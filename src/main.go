package main

import (
	"os"

	"github.com/zzzming/blockchain/src/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	defer os.Exit(0)

	defaultConfigFile := "../config/blockchain.yaml"
	configFile := cmd.AssignString(os.Getenv("BLOCKCHAIN_CONFIG"), defaultConfigFile)
	cmd, err := cmd.NewCommandLine(configFile)
	if err != nil {
		panic(err)
	}
	cmd.Run()
}
