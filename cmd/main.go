package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/schollz/progressbar/v3"
	"github.com/wojtyniak/xde/comparer"
)

const (
	CHUNK_SIZE  int = 4096
	BUFFER_SIZE int = 131072
)

var (
	J int = 2 //runtime.NumCPU() / 2
)

func main() {
	flag.Parse()
	dirnames := flag.Args()
	n, pathBuckets := findPossibleDuplicates(dirnames)

	// Progress bars
	// Yes, I know they're not supposed to be used like that but it works
	// well-enough
	filesBar := progressbar.Default(n)
	bytesBar := progressbar.DefaultBytes(-1, "Bytes scanned")

	// Start comparison
	ctx := context.Background()
	ctx = context.WithValue(ctx, "CHUNK_SIZE", CHUNK_SIZE)
	ctx = context.WithValue(ctx, "BUFFER_SIZE", BUFFER_SIZE)
	ctx = context.WithValue(ctx, "J", J)
	ctx = context.WithValue(ctx, "filesBarAdd", filesBar.Add)
	ctx = context.WithValue(ctx, "bytesBarAdd", bytesBar.Add)
	duplicates := comparer.FindDuplicates(ctx, pathBuckets)

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
