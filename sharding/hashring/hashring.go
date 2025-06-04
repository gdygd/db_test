package hashring

import (
	"crypto/sha1"
	"sort"
	"strconv"
)

// 가상 노드 수 (해시 분산을 높이기 위해)
const virtualNodes = 10

var hashRing []uint32
var nodeMap = map[uint32]string{} // 해시값 → 포트

// 해시 계산 함수
func hashKey(key string) uint32 {
	h := sha1.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return (uint32(bs[0]) << 24) | (uint32(bs[1]) << 16) | (uint32(bs[2]) << 8) | uint32(bs[3])
}

// 해시 링 초기화
func InitHashRing(ports []string) {
	for _, port := range ports {
		for i := 0; i < virtualNodes; i++ {
			vkey := port + "-" + strconv.Itoa(i)
			hash := hashKey(vkey)
			hashRing = append(hashRing, hash)
			nodeMap[hash] = port
		}
	}
	sort.Slice(hashRing, func(i, j int) bool {
		return hashRing[i] < hashRing[j]
	})
}

// 실제 해시로 포트 선택
func GetPortFromHashRing(key string) string {
	h := hashKey(key)
	for _, ringHash := range hashRing {
		if h <= ringHash {
			return nodeMap[ringHash]
		}
	}
	return nodeMap[hashRing[0]] // wrap-around
}
