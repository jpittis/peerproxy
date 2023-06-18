package peerproxy

import (
	"fmt"
	"io"
	"os"
	"sync"

	"go.etcd.io/raft/v3/raftpb"
)

type Event struct {
	Upstream string
	Path     string
	SrcID    string
	DstID    string
	Message  *raftpb.Message
}

func NewRecorder(destination string) (*Recorder, error) {
	var target io.Writer
	var file *os.File
	if destination == "stdout" {
		target = os.Stdout
	} else if destination == "stderr" {
		target = os.Stderr
	} else {
		var err error
		file, err = os.Create(destination)
		if err != nil {
			return nil, err
		}
		target = file
	}
	return &Recorder{
		target: target,
		file:   file,
	}, nil
}

type Recorder struct {
	sync.Mutex
	target io.Writer
	file   *os.File
}

func (r *Recorder) Record(event *Event) {
	r.Lock()
	defer r.Unlock()
	fmt.Fprintf(r.target, "%+v\n", event)
}

func (r *Recorder) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}
