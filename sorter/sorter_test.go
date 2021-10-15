package sorter

import (
	"testing"
)

var (
	data = [][]byte{{1, 2, 3, 4, 5}, {1, 2, 3, 4, 5}, {1, 2, 3, 4}}
)

func TestSortChunks(t *testing.T) {
	sorted := SortChunks(data)
	if len(sorted) != len(data) {
		t.Errorf("Wrong length of sorted chunks: %d", len(sorted))
	}
	if len(sorted[0]) != 2 {
		t.Errorf("Wrong first bucket length: %d", len(sorted[0]))
	}
	if sorted[0][0] != 0 || sorted[0][1] != 1 {
		t.Errorf("Wrong first bucket: %v", sorted[0])
	}
	if len(sorted[1]) != 0 {
		t.Errorf("Wrong second bucket lenght: %d", len(sorted[1]))
	}
	if len(sorted[2]) != 1 {
		t.Errorf("Wrong length of third bucket: %d", len(sorted[2]))
	}
	if sorted[2][0] != 2 {
		t.Errorf("Wrong third bucket: %v", sorted[1])
	}
}
