package main

import (
	"context"
	"fmt"
	"os"

	"path/filepath"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"

	transport "github.com/Jille/raft-grpc-transport"
	boltdb "github.com/hashicorp/raft-boltdb"
)

func NewRaft(ctx context.Context, myID, myAddress string, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	raftDir := "/tmp"
	baseDir := filepath.Join(raftDir, myID)

	fmt.Printf("BaseDir is %v\n", baseDir)

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithInsecure()})

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	raftBootstrap := true
	if raftBootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(myID),
					Address:  raft.ServerAddress(myAddress),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return r, tm, nil
}

func Join(x *raft.Raft, nodeID string, nodeAddress string) error {
	if x.State() != raft.Leader {
		return fmt.Errorf("This node is not a Leader.")
	}

	configFuture := x.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return fmt.Errorf("Failed to get raft configuration: %s", err.Error())
	}

	// This must be run on the leader or it will fail.
	f := x.AddVoter(
		raft.ServerID(nodeID),
		raft.ServerAddress(nodeAddress), 0, 0,
	)

	if f.Error() != nil {
		return fmt.Errorf("Failed to add voter: %s", f.Error().Error())
	}

	return nil
}

func Remove(x *raft.Raft, nodeID string) error {
	if x.State() != raft.Leader {
		return fmt.Errorf("This node is not a Leader.")
	}

	/*configFuture := x.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return fmt.Errorf("Failed to get raft configuration: %s", err.Error())
	}*/

	future := x.RemoveServer(raft.ServerID(nodeID), 0, 0)
	if err := future.Error(); err != nil {
		return fmt.Errorf("Failed to remove existing node %s: %s", nodeID, err.Error())
	}

	return nil
}
