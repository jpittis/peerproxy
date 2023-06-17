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
}

func NewReverseProxy(ln *Listener) *ReverseProxy {
	target := &url.URL{Scheme: "http", Host: ln.UpstreamAddr}
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	p := &ReverseProxy{
		ln:           ln,
		reverseProxy: reverseProxy,
	}
	p.reverseProxy.ModifyResponse = p.modifyResponse
	return p
}

func (p *ReverseProxy) ListenAndServe() error {
	log.Printf("Forwarding from %s to %s", p.ln.ListenerAddr, p.ln.UpstreamAddr)
	return http.ListenAndServe(p.ln.ListenerAddr, p.reverseProxy)
}

func (p *ReverseProxy) modifyResponse(resp *http.Response) (err error) {
	r, w := io.Pipe()
	var dec codec.Decoder
	var enc codec.Encoder
	if strings.HasPrefix(resp.Request.URL.Path, messagePathPrefix) {
		dec = codec.NewMessageDecoder(resp.Body)
		enc = codec.NewMessageEncoder(w)
	} else if strings.HasPrefix(resp.Request.URL.Path, msgappPathPrefix) {
		dec = codec.NewMsgAppDecoder(resp.Body, types.ID(0), types.ID(0))
		enc = codec.NewMsgAppEncoder(w)
	} else {
		// Don't modify payloads for other request paths.
		return nil
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
