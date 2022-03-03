package cluster

import (
	"context"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
)

type NodeServer struct {
	NodeID         string
	NodeStatus     int32
	NodeAddress    string
	PeersID        []string
	PeersAddresses []string
	logger         hclog.Logger
}

func NewNodeServer(address string, l hclog.Logger) *NodeServer {
	id := uuid.New()
	return &NodeServer{
		id.String(),
		0,
		address,
		[]string{},
		[]string{},
		l,
	}
}

func (server *NodeServer) GetInfo(ctx context.Context, r *proto.EmptyRequest) (*proto.NodeInfos, error) {
	server.logger.Info("Handle request for GetInfo")
	return &proto.NodeInfos{
		NodeID:         server.NodeID,
		NodeStatus:     proto.NodeInfos_NodeState(server.NodeStatus),
		NodeAddress:    server.NodeAddress,
		PeersID:        server.PeersID,
		PeersAddresses: server.PeersAddresses,
	}, nil
}

// NodeActions define what a node can do
//type NodeActions interface {
//	getId() string
//	getStatus() string
//	sendMessage(msg string, receiverAddress string)
//	receiveMessage()
//	PrintInfo()
//	Run()
//}
