// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package slices

import (
	"constraints"
	"unsafe"
)

func branchEmilinatable[T constraints.Ordered]() bool {
	var t T
	return unsafe.Sizeof(t) <= unsafe.Sizeof(uintptr(0))
}