package xversion

import "runtime/debug"

var Version string

func init() {
	Version = "unknown version"
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, value := range info.Deps {
			if value.Path == "github.com/coder2z/g-server" {
				Version = value.Version
			}
		}
	}
}
