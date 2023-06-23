package peerproxy

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"

	"go.etcd.io/etcd/client/pkg/v3/types"
)

func CalculateMemberIDToNameMap(config *Config) (map[types.ID]string, error) {
	mapping := map[types.ID]string{}
	for _, ln := range config.Listeners {
		URLs, err := types.NewURLs([]string{fmt.Sprintf("http://%s", ln.ListenerAddr)})
		if err != nil {
			return nil, err
		}
		id := CalculateMemberID(URLs, config.ClusterName)
		mapping[id] = ln.Name
	}
	return mapping, nil
}

func CalculateMemberID(peerURLs types.URLs, clusterName string) types.ID {
	URLs := peerURLs.StringSlice()
	sort.Strings(URLs)
	joinedURLs := strings.Join(URLs, "")
	sum := sha1.Sum([]byte(joinedURLs + clusterName))
	return types.ID(binary.BigEndian.Uint64(sum[0:8]))
}
