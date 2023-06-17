package codec

import "go.etcd.io/raft/v3/raftpb"

type Decoder interface {
	Decode() (raftpb.Message, error)
}

type Encoder interface {
	Encode(*raftpb.Message) error
}
