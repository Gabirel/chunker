package chunker

import (
	"crypto/sha256"
	"github.com/restic/chunker"
	"io"
)

type Chunking struct {
	pol    chunker.Pol
	chnker *chunker.Chunker
}

// ChunkInfo stores one chunk information and its digest
type ChunkInfo struct {
	Chunk  *chunker.Chunk
	Digest []byte
}

// New returns a new Chunking
// min chunk size : 2KB
// max chunk size : 16KB
// average chunk size : 8KB
func New(rd io.Reader) *Chunking {
	var (
		// min chunk size is set to 2KB (1 << 11)
		minBits = 11

		// max chunk size is set to 16KB (1 << 14)
		maxBits = 14

		// average chunk size is set to 8KB (1 << 13)
		averageBits = 13
	)

	return NewWithBoundaries(rd, minBits, maxBits, averageBits)
}

// NewWithBoundaries returns a new Chunking
// this function can set min, max and average chunk size
func NewWithBoundaries(rd io.Reader, minBits, maxBits, averageBits int) *Chunking {
	// polynomial used for all chunk algorithm
	// The value comes from chunker_test.go
	pol := chunker.Pol(0x3DA3358B4DC173)

	minChunkSize := uint(1 << minBits)
	maxChunkSize := uint(1 << maxBits)

	chnker := chunker.NewWithBoundaries(rd, pol, minChunkSize, maxChunkSize)
	chnker.SetAverageBits(averageBits)
	return &Chunking{
		pol:    pol,
		chnker: chnker,
	}
}

func (c *Chunking) Chunking() (*ChunkInfo, error) {

	chunk, err := c.chnker.Next(nil)
	if err == io.EOF {
		return nil, err
	}
	return &ChunkInfo{
		Chunk:  &chunk,
		Digest: hashDigest(chunk.Data),
	}, nil
}

func hashDigest(d []byte) []byte {
	h := sha256.New()
	h.Write(d)
	return h.Sum(nil)
}
