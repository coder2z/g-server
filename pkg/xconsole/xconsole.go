package xconsole

import (
	"fmt"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/pkg/xcolor"
	"os"
)

var (
	debug bool
	is    bool
)

func getDebug() bool {
	if !is {
		debug = xcast.ToBool(os.Getenv("app.debug"))
		is = true
	}
	return debug
}

// Yellow ...
func Yellow(msg string) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Yellow(fmt.Sprintf("%v\n", msg)))
}

// Redf ...
func Yellowf(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Yellowf(fmt.Sprintf("%-40v", msg), arg))
}

// Red ...
func Red(msg string) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Red(fmt.Sprintf("%v\n", msg)))
}

// Redf ...
func Redf(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Redf(fmt.Sprintf("%-40v", msg), arg))
}

// Blue ...
func Blue(msg string) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Blue(fmt.Sprintf("%v\n", msg)))
}

// Greenf ...
func Bluef(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Bluef(fmt.Sprintf("%-40v", msg), arg))
}

// Green ...
func Green(msg string) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Green(fmt.Sprintf("%v\n", msg)))
}

// Greenf ...
func Greenf(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Print(xcolor.Greenf(fmt.Sprintf("%-40v", msg), arg))
}
