package slice_plus

import (
	"reflect"
	"sort"
	"strconv"
)

func Integer2String(i interface{}) string {
	switch i.(type) {
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(i)
		return strconv.FormatInt(v.Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(i)
		return strconv.FormatUint(v.Uint(), 10)
	default:
		return ""
	}
}

func SliceInteger2String(arr ...interface{}) []string {
	s := make([]string, len(arr))

	for i, n := range arr {
		s[i] = Integer2String(n)
	}
	return s
}

// StrInSlice 判断string是否在slice中
func StrInSlice(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// IntInSlice 判断int是否在slice中
func IntInSlice(num int, arr []int) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

// Int64InSlice 判断int是否在slice中
func Int64InSlice(num int64, arr []int64) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

// Interface2String interface的slice转string的slice
func Interface2String(arr []interface{}) []string {
	length := len(arr)
	strSlice := make([]string, length)
	if length == 0 {
		return strSlice
	}
	for k, v := range arr {
		strSlice[k] = v.(string)
	}
	return strSlice
}

// StringToInt string的slice转int
func StringToInt(arr []string) []int {
	length := len(arr)
	intSlice := make([]int, length)
	if length == 0 {
		return intSlice
	}
	for k, v := range arr {
		var ok error
		intSlice[k], ok = strconv.Atoi(v)
		if ok != nil {
			intSlice[k] = 0
		}
	}
	return intSlice
}

// StringToUint32 string的slice转uint32
func StringToUint32(arr []string) []uint32 {
	length := len(arr)
	intSlice := make([]uint32, length)
	if length == 0 {
		return intSlice
	}
	for k, v := range arr {
		tmp, _ := strconv.ParseInt(v, 10, 32)
		intSlice[k] = uint32(tmp)
	}
	return intSlice
}

func StringToUint64(arr []string) []uint64 {
	n := make([]uint64, len(arr))
	for i, v := range arr {
		tmp, _ := strconv.ParseUint(v, 10, 64)
		n[i] = tmp
	}
	return n
}

func Uint32ToInt32(arr []uint32) []int32 {
	intSlice := make([]int32, len(arr))
	for i, v := range arr {
		intSlice[i] = int32(v)
	}
	return intSlice
}

func StrSliceRemove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func ReverseStringSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// 用于uid去重
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func RemoveRep(slice []uint32) []uint32 {
	slc := make([]uint32, len(slice))
	copy(slc, slice)
	if len(slc) <= 1 {
		return slc
	}
	sort.Sort(Uint32Slice(slc))

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func RemoveRepUint64(slc []uint64) []uint64 {
	if len(slc) <= 1 {
		return slc
	}
	sort.Sort(Uint64Slice(slc))

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

func RemoveRepInt(slc []int) []int {
	if len(slc) <= 1 {
		return slc
	}
	sort.Ints(slc)

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

func ReOrderSlice(org []uint32, dst []uint32) {
	if len(org) == 0 || len(dst) == 0 {
		return
	}

	var idx int
	for _, did := range dst {
		for i := idx; i < len(org); i++ {
			if did == org[i] {
				org[idx], org[i] = org[i], org[idx]
				idx++
				break
			}
		}
	}
}

// SliceUint64ToStr  []uint64 to []string
func SliceUint64ToStr(arr ...uint64) []string {
	s := make([]string, len(arr))

	for i, n := range arr {
		s[i] = Integer2String(n)
	}
	return s
}

// SliceUint32ToStr ...
func SliceUint32ToStr(arr ...uint32) []string {
	s := make([]string, len(arr))

	for i, n := range arr {
		s[i] = Integer2String(n)
	}
	return s
}
