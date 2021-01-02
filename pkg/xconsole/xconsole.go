package xconsole

import (
	"fmt"
	"github.com/myxy99/component/pkg/xcolor"
)

// Yellow ...
func Yellow(msg string) {
	fmt.Println(xcolor.Yellow(fmt.Sprintf("%v\n", msg)))
}

// Redf ...
func Yellowf(msg string, arg interface{}) {
	fmt.Println(xcolor.Yellowf(fmt.Sprintf("%-40v", msg), arg))
}

// Red ...
func Red(msg string) {
	fmt.Println(xcolor.Red(fmt.Sprintf("%v\n", msg)))
}

// Redf ...
func Redf(msg string, arg interface{}) {
	fmt.Println(xcolor.Redf(fmt.Sprintf("%-40v", msg), arg))
}

// Blue ...
func Blue(msg string) {
	fmt.Println(xcolor.Blue(fmt.Sprintf("%v\n", msg)))
}

// Greenf ...
func Bluef(msg string, arg interface{}) {
	fmt.Println(xcolor.Bluef(fmt.Sprintf("%-40v", msg), arg))
}

// Green ...
func Green(msg string) {
	fmt.Println(xcolor.Green(fmt.Sprintf("%v\n", msg)))
}

// Greenf ...
func Greenf(msg string, arg interface{}) {
	fmt.Println(xcolor.Greenf(fmt.Sprintf("%-40v", msg), arg))
}
