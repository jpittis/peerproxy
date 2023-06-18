package peerproxy

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/jpittis/peerproxy/pkg/codec"
	"go.etcd.io/etcd/client/pkg/v3/types"
)

const (
	msgappPathPrefix  = "/raft/stream/msgapp"
	messagePathPrefix = "/raft/stream/message"
)

type ReverseProxy struct {
	ln           *Listener
	reverseProxy *httputil.ReverseProxy
	recorder     *Recorder
}

func NewReverseProxy(ln *Listener, recorder *Recorder) *ReverseProxy {
	target := &url.URL{Scheme: "http", Host: ln.UpstreamAddr}
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	p := &ReverseProxy{
		ln:           ln,
		reverseProxy: reverseProxy,
		recorder:     recorder,
	}
	p.reverseProxy.ModifyResponse = p.modifyResponse
	return p
}

func (p *ReverseProxy) ListenAndServe() error {
	log.Printf("Forwarding from %s to %s", p.ln.ListenerAddr, p.ln.UpstreamAddr)
	return http.ListenAndServe(p.ln.ListenerAddr, p.reverseProxy)
}

func (p *ReverseProxy) modifyResponse(resp *http.Response) error {
	// Early return (noop proxy) if this isn't a rafthttp message.
	if !strings.HasPrefix(resp.Request.URL.Path, messagePathPrefix) &&
		!strings.HasPrefix(resp.Request.URL.Path, messagePathPrefix) {
		return nil
	}

	// The source of the decoded messages is the destination of the HTTP request, and the
	// destination is the origin of the HTTP request because in rafthttp, the destination
	// initiates the HTTP request.
	srcID, err := types.IDFromString(resp.Request.Header.Get("X-Raft-To"))
	if err != nil {
		return err
	}
	dstID, err := types.IDFromString(resp.Request.Header.Get("X-Server-From"))
	if err != nil {
		return err
	}

	r, w := io.Pipe()
	var dec codec.Decoder
	var enc codec.Encoder
	if strings.HasPrefix(resp.Request.URL.Path, messagePathPrefix) {
		dec = codec.NewMessageDecoder(resp.Body)
		enc = codec.NewMessageEncoder(w)
	} else if strings.HasPrefix(resp.Request.URL.Path, msgappPathPrefix) {
		// Decode the message as if we're the destination receiving it.
		dec = codec.NewMsgAppDecoder(resp.Body, dstID, srcID)
		enc = codec.NewMsgAppEncoder(w)
	}

	body := resp.Body
	go func() {
		defer body.Close()
		for {
			msg, err := dec.Decode()
			if err != nil {
				log.Printf("error: %+v", err)
				return
			}
			p.recorder.Record(&Event{
				Upstream: p.ln.UpstreamAddr,
				Path:     resp.Request.URL.Path,
				SrcID:    srcID.String(),
				DstID:    dstID.String(),
				Message:  &msg,
			})
			err = enc.Encode(&msg)
			if err != nil {
				log.Printf("error: %+v", err)
				return
			}
		}
	}()
	resp.Body = r

	return nil
}
