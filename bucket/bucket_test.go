package bucket

import "testing"

type Sizeable struct {
	size int
}

func (s Sizeable) Size() int {
	return s.size
}

var testData = []Sizeable{{1}, {2}, {15}, {3}, {1}, {15}, {15}}

func TestBucketing(t *testing.T) {
	testSizers := make([]Sizer, len(testData))
	for i := range testData {
		testSizers[i] = testData[i]
	}
	buckets := bucketBySize(testSizers)
	if len(buckets) != 2 {
		t.Errorf("Wrong bucket length: %d", len(buckets))
	}

	for _, b := range buckets {
		switch b[0].Size() {
		case 1:
			if len(b) != 2 {
				t.Errorf("Wrong number of size 1: %d", len(b))
			}
		case 15:
			if len(b) != 3 {
				t.Errorf("Wrong number of size 15: %d", len(b))
			}
		default:
			t.Errorf("Wrong case or size")
		}
	}
}
