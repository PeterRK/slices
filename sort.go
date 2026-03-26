// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"cmp"
	"unsafe"
)

// Sort sorts a slice of any ordered type in ascending order.
// When sorting floating-point numbers, NaNs are ordered before other values.
func Sort[E cmp.Ordered](list []E) {
	if !tryBlockIntroSort(list) {
		sortFast(list)
	}
}

// SortStableFunc sorts the slice x while keeping the original order of equal
// elements, using cmp to compare elements.
func SortStable[E cmp.Ordered](list []E) {
	sortStable(list, true)
}

// PartlySort moves the smallest k elements to list[:k] and sorts that prefix.
func PartlySort[E cmp.Ordered](list []E, k int) {
	partlySort(list, k)
}

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[E cmp.Ordered](list []E) bool {
	return isSorted(list)
}

// Min returns the minimal value in x. It panics if x is empty.
// For floating-point numbers, Min propagates NaNs (any NaN value in x
// forces the output to be NaN).
func Min[E cmp.Ordered](list []E) E {
	return findMin(list)
}

// Max returns the maximal value in x. It panics if x is empty.
// For floating-point E, Max propagates NaNs (any NaN value in x
// forces the output to be NaN).
func Max[E cmp.Ordered](list []E) E {
	return findMax(list)
}

// BinarySearch searches for target in a sorted slice and returns the position
// where target is found, or the position where target would appear in the
// sort order; it also returns a bool saying whether the target is really found
// in the slice. The slice must be sorted in increasing order.
func BinarySearch[E cmp.Ordered](x []E, target E) (int, bool) {
	return binarySearch(x, target)
}

type lessFunc[E any] func(a, b E) bool
type refLessFunc[E any] func(a, b *E) bool

// Order record the way of comparison, will never changed by its methods.
//
// .Less is a comparison function with value input.
// .RefLess is a comparison function with pointer input.
// At least one of them should be set before use.
// If both of them are set, they must have the same behavior.
type Order[E any] struct {
	Less    func(a, b E) bool
	RefLess func(a, b *E) bool
}

func isSmallUnit[E any]() bool {
	var elem E
	var word uintptr
	return unsafe.Sizeof(elem) <= unsafe.Sizeof(word)*2
}

// The general version of BinarySearch.
func (od *Order[E]) BinarySearch(list []E, target E) (int, bool) {
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		return refLessFunc[E](od.RefLess).binarySearch(list, target)
	}
	return lessFunc[E](od.Less).binarySearch(list, target)
}

// The general version of Min.
func (od *Order[E]) Min(list []E) E {
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		return refLessFunc[E](od.RefLess).findMin(list)
	}
	return lessFunc[E](od.Less).findMin(list)
}

// The general version of Max.
func (od *Order[E]) Max(list []E) E {
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		return refLessFunc[E](od.RefLess).findMax(list)
	}
	return lessFunc[E](od.Less).findMax(list)
}

// The general version of IsSorted.
func (od *Order[E]) IsSorted(list []E) bool {
	if len(list) < 2 {
		return true
	}
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		return refLessFunc[E](od.RefLess).isSorted(list)
	}
	return lessFunc[E](od.Less).isSorted(list)
}

// The general version of Sort.
// It tends to use faster algorithm with extra memory.
func (od *Order[E]) Sort(list []E) {
	od.SortWithOption(list, false, false)
}

// The general version of SortStable.
// It tends to use faster algorithm with extra memory.
func (od *Order[E]) SortStable(list []E) {
	od.SortWithOption(list, true, false)
}

// PartlySort moves the smallest k elements to list[:k] and sorts that prefix.
func (od *Order[E]) PartlySort(list []E, k int) {
	if len(list) < 2 || k <= 0 {
		return
	}
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		refLessFunc[E](od.RefLess).partlySort(list, k)
		return
	}
	lessFunc[E](od.Less).partlySort(list, k)
}

var cacheInfo = struct {
	lineSize  int
	available int
}{
	lineSize:  64,         //most common cache line size
	available: 256 * 1024, //available bytes for sort
}

// The general sort function.
// Guarantee stability when stable flag is set.
// Avoid allocating O(n) size extra memory when inplace flag is set.
func (od *Order[E]) SortWithOption(list []E, stable, inplace bool) {
	if len(list) < 2 {
		return
	}
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		elemSize := int(unsafe.Sizeof(list[0]))
		wordSize := int(unsafe.Sizeof(uintptr(0)))
		footprint := elemSize
		if footprint > cacheInfo.lineSize {
			footprint = cacheInfo.lineSize
		}
		footprint += wordSize
		// movement is cheap for small data
		// random access is expensive for big data
		noRefSort := elemSize*len(list) < 1024 ||
			footprint*len(list) > cacheInfo.available
		if stable {
			if inplace {
				refLessFunc[E](od.RefLess).sortStable(list, true)
				return
			}
			if noRefSort {
				refLessFunc[E](od.RefLess).sortStable(list, false)
				return
			}
		} else if elemSize <= wordSize*4 || noRefSort || inplace {
			//slower than ref mode, but no extra allocation
			refLessFunc[E](od.RefLess).sortFast(list)
			return
		}

		// sort by pointer list, fast in cache
		ref := make([]*E, len(list))
		for i := 0; i < len(list); i++ {
			ref[i] = &list[i]
		}
		if stable {
			lessFunc[*E](od.RefLess).sortStable(ref, false)
		} else {
			lessFunc[*E](od.RefLess).sortFast(ref)
		}
		reorder(list, ref)
		return
	}
	if stable {
		lessFunc[E](od.Less).sortStable(list, inplace)
	} else {
		lessFunc[E](od.Less).sortFast(list)
	}
}

func ptrDiff[T any](a, b *T) int {
	diff := int(uintptr(unsafe.Pointer(a)) - uintptr(unsafe.Pointer(b)))
	return diff / int(unsafe.Sizeof(*a))
}

// ref contains pointers to element of list.
// Reorder elements in list according to the pointer order in ref.
// It's unsafe, should not use it elsewhere.
func reorder[E any](list []E, ref []*E) {
	size := len(list)
	for i := 0; i < size; i++ {
		if ref[i] == nil {
			continue
		}
		j := ptrDiff(ref[i], &list[0])
		ref[i] = nil
		if j == i {
			continue
		}
		k, tmp := i, list[i]
		for j != i {
			k, list[k] = j, list[j]
			j = ptrDiff(ref[k], &list[0])
			ref[k] = nil
		}
		list[k] = tmp
	}
}
