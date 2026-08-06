# Fast generic sort for slices in golang

## API for builtin types
```go
func BinarySearch[E cmp.Ordered](list []E, x E) (int, bool)
func IsSorted[E cmp.Ordered](list []E) bool
func PartlySort[E cmp.Ordered](list []E, k int)
func Sort[E cmp.Ordered](list []E)
func SortStable[E cmp.Ordered](list []E)
```

## Fast API for custom types
```go
type Order[E any] struct {
	Less    func(a, b E) bool
	RefLess func(a, b *E) bool
}

func (od *Order[E]) BinarySearch(list []E, x E) (int, bool)
func (od *Order[E]) IsSorted(list []E) bool
func (od *Order[E]) PartlySort(list []E, k int)
func (od *Order[E]) Sort(list []E)
func (od *Order[E]) SortStable(list []E)
func (od *Order[E]) SortWithOption(list []E, stable, inplace bool)
```

## Func API for custom types
```go
func BinarySearchFunc[E any](list []E, x E, less func(a, b E) bool) (int, bool)
func IsSortedFunc[E any](list []E, less func(a, b E) bool) bool
func PartlySortFunc[E any](list []E, k int, less func(a, b E) bool)
func SortFunc[E any](list []E, less func(a, b E) bool)
func SortStableFunc[E any](list []E, less func(a, b E) bool)
```

## [Benchmark](https://gist.github.com/PeterRK/625e8fad081267d00e5f9e9f7a8e2084) Result on Intel Core Ultra 7 155H
This algorithm runs fast in many cases, but pdqsort is too fast for sorted list. Usually, sorted list is handled well enough, won't be the bottleneck. We should pay more attention to general cases.

