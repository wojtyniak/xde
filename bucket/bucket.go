package bucket

import (
	"log"
	"sort"

	"github.com/wojtyniak/xde/xfiles"
)

type Sizer interface {
	Size() int64
}

type Bucket struct {
	Files []*xfiles.File
}

func NewBucket(files []*xfiles.File) *Bucket {
	return &Bucket{files}
}

func bucketBySize(sizers []Sizer) [][]Sizer {
	sort.Slice(sizers, func(i int, j int) bool {
		return sizers[i].Size() < sizers[j].Size()
	})
	sizeToSizers := make(map[int64][]Sizer)
	for _, f := range sizers {
		if sizeToSizers[f.Size()] != nil {
			sizeToSizers[f.Size()] = append(sizeToSizers[f.Size()], f)
			continue
		}
		sizeToSizers[f.Size()] = []Sizer{f}
	}
	bucketedSizers := make([][]Sizer, 0, 1)
	for _, sizers := range sizeToSizers {
		if len(sizers) <= 1 {
			continue
		}
		bucketedSizers = append(bucketedSizers, sizers)
	}
	return bucketedSizers
}

func BucketFilesBySize(files []*xfiles.File) []*Bucket {
	sizers := make([]Sizer, 0, len(files))
	for _, file := range files {
		sizers = append(sizers, file)
	}
	bucketedSizers := bucketBySize(sizers)
	buckets := make([]*Bucket, 0, len(bucketedSizers))

	for _, bucket := range bucketedSizers {
		fileBucket := make([]*xfiles.File, 0, len(bucket))
		for i := range bucket {
			f := bucket[i].(*xfiles.File)
			if !f.IsReadable() {
				log.Printf("File: %s is unreadable", f.Path())
				continue
			}
			fileBucket = append(fileBucket, f)
		}
		buckets = append(buckets, NewBucket(fileBucket))
	}
	return buckets
}
