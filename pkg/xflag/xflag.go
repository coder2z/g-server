/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/12/23 23:24
 */
package xflag

import (
	"fmt"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

type (
	Command         = cobra.Command
	FlagRegisterFun func(*Command)
	CommandNode     struct {
		Name    string
		Command *Command
		Flags   FlagRegisterFun
	}
	storage struct {
		instances sync.Map
	}
)

var (
	defaultFlags = &Command{
		Use:                "app",
		DisableSuggestions: false,
	}
	bucket = &storage{
		instances: sync.Map{},
	}
)

func NewRootCommand(c *Command) {
	defaultFlags = c
}

func Register(fs ...CommandNode) {
	for _, c := range fs {
		c.Flags(c.Command)
		defaultFlags.AddCommand(c.Command)
		bucket.instances.Store(c.Name, c.Command)
	}
}

func Parse() error {
	return defaultFlags.Execute()
}

func NString(nodeName, name string) string {
	v, _ := NStringE(nodeName, name)
	return v
}

func NBool(nodeName, name string) bool {
	v, _ := NBoolE(nodeName, name)
	return v
}

func NTimeDuration(nodeName, name string) time.Duration {
	v, _ := NTimeDurationE(nodeName, name)
	return v
}

func NFloat32(nodeName, name string) float32 {
	v, _ := NFloat32E(nodeName, name)
	return v
}

func NFloat64(nodeName, name string) float64 {
	v, _ := NFloat64E(nodeName, name)
	return v
}

func NInt(nodeName, name string) int {
	v, _ := NIntE(nodeName, name)
	return v
}

func NStringSlice(nodeName, name string) []string {
	v, _ := NStringSliceE(nodeName, name)
	return v
}

func String(name string) string {
	v, _ := StringE(name)
	return v
}

func Bool(name string) bool {
	v, _ := BoolE(name)
	return v
}

func TimeDuration(name string) time.Duration {
	v, _ := TimeDurationE(name)
	return v
}

func Float32(name string) float32 {
	v, _ := Float32E(name)
	return v
}

func Float64(name string) float64 {
	v, _ := Float64E(name)
	return v
}

func Int(name string) int {
	v, _ := IntE(name)
	return v
}

func StringSlice(name string) []string {
	v, _ := StringSliceE(name)
	return v
}

func NStringE(nodeName, name string) (string, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetString(name)
	}
	return "", fmt.Errorf("undefined flag name: %s", name)
}

func NBoolE(nodeName, name string) (bool, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetBool(name)
	}
	return false, fmt.Errorf("undefined flag name: %s", name)
}

func NTimeDurationE(nodeName, name string) (time.Duration, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetDuration(name)
	}
	return 0, fmt.Errorf("undefined flag name: %s", name)
}

func NFloat32E(nodeName, name string) (float32, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetFloat32(name)
	}
	return 0, fmt.Errorf("undefined flag name: %s", name)
}

func NFloat64E(nodeName, name string) (float64, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetFloat64(name)
	}
	return 0, fmt.Errorf("undefined flag name: %s", name)
}

func NIntE(nodeName, name string) (int, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetInt(name)
	}
	return 0, fmt.Errorf("undefined flag name: %s", name)
}

func NStringSliceE(nodeName, name string) ([]string, error) {
	if val, ok := bucket.instances.Load(nodeName); ok {
		return val.(*Command).Flags().GetStringSlice(name)
	}
	return nil, fmt.Errorf("undefined flag name: %s", name)
}

func StringE(name string) (val string, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetString(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}

func BoolE(name string) (val bool, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetBool(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}

func TimeDurationE(name string) (val time.Duration, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetDuration(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}

func Float32E(name string) (val float32, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetFloat32(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}

func Float64E(name string) (val float64, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetFloat64(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}

func IntE(name string) (val int, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetInt(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}

func StringSliceE(name string) (val []string, err error) {
	var ok bool
	bucket.instances.Range(func(key, value interface{}) bool {
		val, err = value.(*Command).Flags().GetStringSlice(name)
		ok = err == nil
		return err != nil
	})
	if !ok {
		err = fmt.Errorf("undefined flag name: %s", name)
	}
	return
}
