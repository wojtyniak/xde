package bucket

import (
	"context"
	"log"

	fc "github.com/wojtyniak/xde/filechunker"
	"github.com/wojtyniak/xde/sorter"
)

type Bucket struct {
	chunkers []*fc.FileChunker
	done     bool
}

func NewBucket(ctx context.Context, paths []string) *Bucket {
	b := new(Bucket)
	b.chunkers = make([]*fc.FileChunker, len(paths))
	errs := 0
	for i, path := range paths {
		chunker, err := fc.NewFileChunker(ctx, path)
		if err != nil {
			log.Printf("Cannot create a chunker for file %s: %s", path, err)
			errs++
			continue
		}
		b.chunkers[i] = chunker
	}
	b.chunkers = b.chunkers[:len(b.chunkers)-errs]
	return b
}

func (b *Bucket) Close() {
	b.done = true
	for _, c := range b.chunkers {
		c.Close()
	}
}

func (b *Bucket) Done() bool {
	return b.done
}

func (b *Bucket) getNextChunks() [][]byte {
	chunks := make([][]byte, len(b.chunkers))
	for i, c := range b.chunkers {
		chunks[i] = c.NextChunk()
		if chunks[i] == nil {
			return nil
		}
	}
	return chunks
}

func (b *Bucket) subBucket(chunkerIDs []int) *Bucket {
	nb := new(Bucket)
	nb.chunkers = make([]*fc.FileChunker, len(chunkerIDs))
	for i, c := range chunkerIDs {
		nb.chunkers[i] = b.chunkers[c]
	}
	return nb
}

func (b *Bucket) splitByChunks(sorted [][]int) []*Bucket {
	buckets := make([]*Bucket, 0, len(sorted))
	for i, s := range sorted {
		if len(s) == 0 {
			// File went into a different bucket
			continue
		}
		if len(s) == 1 {
			// Only one file in the bucket so it's unique
			b.chunkers[sorted[i][0]].Close()
			continue
		}
		buckets = append(buckets, b.subBucket(sorted[i]))
	}
	return buckets
}

func (b *Bucket) Sort() []*Bucket {
	chunks := b.getNextChunks()
	if chunks == nil {
		b.Close()
		return []*Bucket{b}
	}
	sortedChunks := sorter.SortChunks(chunks)
	return b.splitByChunks(sortedChunks)
}

func (b *Bucket) Paths() []string {
	ss := make([]string, len(b.chunkers))
	for i, c := range b.chunkers {
		ss[i] = c.Path()
	}
	return ss
}
