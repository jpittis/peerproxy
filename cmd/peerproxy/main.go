package main

import (
	"log"
	"os"

	"github.com/jpittis/peerproxy/pkg/peerproxy"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Expected config filename as arg")
		os.Exit(1)
	}
	filename := os.Args[1]
	config, err := peerproxy.ParseConfigFile(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	peerProxy := peerproxy.NewPeerProxy(config)
	if err := peerProxy.ListenAndServe(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
