package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	v1 "github.com/Hammond95/FattarielloDB/client/v1"
	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/hashicorp/go-hclog"
)

func main() {
	log := hclog.Default()

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide the address of the server,for the client to connect.")
		os.Exit(1)
	}
	address := args[0]

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewNodeClient(conn)

	g := gin.Default()
	v1.ApplyRouteGroupDefinition(g, client)

	if err := g.Run(":8080"); err != nil {
		log.Error("Failed to run server: %v", err)
	}
}
