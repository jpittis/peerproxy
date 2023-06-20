package peerproxy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var exampleConfig = `
cluster_name: "etcd-cluster-1"
destination: stderr
listeners:
  - name: infra1
    listener_addr: "localhost:12381"
    upstream_addr: "localhost:12380"
    outbound_latency:
      infra2: 30ms
      infra3: 30ms
  - name: infra2
    listener_addr: "localhost:22381"
    upstream_addr: "localhost:22380"
    outbound_latency:
      infra1: 30ms
      infra3: 10ms
  - name: infra3
    listener_addr: "localhost:32381"
    upstream_addr: "localhost:32380"
    outbound_latency:
      infra1: 30ms
      infra2: 10ms
`

var expectedExampleConfig = &Config{
	ClusterName: "etcd-cluster-1",
	Destination: "stderr",
	Listeners: []*Listener{
		{
			Name:         "infra1",
			ListenerAddr: "localhost:12381",
			UpstreamAddr: "localhost:12380",
			OutboundLatency: map[string]time.Duration{
				"infra2": 30 * time.Millisecond,
				"infra3": 30 * time.Millisecond,
			},
		},
		{
			Name:         "infra2",
			ListenerAddr: "localhost:22381",
			UpstreamAddr: "localhost:22380",
			OutboundLatency: map[string]time.Duration{
				"infra1": 30 * time.Millisecond,
				"infra3": 10 * time.Millisecond,
			},
		},
		{
			Name:         "infra3",
			ListenerAddr: "localhost:32381",
			UpstreamAddr: "localhost:32380",
			OutboundLatency: map[string]time.Duration{
				"infra1": 30 * time.Millisecond,
				"infra2": 10 * time.Millisecond,
			},
		},
	},
}

func TestParseExampleConfig(t *testing.T) {
	config, err := ParseConfigBytes([]byte(exampleConfig))
	require.NoError(t, err)
	require.Equal(t, expectedExampleConfig, config)
}

func TestValidateConfig(t *testing.T) {
	require.NoError(t, ValidateConfig(expectedExampleConfig))

	require.ErrorContains(t, ValidateConfig(
		&Config{
			Listeners: []*Listener{
				{Name: "dupe"},
				{Name: "dupe"},
			},
		},
	), "duplicate listener name dupe")

	require.ErrorContains(t, ValidateConfig(
		&Config{
			Listeners: []*Listener{
				{
					Name: "bbq",
					OutboundLatency: map[string]time.Duration{
						"bbq": 100 * time.Millisecond,
					},
				},
			},
		},
	), "listener bbq self reference in outbound latency")

	require.ErrorContains(t, ValidateConfig(
		&Config{
			Listeners: []*Listener{
				{
					Name: "bbq",
					OutboundLatency: map[string]time.Duration{
						"nerp": 100 * time.Millisecond,
					},
				},
			},
		},
	), "outbound latency for listener bbq references nerp which does not exist")
}
