package xapp

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xcfg"
	"testing"
)

func TestAppVersion(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		[app]
		AppName = "testAppName"
		BuildAppVersion = "v1.1.1"
		`))
	_ = xcfg.LoadFromReader(reader, toml.Unmarshal)
	RegisterAppInfoCfg("app")
	PrintVersion()
}
