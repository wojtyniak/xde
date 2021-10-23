# xde

Fast file deduplicator

## Istallation

`go install github.com/wojtyniak/xde/cmd/xde@latest`

## Usage

```
Usage: xde [options] [directory1] [directory2]

Options:
  -buffer-size int
    	Buffer size for the data read from disk in bytes (default 131072)
  -chunk-size int
    	Length of the data being compared at once in bytes (default 4096)
  -j int
    	Number of concurrent jobs running in parallel. Low values are ok since the program is I/O bound. (default 2)
  -q	Don't print stats to stderr and don't show the progress bar
  -w string
    	Write output to the specified file (the file is going to be truncated)
```
