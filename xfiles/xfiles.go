package xfiles

import (
	"io/fs"
	"log"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type File struct {
	path string
	size int64
}

func (f *File) IsReadable() bool {
	err := unix.Access(f.path, unix.O_RDONLY)
	if err != nil {
		return false
	}
	return true
}

func (f *File) Size() int64 {
	return f.size
}

func (f *File) Path() string {
	return f.path
}

func FindFiles(dir string) (files []*File) {
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d == nil {
				// initial fs.Stat on the root directory failed
				return err
			}
			log.Printf("Cannot read directory %s/%s\n", path, d.Name())
			return fs.SkipDir
		}
		if d.IsDir() {
			// skip this entry
			return nil
		}

		info, err := d.Info()
		if err != nil {
			log.Fatalf("Failed to get info for file: %s/%s\n", path, d.Name())
		} else {
			f := &File{path, info.Size()}
			files = append(files, f)
		}
		return nil
	})
	return files
}
