package sorter

import (
	"bytes"

	"github.com/cespare/xxhash"
)

type Comparable interface {
	Len() int
	Equal(i, j int) bool
}

type ByteComparable struct {
	chunks [][]byte
}

func (b *ByteComparable) Equal(i, j int) bool {
	return bytes.Compare(b.chunks[i], b.chunks[j]) == 0
}

func (b *ByteComparable) Len() int {
	return len(b.chunks)
}

type HashComparable struct {
	hashes []uint64
}

func newHashComparable(chunks [][]byte) *HashComparable {
	h := &HashComparable{make([]uint64, len(chunks))}
	for i := range h.hashes {
		h.hashes[i] = xxhash.Sum64(chunks[i])
	}
	return h
}

func (h *HashComparable) Equal(i, j int) bool {
	return h.hashes[i] == h.hashes[j]
}

func (h *HashComparable) Len() int {
	return len(h.hashes)
}
