// switch statements express conditionals across many
// branches.

package main

import "fmt"

func main() {
	// A basic switch.
	i := 2
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	}

	// You can use commas to separate multiple
	// expressions in the same case statement.
	switch time := 9; {
	case time < 12:
		fmt.Println("it's before noon")
	default:
		fmt.Println("it's after noon")
	}

	// switch without an expression is an alternate way
	// to express if/else logic.
	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm an int")
		default:
			fmt.Printf("Don't know type %T\n", t)
		}
	}
	whatAmI(true)
	whatAmI(1)
	whatAmI("hey")
}
