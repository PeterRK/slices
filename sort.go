// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run genzfunc.go

package slices

import (
	"golang.org/x/exp/constraints"
	"math/bits"
)

// Search one E in sorted list, return index that list[index] >= x.
// The retuned index can be len(list), but never negtive.
func BinarySearch[E constraints.Ordered](list []E, x E) int {
	return binarySearch(list, x)
}

func binarySearch[E constraints.Ordered](list []E, x E) int {
	a, b := 0, len(list)
	for a < b {
		m := int(uint(a+b) / 2)
		if less(list[m], x) {
			a = m + 1
		} else {
			b = m
		}
	}
	return a
}

func reverse[E any](list []E) {
	for l, r := 0, len(list)-1; l < r; {
		list[l], list[r] = list[r], list[l]
		l++
		r--
	}
}

func log2Ceil(num uint) int {
	return bits.Len(num)
}

// With small E, double reversion is faster than the BlockSwap rotation.
// BlockSwap rotation needs less swaps, but more branches.
func rotate[E any](list []E, border int) {
	reverse(list[:border])
	reverse(list[border:])
	reverse(list)
}

// It's hoped to be inlined.
func less[E constraints.Ordered](a, b E) bool {
	return a < b
}

// IsSorted reports whether list is sorted.
func IsSorted[E constraints.Ordered](list []E) bool {
	return isSorted(list)
}

func isSorted[E constraints.Ordered](list []E) bool {
	for i := 1; i < len(list); i++ {
		if less(list[i], list[i-1]) {
			return false
		}
	}
	return true
}

// Sort sorts data.
// It makes O(n*log(n)) calls to less function.
// The sort is not guaranteed to be stable.
func Sort[E constraints.Ordered](list []E) {
	sort(list)
}

func sort[E constraints.Ordered](list []E) {
	size := len(list)
	chance := log2Ceil(uint(size)) * 3 / 2
	if size > 50 {
		a, b, c := size/4, size/2, size*3/4
		a, ha := median(list, a-1, a, a+1)
		b, hb := median(list, b-1, b, b+1)
		c, hc := median(list, c-1, c, c+1)
		m, hint := median(list, a, b, c)
		hint &= ha & hb & hc

		pivot := list[m]
		if hint == hintRevered {
			reverse(list)
			hint = hintSorted
		}
		if hint == hintSorted && isSorted(list) {
			return
		}

		l, r := 0, size-1
		for {
			for less(list[l], pivot) {
				l++
			}
			for less(pivot, list[r]) {
				r--
			}
			if l >= r {
				break
			}
			list[l], list[r] = list[r], list[l]
			l++
			r--
		}

		if l > size/2 {
			introSort(list[l:], chance)
			list = list[:l]
		} else {
			introSort(list[:l], chance)
			list = list[l:]
		}
	}
	introSort(list, chance)
}

const (
	hintSorted uint8 = 1 << iota
	hintRevered
)

func median[E constraints.Ordered](list []E, a, b, c int) (int, uint8) {
	// keep stable
	if less(list[b], list[a]) {
		if less(list[c], list[b]) {
			return b, hintRevered //c, b, a
		} else if less(list[c], list[a]) {
			return c, 0 //b, c, a
		} else {
			return a, 0 //b, a, c
		}
	} else {
		if less(list[c], list[a]) {
			return a, 0 //c, a, b
		} else if less(list[c], list[b]) {
			return c, 0 //a, c, b
		} else {
			return b, hintSorted //a, b, c
		}
	}
}

// StableSort sorts data while keeping the original order of equal elements.
// It makes O(n*log(n)) calls to less function.
func SortStable[E constraints.Ordered](list []E) {
	sortStable(list, true)
}

// Avoid allocating O(n) size extra memory when inplace flag is set.
func sortStable[E constraints.Ordered](list []E, inplace bool) {
	if size := len(list); size < 16 {
		simpleSort(list)
	} else if inplace {
		step := 8
		a, b := 0, step
		for b <= size {
			simpleSort(list[a:b])
			a = b
			b += step
		}
		simpleSort(list[a:])

		for step < size {
			a, b = 0, step*2
			for b <= size {
				symmerge(list[a:b], step)
				a = b
				b += step * 2
			}
			if a+step < size {
				symmerge(list[a:], step)
			}
			step *= 2
		}
	} else {
		temp := make([]E, size)
		copy(temp, list)
		mergeSort(temp, list)
	}
}

