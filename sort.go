// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run genzfunc.go

package slices

import (
	"constraints"
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
	if branchEmilinatable[E]() {
		blockIntroSort(list, log2Ceil(uint(len(list)))*2)
	} else {
		sort(list)
	}
}

func sort[E constraints.Ordered](list []E) {
	chance := log2Ceil(uint(len(list))) * 3 / 2
	introSort(list, chance)
}

// StableSort sorts data while keeping the original order of equal elements.
// It makes O(n*log(n)) calls to less function.
func SortStable[E constraints.Ordered](list []E) {
	sortStable(list, true)
}

// Avoid allocating O(n) size extra memory when inplace flag is set.
func sortStable[E constraints.Ordered](list []E, inplace bool) {
	if size := len(list); inplace {
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
	} else if size < 16 {
		simpleSort(list)
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

// no codegen
func sortIndex3[T constraints.Ordered](list []T, a, b, c int) (int, int, int) {
	// keep stable
	if list[a] > list[b] {
		if list[b] > list[c] {
			return c, b, a
		} else if list[a] > list[c] {
			return b, c, a
		} else {
			return b, a, c
		}
	} else {
		if list[a] > list[c] {
			return c, a, b
		} else if list[b] > list[c] {
			return a, c, b
		} else {
			return a, b, c
		}
	}
}

// no codegen
// return int instead of bool
// should be inlined
func cmpGT[T constraints.Ordered](a, b T) int {
	if a > b {
		return 1
	} else {
		return 0
	}
}

// no codegen
// the block partition algorithm from Edelkamp, Stefan, and Armin Weiß.
// "Blockquicksort: Avoiding branch mispredictions in quicksort."
// Journal of Experimental Algorithmics (JEA) 24 (2019): 1-22.
func blockPartition[T constraints.Ordered](list []T, hard bool) int {
	size := len(list)
	x, s := size/2, size-1
	var a, m, b int
	if hard {
		a, m, b = sortIndex3(list, 1, x, s-1)
		if size > 32 {
			y, z := size/4, size/8
			_, a, _ = sortIndex3(list, x-y, x-1, x+z)
			_, b, _ = sortIndex3(list, x-z, x+1, x+y)
			a, m, b = sortIndex3(list, a, m, b)
		}
	} else {
		if size > 128 {
			y, z := size/4, size/8
			_, m, _ = sortIndex3(list, 1, x, s-1)
			_, a, _ = sortIndex3(list, x-y, x-1, x+z)
			_, b, _ = sortIndex3(list, x-z, x+1, x+y)
			a, m, b = sortIndex3(list, a, m, b)
		} else {
			a, m, b = sortIndex3(list, x-1, x, x+1)
		}
		if list[0] > list[s] {
			// may convert descent array to ascent array
			list[0], list[s] = list[s], list[0]
		}
	}

	pivot := list[m]
	list[0], list[a] = list[a], list[0]
	list[s], list[b] = list[b], list[s]

	l, r := 1, s-1
	pattern := 0 // try to detect ascent, descent, constant
	// 0: constant
	// 1: partitioned, maybe ascent
	// 2: reverse partitioned, maybe descent
	// 3: unordered

	// with branch elimination
	// complicatied but fast in some superscalar machine
	const blockSize = 64
	if r-l > blockSize*2-1 {
		// branch elimination may be faster only in unordered pattern
		for pattern != 3 {
			for list[l] < pivot {
				l++
				pattern |= 1
			}
			for list[r] > pivot {
				r--
				pattern |= 1
			}
			if l >= r {
				goto finish
			}
			list[l], list[r] = list[r], list[l]
			if (pattern&2) == 0 && list[l] != list[r] {
				pattern |= 2
			}
			l++
			r--
		}
		var ml, mr struct {
			v [blockSize]uint8
			a int
			b int
		}
		for r-l > blockSize*2-1 {
			if ml.a == ml.b {
				ml.a, ml.b = 0, 0
				for i := 0; i < blockSize; i++ {
					ml.v[ml.b] = uint8(i)
					ml.b += cmpGT(list[l+i], pivot)
				}
			}
			if mr.a == mr.b {
				mr.a, mr.b = 0, 0
				for i := 0; i < blockSize; i++ {
					mr.v[mr.b] = uint8(i)
					mr.b += cmpGT(pivot, list[r-i])
				}
			}
			sz := ml.b - ml.a
			if t := mr.b - mr.a; t < sz {
				sz = t
			}
			for i := 0; i < sz; i++ {
				ll := l + int(ml.v[ml.a])
				ml.a++
				rr := r - int(mr.v[mr.a])
				mr.a++
				list[ll], list[rr] = list[rr], list[ll]
			}
			if ml.a == ml.b {
				l += blockSize
			}
			if mr.a == mr.b {
				r -= blockSize
			}
		}
		if ml.a != ml.b {
			for {
				for list[r] > pivot {
					r--
				}
				ll := l + int(ml.v[ml.a])
				// list[r] <= pivot
				// list[r+1] > pivot
				if ll >= r {
					return r + 1
				}
				list[ll], list[r] = list[r], list[ll]
				r--
				// list[r] ?
				// list[r+1] >= pivot
				if ml.a++; ml.a == ml.b {
					l += blockSize
					if l > r {
						return r + 1
					}
					break
				}
			}
		} else if mr.a != mr.b {
			for {
				for list[l] < pivot {
					l++
				}
				rr := r - int(mr.v[mr.a])
				// list[l] >= pivot
				// list[l-1] < pivot
				if l >= rr {
					return l
				}
				list[l], list[rr] = list[rr], list[l]
				l++
				// list[l] ?
				// list[l-1] <= pivot
				if mr.a++; mr.a == mr.b {
					r -= blockSize
					if l > r {
						return l
					}
					break
				}
			}
		}
	}

	for {
		for list[l] < pivot {
			l++
			pattern |= 1
		}
		for list[r] > pivot {
			r--
			pattern |= 1
		}
		if l >= r {
			break
		}
		list[l], list[r] = list[r], list[l]
		if (pattern&2) == 0 && list[l] != list[r] {
			pattern |= 2
		}
		l++
		r--
	}
finish:
	if pattern == 3 {
		// common case
	} else if pattern == 0 {
		// list[0] <= pivot <= list[s]
		// values in list[1:s-1] are all pivot
		return -1
	} else if pattern == 1 {
		if a < l && l <= b {
			// rollback then check
			list[0], list[a] = list[a], list[0]
			list[s], list[b] = list[b], list[s]
			for i := 0; i < s; i++ {
				if list[i] > list[i+1] {
					return l
				}
			}
			return -1
		}
	} else if pattern == 2 {
		// fix up to make ascent sub-arrays if orign array is descent
		if list[l] == pivot {
			list[s], list[l+1] = list[l+1], list[s]
		} else if list[r] == pivot {
			list[0], list[r-1] = list[r-1], list[0]
		}
	}
	return l
}

// no codegen
func blockIntroSort[T constraints.Ordered](list []T, chance int) {
	// blockPartition doesn't work well on all patterns, fallback if necessary 
	for hard := false; len(list) > 12; {
		if chance--; chance < 0 {
			heapSort(list)
			return
		}
		m := blockPartition(list, hard)
		if m < 0 {
			return
		}
		s := len(list) / 8
		if m < s {
			blockIntroSort(list[:m], chance)
			list = list[m:]
			if hard {
				introSort(list, chance)
				return
			}
			hard = true
		} else {
			blockIntroSort(list[m:], chance)
			list = list[:m]
			if m > 7 * s {
				if hard {
					introSort(list, chance)
					return
				}
				hard = true
			} else {
				hard = false
			}
		}
	}
	simpleSort(list)
}