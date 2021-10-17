package filechunker

import (
	"context"
	"log"
	"os"
)

type FileChunker struct {
	path        string
	ctx         context.Context
	cancel      context.CancelFunc
	chunkChan   chan []byte
	chunkSize   int
	filesBarAdd func(num int) error
	bytesBarAdd func(num int) error
}

func NewFileChunker(ctx context.Context, path string) (*FileChunker, error) {
	fc := new(FileChunker)
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

func (fc *FileChunker) startChunkReader(out chan<- []byte) error {
	f, err := os.Open(fc.path)
	if err != nil {
		return err
	}

	go func() {
		if fc.filesBarAdd != nil {
			defer fc.filesBarAdd(1)
		}
		defer close(out)
		defer fc.Close()

		for {
			b := make([]byte, fc.chunkSize)
			n, err := f.Read(b)
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
