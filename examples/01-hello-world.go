// Our first program prints the classic "hello world" message.
// Here's the full source code.

package main

import "fmt"

func main() {
	fmt.Println("hello world")
}

// To run the program, put the code in a file called hello-world.go
// and use go run.

// $ go run hello-world.go
// hello world

// Sometimes we'll want to build our programs into binaries.
// We can do this using go build.

// $ go build hello-world.go
// $ ls
// hello-world	hello-world.go
// $ ./hello-world
// hello world

// We've just run our first Go program. Now let's move on to
// learning a bit more about the language.
