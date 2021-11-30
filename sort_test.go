// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"constraints"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8}
var strings = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestInts(t *testing.T) {
	data := ints
	Sort(data[:])
	if !IsSorted(data[:]) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestFloat64s(t *testing.T) {
	data := float64s
	Sort(data[:])
	if !IsSorted(data[:]) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestStrings(t *testing.T) {
	data := strings
	Sort(data[:])
	if !IsSorted(data[:]) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
}

func TestHeapSort(t *testing.T) {
	data := ints
	heapSort(data[:])
	if !IsSorted(data[:]) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortStable(t *testing.T) {
	data := ints
	heapSort(data[:])
	if !IsSorted(data[:]) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortLarge_Random(t *testing.T) {
	n := 100000
	if testing.Short() {
		n /= 10
	}
	data := make([]int, n)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100)
	}
	if IsSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sort didn't sort - 100k ints")
	}
}

func TestReverseSortIntSlice(t *testing.T) {
	data1 := ints
	data2 := ints
	naiveOrder[int](false).Sort(data1[:])
	naiveOrder[int](true).Sort(data2[:])
	for i := 0; i < len(ints); i++ {
		if data1[i] != data2[len(ints)-1-i] {
			t.Errorf("reverse sort didn't sort")
		}
		if i > len(ints)/2 {
			break
		}
	}
}

const (
	_Sawtooth = iota
	_Rand
	_Stagger
	_Plateau
	_Shuffle
	_NDist
)

const (
	_Copy = iota
	_Reverse
	_ReverseFirstHalf
	_ReverseSecondHalf
	_Sorted
	_Dither
	_NMode
)

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func testBentleyMcIlroy(t *testing.T, stable, inplace bool) {
	sizes := []int{100, 1023, 1024, 1025}
	if testing.Short() {
		sizes = []int{100, 127, 128, 129}
	}
	dists := []string{"sawtooth", "rand", "stagger", "plateau", "shuffle"}
	modes := []string{"copy", "reverse", "reverse1", "reverse2", "sort", "dither"}
	od := naiveOrder[int](false)
	var tmp1, tmp2 [1025]int
	for _, n := range sizes {
		for m := 1; m < 2*n; m *= 2 {
			for dist := 0; dist < _NDist; dist++ {
				j := 0
				k := 1
				data := tmp1[0:n]
				for i := 0; i < n; i++ {
					switch dist {
					case _Sawtooth:
						data[i] = i % m
					case _Rand:
						data[i] = rand.Intn(m)
					case _Stagger:
						data[i] = (i*m + i) % n
					case _Plateau:
						data[i] = min(i, m)
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = j
						} else {
							k += 2
							data[i] = k
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = data[n-i-1]
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = data[n/2-i-1]
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = data[i]
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = data[i]
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = data[n-(i-n/2)-1]
						}
					case _Sorted:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
						// Ints is known to be correct
						// because mode Sort runs after mode _Copy.
						od.Sort(mdata)
					case _Dither:
						for i := 0; i < n; i++ {
							mdata[i] = data[i] + i%5
						}
					}

					sum, xor := 0, 0
					for i := 0; i < n; i++ {
						sum += mdata[i]
						xor ^= mdata[i]
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])
					ncmp := 0
					xod := Order[int]{Less: func(a, b int) bool {
						ncmp++
						return a < b
					}}
					xod.SortWithOption(mdata[:n], stable, inplace)
					// Uncomment if you are trying to improve the number of compares.
					//t.Logf("%s: ncmp=%d", desc, ncmp)
					for i := 0; i < n; i++ {
						sum -= mdata[i]
						xor ^= mdata[i]
					}
					if !od.IsSorted(mdata) || sum != 0 || xor != 0 {
						t.Fatalf("%s: ints not sorted\n\t%v", desc, mdata)
					}
				}
			}
		}
	}
}

func TestSortBM(t *testing.T)              { testBentleyMcIlroy(t, false, false) }
func TestSortStableBM(t *testing.T)        { testBentleyMcIlroy(t, true, false) }
func TestSortStableInplaceBM(t *testing.T) { testBentleyMcIlroy(t, true, true) }

type intPair struct {
	a, b int
}

