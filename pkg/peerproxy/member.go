package peerproxy

import (
	"crypto/sha1"
	"encoding/binary"
	"sort"
	"strings"

	"go.etcd.io/etcd/pkg/types"
)

func CalculateMemberID(peerURLs types.URLs, clusterName string) types.ID {
	URLs := peerURLs.StringSlice()
	sort.Strings(URLs)
	joinedURLs := strings.Join(URLs, "")
	sum := sha1.Sum([]byte(joinedURLs + clusterName))
	return types.ID(binary.BigEndian.Uint64(sum[0:8]))
}
