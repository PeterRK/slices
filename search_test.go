// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "testing"

var data = []int{0: -10, 1: -5, 2: 0, 3: 1, 4: 2, 5: 3, 6: 5, 7: 7, 8: 11, 9: 100, 10: 100, 11: 100, 12: 1000, 13: 10000}
var fdata = []float64{0: -3.14, 1: 0, 2: 1, 3: 2, 4: 1000.7}
var sdata = []string{0: "f", 1: "foo", 2: "foobar", 3: "x"}

var tests = []struct {
	name   string
	result int
	i      int
}{
	{"SearchInts", BinarySearch(data, 11), 8},
	{"SearchFloat64s", BinarySearch(fdata, 2.1), 4},
	{"SearchStrings", BinarySearch(sdata, ""), 0},
	{"IntSlice.Search", naiveOrder[int](false).BinarySearch(data, 0), 2},
	{"Float64Slice.Search", naiveOrder[float64](false).BinarySearch(fdata, 2.0), 3},
	{"StringSlice.Search", naiveOrder[string](false).BinarySearch(sdata, "x"), 3},
}

func TestBinarySearch(t *testing.T) {
	for _, e := range tests {
		if e.result != e.i {
			t.Errorf("%s: expected index %d; got %d", e.name, e.i, e.result)
		}
	}
}

func BenchmarkBinarySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BinarySearch(data, 11)
		BinarySearch(fdata, 2.1)
		BinarySearch(sdata, "")
	}
}
