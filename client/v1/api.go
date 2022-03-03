package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

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
}