### Compared to generic sort in Go's generic slices package
```
name                 exp time/op     new time/op      delta
Int/Small-1K         17.31µs ± 0%    15.38µs ± 1%   -11.17%  (p=0.002 n=6+6)
Int/Small-10K        205.0µs ± 1%    188.2µs ± 1%    -8.17%  (p=0.002 n=6+6)
Int/Small-100K       2.465ms ± 0%    2.294ms ± 0%    -6.92%  (p=0.002 n=6+6)
Int/Small-1M         28.99ms ± 1%    27.16ms ± 1%    -6.33%  (p=0.002 n=6+6)
Int/Random-1K        30.84µs ± 0%    25.89µs ± 1%   -16.05%  (p=0.002 n=6+6)
Int/Random-10K       400.6µs ± 0%    340.2µs ± 1%   -15.07%  (p=0.002 n=6+6)
Int/Random-100K      4.971ms ± 1%    4.263ms ± 1%   -14.23%  (p=0.002 n=6+6)
Int/Random-1M        59.57ms ± 1%    51.57ms ± 1%   -13.44%  (p=0.002 n=6+6)
Int/Constant-1K      470.8ns ± 1%    314.9ns ± 0%   -33.12%  (p=0.002 n=6+6)
Int/Constant-10K     4.386µs ± 1%    3.144µs ± 1%   -28.32%  (p=0.002 n=6+6)
Int/Constant-100K    43.78µs ± 1%    28.48µs ± 1%   -34.94%  (p=0.002 n=6+6)
Int/Constant-1M      434.2µs ± 1%    280.0µs ± 1%   -35.50%  (p=0.002 n=6+6)
Int/Ascent-1K        477.7ns ± 2%    316.2ns ± 1%   -33.81%  (p=0.002 n=6+6)
Int/Ascent-10K       4.423µs ± 0%    3.110µs ± 2%   -29.69%  (p=0.002 n=6+6)
Int/Ascent-100K      44.26µs ± 2%    28.12µs ± 0%   -36.46%  (p=0.002 n=6+6)
Int/Ascent-1M        440.5µs ± 1%    280.4µs ± 1%   -36.35%  (p=0.002 n=6+6)
Int/Descent-1K       694.8ns ± 0%    522.1ns ± 1%   -24.84%  (p=0.002 n=6+6)
Int/Descent-10K      6.883µs ± 0%    5.415µs ± 2%   -21.32%  (p=0.002 n=6+6)
Int/Descent-100K     69.83µs ± 3%    52.27µs ± 1%   -25.14%  (p=0.002 n=6+6)
Int/Descent-1M       836.6µs ± 1%    667.4µs ± 1%   -20.23%  (p=0.002 n=6+6)
Int/Mixed-1K         11.73µs ± 2%    10.31µs ± 1%   -12.12%  (p=0.002 n=6+6)
Int/Mixed-10K        135.5µs ± 2%    126.0µs ± 1%    -6.95%  (p=0.002 n=6+6)
Int/Mixed-100K       1.613ms ± 0%    1.481ms ± 1%    -8.17%  (p=0.002 n=6+6)
Int/Mixed-1M         18.70ms ± 1%    17.50ms ± 1%    -6.41%  (p=0.002 n=6+6)
Hybrid/5%            2.645ms ± 1%    2.162ms ± 1%   -18.25%  (p=0.002 n=6+6)
Hybrid/10%           4.680ms ± 1%    3.883ms ± 2%   -17.03%  (p=0.002 n=6+6)
Hybrid/20%           8.728ms ± 1%    7.312ms ± 2%   -16.23%  (p=0.002 n=6+6)
Hybrid/30%           12.73ms ± 1%    10.70ms ± 1%   -15.98%  (p=0.002 n=6+6)
Hybrid/50%           20.78ms ± 0%    17.54ms ± 1%   -15.58%  (p=0.002 n=6+6)
Float/1K             37.19µs ± 0%    27.23µs ± 1%   -26.79%  (p=0.002 n=6+6)
Float/10K            486.7µs ± 0%    360.0µs ± 1%   -26.03%  (p=0.002 n=6+6)
Float/100K           6.022ms ± 0%    4.508ms ± 1%   -25.15%  (p=0.002 n=6+6)
Float/1M             72.13ms ± 1%    53.31ms ± 1%   -26.09%  (p=0.002 n=6+6)
Str/1K               77.99µs ± 0%    70.09µs ± 1%   -10.13%  (p=0.002 n=6+6)
Str/10K             1028.2µs ± 1%    930.1µs ± 2%    -9.55%  (p=0.002 n=6+6)
Str/100K             13.39ms ± 3%    12.29ms ± 3%    -8.20%  (p=0.002 n=6+6)
Str/1M               155.1ms ± 1%    140.7ms ± 1%    -9.28%  (p=0.002 n=6+6)
Struct/1K            76.39µs ± 1%    52.45µs ± 3%   -31.34%  (p=0.002 n=6+6)
Struct/10K          1016.3µs ± 1%    740.1µs ± 5%   -27.18%  (p=0.002 n=6+6)
Struct/100K         12.732ms ± 1%    7.821ms ± 3%   -38.58%  (p=0.002 n=6+6)
Struct/1M           152.25ms ± 1%    97.93ms ± 3%   -35.68%  (p=0.002 n=6+6)
Stable/1K           134.50µs ± 2%    67.14µs ± 2%   -50.08%  (p=0.002 n=6+6)
Stable/10K          2040.6µs ± 1%    998.3µs ± 12%  -51.08%  (p=0.002 n=6+6)
Stable/100K         30.071ms ± 1%    9.777ms ± 2%   -67.49%  (p=0.002 n=6+6)
Stable/1M            399.2ms ± 0%    120.1ms ± 10%  -69.91%  (p=0.002 n=6+6)
Pointer/1K           47.89µs ± 2%    44.43µs ± 1%    -7.22%  (p=0.002 n=6+6)
Pointer/10K          642.2µs ± 1%    609.6µs ± 1%    -5.09%  (p=0.002 n=6+6)
Pointer/100K         8.422ms ± 1%    8.084ms ± 1%    -4.01%  (p=0.002 n=6+6)
Pointer/1M           130.2ms ± 2%    128.7ms ± 0%         ~  (p=0.394 n=6+6)
```

