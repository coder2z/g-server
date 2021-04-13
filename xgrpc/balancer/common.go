package xbalancer

import (
	"google.golang.org/grpc/resolver"
	"strconv"
)

const (
	WeightKey = "weight"
)

func GetWeight(addr resolver.Address) int {
	if addr.Attributes == nil {
		return 1
	}
	values := addr.Attributes.Value(WeightKey).([]string)
	if len(values) > 0 {
		weight, err := strconv.Atoi(values[0])
		if err == nil {
			return weight
		}
	}
	return 1
}
