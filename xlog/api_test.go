/**
 * @Author: yangon
 * @Description
 * @Date: 2021/3/5 17:35
 **/
package xlog

import (
	"errors"
	"testing"
)

func TestAuto(t *testing.T) {
	Info("test",FieldErr(errors.New("error")))
	Warn("test",FieldErr(errors.New("error")))
	Error("test",FieldErr(errors.New("error")))
}