/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 15:27
 */
package xreminder

import (
	"github.com/robfig/cron/v3"
)

type options struct {
	Spec []string `json:"spec,omitempty" yaml:"spec"`
}

func NewReminderOptions() *options {
	return &options{
		Spec: []string{"* * * * * *"},
	}
}

type cronServer struct {
	o *options
	c *cron.Cron
}

func (r *cronServer) Run(stopCh <-chan struct{}, f cron.Job) {
	for _, v := range r.o.Spec {
		_, _ = r.c.AddJob(v, f)
	}
	r.c.Start()
	<-stopCh
	r.c.Stop()
}

func NewReminderClient(o *options) *cronServer {
	return &cronServer{c: cron.New(), o: o}
}
