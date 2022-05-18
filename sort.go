// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"golang.org/x/exp/constraints"
	"unsafe"
)

// Sort sorts a slice of any ordered type in ascending order.
// Sort may fail to sort correctly when sorting slices of floating-point
// numbers containing Not-a-number (NaN) values.
// Use slices.SortFunc(x, func(a, b float64) bool {return a < b || (math.IsNaN(a) && !math.IsNaN(b))})
// instead if the input may contain NaNs.
func Sort[E constraints.Ordered](x []E) {
	sortFast(x)
}

// SortStable sorts the slice x while keeping the original order of equal
func SortStable[E constraints.Ordered](x []E) {
	sortStable(x, true)
}

// SortFunc sorts the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
//
// SortFunc requires that less is a strict weak ordering.
// See https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.
func SortFunc[E any](x []E, less func(a, b E) bool) {
	lessFunc[E](less).sortFast(x)
}

// SortStable sorts the slice x while keeping the original order of equal
// elements, using less to compare elements.
func SortStableFunc[E any](x []E, less func(a, b E) bool) {
	lessFunc[E](less).sortStable(x, true)
}

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[E constraints.Ordered](x []E) bool {
	return isSorted(x)
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function.
func IsSortedFunc[E any](x []E, less func(a, b E) bool) bool {
	return lessFunc[E](less).isSorted(x)
}

// BinarySearch searches for target in a sorted slice and returns the position
// where target is found, or the position where target would appear in the
// sort order; it also returns a bool saying whether the target is really found
// in the slice. The slice must be sorted in increasing order.
func BinarySearch[E constraints.Ordered](x []E, target E) (int, bool) {
	return binarySearch(x, target)
}

// BinarySearchFunc works like BinarySearch, but uses a custom comparison
// function. The slice must be sorted in increasing order, where "increasing" is
// defined by the less function.
func BinarySearchFunc[E any](x []E, target E, less func(a, b E) bool) (int, bool) {
	return lessFunc[E](less).binarySearch(x, target)
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
