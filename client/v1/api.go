package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Hammond95/FattarielloDB/proto"
)

func ApplyRouteGroupDefinition(router *gin.Engine, client proto.FattarielloClient) {
	apiV1 := *router.Group("/v1")

	apiV1.GET("/info", func(ctx *gin.Context) {
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

	apiV1.POST("/sendMessage", func(ctx *gin.Context) {
		// Parse JSON
		var json struct {
			Message                string `json:"message" binding:"required"`
			DestinationNodeAddress string `json:"destinationNodeAddress" binding:"required"`
		}

		if ctx.Bind(&json) == nil {
			response, err := client.SendMessage(
				ctx,
				&proto.SendMessageRequest{
					Message:                []byte(json.Message),
					DestinationNodeAddress: json.DestinationNodeAddress,
				},
			)
			if err == nil {
				ctx.JSON(
					http.StatusAccepted,
					gin.H{
						"AckMessage": response.AckMessage,
					},
				)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
	})

	raftGroup := *apiV1.Group("/raft")
	raftGroup.GET("/stats", func(ctx *gin.Context) {
		var stats RaftStats

		req := &proto.EmptyRequest{}

		if response, err := client.RaftStats(ctx, req); err == nil {
			json.Unmarshal([]byte(response.Body), &stats)
			ctx.JSON(http.StatusOK, stats)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	raftGroup.POST("/join", func(ctx *gin.Context) {
		var json struct {
			NodeID      string `json:"nodeID" binding:"required"`
			NodeAddress string `json:"nodeAddress" binding:"required"`
		}

		if ctx.Bind(&json) == nil {
			response, err := client.RaftJoin(
				ctx,
				&proto.NodeMinimalInfos{
					NodeID:      json.NodeID,
					NodeAddress: json.NodeAddress,
				},
			)
			if err == nil {
				ctx.JSON(
					http.StatusAccepted,
					gin.H{
						"AckMessage": response.AckMessage,
					},
				)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
	})

	raftGroup.POST("/remove", func(ctx *gin.Context) {
		var json struct {
			NodeID      string `json:"nodeID" binding:"required"`
			NodeAddress string `json:"nodeAddress" binding:"required"`
		}

		if ctx.Bind(&json) == nil {
			response, err := client.RaftRemove(
				ctx,
				&proto.NodeMinimalInfos{
					NodeID:      json.NodeID,
					NodeAddress: json.NodeAddress,
				},
			)
			if err == nil {
				ctx.JSON(
					http.StatusAccepted,
					gin.H{
						"AckMessage": response.AckMessage,
					},
				)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
	})
}
