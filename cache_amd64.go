// Code generated from sort.go using genzfunc.go; DO NOT EDIT.

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package slices

//go:noescape
func cpuid(eax, ecx uint32) (a, b, c, d uint32)

func init() {
	a, b, c, d := cpuid(0, 0)
	level := a
	if level < 1 {
		return
	}
	venderBytes := [12]byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(d), byte(d >> 8), byte(d >> 16), byte(d >> 24),
		byte(c), byte(c >> 8), byte(c >> 16), byte(c >> 24),
	}
	vender := string(venderBytes[:])

	extLv, _, _, _ := cpuid(0x80000000, 0)
	if extLv < 0x80000006 {
		return
	}

	_, b, _, _ = cpuid(1, 0)
	cores := int(b>>16) & 0xff

	_, _, c, d = cpuid(0x80000006, 0)
	cacheLineSize := int(c & 0xff)
	l2CacheSize := int(c>>16) * 1024

	if cacheLineSize <= 0 || l2CacheSize <= 0 {
		return
	}
	cacheInfo.lineSize = cacheLineSize
	cacheInfo.available = l2CacheSize

	cmd := uint32(0)
	if vender == "GenuineIntel" {
		if level >= 4 {
			cmd = 4
		}
	} else if vender == "AuthenticAMD" {
		if extLv >= 0x8000001D {
			cmd = 0x8000001D
		}
	}

	if cmd == 0 {
		return
	}
	a, b, c, _ = cpuid(cmd, 3)
	cacheLevel := (a >> 5) & 7
	if cacheLevel != 3 {
		return
	}
	lineSize := int(b&0xfff) + 1
	partitions := int((b>>12)&0x3ff) + 1
	associativity := int((b>>22)&0x3ff) + 1
	sets := int(c) + 1
	l3CacheSize := lineSize * partitions * associativity * sets

	if cores <= 0 || l3CacheSize <= 0 {
		return
	}
	l3CachePerCore := l3CacheSize / cores
	l3Available := l3CachePerCore / 2
	if l3Available > cacheInfo.available {
		cacheInfo.available = l3Available
	}
}