// A variant of insertion sort for short list.
func simpleSort[E constraints.Ordered](list []E) {
	if len(list) < 2 {
		return
	}
	for i := 1; i < len(list); i++ {
		curr := list[i]
		if less(curr, list[0]) {
			for j := i; j > 0; j-- {
				list[j] = list[j-1]
			}
			list[0] = curr
		} else {
			pos := i
			for ; less(curr, list[pos-1]); pos-- {
				list[pos] = list[pos-1]
			}
			list[pos] = curr
		}
	}
}

func heapSort[E constraints.Ordered](list []E) {
	for idx := len(list)/2 - 1; idx >= 0; idx-- {
		heapDown(list, idx)
	}
	for end := len(list) - 1; end > 0; end-- {
		list[0], list[end] = list[end], list[0]
		heapDown(list[:end], 0)
	}
}

func heapDown[E constraints.Ordered](list []E, pos int) {
	curr := list[pos]
	kid, last := pos*2+1, len(list)-1
	for kid < last {
		if less(list[kid], list[kid+1]) {
			kid++
		}
		if !less(curr, list[kid]) {
			break
		}
		list[pos] = list[kid]
		pos, kid = kid, kid*2+1
	}
	if kid == last && less(curr, list[kid]) {
		list[pos], pos = list[kid], kid
	}
	list[pos] = curr
}

// Sort 5 elemnt in list with 7 comparison.
func sortIndex5[E constraints.Ordered](list []E,
	a, b, c, d, e int) (int, int, int, int, int) {
	if less(list[b], list[a]) {
		a, b = b, a
	}
	if less(list[d], list[c]) {
		c, d = d, c
	}
	if less(list[c], list[a]) {
		a, c = c, a
		b, d = d, b
	}
	if less(list[c], list[e]) {
		if less(list[d], list[e]) {
			if less(list[b], list[d]) {
				if less(list[c], list[b]) {
					return a, c, b, d, e
				} else {
					return a, b, c, d, e
				}
			} else if less(list[b], list[e]) {
				return a, c, d, b, e
			} else {
				return a, c, d, e, b
			}
		} else {
			if less(list[b], list[e]) {
				if less(list[c], list[b]) {
					return a, c, b, e, d
				} else {
					return a, b, c, e, d
				}
			} else if less(list[b], list[d]) {
				return a, c, e, b, d
			} else {
				return a, c, e, d, b
			}
		}
	} else {
		if less(list[b], list[c]) {
			if less(list[e], list[a]) {
				return e, a, b, c, d
			} else if less(list[e], list[b]) {
				return a, e, b, c, d
			} else {
				return a, b, e, c, d
			}
		} else {
			if less(list[a], list[e]) {
				a, e = e, a
			}
			if less(list[d], list[b]) {
				b, d = d, b
			}
			return e, a, c, b, d
		}
	}
}

// triPartition divides list into 3 segments.
// Eents before list[l] are all not greater than it.
// Eents after list[r] are all not less than it.
func triPartition[E constraints.Ordered](list []E) (l, r int) {
	size := len(list)
	m, s := size/2, size/4
	// Get a guide to avoid skewness.
	x, l, _, r, y := sortIndex5(list, m-s, m-1, m, m+1, m+s)

	s = size - 1
	pivotL, pivotR := list[l], list[r]
	list[l], list[r] = list[0], list[s]
	list[1], list[x] = list[x], list[1]
	list[s-1], list[y] = list[y], list[s-1]

	//  | less than pivotL | between pivotL and pivotR | greater than pivotR |
	// 0|                  |l        k -- untested -- r|                     |s

	l, r = 2, s-2
	for {
		for less(list[l], pivotL) {
			l++
		}
		for less(pivotR, list[r]) {
			r--
		}
		if less(pivotR, list[l]) {
			list[l], list[r] = list[r], list[l]
			r--
			if less(list[l], pivotL) {
				l++
				continue
			}
		}
		break
	}

	for k := l + 1; k <= r; k++ {
		if less(pivotR, list[k]) {
			for less(pivotR, list[r]) {
				r--
			}
			if k >= r {
				break
			}
			if less(list[r], pivotL) {
				list[l], list[k], list[r] = list[r], list[l], list[k]
				l++
			} else {
				list[k], list[r] = list[r], list[k]
			}
			r--
		} else if less(list[k], pivotL) {
			list[k], list[l] = list[l], list[k]
			l++
		}
	}

	l--
	r++
	list[0], list[l] = list[l], pivotL
	list[s], list[r] = list[r], pivotR
	return l, r
}

