package example

import "github.com/cespare/xxhash"

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

type Member string

func (m Member) String() string {
	return string(m)
}
