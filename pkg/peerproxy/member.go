package peerproxy

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"

	"go.etcd.io/etcd/pkg/types"
)

func CalculateMemberNameToIDMap(config *Config) (map[string]types.ID, error) {
	mapping := map[string]types.ID{}
	for _, ln := range config.Listeners {
		URLs, err := types.NewURLs([]string{fmt.Sprintf("http://%s", ln.ListenerAddr)})
		if err != nil {
			return nil, err
		}
		mapping[ln.Name] = CalculateMemberID(URLs, config.ClusterName)
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
