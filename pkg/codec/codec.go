package codec

import (
	"io"

	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/raft/v3/raftpb"
)

type Decoder interface {
	Decode() (raftpb.Message, error)
}

type Encoder interface {
	Encode(*raftpb.Message) error
}

func NewMessageDecoder(r io.Reader) Decoder {
	return &messageDecoder{r}
}

func NewMessageEncoder(w io.Writer) Encoder {
	return &messageEncoder{w}
}

func NewMsgAppDecoder(r io.Reader, local, remote types.ID) Decoder {
	return newMsgAppV2Decoder(r, local, remote)
}

func NewMsgAppEncoder(w io.Writer) Encoder {
	return newMsgAppV2Encoder(w)
}
