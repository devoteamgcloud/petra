package main

import (
	"fmt"
	"os"

	server "github.com/arthur-laurentdka/petra/server/cmd"
)

func main() {
	err := server.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
