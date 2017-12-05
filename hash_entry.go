package protohash

import "bytes"

type hashEntry struct {
	khash []byte
	vhash []byte
}

type byKHash []hashEntry

func (h byKHash) Len() int {
	return len(h)
}

func (h byKHash) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h byKHash) Less(i, j int) bool {
	return bytes.Compare(h[i].khash[:], h[j].khash[:]) < 0
}
