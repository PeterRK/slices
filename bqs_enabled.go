// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package slices

import (
	"cmp"
	"unsafe"
)

const bqsSize = 1024

func tryBlockIntroSort[E cmp.Ordered](list []E) bool {
	var elem E
	var word uintptr
	if unsafe.Sizeof(elem) > unsafe.Sizeof(word) ||
		unsafe.Sizeof(elem) < 2 || len(list) < bqsSize {
		return false
	}
	chance := log2Ceil(uint(len(list))) * 2
	blockIntroSort(list, chance)
	return true
}

func blockIntroSort[E cmp.Ordered](list []E, chance int) {
	for len(list) >= bqsSize {
		if chance--; chance < 0 {
			heapSort(list)
			return
		}
		m := blockPartition(list)
		if m < 0 {
			return
		}
		blockIntroSort(list[m:], chance)
		list = list[:m]
	}
	introSort(list, chance)
}

func compGE[E cmp.Ordered](a, b E) int {
	if a < b {
		return 0
	} else {
		return 1
	}
}

func blockPartition[E cmp.Ordered](list []E) int {
	size := len(list) // size >= 16

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
		return -1
	}

	l, r := 0, size-1

	const blockSize = 64
	var ml, mr struct {
		v [blockSize]uint8
		a int
		b int
	}
	for r-l >= blockSize*2 {
		if ml.a == ml.b {
			ml.a, ml.b = 0, 0
			for i := 0; i < blockSize; i++ {
				ml.v[ml.b] = uint8(i)
				ml.b += compGE(list[l+i], pivot)
			}
		}
		if mr.a == mr.b {
			mr.a, mr.b = 0, 0
			for i := 0; i < blockSize; i++ {
				mr.v[mr.b] = uint8(i)
				mr.b += compGE(pivot, list[r-i])
			}
		}
		sz := min(ml.b-ml.a, mr.b-mr.a)
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
	}
	if mr.a != mr.b {
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

	for {
		for list[l] < pivot {
			l++
		}
		for list[r] > pivot {
			r--
		}
		if l >= r {
			break
		}
		list[l], list[r] = list[r], list[l]
		l++
		r--
	}
	return l
}
