package main

import (
	"io/fs"
	"log"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func scanDirectory(sizeToPath map[int64][]string, dir string) {
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d == nil {
				// initial fs.Stat on the root directory failed
				panic(err)
			}
			log.Printf("Cannot read directory %s/%s\n", path, d.Name())
			return fs.SkipDir
		}
		if d.IsDir() {
			// don't log directories
			return nil
		}
		info, err := d.Info()
		if err != nil {
			log.Fatalf("Failed to get info for file: %s, skipping", path)
			return nil
		}
		size := info.Size()
		if size == 0 {
			// Don't process empty files
			return nil
		}

		err = unix.Access(path, unix.O_RDONLY)
		if err != nil {
			log.Printf("Cannot read file: %s, skipping", path)
			return nil
		}
		if len(sizeToPath[size]) == 0 {
			sizeToPath[size] = []string{path}
		} else {
			sizeToPath[size] = append(sizeToPath[size], path)
		}
		return nil
	})

}

func findPossibleDuplicates(dirs []string) (n int, paths [][]string) {
	sizeToPath := make(map[int64][]string)
	for _, dir := range dirs {
		scanDirectory(sizeToPath, dir)
	}

	paths = make([][]string, 0, 16)
	for key := range sizeToPath {
		if len(sizeToPath[key]) > 1 {
			paths = append(paths, sizeToPath[key])
			n += len(sizeToPath[key])
		}
	}
	return n, paths
}
