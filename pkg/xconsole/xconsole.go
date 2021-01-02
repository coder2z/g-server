package xconsole

import (
	"fmt"
	"github.com/myxy99/component/pkg/xcolor"
)

// Yellow ...
func Yellow(msg string) {
	fmt.Println(xcolor.Yellow(msg))
}

// Redf ...
func Yellowf(msg string, arg interface{}) {
	fmt.Println(xcolor.Yellowf(fmt.Sprintf("%v\t\t\t", msg), arg))
}

// Red ...
func Red(msg string) {
	fmt.Println(xcolor.Red(msg))
}

// Redf ...
func Redf(msg string, arg interface{}) {
	fmt.Println(xcolor.Redf(fmt.Sprintf("%v\t\t\t", msg), arg))
}

// Blue ...
func Blue(msg string) {
	fmt.Println(xcolor.Blue(msg))
}

// Greenf ...
func Bluef(msg string, arg interface{}) {
	fmt.Println(xcolor.Bluef(fmt.Sprintf("%v\t\t\t", msg), arg))
}

// Green ...
func Green(msg string) {
	fmt.Println(xcolor.Green(msg))
}

// Greenf ...
func Greenf(msg string, arg interface{}) {
	fmt.Println(xcolor.Greenf(fmt.Sprintf("%v\t\t\t", msg), arg))
}
