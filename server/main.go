package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/Hammond95/FattarielloDB/server/cluster"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide an address for the node to run.")
		os.Exit(1)
	}

	log := hclog.Default()

	// create a TCP socket for inbound server connections
	address := args[0]
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("Created listener at %v.", address))

	// create a new gRPC server, use WithInsecure to allow http connections
	grpcServer := grpc.NewServer()
	log.Info("Created gRPC Server.")

	// create an instance of the Node server
	server := cluster.NewNodeServer(address, log)
	log.Info("Created NodeServer.")

	// register the Node Server
	proto.RegisterNodeServer(grpcServer, server)
	log.Info("Registered NodeServer.")

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)
	log.Info("Registered reflection.")

	if e := grpcServer.Serve(listener); e != nil {
		panic(e)
	}
}
