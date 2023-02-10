// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64

package slices

import (
	"golang.org/x/exp/constraints"
)

func tryBlockIntroSort[E constraints.Ordered](x []E) bool {
	return false
}
