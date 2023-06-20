package peerproxy

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ClusterName string      `yaml:"cluster_name"`
	Destination string      `yaml:"destination"`
	Listeners   []*Listener `yaml:"listeners"`
}

type Listener struct {
	Name            string                   `yaml:"name"`
	ListenerAddr    string                   `yaml:"listener_addr"`
	UpstreamAddr    string                   `yaml:"upstream_addr"`
	OutboundLatency map[string]time.Duration `yaml:"outbound_latency"`
}

func ParseConfigBytes(buf []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(buf, config); err != nil {
		return nil, err
	}
	if err := ValidateConfig(config); err != nil {
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

func ValidateConfig(config *Config) error {
	listenerNames := map[string]struct{}{}
	for _, ln := range config.Listeners {
		if _, alreadyPresent := listenerNames[ln.Name]; alreadyPresent {
			return fmt.Errorf("duplicate listener name %s", ln.Name)
		}
		listenerNames[ln.Name] = struct{}{}
	}

	for _, ln := range config.Listeners {
		for name := range ln.OutboundLatency {
			if name == ln.Name {
				return fmt.Errorf("listener %s self reference in outbound latency", ln.Name)
			} else if _, listenerExists := listenerNames[name]; !listenerExists {
				return fmt.Errorf(
					"outbound latency for listener %s references %s which does not exist",
					ln.Name, name,
				)
			}
		}
	}

	return nil
}
