package main

import (
	"os"

	"github.com/Hammond95/FattarielloDB/cluster"
)

func main() {

	// argsWithProg := os.Args
	args := os.Args[1:]

	address := args[0]

	n := cluster.NewNode(address)

	n.printInfo()
}
