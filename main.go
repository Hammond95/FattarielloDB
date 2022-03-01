package main

import (
	"fmt"
	"os"

	"github.com/Hammond95/FattarielloDB/cluster"
)

func main() {

	// argsWithProg := os.Args
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide an address for the node to run.")
		os.Exit(1)
	}

	address := args[0]
	n := cluster.NewNode(address)
	n.PrintInfo()

	n.Run()
}
