// +build windows

package xsignals

import "os"

var shutdownSignals = []os.Signal{os.Interrupt}
