/**
 * @Author: yangon
 * @Description
 * @Date: 2021/3/5 18:06
 **/
package xstring

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	t.Log(GenerateID(),randInstance.Int63())
}
