package sorter

import (
	"bytes"
)

func sortTwo(chunks [][]byte) [][]int {
	if bytes.Compare(chunks[0], chunks[1]) == 0 {
		return [][]int{{0, 1}}
	} else {
		return [][]int{{0}, {1}}
	}
}

func sortThree(chunks [][]byte) [][]int {
	ab := bytes.Compare(chunks[0], chunks[1]) == 0
	ac := bytes.Compare(chunks[0], chunks[2]) == 0
	bc := bytes.Compare(chunks[1], chunks[2]) == 0

	switch {
	case ab && ac:
		return [][]int{{0, 1, 2}, {}, {}}
	case ab && !bc:
		return [][]int{{0, 1}, {}, {2}}
	case !ab && ac:
		return [][]int{{0, 2}, {1}, {}}
	case !ab && bc:
		return [][]int{{0}, {1, 2}, {}}
	default:
		return [][]int{{0}, {1}, {2}}
	}
}

func sortHash(chunks [][]byte) [][]int {
	hc := newHashComparable(chunks)
	return sort(hc)
}

func sortBytes(chunks [][]byte) [][]int {
	bc := &ByteComparable{chunks}
	return sort(bc)
}

func sort(c Comparable) [][]int {
	matches := make([]int, c.Len())
	for i := range matches {
		matches[i] = i
	}
	for i := 0; i < c.Len()-1; i++ {
		if matches[i] < i {
			continue
		}
		for j := i + 1; j < c.Len(); j++ {
			if matches[j] <= i {
				continue
			} else if c.Equal(i, j) {
				matches[j] = i
			}
		}
	}
	assigments := make([][]int, c.Len())
	for a := range assigments {
		assigments[a] = make([]int, 0)
	}
	for i := range matches {
		assigments[matches[i]] = append(assigments[matches[i]], i)
	}
	return assigments
}

func SortChunks(chunks [][]byte) [][]int {
	n := len(chunks)
	switch {
	case n == 2:
		return sortTwo(chunks)
	case n == 3:
		return sortThree(chunks)
	default:
		return sortHash(chunks)
	}
}