func introSort[E constraints.Ordered](list []E, chance int) {
	for len(list) > 14 {
		if chance--; chance < 0 {
			heapSort(list)
			return
		}
		// Dual pivot quicksort need less memory access, witch makes it faster
		// than single pivot version in many cases, but not always.
		l, r := triPartition(list)
		introSort(list[:l], chance)
		introSort(list[r+1:], chance)
		if !less(list[l], list[r]) {
			return // All emelents in the middle segemnt are equal.
		}
		list = list[l+1 : r]
	}
	simpleSort(list)
}

// symmerge merges the two sorted subsequences data[a:m] and data[m:b] using
// the symmerge algorithm from Pok-Son Kim and Arne Kutzner, "Stable Minimum
// Storage Merging by Symmetric Comparisons", in Susanne Albers and Tomasz
// Radzik, editors, Algorithms - ESA 2004, volume 3221 of Lecture Notes in
// Computer Science, pages 714-723. Springer, 2004.
func symmerge[E constraints.Ordered](list []E, border int) {
	size := len(list)

	// Avoid unnecessary recursions of symmerge by direct insertion.
	if border == 1 {
		curr := list[0]
		a, b := 1, size
		for a < b {
			m := int(uint(a+b) / 2)
			if less(list[m], curr) {
				a = m + 1
			} else {
				b = m
			}
		}
		for i := 1; i < a; i++ {
			list[i-1] = list[i]
		}
		list[a-1] = curr
		return
	}

	// Avoid unnecessary recursions of symmerge by direct insertion.
	if border == size-1 {
		curr := list[border]
		a, b := 0, border
		for a < b {
			m := int(uint(a+b) / 2)
			if less(curr, list[m]) {
				b = m
			} else {
				a = m + 1
			}
		}
		for i := border; i > a; i-- {
			list[i] = list[i-1]
		}
		list[a] = curr
		return
	}

	// Divide list into 3 segments, then handle non-empty ones recursively.
	half := size / 2
	n := border + half
	a, b := 0, border
	if border > half {
		a, b = n-size, half
	}
	// Part of the small piece should be moved to another side.
	// |            |half         |
	// |===|border  |             |
	// |===         |***|n        |
	// |a  |b       |   |         |
	// Keep x-0 == n-y, then x+y == n.
	// It's easy to see the binary search below works
	// when left piece is the small one.
	// Size ceil of left and center segments is border+half.
	// |            |half         |
	// |            |     |border |
	// |    |*******|     |=======|
	// |    |a      |b    |       |
	// When right piece is the small one, size ceil of right and center
	// is (size-border)+(size-half) = size*2-(border+half).
	// size - ceil = (border+half) - size = n - size
	// Keep x-(n-size) == size-y, then x+y == n.
	// Now binary search code can be shared for both cases.
	p := n - 1
	for a < b {
		m := int(uint(a+b) / 2)
		if less(list[p-m], list[m]) { //p-m == (n-m)-1
			b = m
		} else {
			a = m + 1
		}
	}
	b = n - a
	// list[a] > list[b-1] && list[a] <= list[b] && list[b-1] >= list[a-1]
	if a < border && border < b {
		rotate(list[a:b], border-a)
	}
	if 0 < a && a < half {
		symmerge(list[:half], a)
	}
	if half < b && b < size {
		symmerge(list[half:], b-half)
	}
}

func mergeSort[E constraints.Ordered](a, b []E) {
	if size := len(a); size < 12 {
		if size == 0 {
			return
		}
		b[0] = a[0]
		for i := 1; i < size; i++ {
			if curr := a[i]; less(curr, b[0]) {
				for j := i; j > 0; j-- {
					b[j] = b[j-1]
				}
				b[0] = curr
			} else {
				pos := i
				for ; less(curr, b[pos-1]); pos-- {
					b[pos] = b[pos-1]
				}
				b[pos] = curr
			}
		}
	} else {
		half := size / 2
		mergeSort(b[:half], a[:half])
		mergeSort(b[half:], a[half:])

		i, j, k := 0, half, 0
		for ; i < half && j < size; k++ {
			if less(a[j], a[i]) {
				b[k] = a[j]
				j++
			} else {
				b[k] = a[i]
				i++
			}
		}
		for ; i < half; k++ {
			b[k] = a[i]
			i++
		}
		for ; j < size; k++ {
			b[k] = a[j]
			j++
		}
	}
}
