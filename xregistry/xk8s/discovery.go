package xk8s

import (
	"context"
	"fmt"
	"github.com/coder2z/g-server/xregistry"
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xstring"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"sync"
)

type k8sDiscovery struct {
	clients   *kubernetes.Clientset
	namespace string

	closeOnce sync.Once
	closeCh   chan struct{}
}

func newDiscovery(namespace string) (xregistry.Discovery, error) {
	if namespace == "" {
		namespace = "default"
	}
	conf, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clients, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}
	return &k8sDiscovery{
		clients:   clients,
		namespace: namespace,
		closeCh:   make(chan struct{}),
	}, nil
}

// target: service-name:port or service-name:port-name
func (d *k8sDiscovery) Discover(target string) (<-chan []xregistry.Instance, error) {
	service, port := parse(target)
	if service == "" || port == "" {
		return nil, fmt.Errorf("target not valid: %s", target)
	}
	ch := make(chan []xregistry.Instance)
	return ch, d.watch(ch, service, port)
}

func parse(target string) (service, port string) {
	ss := strings.Split(target, ":")
	if len(ss) == 2 {
		service, port = ss[0], ss[1]
	}
	return
}

func (d *k8sDiscovery) watch(ch chan<- []xregistry.Instance, service, port string) error {
	watcher, err := d.clients.CoreV1().Endpoints(d.namespace).
		Watch(context.Background(), metaV1.ListOptions{FieldSelector: fmt.Sprintf("%s=%s", "metadata.name", service)})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-d.closeCh:
				return
			case <-watcher.ResultChan():
				endpoints, err := d.clients.CoreV1().Endpoints(d.namespace).
					List(context.Background(), metaV1.ListOptions{FieldSelector: fmt.Sprintf("%s=%s", "metadata.name", service)})
				if err != nil {
					continue
				}
				var i []xregistry.Instance
				for _, endpoint := range endpoints.Items {
					for _, subset := range endpoint.Subsets {
						realPort := port
						for _, p := range subset.Ports {
							if p.Name == port {
								realPort = fmt.Sprint(p.Port)
								break
							}
						}
						for _, addr := range subset.Addresses {
							ins := xregistry.Instance{Address: fmt.Sprintf("%s:%s", addr.IP, realPort)}
							i = append(i, ins)
						}
						xlog.Info("service discovery k8s", xlog.FieldValue(xstring.Json(i)))
					}
				}
				ch <- i
			}
		}
	}()
	return nil
}

func (d *k8sDiscovery) Close() {
	d.closeOnce.Do(func() {
		xconsole.Red("service registration shutdown")
		close(d.closeCh)
	})
}
