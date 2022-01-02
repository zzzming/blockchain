package main

import (
	"os"

	"github.com/zzzming/mbt/src/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	defer os.Exit(0)

	defaultConfigFile := "../config/mbt.yaml"
	configFile := cmd.AssignString(os.Getenv("BLOCKCHAIN_CONFIG"), defaultConfigFile)
	cmd, err := cmd.NewCommandLine(configFile)
	if err != nil {
		panic(err)
	}
	cmd.Run()
}