func testStability(t *testing.T, ref, inplace bool) {
	n, m := 100000, 1000
	if testing.Short() {
		n, m = 1000, 100
	}
	data := make([]intPair, n)

	od := Order[intPair]{Less: func(x, y intPair) bool {
		return x.a < y.a
	}}
	if ref {
		od = Order[intPair]{RefLess: func(x, y *intPair) bool {
			return x.a < y.a
		}}
	}
	xod := Order[intPair]{Less: func(x, y intPair) bool {
		if x.a < y.a {
			return true
		} else if x.a == y.a {
			return x.b < y.b
		}
		return false
	}}

	// random distribution
	for i := 0; i < len(data); i++ {
		data[i].a = rand.Intn(m)
		data[i].b = i
	}
	if od.IsSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	od.SortWithOption(data, true, inplace)
	if !xod.IsSorted(data) {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}

	// already sorted
	for i := 0; i < len(data); i++ {
		data[i].b = i
	}
	od.SortWithOption(data, true, inplace)
	if !xod.IsSorted(data) {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}

	// sorted reversed
	for i := 0; i < len(data); i++ {
		data[i].a = (len(data) - i) / m
		data[i].b = i
	}
	od.SortWithOption(data, true, inplace)
	if !xod.IsSorted(data) {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}
}

func TestStability(t *testing.T) {
	testStability(t, false, false)
	testStability(t, true, false)
}
func TestStabilityInplace(t *testing.T) {
	testStability(t, false, true)
}

var countOpsSizes = []int{1e2, 3e2, 1e3, 3e3, 1e4, 3e4, 1e5, 3e5, 1e6}

func countOps(t *testing.T, stable, inplace bool) {
	sizes := countOpsSizes
	if testing.Short() {
		sizes = sizes[:5]
	}
	if !testing.Verbose() {
		t.Skip("Counting skipped as non-verbose mode.")
	}
	ncmp := 0
	od := Order[int]{Less: func(a, b int) bool {
		ncmp++
		return a < b
	}}
	for _, n := range sizes {
		data := make([]int, n)
		for i := 0; i < n; i++ {
			data[i] = rand.Intn(n)
		}
		ncmp = 0
		name := "Sort"
		if stable {
			name = "StableSort"
		}
		if inplace {
			name += "(inplace)"
		}
		od.SortWithOption(data, stable, inplace)
		t.Logf("%s %8d elements: %10d Less", name, n, ncmp)
	}
}

func TestCountSortOps(t *testing.T)              { countOps(t, false, false) }
func TestCountSortStableOps(t *testing.T)        { countOps(t, true, false) }
func TestCountSortStableInplaceOps(t *testing.T) { countOps(t, true, true) }

type object struct {
	val int
}
type smallObject struct {
	object
	pad [15]byte
}
type bigObject struct {
	object
	pad [199]byte
}

func (o object) String() string {
	return strconv.Itoa(o.val)
}

func testSortObject(t *testing.T, stable, inplace bool) {
	n := 100000
	if testing.Short() {
		n = 10000
	}
	data1 := make([]smallObject, n)
	data2 := make([]bigObject, n)
	for i := 0; i < n; i++ {
		val := rand.Intn(n)
		data1[i].val = val
		data2[i].val = val
	}
	od1 := Order[smallObject]{Less: func(a, b smallObject) bool {
		return a.val < b.val
	}, RefLess: func(a, b *smallObject) bool {
		return a.val < b.val
	}}

	od1.SortWithOption(data1, stable, inplace)
	if !od1.IsSorted(data1) {
		t.Errorf("small objects didn't sort")
	}

	od2 := Order[bigObject]{Less: func(a, b bigObject) bool {
		return a.val < b.val
	}, RefLess: func(a, b *bigObject) bool {
		return a.val < b.val
	}}
	od2.SortWithOption(data2, stable, inplace)
	if !od2.IsSorted(data2) {
		t.Errorf("big objects didn't sort")
	}
}

func TestSortObject(t *testing.T)              { testSortObject(t, false, false) }
func TestSortObjectInplace(t *testing.T)       { testSortObject(t, false, true) }
func TestSortObjectStable(t *testing.T)        { testSortObject(t, true, false) }
func TestSortObjectStableInplace(t *testing.T) { testSortObject(t, true, true) }

func benchString(b *testing.B, size int, stable, inplace bool) {
	b.StopTimer()
	unsorted := make([]string, size)
	for i := 0; i < size; i++ {
		unsorted[i] = strconv.Itoa(rand.Int())
	}
	data := make([]string, size)
	for i := 0; i < b.N; i++ {
		copy(data, unsorted)
		b.StartTimer()
		if stable {
			sortStable(data, inplace)
		} else {
			sort(data)
		}
		b.StopTimer()
	}
}

func BenchmarkSortString1K(b *testing.B)              { benchString(b, 1<<10, false, false) }
func BenchmarkSortStableString1K(b *testing.B)        { benchString(b, 1<<10, true, false) }
func BenchmarkSortStableInplaceString1K(b *testing.B) { benchString(b, 1<<10, true, true) }

