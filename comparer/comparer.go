package comparer

import (
	"bufio"
	"context"
	"sync"

	"github.com/wojtyniak/xde/bucket"
)

func processBucket(initialBucket *bucket.Bucket) []*bucket.Bucket {
	buckets := make([]*bucket.Bucket, 0, 10)
	queue := []*bucket.Bucket{initialBucket}
	for len(queue) > 0 {
		subBuckets := queue[0].Sort()
		queue = queue[1:]
		for _, b := range subBuckets {
			if b.Done() {
				buckets = append(buckets, b)
				continue
			}
			queue = append(queue, b)
		}
	}
	return buckets
}

func worker(ctx context.Context, in chan *bucket.Bucket, out chan []*bucket.Bucket) {
	defer close(out)
	for b := range in {
		buckets := processBucket(b)
		select {
		case out <- buckets:
		case <-ctx.Done():
			<-in
			return
		}
	}
}

func mergeChannels(outs []chan []*bucket.Bucket) chan []*bucket.Bucket {
	out := make(chan []*bucket.Bucket, len(outs))
	wg := sync.WaitGroup{}
	wg.Add(len(outs))

	for _, o := range outs {
		go func(o chan []*bucket.Bucket) {
			for b := range o {
				out <- b
			}
			wg.Done()
		}(o)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func startWorkers(ctx context.Context, n int) (in chan *bucket.Bucket, out chan []*bucket.Bucket) {
	in = make(chan *bucket.Bucket)
	if n == 1 {
		out := make(chan []*bucket.Bucket)
		go worker(ctx, in, out)
		return in, out
	}

	outChannels := make([]chan []*bucket.Bucket, n)
	for i := 0; i < n; i++ {
		outChannels[i] = make(chan []*bucket.Bucket)
		go worker(ctx, in, outChannels[i])
	}
	out = mergeChannels(outChannels)

	return in, out
}

func startFeeder(ctx context.Context, pathBuckets [][]string, in chan<- *bucket.Bucket) {
	defer close(in)
	bufferSize := ctx.Value("BUFFER_SIZE").(int)
	chunkSize := ctx.Value("CHUNK_SIZE").(int)
	bufferPool := &sync.Pool{New: func() interface{} {
		b := bufio.NewReaderSize(nil, bufferSize)
		return b
	}}
	chunkPool := &sync.Pool{New: func() interface{} {
		b := make([]byte, chunkSize)
		return b
	}}
	for _, pathBucket := range pathBuckets {
		bucket := bucket.NewBucket(ctx, pathBucket, bufferPool, chunkPool)
		select {
		case in <- bucket:
		case <-ctx.Done():
			return
		}
	}
}

func FindDuplicates(ctx context.Context, pathBuckets [][]string) [][]string {
	j := ctx.Value("J").(int)
	// chunk_size := ctx.Value("CHUNK_SIZE").(int)
	in, out := startWorkers(ctx, j)
	go startFeeder(ctx, pathBuckets, in)

	pathBucket := make([][]string, 0, 10)
	for buckets := range out {
		for _, bucket := range buckets {
			pathBucket = append(pathBucket, bucket.Paths())
		}
	}
	return pathBucket
}
