package filechunker

import (
	"bufio"
	"context"
	"log"
	"os"
	"sync"
)

type FileChunker struct {
	path        string
	ctx         context.Context
	cancel      context.CancelFunc
	chunkChan   chan []byte
	chunkSize   int
	bufferPool  *sync.Pool
	chunkPool   *sync.Pool
	filesBarAdd func(num int) error
	bytesBarAdd func(num int) error
}

func NewFileChunker(ctx context.Context, path string, bufferPool, chunkPool *sync.Pool) (*FileChunker, error) {
	fc := new(FileChunker)
	fc.bufferPool = bufferPool
	fc.chunkPool = chunkPool
	fc.path = path
	fc.ctx, fc.cancel = context.WithCancel(ctx)
	chunkChan := make(chan []byte, 1)
	fc.chunkSize = ctx.Value("CHUNK_SIZE").(int)
	fab := ctx.Value("filesBarAdd")
	if fab != nil {
		fc.filesBarAdd = fab.(func(num int) error)
	}
	bba := ctx.Value("bytesBarAdd")
	if bba != nil {
		fc.bytesBarAdd = bba.(func(num int) error)
	}
	err := fc.startChunkReader(chunkChan)
	if err != nil {
		return nil, err
	}
	fc.chunkChan = chunkChan

	return fc, nil
}

func (fc *FileChunker) Path() string {
	return fc.path
}

func (fc *FileChunker) Close() {
	fc.cancel()
}

func clearChunk(chunk []byte) {
	for i := 0; i < cap(chunk); i++ {
		chunk[i] = 0
	}
}

func (fc *FileChunker) startChunkReader(out chan<- []byte) error {
	f, err := os.Open(fc.path)
	br := fc.bufferPool.Get().(*bufio.Reader)
	br.Reset(f)
	if err != nil {
		return err
	}

	go func() {
		if fc.filesBarAdd != nil {
			defer fc.filesBarAdd(1)
		}
		defer fc.bufferPool.Put(br)
		defer close(out)
		defer f.Close()

		for {
			// b := make([]byte, fc.chunkSize)
			b := (fc.chunkPool.Get().([]byte))
			clearChunk(b)
			n, err := br.Read(b)
			if fc.bytesBarAdd != nil {
				go fc.bytesBarAdd(n)
			}
			if err != nil && err.Error() != "EOF" {
				log.Printf("Error while reading file %s: %s", fc.path, err)
				return
			}
			if n == 0 {
				// Reached EOF
				return
			}
			select {
			case out <- b:
			case <-fc.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (fc *FileChunker) NextChunk() []byte {
	// This function doesn't check for ctx.Done as it's handled by the goroutine
	// in the background. Returning nil should be handled by the user.
	return <-fc.chunkChan
}
