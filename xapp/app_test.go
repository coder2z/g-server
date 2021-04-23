package xapp

import (
	"os"
	"testing"
)

func TestAppVersion(t *testing.T) {
	os.Setenv("SERVER_APP_ID","1111")
	PrintVersion()
}
