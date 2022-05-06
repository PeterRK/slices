# Fast generic sort for slices in golang

## API for builtin types
```go
func BinarySearch[E constraints.Ordered](list []E, x E) int
func IsSorted[E constraints.Ordered](list []E) bool
func Sort[E constraints.Ordered](list []E)
func SortStable[E constraints.Ordered](list []E)
```

## Fast API for custom types
```go
type Order[E any] struct {
	Less    func(a, b E) bool
	RefLess func(a, b *E) bool
}

func (od *Order[E]) BinarySearch(list []E, x E) int
func (od *Order[E]) IsSorted(list []E) bool
func (od *Order[E]) Sort(list []E)
func (od *Order[E]) SortStable(list []E)
func (od *Order[E]) SortWithOption(list []E, stable, inplace bool)
```

## Func API for custom types
```go
func BinarySearchFunc[E any](list []E, x E, less func(a, b E) bool) int
func IsSortedFunc[E any](list []E, less func(a, b E) bool) bool
func SortFunc[E any](list []E, less func(a, b E) bool)
func SortStableFunc[E any](list []E, less func(a, b E) bool)
```

## [Benchmark](https://gist.github.com/PeterRK/7faec0e10e82effb57c2b07e890a6f45) Result on Xeon-8374C
This algorithm runs fast in many cases, but pdqsort is too fast for sorted list. Usually, sorted list is handled well enough, won't be the bottleneck. We should pay more attention to general cases. 

