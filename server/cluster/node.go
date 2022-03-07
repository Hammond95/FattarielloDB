package cluster

import (
	"context"
	"fmt"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
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

func (server *NodeServer) SendMessage(ctx context.Context, r *proto.SendMessageRequest) (*proto.AckResponse, error) {
	server.logger.Info("Sending message to " + fmt.Sprint(r.DestinationNodeAddress))

	conn, err := grpc.Dial(r.DestinationNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewNodeClient(conn)

	ackResponse, err := client.ReceiveMessage(
		ctx,
		&proto.ReceiveMessageRequest{
			SenderNodeID:      server.NodeID,
			SenderNodeAddress: server.NodeAddress,
			Message:           r.Message,
		},
	)

	return ackResponse, err
}

func (server *NodeServer) ReceiveMessage(ctx context.Context, r *proto.ReceiveMessageRequest) (*proto.AckResponse, error) {
	server.logger.Info("Receiving message from " + fmt.Sprint(r.SenderNodeID))

	// TODO: Write message in storage
	server.logger.Info("[FAKE] Wrote message in storage.")

	ackMessage := "[FAKE] Wrote message in storage."
	var err error = nil

	return &proto.AckResponse{AckMessage: ackMessage}, err
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
