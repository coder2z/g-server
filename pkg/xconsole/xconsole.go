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
	fmt.Println(xcolor.Yellow(fmt.Sprintf("%v\n", msg)))
}

// Redf ...
func Yellowf(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Yellowf(fmt.Sprintf("%-40v", msg), arg))
}

// Red ...
func Red(msg string) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Red(fmt.Sprintf("%v\n", msg)))
}

// Redf ...
func Redf(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Redf(fmt.Sprintf("%-40v", msg), arg))
}

// Blue ...
func Blue(msg string) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Blue(fmt.Sprintf("%v\n", msg)))
}

// Greenf ...
func Bluef(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Bluef(fmt.Sprintf("%-40v", msg), arg))
}

// Green ...
func Green(msg string) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Green(fmt.Sprintf("%v\n", msg)))
}

// Greenf ...
func Greenf(msg string, arg interface{}) {
	if !getDebug() {
		return
	}
	fmt.Println(xcolor.Greenf(fmt.Sprintf("%-40v", msg), arg))
}
