package xcode

import (
	"github.com/coder2z/g-saber/xjson"
	"github.com/coder2z/g-saber/xlog"
	"net/http"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

// CodeBreakUp 低于10000均为系统错误码，业务错误码请使用10000以上
const (
	CodeBreakUp uint32 = 9999
	SystemType         = iota
	BusinessType
)

var (
	aid            uint32
	_codesSystem   sync.Map
	_codesBusiness sync.Map
	// OK ...
	OK      = add(SystemType, uint32(codes.OK), "OK")
	Unknown = add(SystemType, uint32(codes.Unknown), "UNKNOWN")
)

type CodeInfo struct {
	CodeT   uint
	Code    uint32
	Message string
}

func XCodeSystemCodeHttp(w http.ResponseWriter, r *http.Request) {
	var res = make([]*spbStatus, 0)
	_codesSystem.Range(func(key, val interface{}) bool {
		res = append(res, val.(*spbStatus))
		return true
	})
	_ = xjson.NewEncoder(w).Encode(res)
}

func XCodeBusinessCodeHttp(w http.ResponseWriter, r *http.Request) {
	var res = make([]*spbStatus, 0)
	_codesBusiness.Range(func(key, val interface{}) bool {
		res = append(res, val.(*spbStatus))
		return true
	})
	_ = xjson.NewEncoder(w).Encode(res)
}

func SystemCode(code uint32) (spb *spbStatus) {
	_codesSystem.Range(func(key, val interface{}) bool {
		if code == key.(uint32) {
			spb = val.(*spbStatus)
			return false
		}
		return true
	})
	if spb == nil {
		spb = Unknown
	}
	return spb
}

func BusinessCode(code uint32) (spb *spbStatus) {
	_codesBusiness.Range(func(key, val interface{}) bool {
		if code == key.(uint32) {
			spb = val.(*spbStatus)
			return false
		}
		return true
	})
	if spb == nil {
		spb = Unknown
	}
	return spb
}

func SystemCodeAdd(code uint32, message string) *spbStatus {
	if code > CodeBreakUp {
		xlog.Panic("customize code must less than 9999")
	}

	return add(SystemType, aid*10000+code, message)
}

func BusinessCodeAdd(code uint32, message string) *spbStatus {
	if code < CodeBreakUp {
		xlog.Panic("customize code must less than 9999")
	}
	return add(BusinessType, code, message)
}

func CodeAdds(data []CodeInfo) {
	for _, datum := range data {
		_ = add(datum.CodeT, datum.Code, datum.Message)
	}
}

func add(codeT uint, code uint32, message string) *spbStatus {
	s := &spbStatus{
		&spb.Status{
			Code:    int32(code),
			Message: message,
			Details: make([]*any.Any, 0),
		},
	}
	if codeT == SystemType {
		_codesSystem.Store(code, s)
	}
	if codeT == BusinessType {
		_codesBusiness.Store(code, s)
	}
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
