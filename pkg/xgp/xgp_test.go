package xgp

import (
	"github.com/myxy99/component/pkg/xconsole"
	"testing"
	"time"
)

func TestXgp(t *testing.T) {
	pool := NewGoPool(3, 5*time.Minute)
	pool.Go(func() {
		xconsole.Red("test1")
		time.Sleep(5 * time.Second)
		xconsole.Red("test1 结束")
	}, func(err error) {

	})
	pool.Go(func() {
		xconsole.Red("test2")
		select {}
	}, func(err error) {

	})
	pool.Go(func() {
		xconsole.Red("test3")
		select {}
	}, func(err error) {

	})
	time.Sleep(6 * time.Second)
	xconsole.Red("test4 加入")
	pool.Go(func() {
		xconsole.Red("test4")
		select {}
	}, func(err error) {

	})
	xconsole.Red("test5 加入")
	pool.Go(func() {
		xconsole.Red("test5")
		select {}
	}, func(err error) {

	})


	select {

	}
}
