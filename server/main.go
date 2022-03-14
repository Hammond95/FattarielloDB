package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/Hammond95/FattarielloDB/server/cluster"
	raft "github.com/Hammond95/FattarielloDB/server/cluster/raft"
	//"github.com/Hammond95/FattarielloDB/server/application"
)

var (
	address = flag.String("address", ":8888", "TCP host+port for this node")
	//raftDir       = flag.String("raft_data_dir", "data/", "Raft data dir")
	//raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()
	log := hclog.Default()

	_, port, err := net.SplitHostPort(*address)
	if err != nil {
		log.Error("failed to parse local address (%q): %v", *&address, err)
		os.Exit(1)
	}

	fmt.Println("port used is: " + port)

	ctx := context.Background()

	// create a TCP socket for inbound server connections
	listener, err := net.Listen("tcp", *address)
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("Created listener at %v.", *address))

	raftId := uuid.New()
	fattarielloFSM := &fattariello{}

	r, tm, err := raft.NewRaft(ctx, raftId, *address, fattarielloFSM)
	if err != nil {
		log.Error("failed to start raft: %v", err)
		os.Exit(1)
	}

	// create a new gRPC server, use WithInsecure to allow http connections
	grpcServer := grpc.NewServer()
	log.Info("Created gRPC Server.")

	// create an instance of the Node server
	server := cluster.NewNodeServer(*address, log)
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
