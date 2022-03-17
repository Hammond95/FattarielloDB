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
)

var (
	address = flag.String("address", ":8888", "TCP host+port for this node")
	//raftAddress = flag.String("raftAddress", ":7777", "TCP host+port for the raft server.")
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

	//raftListener, err := net.Listen("tcp", *raftAddress)
	//if err != nil {
	//	panic(err)
	//}
	//log.Info(fmt.Sprintf("Created raftListener at %v.", *raftAddress))

	//fmt.Printf("raftListener = %v\n", raftListener)

	raftId := uuid.New()

	// create an instance of the Fattariello server
	server := NewFattarielloServer(raftId.String(), *address, log)
	log.Info("Created FattarielloServer.")

	fattarielloFSM := &FattarielloDB{}

	raftDoh, tm, err := NewRaft(ctx, raftId.String(), *address, fattarielloFSM)
	if err != nil {
		log.Error("failed to start raft: %v", err)
		os.Exit(1)
	}

	fmt.Printf("RAFT OBJ: %v, TransportManager: %v", raftDoh, tm)

	// create a new gRPC server, use WithInsecure to allow http connections
	grpcServer := grpc.NewServer()
	log.Info("Created gRPC Server.")

	// register the Node Server
	proto.RegisterFattarielloServer(
		grpcServer,
		&RpcInterface{
			fattarielloServer: server,
			raft:              raftDoh,
		})
	log.Info("Registered Fattariello Server.")

	tm.Register(grpcServer)
	log.Info("Registered TransportManager to gRPC Server.")
	//leaderhealth.Setup(r, s, []string{"Example"})
	// == grpcServer.RegisterService(&raftDoh, interface for raft with all gRPC methods)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)
	log.Info("Registered reflection.")

	if e := grpcServer.Serve(listener); e != nil {
		panic(e)
	}
}
