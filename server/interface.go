package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/hashicorp/raft"
)

type RpcInterface struct {
	fattarielloServer *FattarielloServer
	raft              *raft.Raft
}

func (x RpcInterface) GetInfo(ctx context.Context, r *proto.EmptyRequest) (*proto.NodeInfos, error) {
	x.fattarielloServer.logger.Info("Running 'func (x RpcInterface) GetInfo'")
	x.fattarielloServer.UpdateNodeStatus(x.raft)
	return x.fattarielloServer.GetInfo(ctx, r)
}

func (x RpcInterface) SendMessage(ctx context.Context, r *proto.SendMessageRequest) (*proto.AckResponse, error) {
	return x.fattarielloServer.SendMessage(ctx, r)
}

func (x RpcInterface) ReceiveMessage(ctx context.Context, r *proto.ReceiveMessageRequest) (*proto.AckResponse, error) {
	return x.fattarielloServer.ReceiveMessage(ctx, r)
}

func (x RpcInterface) RaftStats(ctx context.Context, r *proto.EmptyRequest) (*proto.JSONResponse, error) {
	raftStats := x.raft.Stats()

	statsJSONString, err := json.Marshal(raftStats)

	return &proto.JSONResponse{Body: string(statsJSONString)}, err
}

func (x RpcInterface) RaftJoin(ctx context.Context, r *proto.NodeMinimalInfos) (*proto.AckResponse, error) {
	err := Join(*x.raft, r.NodeID, r.NodeAddress)
	ackMessage := fmt.Sprintf("Node %s running at %s, successfully joined to raft cluster!", r.NodeID, r.NodeAddress)
	return &proto.AckResponse{AckMessage: ackMessage}, err
}

func (x RpcInterface) RaftRemove(ctx context.Context, r *proto.NodeMinimalInfos) (*proto.AckResponse, error) {
	err := Remove(*x.raft, r.NodeID)
	ackMessage := fmt.Sprintf("Node %s at %s, Removed from the raft cluster!", r.NodeID, r.NodeAddress)
	return &proto.AckResponse{AckMessage: ackMessage}, err
}
