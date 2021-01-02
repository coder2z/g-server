/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 13:42
 */
package xtrace

import (
	"strings"
)

// MetadataReaderWriter ...
type MetadataReaderWriter struct {
	MD map[string][]string
}

// Set ...
func (w MetadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

// ForeachKey ...
func (w MetadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, md := range w.MD {
		for _, v := range md {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