func benchInt(b *testing.B, size int, stable, inplace bool) {
	b.StopTimer()
	unsorted := make([]int, size)
	for i := 0; i < size; i++ {
		unsorted[i] = rand.Int()
	}
	data := make([]int, size)
	for i := 0; i < b.N; i++ {
		copy(data, unsorted)
		b.StartTimer()
		if stable {
			sortStable(data, inplace)
		} else {
			sort(data)
		}
		b.StopTimer()
	}
}

func BenchmarkSortInt1K(b *testing.B)              { benchInt(b, 1<<10, false, false) }
func BenchmarkSortStableInt1K(b *testing.B)        { benchInt(b, 1<<10, true, false) }
func BenchmarkSortStableInplaceInt1K(b *testing.B) { benchInt(b, 1<<10, true, true) }

func BenchmarkSortInt64K(b *testing.B)              { benchInt(b, 1<<16, false, false) }
func BenchmarkSortStableInt64K(b *testing.B)        { benchInt(b, 1<<16, true, false) }
func BenchmarkSortStableInplaceInt64K(b *testing.B) { benchInt(b, 1<<16, true, true) }

func benchSmallObject(b *testing.B, size int, stable, inplace bool) {
	b.StopTimer()
	od := Order[smallObject]{Less: func(a, b smallObject) bool {
		return a.val < b.val
	}, RefLess: func(a, b *smallObject) bool {
		return a.val < b.val
	}}
	unsorted := make([]int, size)
	for i := range unsorted {
		unsorted[i] = rand.Int()
	}
	data := make([]smallObject, size)
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(data); i++ {
			data[i].val = unsorted[i]
		}
		b.StartTimer()
		od.SortWithOption(data, stable, inplace)
		b.StopTimer()
	}
}

func BenchmarkSortSmallObject64K(b *testing.B)   { benchSmallObject(b, 1<<16, false, false) }
func BenchmarkStableSmallObject64K(b *testing.B) { benchSmallObject(b, 1<<16, true, false) }

func benchBigObject(b *testing.B, size int, stable, inplace bool) {
	b.StopTimer()
	od := Order[bigObject]{Less: func(a, b bigObject) bool {
		return a.val < b.val
	}, RefLess: func(a, b *bigObject) bool {
		return a.val < b.val
	}}
	unsorted := make([]int, size)
	for i := range unsorted {
		unsorted[i] = rand.Int()
	}
	data := make([]bigObject, size)
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(data); i++ {
			data[i].val = unsorted[i]
		}
		b.StartTimer()
		od.SortWithOption(data, stable, inplace)
		b.StopTimer()
	}
}

func BenchmarkSortBigObject64K(b *testing.B)              { benchBigObject(b, 1<<16, false, false) }
func BenchmarkSortInplaceBigObject64K(b *testing.B)       { benchBigObject(b, 1<<16, false, true) }
func BenchmarkSortStableBigObject64K(b *testing.B)        { benchBigObject(b, 1<<16, true, false) }
func BenchmarkSortStableInplaceBigObject64K(b *testing.B) { benchBigObject(b, 1<<16, true, true) }

func benchIntArray(b *testing.B, size int, stable, inplace bool) {
	b.StopTimer()
	od := Order[int]{Less: func(a, b int) bool { return a > b }}
	data := make([]int, size)
	x := ^uint32(0)
	for i := 0; i < b.N; i++ {
		for n := size - 3; n <= size+3; n++ {
			for i := 0; i < len(data); i++ {
				x += x
				x ^= 1
				if int32(x) < 0 {
					x ^= 0x88888eef
				}
				data[i] = int(x % uint32(n/5))
			}
			b.StartTimer()
			od.SortWithOption(data, stable, inplace)
			b.StopTimer()
		}
	}
}

func BenchmarkSort1e2(b *testing.B)       { benchIntArray(b, 1e2, false, false) }
func BenchmarkSortStable1e2(b *testing.B) { benchIntArray(b, 1e2, true, false) }
func BenchmarkSort1e4(b *testing.B)       { benchIntArray(b, 1e4, false, false) }
func BenchmarkSortStable1e4(b *testing.B) { benchIntArray(b, 1e4, true, false) }
func BenchmarkSort1e6(b *testing.B)       { benchIntArray(b, 1e6, false, false) }
func BenchmarkSortStable1e6(b *testing.B) { benchIntArray(b, 1e6, true, false) }
