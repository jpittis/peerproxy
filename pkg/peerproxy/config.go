package peerproxy

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ClusterName string      `yaml:"cluster_name"`
	Destination string      `yaml:"destination"`
	Listeners   []*Listener `yaml:"listeners"`
}

type Listener struct {
	Name         string `yaml:"name"`
	ListenerAddr string `yaml:"listener_addr"`
	UpstreamAddr string `yaml:"upstream_addr"`
}

func ParseConfigBytes(buf []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(buf, config); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseConfigFile(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfigBytes(buf)
}
