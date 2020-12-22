/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/22 19:03
 **/
package xcast

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
)

// ToSliceStringMapE casts an empty interface to a []interface{}.
func ToSliceStringMapE(i interface{}) ([]map[string]interface{}, error) {
	var s = make([]map[string]interface{}, 0)

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s = append(s, cast.ToStringMap(u))
		}
		return s, nil
	case []map[string]interface{}:
		s = append(s, v...)
		return s, nil
	default:
		return s, fmt.Errorf("unable to Cast %#v of type %v to []map[string]interface{}", i, reflect.TypeOf(i))
	}
}

func ToSliceStringMap(i interface{}) []map[string]interface{} {
	v, _ := ToSliceStringMapE(i)
	return v
}
