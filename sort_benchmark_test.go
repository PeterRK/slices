// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"fmt"
	"math"
	"math/rand"
	std "slices"
	"strconv"
	"testing"
)

func randomInts(list []int) {
	size := len(list)
	for i := 0; i < size; i++ {
		list[i] = rand.Intn(size)
	}
}

func constantInts(list []int) {
	size := len(list)
	v := rand.Int()
	for i := 0; i < size; i++ {
		list[i] = v
	}
}

func descentInts(list []int) {
	size := len(list)
	v := rand.Int()
	for i := 0; i < size; i++ {
		list[i] = v + size - i
	}
}

func ascentInts(list []int) {
	size := len(list)
	v := rand.Int()
	for i := 0; i < size; i++ {
		list[i] = v + i
	}
}

func smallInts(list []int) {
	size := len(list)
	limit := int(math.Sqrt(float64(size)))
	if limit < 10 {
		limit = 10
	}
	for i := 0; i < size; i++ {
		list[i] = rand.Intn(limit)
	}
}

func mixedInts(list []int) {
	size := len(list)
	m := size / 5
	smallInts(list[:m])
	constantInts(list[m : m*2])
	ascentInts(list[m*2 : m*3])
	descentInts(list[m*3 : m*4])
	randomInts(list[m*4:])
}

type genFunc struct {
	name string
	fn   func([]int)
}

var pattern = []genFunc{
	genFunc{"Small", smallInts},
	genFunc{"Random", randomInts},
	genFunc{"Constant", constantInts},
	genFunc{"Ascent", ascentInts},
	genFunc{"Descent", descentInts},
	genFunc{"Mixed", mixedInts},
}

type sizeClass struct {
	name string
	size int
}

var level = []sizeClass{
	sizeClass{"1K", 1000},
	sizeClass{"10K", 10_000},
	sizeClass{"100K", 100_000},
	sizeClass{"1M", 1000_000},
}

func benchmarkInt(b *testing.B, sort func([]int)) {
	for _, gen := range pattern {
		for _, sc := range level {
			b.Run(gen.name+"-"+sc.name, func(b *testing.B) {
				b.StopTimer()
				rand.Seed(0)
				list := make([]int, sc.size)
				for i := 0; i < b.N; i++ {
					gen.fn(list)
					b.StartTimer()
					sort(list)
					b.StopTimer()
				}
			})
		}
	}
}

func BenchmarkIntNew(b *testing.B) {
	benchmarkInt(b, Sort[int])
}

func BenchmarkIntStd(b *testing.B) {
	benchmarkInt(b, std.Sort[[]int, int])
}

func benchmarkHybrid(b *testing.B, sort func([]int)) {
	n := 10000
	for _, m := range []int{5, 10, 20, 30, 50} {
		b.Run(fmt.Sprintf("%d%%", m), func(b *testing.B) {
			b.StopTimer()
			rand.Seed(0)
			var all [100][]int
			for j := 0; j < 100; j++ {
				all[j] = make([]int, n)
			}
			a := (102 - m) / 3
			for i := 0; i < b.N; i++ {
				for j := 0; j < m; j++ {
					randomInts(all[j])
				}
				for j := m; j < m+a; j++ {
					ascentInts(all[j])
				}
				for j := m + a; j < m+a*2; j++ {
					descentInts(all[j])
				}
				for j := m + a*2; j < 100; j++ {
					constantInts(all[j])
				}
				b.StartTimer()
				for j := 0; j < 100; j++ {
					sort(all[j])
				}
				b.StopTimer()
			}
		})
	}
}

func BenchmarkHybridNew(b *testing.B) {
	benchmarkHybrid(b, Sort[int])
}

func BenchmarkHybridStd(b *testing.B) {
	benchmarkHybrid(b, std.Sort[[]int, int])
}

func benchmarkFloat(b *testing.B, sort func([]float64)) {
	for _, sc := range level {
		b.Run(sc.name, func(b *testing.B) {
			b.StopTimer()
			rand.Seed(0)
			list := make([]float64, sc.size)
			for i := 0; i < b.N; i++ {
				for j := 0; j < sc.size; j++ {
					list[j] = rand.Float64()
				}
				b.StartTimer()
				sort(list)
				b.StopTimer()
			}
		})
	}
}

