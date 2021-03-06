package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/wojtyniak/xde/comparer"
)

const ()

var (
	j          int
	chunkSize  int
	bufferSize int
	writeOut   string
	quiet      bool
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
	flag.StringVar(&writeOut, "w", "", "Write output to the specified file (the file is going to be truncated)")
	flag.BoolVar(&quiet, "q", false, "Don't print stats to stderr and don't show the progress bar")

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

	var out io.Writer
	if writeOut != "" {
		outFile, err := os.Create(writeOut)
		if err != nil {
			log.Fatalln(err)
		}
		outBuf := bufio.NewWriter(outFile)
		out = outBuf
		defer func() {
			errors := false
			err := outBuf.Flush()
			if err != nil {
				log.Println(err)
				errors = true
			}

			err = outFile.Close()
			if err != nil {
				log.Fatalln(err)
			}
			if errors {
				log.Fatalf("Writing to the file: %s failed\n", writeOut)
			}
		}()
	} else {
		out = os.Stdout
	}

	startTime := time.Now()
	n, pathBuckets := findPossibleDuplicates(dirnames)
	if !quiet {
		fmt.Fprintf(os.Stdout, "Found %d possibly duplicated files in %d buckets\n", n, len(pathBuckets))
	}

	// Progress bars
	// Yes, I know they're not supposed to be used like that but it works
	// well-enough

	// Start comparison
	ctx := context.Background()
	ctx = context.WithValue(ctx, "CHUNK_SIZE", chunkSize)
	ctx = context.WithValue(ctx, "BUFFER_SIZE", bufferSize)
	if !quiet {
		filesBar := progressbar.NewOptions(n, progressbar.OptionClearOnFinish())
		ctx = context.WithValue(ctx, "filesBarAdd", filesBar.Add)
	}

	duplicates := comparer.FindDuplicates(ctx, pathBuckets, j)
	endTime := time.Now()

	n = 0
	for i, dups := range duplicates {
		n += len(dups)
		for _, d := range dups {
			fmt.Fprintln(out, d)

		}
		if i < len(duplicates)-1 {
			fmt.Fprintln(out) // Separate entries
		}
	}
	if !quiet {
		fmt.Fprintf(os.Stderr, "Found %d duplicated files in %d buckets in %s \n", n, len(duplicates), endTime.Sub(startTime))
	}
}
