/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 12:55
 */
package xcode

import (
	"encoding/json"
	"github.com/myxy99/component/xgovern"
	"github.com/myxy99/component/xlog"
	"net/http"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

// EcodeNum 低于10000均为系统错误码，业务错误码请使用10000以上
const EcodeNum int32 = 9999

var (
	aid              int
	maxCustomizeCode = 9999
	_codes           sync.Map
	// OK ...
	OK = add(int(codes.OK), "OK")
)

func init() {
	xgovern.HandleFunc("/status/code/list", func(w http.ResponseWriter, r *http.Request) {
		var res = make(map[int]*spbStatus)
		_codes.Range(func(key, val interface{}) bool {
			code := key.(int)
			res[code] = val.(*spbStatus)
			return true
		})
		_ = json.NewEncoder(w).Encode(res)
	})
}

// Add ...
func Add(code int, message string) *spbStatus {
	if code > maxCustomizeCode {
		xlog.Panic("customize code must less than 9999", xlog.Any("code", code))
	}

	return add(aid*10000+code, message)
}

func add(code int, message string) *spbStatus {
	s := &spbStatus{
		&spb.Status{
			Code:    int32(code),
			Message: message,
			Details: make([]*any.Any, 0),
		},
	}
	_codes.Store(code, s)
	return s
}

// ExtractCodes cause from error to ecode.
func ExtractCodes(e error) *spbStatus {
	if e == nil {
		return OK
	}
	gst, _ := status.FromError(e)
	return &spbStatus{
		&spb.Status{
			Code:    int32(gst.Code()),
			Message: gst.Message(),
			Details: make([]*any.Any, 0),
		},
	}
}
