package main

import (
	"context"
	"encoding/json"

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
