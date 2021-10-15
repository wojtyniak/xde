package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/wojtyniak/xde/comparer"
)

const (
	CHUNK_SIZE int = 4096
)

var (
	J int = runtime.NumCPU() / 2
)

func main() {
	flag.Parse()
	dirnames := flag.Args()

	pathBuckets := findPossibleDuplicates(dirnames)

	// Start comparison
	ctx := context.Background()
	ctx = context.WithValue(ctx, "CHUNK_SIZE", CHUNK_SIZE)
	ctx = context.WithValue(ctx, "J", J)
	duplicates := comparer.FindDuplicates(ctx, pathBuckets)

	for _, dups := range duplicates {
		fmt.Println()
		for _, d := range dups {
			fmt.Println(d)
		}
	}
}
