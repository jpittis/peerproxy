package peerproxy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var exampleConfig = `
destination: stdout
listeners:
  - listener_addr: "localhost:12381"
    upstream_addr: "localhost:12380"
  - listener_addr: "localhost:22381"
    upstream_addr: "localhost:22380"
  - listener_addr: "localhost:32381"
    upstream_addr: "localhost:32380"
`

var expectedExampleConfig = &Config{
	Destination: "stdout",
	Listeners: []*Listener{
		{ListenerAddr: "localhost:12381", UpstreamAddr: "localhost:12380"},
		{ListenerAddr: "localhost:22381", UpstreamAddr: "localhost:22380"},
		{ListenerAddr: "localhost:32381", UpstreamAddr: "localhost:32380"},
	},
}

func TestParseExampleConfig(t *testing.T) {
	config, err := ParseConfigBytes([]byte(exampleConfig))
	require.NoError(t, err)
	require.Equal(t, expectedExampleConfig, config)
}
