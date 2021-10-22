package filechunker

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"sync"
)

type FileChunker struct {
	ctx         context.Context
	cancel      context.CancelFunc
	path        string
	chunkChan   chan []byte
	bufferPool  *sync.Pool
	chunkPool   *sync.Pool
	filesBarAdd func(num int) error
}

func NewFileChunker(ctx context.Context, path string, bufferPool, chunkPool *sync.Pool) (*FileChunker, error) {
	fc := new(FileChunker)
	fc.bufferPool = bufferPool
	fc.chunkPool = chunkPool
	fc.path = path
	fc.ctx, fc.cancel = context.WithCancel(ctx)
	fc.chunkChan = make(chan []byte, 1)

	fab := ctx.Value("filesBarAdd")
	if fab != nil {
		fc.filesBarAdd = fab.(func(num int) error)
	}

	err := fc.startChunkReader(fc.chunkChan)
	if err != nil {
		return nil, err
	}

	return fc, nil
}

func (fc *FileChunker) Path() string {
	return fc.path
}

func (fc *FileChunker) Close() {
	fc.cancel()
}

func clearChunk(chunk []byte) {
	for i := 0; i < len(chunk); i++ {
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
		defer f.Close()
		defer close(out)

		for {
			b := (fc.chunkPool.Get().([]byte))
			clearChunk(b)
			n, err := br.Read(b)
			if err != nil && err != io.EOF {
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
