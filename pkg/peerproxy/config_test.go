package peerproxy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var exampleConfig = `
cluster_name: "etcd-cluster-1"
destination: stdout
listeners:
  - name: infra1
    listener_addr: "localhost:12381"
    upstream_addr: "localhost:12380"
  - name: infra2
    listener_addr: "localhost:22381"
    upstream_addr: "localhost:22380"
  - name: infra3
    listener_addr: "localhost:32381"
    upstream_addr: "localhost:32380"
`

var expectedExampleConfig = &Config{
	ClusterName: "etcd-cluster-1",
	Destination: "stdout",
	Listeners: []*Listener{
		{Name: "infra1", ListenerAddr: "localhost:12381", UpstreamAddr: "localhost:12380"},
		{Name: "infra2", ListenerAddr: "localhost:22381", UpstreamAddr: "localhost:22380"},
		{Name: "infra3", ListenerAddr: "localhost:32381", UpstreamAddr: "localhost:32380"},
	},
}

func TestParseExampleConfig(t *testing.T) {
	config, err := ParseConfigBytes([]byte(exampleConfig))
	require.NoError(t, err)
	require.Equal(t, expectedExampleConfig, config)
}
