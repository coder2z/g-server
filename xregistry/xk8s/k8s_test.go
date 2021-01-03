package xk8s

import (
	"google.golang.org/grpc"
	"testing"
)

func TestK8s(t *testing.T) {

	err := RegisterBuilder() //发现
	if err != nil {
		t.Failed()
		return
	}

	conn, err := grpc.Dial("k8s://namespaces/servicename:portname")
	t.Log(conn)
}