func BenchmarkFloatNew(b *testing.B) {
	benchmarkFloat(b, Sort[float64])
}

func BenchmarkFloatStd(b *testing.B) {
	benchmarkFloat(b, std.Sort[[]float64, float64])
}

func benchmarkString(b *testing.B, sort func([]string)) {
	for _, sc := range level {
		b.Run(sc.name, func(b *testing.B) {
			b.StopTimer()
			rand.Seed(0)
			list := make([]string, sc.size)
			for i := 0; i < b.N; i++ {
				for j := 0; j < sc.size; j++ {
					list[j] = strconv.Itoa(rand.Intn(sc.size))
				}
				b.StartTimer()
				sort(list)
				b.StopTimer()
			}
		})
	}
}

func BenchmarkStrNew(b *testing.B) {
	benchmarkString(b, Sort[string])
}

func BenchmarkStrStd(b *testing.B) {
	benchmarkString(b, std.Sort[[]string, string])
}

func benchmarkStruct(b *testing.B, sort func([]smallObject)) {
	for _, sc := range level {
		b.Run(sc.name, func(b *testing.B) {
			b.StopTimer()
			rand.Seed(0)
			list := make([]smallObject, sc.size)
			for i := 0; i < b.N; i++ {
				for j := 0; j < sc.size; j++ {
					list[j].val = rand.Intn(sc.size)
				}
				b.StartTimer()
				sort(list)
				b.StopTimer()
			}
		})
	}
}

func BenchmarkStructNew(b *testing.B) {
	order := Order[smallObject]{
		Less: func(a, b smallObject) bool {
			return a.val < b.val
		}, RefLess: func(a, b *smallObject) bool {
			return a.val < b.val
		}}
	benchmarkStruct(b, func(list []smallObject) {
		order.Sort(list)
	})
}

func BenchmarkStructStd(b *testing.B) {
	benchmarkStruct(b, func(list []smallObject) {
		std.SortFunc[[]smallObject, smallObject](list, func(a, b smallObject) int {
			switch {
			case a.val < b.val:
				return -1
			case a.val > b.val:
				return 1
			default:
				return 0
			}
		})
	})
}

func BenchmarkStableNew(b *testing.B) {
	order := Order[smallObject]{
		Less: func(a, b smallObject) bool {
			return a.val < b.val
		}, RefLess: func(a, b *smallObject) bool {
			return a.val < b.val
		}}
	benchmarkStruct(b, func(list []smallObject) {
		order.SortStable(list)
	})
}

func BenchmarkStableStd(b *testing.B) {
	benchmarkStruct(b, func(list []smallObject) {
		std.SortStableFunc[[]smallObject, smallObject](list, func(a, b smallObject) int {
			switch {
			case a.val < b.val:
				return -1
			case a.val > b.val:
				return 1
			default:
				return 0
			}
		})
	})
}

func benchmarkPointer(b *testing.B, sort func([]*smallObject)) {
	for _, sc := range level {
		b.Run(sc.name, func(b *testing.B) {
			b.StopTimer()
			rand.Seed(0)
			data := make([]smallObject, sc.size)
			list := make([]*smallObject, sc.size)
			for i := 0; i < b.N; i++ {
				for j := 0; j < sc.size; j++ {
					data[j].val = rand.Intn(sc.size)
					list[j] = &data[j]
				}
				b.StartTimer()
				sort(list)
				b.StopTimer()
			}
		})
	}
}

func BenchmarkPointerNew(b *testing.B) {
	order := Order[*smallObject]{
		Less: func(a, b *smallObject) bool {
			return a.val < b.val
		},
	}
	benchmarkPointer(b, order.Sort)
}

func BenchmarkPointerStd(b *testing.B) {
	benchmarkPointer(b, func(list []*smallObject) {
		std.SortFunc(list, func(a, b *smallObject) int {
			switch {
			case a.val < b.val:
				return -1
			case a.val > b.val:
				return 1
			default:
				return 0
			}
		})
	})
}