### Compared to non-generic sort in stdlib
```
name                 std time/op     new time/op      delta
Int/Small-1K         17.36µs ± 1%    15.38µs ± 1%   -11.42%  (p=0.002 n=6+6)
Int/Small-10K        205.9µs ± 1%    188.2µs ± 1%    -8.58%  (p=0.002 n=6+6)
Int/Small-100K       2.472ms ± 1%    2.294ms ± 0%    -7.19%  (p=0.002 n=6+6)
Int/Small-1M         29.22ms ± 1%    27.16ms ± 1%    -7.04%  (p=0.002 n=6+6)
Int/Random-1K        31.06µs ± 1%    25.89µs ± 1%   -16.62%  (p=0.002 n=6+6)
Int/Random-10K       404.2µs ± 1%    340.2µs ± 1%   -15.84%  (p=0.002 n=6+6)
Int/Random-100K      4.998ms ± 1%    4.263ms ± 1%   -14.70%  (p=0.002 n=6+6)
Int/Random-1M        59.57ms ± 1%    51.57ms ± 1%   -13.43%  (p=0.002 n=6+6)
Int/Constant-1K      473.3ns ± 1%    314.9ns ± 0%   -33.48%  (p=0.002 n=6+6)
Int/Constant-10K     4.369µs ± 1%    3.144µs ± 1%   -28.04%  (p=0.002 n=6+6)
Int/Constant-100K    43.62µs ± 1%    28.48µs ± 1%   -34.71%  (p=0.002 n=6+6)
Int/Constant-1M      435.0µs ± 1%    280.0µs ± 1%   -35.63%  (p=0.002 n=6+6)
Int/Ascent-1K        479.1ns ± 1%    316.2ns ± 1%   -33.99%  (p=0.002 n=6+6)
Int/Ascent-10K       4.451µs ± 1%    3.110µs ± 2%   -30.13%  (p=0.002 n=6+6)
Int/Ascent-100K      44.15µs ± 1%    28.12µs ± 0%   -36.30%  (p=0.002 n=6+6)
Int/Ascent-1M        440.8µs ± 1%    280.4µs ± 1%   -36.39%  (p=0.002 n=6+6)
Int/Descent-1K       700.1ns ± 1%    522.1ns ± 1%   -25.42%  (p=0.002 n=6+6)
Int/Descent-10K      6.896µs ± 2%    5.415µs ± 2%   -21.47%  (p=0.002 n=6+6)
Int/Descent-100K     69.01µs ± 2%    52.27µs ± 1%   -24.25%  (p=0.002 n=6+6)
Int/Descent-1M       836.1µs ± 1%    667.4µs ± 1%   -20.19%  (p=0.002 n=6+6)
Int/Mixed-1K         11.80µs ± 2%    10.31µs ± 1%   -12.64%  (p=0.002 n=6+6)
Int/Mixed-10K        136.1µs ± 1%    126.0µs ± 1%    -7.38%  (p=0.002 n=6+6)
Int/Mixed-100K       1.611ms ± 1%    1.481ms ± 1%    -8.07%  (p=0.002 n=6+6)
Int/Mixed-1M         18.76ms ± 2%    17.50ms ± 1%    -6.71%  (p=0.002 n=6+6)
Hybrid/5%            2.635ms ± 1%    2.162ms ± 1%   -17.94%  (p=0.002 n=6+6)
Hybrid/10%           4.677ms ± 1%    3.883ms ± 2%   -16.99%  (p=0.002 n=6+6)
Hybrid/20%           8.761ms ± 1%    7.312ms ± 2%   -16.54%  (p=0.002 n=6+6)
Hybrid/30%           12.71ms ± 1%    10.70ms ± 1%   -15.85%  (p=0.002 n=6+6)
Hybrid/50%           20.82ms ± 1%    17.54ms ± 1%   -15.72%  (p=0.002 n=6+6)
Float/1K             37.43µs ± 1%    27.23µs ± 1%   -27.25%  (p=0.002 n=6+6)
Float/10K            488.0µs ± 1%    360.0µs ± 1%   -26.22%  (p=0.002 n=6+6)
Float/100K           6.026ms ± 1%    4.508ms ± 1%   -25.20%  (p=0.002 n=6+6)
Float/1M             72.10ms ± 1%    53.31ms ± 1%   -26.06%  (p=0.002 n=6+6)
Str/1K               78.28µs ± 1%    70.09µs ± 1%   -10.46%  (p=0.002 n=6+6)
Str/10K             1028.2µs ± 1%    930.1µs ± 2%    -9.54%  (p=0.002 n=6+6)
Str/100K             13.44ms ± 4%    12.29ms ± 3%    -8.57%  (p=0.002 n=6+6)
Str/1M               155.1ms ± 1%    140.7ms ± 1%    -9.34%  (p=0.002 n=6+6)
Struct/1K            86.60µs ± 1%    52.45µs ± 3%   -39.43%  (p=0.002 n=6+6)
Struct/10K          1135.5µs ± 3%    740.1µs ± 5%   -34.83%  (p=0.002 n=6+6)
Struct/100K         14.107ms ± 1%    7.821ms ± 3%   -44.56%  (p=0.002 n=6+6)
Struct/1M           170.63ms ± 2%    97.93ms ± 3%   -42.61%  (p=0.002 n=6+6)
Stable/1K           291.00µs ± 1%    67.14µs ± 2%   -76.93%  (p=0.002 n=6+6)
Stable/10K          4846.8µs ± 3%    998.3µs ± 12%  -79.40%  (p=0.002 n=6+6)
Stable/100K         77.825ms ± 2%    9.777ms ± 2%   -87.44%  (p=0.002 n=6+6)
Stable/1M           1140.9ms ± 1%    120.1ms ± 10%  -89.47%  (p=0.002 n=6+6)
Pointer/1K           60.16µs ± 2%    44.43µs ± 1%   -26.14%  (p=0.002 n=6+6)
Pointer/10K          814.0µs ± 1%    609.6µs ± 1%   -25.12%  (p=0.002 n=6+6)
Pointer/100K        10.833ms ± 1%    8.084ms ± 1%   -25.38%  (p=0.002 n=6+6)
Pointer/1M           176.4ms ± 2%    128.7ms ± 0%   -27.05%  (p=0.002 n=6+6)
```
