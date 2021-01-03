package xjson

import "testing"

type Info struct {
	Name string `json:"name"`
}

func TestJson(t *testing.T) {
	var info = Info{
		Name: "test_name",
	}
	data, err := Marshal(info)
	t.Log(data, err)

	var info2 Info
	err = Unmarshal(data, &info2)
	t.Log(info2, err)
}
