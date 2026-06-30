// In Go it's idiomatic to communicate errors via an
// explicit, separate return value. This contrasts with
// the exceptions used in languages like Java and Ruby
// and the overloaded single result / error value
// sometimes used in C.

package main

import "errors"
import "fmt"

// By convention, errors are the last return value and
// have type error, a built-in interface.
func f(arg int) (int, error) {
	if arg == 42 {
		// errors.New constructs a basic error value
		// with the given error message.
		return -1, errors.New("can't work with 42")
	}

	// A nil value in the error position indicates that
	// there was no error.
	return arg + 3, nil
}

func main() {
	// The two loops below test out each of our
	// error-returning functions. Note that the use of
	// an inline error check on the if line is a common
	// idiom in Go code.
	for _, i := range []int{7, 42} {
		if r, err := f(i); err != nil {
			fmt.Println("f failed:", err)
		} else {
			fmt.Println("f worked:", r)
		}
	}
}
