/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 12:56
 */
package xcode

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/myxy99/component/pkg/xjson"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"reflect"
)

type spbStatus struct {
	*spb.Status
}

// GetCodeAsInt ...
func (s *spbStatus) GetCodeAsInt() int {
	return int(s.Code)
}

// GetCodeAsUint32 ...
func (s *spbStatus) GetCodeAsUint32() uint32 {
	return uint32(s.Code)
}

func (s *spbStatus) Error() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", codes.Code(s.GetCode()), s.GetMessage())
}

// GetCodeAsBool ...
func (s *spbStatus) IsOk() bool {
	return s.CauseCode() == 0
}

// GetMessage ...
func (s *spbStatus) GetMessage(exts ...interface{}) string {
	if len(exts)%2 != 0 {
		panic("parameter must be odd")
	}

	var buf bytes.Buffer
	buf.WriteString(s.Message)

	if len(exts) > 0 {
		buf.WriteByte(',')
	}
	for i := 0; i < len(exts); i++ {
		buf.WriteString(fmt.Sprintf("%v", exts[i]))
		buf.WriteByte(':')
		buf.WriteString(fmt.Sprintf("%v", exts[i+1]))
		i++
	}
	return buf.String()
}

func (s *spbStatus) SetMsg(msg string) *spbStatus {
	s.Message = msg
	return s
}

func (s *spbStatus) SetMsgf(format string, age ...interface{}) *spbStatus {
	s.Message = fmt.Sprintf(format, age)
	return s
}

// GetDetailMessage ...
func (s *spbStatus) GetDetailMessage(exts ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString(s.GetMessage(exts...))
	for _, detail := range s.Details {
		buf.WriteByte('\n')
		buf.WriteString(detail.String())
	}
	return buf.String()
}

// String ...
func (s *spbStatus) String() string {
	bs, _ := xjson.Marshal(s)
	return string(bs)
}

// CauseCode ...
func (s *spbStatus) CauseCode() int {
	return int(s.Code)
}

// Proto ...
func (s *spbStatus) Proto() *spb.Status {
	if s == nil {
		return nil
	}
	return proto.Clone(s.Status).(*spb.Status)
}

// MustWithDetails ...
func (s *spbStatus) MustWithDetails(details ...interface{}) *spbStatus {
	status, err := s.WithDetails(details...)
	if err != nil {
		panic(err)
	}
	return status
}

// WithDetails returns a new status with the provided details messages appended to the status.
// If any errors are encountered, it returns nil and the first error encountered.
func (s *spbStatus) WithDetails(details ...interface{}) (*spbStatus, error) {
	if s.CauseCode() == 0 {
		return nil, errors.New("no error details for status with code OK")
	}
	p := s.Proto()
	for _, detail := range details {
		if pms, ok := detail.(proto.Message); ok {
			message, err := marshalAnyProtoMessage(pms)
			if err != nil {
				return nil, err
			}
			p.Details = append(p.Details, message)
		} else {
			a, err := marshalAny(detail)
			if err != nil {
				return nil, err
			}
			p.Details = append(p.Details, a)
		}
	}
	return &spbStatus{Status: p}, nil
}

func marshalAny(obj interface{}) (*any.Any, error) {
	typ := reflect.TypeOf(obj)
	val := fmt.Sprintf("%+v", obj)

	return &any.Any{TypeUrl: typ.Name(), Value: []byte(val)}, nil
}

func marshalAnyProtoMessage(pb proto.Message) (*any.Any, error) {
	value, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return &any.Any{TypeUrl: proto.MarshalTextString(pb), Value: value}, nil
}
