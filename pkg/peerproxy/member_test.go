package peerproxy

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/pkg/v3/types"
)

func TestCalculateMemberID(t *testing.T) {
	URLs, err := types.NewURLs([]string{"http://127.0.0.1:12381"})
	require.NoError(t, err)
	id := CalculateMemberID(URLs, "etcd-cluster-1")
	expectedID, err := types.IDFromString("5ec957d7b9f69bbb")
	require.NoError(t, err)
	require.Equal(t, expectedID, id)
}