### Compared to generic sort in golang.org/x/exp/slices
```
name               exp time/op  new time/op   delta
Int/Small-1K       18.9µs ± 0%   19.3µs ± 1%   +1.89%  (p=0.000 n=9+10)
Int/Small-10K       311µs ± 0%    292µs ± 0%   -6.06%  (p=0.000 n=10+10)
Int/Small-100K     4.49ms ± 0%   4.09ms ± 0%   -8.88%  (p=0.000 n=10+10)
Int/Small-1M       59.2ms ± 0%   52.9ms ± 0%  -10.59%  (p=0.000 n=10+8)
Int/Random-1K      49.5µs ± 0%   40.0µs ± 0%  -19.23%  (p=0.000 n=10+10)
Int/Random-10K      614µs ± 0%    496µs ± 0%  -19.20%  (p=0.000 n=10+9)
Int/Random-100K    7.51ms ± 1%   6.12ms ± 0%  -18.54%  (p=0.000 n=9+10)
Int/Random-1M      89.3ms ± 0%   73.0ms ± 0%  -18.25%  (p=0.000 n=8+10)
Int/Constant-1K    1.35µs ± 0%   1.78µs ± 0%  +31.76%  (p=0.000 n=10+10)
Int/Constant-10K   11.0µs ± 0%   13.6µs ± 0%  +22.91%  (p=0.000 n=10+10)
Int/Constant-100K  67.4µs ± 0%   96.2µs ± 0%  +42.74%  (p=0.000 n=10+10)
Int/Constant-1M     623µs ± 0%    924µs ± 0%  +48.36%  (p=0.000 n=10+10)
Int/Descent-1K     2.13µs ± 1%   3.06µs ± 0%  +43.84%  (p=0.000 n=10+10)
Int/Descent-10K    16.0µs ± 0%   20.5µs ± 0%  +27.79%  (p=0.000 n=10+9)
Int/Descent-100K    112µs ± 0%    157µs ± 0%  +39.62%  (p=0.000 n=10+8)
Int/Descent-1M     1.09ms ± 0%   1.52ms ± 0%  +40.27%  (p=0.000 n=9+9)
Int/Ascent-1K      1.35µs ± 0%   2.36µs ± 0%  +74.50%  (p=0.000 n=10+10)
Int/Ascent-10K     10.9µs ± 0%   16.7µs ± 0%  +53.30%  (p=0.000 n=10+9)
Int/Ascent-100K    67.2µs ± 0%  118.8µs ± 0%  +76.81%  (p=0.000 n=10+10)
Int/Ascent-1M       622µs ± 0%   1139µs ± 0%  +83.12%  (p=0.000 n=10+10)
Int/Mixed-1K       22.7µs ± 0%   21.4µs ± 0%   -5.49%  (p=0.000 n=10+10)
Int/Mixed-10K       259µs ± 0%    222µs ± 0%  -13.96%  (p=0.000 n=10+10)
Int/Mixed-100K     3.11ms ± 0%   2.55ms ± 0%  -17.79%  (p=0.000 n=10+10)
Int/Mixed-1M       37.3ms ± 0%   29.9ms ± 0%  -19.92%  (p=0.000 n=9+10)
Float/1K           51.9µs ± 0%   43.2µs ± 0%  -16.85%  (p=0.000 n=10+9)
Float/10K           657µs ± 0%    548µs ± 0%  -16.56%  (p=0.000 n=10+10)
Float/100K         8.09ms ± 0%   6.81ms ± 0%  -15.82%  (p=0.000 n=10+10)
Float/1M           96.3ms ± 0%   81.4ms ± 0%  -15.54%  (p=0.000 n=10+9)
Str/1K              113µs ± 0%    103µs ± 0%   -9.52%  (p=0.000 n=10+10)
Str/10K            1.50ms ± 0%   1.37ms ± 0%   -8.50%  (p=0.000 n=10+9)
Str/100K           19.7ms ± 0%   18.5ms ± 0%   -6.28%  (p=0.000 n=8+9)
Str/1M              276ms ± 2%    267ms ± 3%   -3.09%  (p=0.000 n=10+10)
Struct/1K           112µs ± 0%     79µs ± 0%  -29.90%  (p=0.000 n=10+8)
Struct/10K         1.45ms ± 0%   1.20ms ± 0%  -17.75%  (p=0.000 n=10+10)
Struct/100K        18.1ms ± 0%   14.5ms ± 0%  -20.06%  (p=0.000 n=10+10)
Struct/1M           219ms ± 0%    171ms ± 1%  -21.76%  (p=0.000 n=10+10)
Stable/1K           203µs ± 0%     90µs ± 0%  -55.92%  (p=0.000 n=10+9)
Stable/10K         3.02ms ± 0%   1.44ms ± 0%  -52.21%  (p=0.000 n=9+10)
Stable/100K        43.0ms ± 0%   17.6ms ± 1%  -59.19%  (p=0.000 n=10+10)
Stable/1M           614ms ± 0%    217ms ± 1%  -64.69%  (p=0.000 n=10+8)
Pointer/1K         75.3µs ± 0%   69.1µs ± 0%   -8.29%  (p=0.000 n=10+9)
Pointer/10K        1.03ms ± 0%   0.96ms ± 0%   -6.38%  (p=0.000 n=9+10)
Pointer/100K       14.3ms ± 0%   13.6ms ± 0%   -4.50%  (p=0.000 n=10+9)
Pointer/1M          216ms ± 6%    223ms ± 2%     ~     (p=0.052 n=10+10)
```
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       37.1µs ± 0%  19.3µs ± 1%  -47.94%  (p=0.000 n=10+10)
Int/Small-10K       617µs ± 1%   292µs ± 0%  -52.63%  (p=0.000 n=9+10)
Int/Small-100K     8.79ms ± 0%  4.09ms ± 0%  -53.48%  (p=0.000 n=10+10)
Int/Small-1M        114ms ± 0%    53ms ± 0%  -53.71%  (p=0.000 n=9+8)
Int/Random-1K      87.5µs ± 0%  40.0µs ± 0%  -54.34%  (p=0.000 n=10+10)
Int/Random-10K     1.12ms ± 0%  0.50ms ± 0%  -55.68%  (p=0.000 n=10+9)
Int/Random-100K    13.8ms ± 0%   6.1ms ± 0%  -55.75%  (p=0.000 n=10+10)
Int/Random-1M       165ms ± 0%    73ms ± 0%  -55.73%  (p=0.000 n=10+10)
Int/Constant-1K    10.0µs ± 0%   1.8µs ± 0%  -82.19%  (p=0.000 n=10+10)
Int/Constant-10K   53.9µs ± 1%  13.6µs ± 0%  -74.82%  (p=0.000 n=10+10)
Int/Constant-100K   487µs ± 0%    96µs ± 0%  -80.22%  (p=0.000 n=10+10)
Int/Constant-1M    4.81ms ± 0%  0.92ms ± 0%  -80.77%  (p=0.000 n=10+10)
Int/Descent-1K     34.7µs ± 0%   3.1µs ± 0%  -91.18%  (p=0.000 n=10+10)
Int/Descent-10K     342µs ± 0%    20µs ± 0%  -94.01%  (p=0.000 n=10+9)
Int/Descent-100K   4.15ms ± 0%  0.16ms ± 0%  -96.21%  (p=0.000 n=10+8)
Int/Descent-1M     49.4ms ± 1%   1.5ms ± 0%  -96.91%  (p=0.000 n=10+9)
Int/Ascent-1K      33.0µs ± 0%   2.4µs ± 0%  -92.85%  (p=0.000 n=9+10)
Int/Ascent-10K      326µs ± 0%    17µs ± 0%  -94.87%  (p=0.000 n=10+9)
Int/Ascent-100K    3.99ms ± 0%  0.12ms ± 0%  -97.02%  (p=0.000 n=9+10)
Int/Ascent-1M      48.0ms ± 0%   1.1ms ± 0%  -97.63%  (p=0.000 n=10+10)
Int/Mixed-1K       50.5µs ± 1%  21.4µs ± 0%  -57.51%  (p=0.000 n=10+10)
Int/Mixed-10K       602µs ± 0%   222µs ± 0%  -63.05%  (p=0.000 n=10+10)
Int/Mixed-100K     7.37ms ± 0%  2.55ms ± 0%  -65.38%  (p=0.000 n=10+10)
Int/Mixed-1M       88.2ms ± 0%  29.9ms ± 0%  -66.13%  (p=0.000 n=10+10)
Float/1K            100µs ± 0%    43µs ± 0%  -57.04%  (p=0.000 n=10+9)
Float/10K          1.30ms ± 0%  0.55ms ± 0%  -57.67%  (p=0.000 n=10+10)
Float/100K         16.0ms ± 0%   6.8ms ± 0%  -57.54%  (p=0.000 n=10+10)
Float/1M            192ms ± 0%    81ms ± 0%  -57.54%  (p=0.000 n=9+9)
Str/1K              140µs ± 0%   103µs ± 0%  -26.41%  (p=0.000 n=10+10)
Str/10K            1.84ms ± 0%  1.37ms ± 0%  -25.45%  (p=0.000 n=10+9)
Str/100K           24.3ms ± 1%  18.5ms ± 0%  -23.93%  (p=0.000 n=10+9)
Str/1M              345ms ± 5%   267ms ± 3%  -22.48%  (p=0.000 n=9+10)
Struct/1K           144µs ± 1%    79µs ± 0%  -45.42%  (p=0.000 n=10+8)
Struct/10K         1.82ms ± 1%  1.20ms ± 0%  -34.19%  (p=0.000 n=10+10)
Struct/100K        22.4ms ± 0%  14.5ms ± 0%  -35.40%  (p=0.000 n=10+10)
Struct/1M           269ms ± 1%   171ms ± 1%  -36.30%  (p=0.000 n=10+10)
Stable/1K           501µs ± 1%    90µs ± 0%  -82.13%  (p=0.000 n=10+9)
Stable/10K         8.34ms ± 0%  1.44ms ± 0%  -82.71%  (p=0.000 n=10+10)
Stable/100K         133ms ± 0%    18ms ± 1%  -86.82%  (p=0.000 n=10+10)
Stable/1M           1.93s ± 1%   0.22s ± 1%  -88.75%  (p=0.000 n=10+8)
Pointer/1K         96.8µs ± 0%  69.1µs ± 0%  -28.65%  (p=0.000 n=10+9)
Pointer/10K        1.32ms ± 0%  0.96ms ± 0%  -27.00%  (p=0.000 n=10+10)
Pointer/100K       18.2ms ± 1%  13.6ms ± 0%  -25.13%  (p=0.000 n=10+9)
Pointer/1M          267ms ± 1%   223ms ± 2%  -16.79%  (p=0.000 n=8+10)
```

