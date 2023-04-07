package main

import (
	"fmt"
	"os"
)

func main() {
	err := Run()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
