// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64

package slices

import (
	"cmp"
)

func tryBlockIntroSort[E cmp.Ordered](x []E) bool {
	return false
}
