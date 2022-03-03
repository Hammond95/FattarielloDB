package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

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
	g.GET("/info", func(ctx *gin.Context) {
		req := &proto.EmptyRequest{}
		if response, err := client.GetInfo(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"NodeID":         fmt.Sprint(response.NodeID),
				"NodeStatus":     fmt.Sprint(response.NodeStatus.Enum()),
				"NodeAddress":    fmt.Sprint(response.NodeAddress),
				"PeersID":        fmt.Sprint(response.PeersID),
				"PeersAddresses": fmt.Sprint(response.PeersAddresses),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	if err := g.Run(":8080"); err != nil {
		log.Error("Failed to run server: %v", err)
	}
}
