package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/wojtyniak/xde/comparer"
)

const ()

var (
	j          int
	chunkSize  int
	bufferSize int
)

func init() {
	const (
		defaultJ          = 2
		defaultChunkSize  = 4096
		defaultBufferSize = 131072
	)
	flag.IntVar(&j, "j", defaultJ, "Number of concurrent jobs running in parallel. Low values are ok since the program is I/O bound.")
	flag.IntVar(&chunkSize, "chunk-size", defaultChunkSize, "Length of the data being compared at once in bytes")
	flag.IntVar(&bufferSize, "buffer-size", defaultBufferSize, "Buffer size for the data read from disk in bytes")

	flag.Usage = func() {
		fmt.Printf("Usage: xde [options] [directory1] [directory2]\n\nOptions:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	dirnames := flag.Args()
	if len(dirnames) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	n, pathBuckets := findPossibleDuplicates(dirnames)

	// Progress bars
	// Yes, I know they're not supposed to be used like that but it works
	// well-enough
	filesBar := progressbar.Default(n)
	bytesBar := progressbar.DefaultBytes(-1, "Bytes scanned")

	// Start comparison
	ctx := context.Background()
	ctx = context.WithValue(ctx, "CHUNK_SIZE", chunkSize)
	ctx = context.WithValue(ctx, "BUFFER_SIZE", bufferSize)
	ctx = context.WithValue(ctx, "filesBarAdd", filesBar.Add)
	ctx = context.WithValue(ctx, "bytesBarAdd", bytesBar.Add)
	duplicates := comparer.FindDuplicates(ctx, pathBuckets, j)

	for i, dups := range duplicates {
		for _, d := range dups {
			fmt.Println(d)
		}
		if i < len(duplicates)-1 {
			fmt.Println()
		}
	}
	fmt.Println(len(duplicates))
}
