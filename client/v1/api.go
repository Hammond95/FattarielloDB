package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Hammond95/FattarielloDB/proto"
)

func ApplyRouteGroupDefinition(router *gin.Engine, client proto.NodeClient) {
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
}
