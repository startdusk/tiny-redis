package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc      HashFunc
	nodeHashes    []int
	nodeHashesMap map[int]string
}

func NewNodeMap(fn HashFunc) *NodeMap {
	m := NodeMap{
		hashFunc:      fn,
		nodeHashesMap: make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return &m
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashes) == 0
}

func (m *NodeMap) Add(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashes = append(m.nodeHashes, hash)
		m.nodeHashesMap[hash] = key
	}
	sort.Ints(m.nodeHashes)
}

func (m *NodeMap) Pick(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	index := sort.Search(len(m.nodeHashes), func(i int) bool {
		return m.nodeHashes[i] >= hash
	})
	if index == len(m.nodeHashes) {
		index = 0
	}
	hash = m.nodeHashes[index]
	return m.nodeHashesMap[hash]
}
