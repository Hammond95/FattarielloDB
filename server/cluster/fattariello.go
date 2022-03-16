package cluster

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/Hammond95/FattarielloDB/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
)

type FattarielloDB struct {
	mtx  sync.RWMutex
	fatt [3]string
}

type snapshot struct {
	fatt []string
}

var _ raft.FSM = &FattarielloDB{}

type FattarielloServer struct {
	NodeID         string
	NodeStatus     int32
	NodeAddress    string
	PeersID        []string
	PeersAddresses []string
	DB             FattarielloDB
	logger         hclog.Logger
}

func NewFattarielloServer(id string, address string, l hclog.Logger) *FattarielloServer {
	return &FattarielloServer{
		id,
		0,
		address,
		[]string{},
		[]string{},
		FattarielloDB{},
		l,
	}
}

func cloneWords(words [3]string) []string {
	var ret [3]string
	copy(ret[:], words[:])
	return ret[:]
}

func (f *FattarielloDB) Apply(l *raft.Log) interface{} {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	w := string(l.Data)
	for i := 0; i < len(f.fatt); i++ {
		copy(f.fatt[i+1:], f.fatt[i:])
		f.fatt[i] = w
	}
	return nil
}

func (f *FattarielloDB) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	return &snapshot{cloneWords(f.fatt)}, nil
}

func (f *FattarielloDB) Restore(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	fatt := strings.Split(string(b), "\n")
	copy(f.fatt[:], fatt)
	return nil
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte(strings.Join(s.fatt, "\n")))
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("sink.Write(): %v", err)
	}
	return sink.Close()
}

func (s *snapshot) Release() {
}

func (server *FattarielloServer) GetInfo(ctx context.Context, r *proto.EmptyRequest) (*proto.NodeInfos, error) {
	server.logger.Info("Handle request for GetInfo")
	return &proto.NodeInfos{
		NodeID:         server.NodeID,
		NodeStatus:     proto.NodeInfos_NodeState(server.NodeStatus),
		NodeAddress:    server.NodeAddress,
		PeersID:        server.PeersID,
		PeersAddresses: server.PeersAddresses,
	}, nil
}

func (server *FattarielloServer) SendMessage(ctx context.Context, r *proto.SendMessageRequest) (*proto.AckResponse, error) {
	server.logger.Info("Sending message to " + fmt.Sprint(r.DestinationNodeAddress))

	conn, err := grpc.Dial(r.DestinationNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewFattarielloClient(conn)

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

func (server *FattarielloServer) ReceiveMessage(ctx context.Context, r *proto.ReceiveMessageRequest) (*proto.AckResponse, error) {
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
