// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "constraints"

func (lt refLessFunc[E]) binarySearch(list []E, x E) int {
	a, b := 0, len(list)
	for a < b {
		m := int(uint(a+b) / 2)
		if lt(&list[m], &x) {
			a = m + 1
		} else {
			b = m
		}
	}
	return a
}

func naiveRefLess[E constraints.Ordered](a, b *E) bool {
	return *a < *b
}

// Search one E in sorted list, return index that list[index] >= x.
// The retuned index can be len(list), but never negtive.
func BinarySearch[E constraints.Ordered](list []E, x E) int {
	return refLessFunc[E](naiveRefLess[E]).binarySearch(list, x)
}
