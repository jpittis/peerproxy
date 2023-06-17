package peerproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	return nil
}
