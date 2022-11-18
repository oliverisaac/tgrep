package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	fmt.Printf("Hello world from %s\n", path.Base(os.Args[0]))
}
