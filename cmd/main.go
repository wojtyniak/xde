package main

import (
	"flag"
	"fmt"

	"github.com/wojtyniak/xde/bucket"
	"github.com/wojtyniak/xde/xfiles"
)

func main() {
	flag.Parse()
	dirnames := flag.Args()

	fileSlices := make([][]*xfiles.File, len(dirnames))
	noOfFiles := 0
	for i, dir := range dirnames {
		fileSlices[i] = xfiles.FindFiles(dir)
		noOfFiles += len(fileSlices[i])
	}

	files := make([]*xfiles.File, 0, noOfFiles)
	for _, fileSlice := range fileSlices {
		files = append(files, fileSlice...)
	}

	buckets := bucket.BucketFilesBySize(files)

	fmt.Println(len(buckets))
}
