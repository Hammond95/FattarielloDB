package main

import (
	"context"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/Hammond95/FattarielloDB/server/cluster"
	"github.com/hashicorp/raft"
)

type RpcInterface struct {
	fattarielloServer *cluster.FattarielloServer
	raft              *raft.Raft
}

func (x RpcInterface) GetInfo(ctx context.Context, r *proto.EmptyRequest) (*proto.NodeInfos, error) {
	return x.fattarielloServer.GetInfo(ctx, r)
}

func (x RpcInterface) SendMessage(ctx context.Context, r *proto.SendMessageRequest) (*proto.AckResponse, error) {
	return x.fattarielloServer.SendMessage(ctx, r)
}

func (x RpcInterface) ReceiveMessage(ctx context.Context, r *proto.ReceiveMessageRequest) (*proto.AckResponse, error) {
	return x.fattarielloServer.ReceiveMessage(ctx, r)
}
