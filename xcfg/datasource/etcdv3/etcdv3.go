package etcdv3

import (
	"context"
	"github.com/coder2z/component/xcfg"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/pkg/errors"
)

type etcdDataSource struct {
	propertyKey         string
	lastUpdatedRevision int64
	client              *clientv3.Client
	// cancel is the func, call cancel will stop watching on the propertyKey
	cancel context.CancelFunc
	// closed indicate whether continuing to watch on the propertyKey
	// closed util.AtomicBool

	changed chan struct{}
}

// NewDataSource new a etcdDataSource instance.
// client is the etcd client, it must be useful and should be release by User.
func NewDataSource(client *clientv3.Client, key string,watch bool) xcfg.DataSource {
	ds := &etcdDataSource{
		client:      client,
		propertyKey: key,
	}
	if watch {
		go ds.watch()
	}
	return ds
}

// ReadConfig ...
func (s *etcdDataSource) ReadConfig() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.client.Get(ctx, s.propertyKey)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, errors.New("empty response")
	}
	s.lastUpdatedRevision = resp.Header.GetRevision()
	return resp.Kvs[0].Value, nil
}

// IsConfigChanged ...
func (s *etcdDataSource) IsConfigChanged() <-chan struct{} {
	return s.changed
}

func (s *etcdDataSource) handle(resp *clientv3.WatchResponse) {
	if resp.CompactRevision > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.CompactRevision
	}
	if resp.Header.GetRevision() > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.Header.GetRevision()
	}

	if err := resp.Err(); err != nil {
		return
	}

	for _, ev := range resp.Events {
		if ev.Type == mvccpb.PUT || ev.Type == mvccpb.DELETE {
			select {
			case s.changed <- struct{}{}:
			default:
			}
		}
	}
}

func (s *etcdDataSource) watch() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	rch := s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))
	for {
		for resp := range rch {
			s.handle(&resp)
		}
		time.Sleep(time.Second)

		ctx, cancel = context.WithCancel(context.Background())
		if s.lastUpdatedRevision > 0 {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))
		} else {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify())
		}
		s.cancel = cancel
	}
}

// Close ...
func (s *etcdDataSource) Close() error {
	s.cancel()
	return nil
}
