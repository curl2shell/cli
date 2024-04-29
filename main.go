package main

import (
	"fmt"
	"os"

	"cli/c2b"
)

func main() {
	args := os.Args[1:]

	err := c2b.Run(
		args,
		os.Getenv,
	)
	if err != nil {
		fmt.Println("error", err)
	}
}
