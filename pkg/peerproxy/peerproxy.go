package peerproxy

import (
	"log"
	"sync"
)

type PeerProxy struct {
	config *Config
}

func NewPeerProxy(config *Config) *PeerProxy {
	return &PeerProxy{
		config: config,
	}
}

func (p *PeerProxy) ListenAndServe() error {
	log.Println("Starting peerpoxy with config", p.config)

	memberIDToNameMap, err := CalculateMemberIDToNameMap(p.config)
	if err != nil {
		return err
	}

	recorder, err := NewRecorder(p.config.Destination)
	if err != nil {
		return err
	}
	defer recorder.Close()

	var wg sync.WaitGroup
	wg.Add(len(p.config.Listeners))
	for _, ln := range p.config.Listeners {
		reverseProxy := NewReverseProxy(ln, recorder, memberIDToNameMap)
		go func() {
			if err := reverseProxy.ListenAndServe(); err != nil {
				log.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
