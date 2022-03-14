package application

import (
	"github.com/hashicorp/raft"
	"sync"
)

type fattariello struct {
	mtx  sync.RWMutex
	fatt [3]string
}

type snapshot struct {
	fatt []string
}

var _ raft.FSM = &fattariello{}

func cloneWords(words [3]string) []string {
	var ret [3]string
	copy(ret[:], words[:])
	return ret[:]
}

func (f *fattariello) Apply(l *raft.Log) interface{} {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	w := string(l.Data)
	for i := 0; i < len(f.fatt); i++ {
		copy(f.fatt[i+1:], f.fatt[i:])
		f.fatt[i] = w
	}
	return nil
}

func (f *fattariello) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	return &snapshot{cloneWords(f.fatt)}, nil
}

func (f *fattariello) Restore(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	fatt := strings.Split(string(b), "\n")
	copy(f.fatt[:], fatt)
	return nil
}
