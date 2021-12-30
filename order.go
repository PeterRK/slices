// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "unsafe"

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
func (od *Order[E]) BinarySearch(list []E, x E) int {
	if od.RefLess == nil {
		if od.Less == nil {
			panic("uninitialized Order")
		}
	} else if od.Less == nil || !isSmallUnit[E]() {
		return refLessFunc[E](od.RefLess).binarySearch(list, x)
	}
	return lessFunc[E](od.Less).binarySearch(list, x)
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
		elemSize := unsafe.Sizeof(list[0])
		wordSize := unsafe.Sizeof(uintptr(0))
		big := int(elemSize + wordSize) * len(list) > 256*1024
		if stable {
			if inplace {
				refLessFunc[E](od.RefLess).sortStable(list, true)
				return
			}
			if elemSize <= wordSize*2 || big {
				refLessFunc[E](od.RefLess).sortStable(list, false)
				return
			}
		} else if elemSize <= wordSize*4 || big || inplace {
			//slower than ref mode, but no extra allocation
			refLessFunc[E](od.RefLess).sort(list)
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
			lessFunc[*E](od.RefLess).sort(ref)
		}
		reorder(list, ref)
		return
	}
	if stable {
		lessFunc[E](od.Less).sortStable(list, inplace)
	} else {
		lessFunc[E](od.Less).sort(list)
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

// The general version of BinarySearch.
func BinarySearchFunc[E any](list []E, x E, less func(a, b E) bool) int {
	return lessFunc[E](less).binarySearch(list, x)
}

// The general version of IsSorted.
func IsSortedFunc[E any](list []E, less func(a, b E) bool) bool {
	return lessFunc[E](less).isSorted(list)
}

// The general version of Sort.
func SortFunc[E any](list []E, less func(a, b E) bool) {
	lessFunc[E](less).sort(list)
}

// The general version of SortStable.
func SortStableFunc[E any](list []E, less func(a, b E) bool) {
	lessFunc[E](less).sortStable(list, true)
}
