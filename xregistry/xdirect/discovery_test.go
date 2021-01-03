package xdirect

import (
	"google.golang.org/grpc"
	"testing"
)

func TestDirect(t *testing.T) {

	err := RegisterBuilder() //发现
	if err != nil {
		t.Failed()
		return
	}

	conn, err := grpc.Dial("direct://namespaces/127.0.0.1:8000,127.0.0.1:8001")

	t.Log(conn)
}
