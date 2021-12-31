package main

import (
	"os"

	"github.com/zzzming/mbt/src/cmd"
)

func main() {
	defer os.Exit(0)

	cmd := cmd.CommandLine{}
	cmd.Run()
}
