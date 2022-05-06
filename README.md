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

## [Benchmark](https://gist.github.com/PeterRK/7faec0e10e82effb57c2b07e890a6f45) Result on Xeon-8372C
This algorithm runs fast in many cases, but pdqsort is too fast for sorted list. Usually, sorted list is handled well enough, won't be the bottleneck. We should pay more attention to general cases. 

### Compared to generic sort in golang.org/x/exp/slices
```
name               exp time/op  new time/op   delta
Int/Small-1K       16.3µs ± 1%   15.9µs ± 1%   -2.46%  (p=0.000 n=10+10)
Int/Small-10K       291µs ± 0%    263µs ± 0%   -9.49%  (p=0.000 n=10+10)
Int/Small-100K     4.23ms ± 0%   3.75ms ± 0%  -11.21%  (p=0.000 n=9+10)
Int/Small-1M       55.8ms ± 0%   48.7ms ± 0%  -12.65%  (p=0.000 n=10+9)
Int/Random-1K      44.7µs ± 1%   36.7µs ± 1%  -18.03%  (p=0.000 n=10+10)
Int/Random-10K      574µs ± 0%    472µs ± 0%  -17.79%  (p=0.000 n=10+9)
Int/Random-100K    7.05ms ± 0%   5.84ms ± 0%  -17.10%  (p=0.000 n=10+10)
Int/Random-1M      83.9ms ± 0%   69.6ms ± 0%  -17.09%  (p=0.000 n=10+10)
Int/Constant-1K     908ns ± 1%   1362ns ± 1%  +49.97%  (p=0.000 n=10+10)
Int/Constant-10K   6.74µs ± 3%  10.29µs ± 5%  +52.76%  (p=0.000 n=10+10)
Int/Constant-100K  59.4µs ± 1%   89.5µs ± 1%  +50.75%  (p=0.000 n=10+10)
Int/Constant-1M     581µs ± 0%    887µs ± 1%  +52.53%  (p=0.000 n=10+10)
Int/Descent-1K     1.50µs ± 1%   2.23µs ± 2%  +47.91%  (p=0.000 n=10+10)
Int/Descent-10K    11.3µs ± 2%   16.0µs ± 3%  +42.07%  (p=0.000 n=10+10)
Int/Descent-100K    102µs ± 1%    147µs ± 1%  +43.95%  (p=0.000 n=10+10)
Int/Descent-1M     1.02ms ± 0%   1.46ms ± 0%  +44.06%  (p=0.000 n=10+10)
Int/Ascent-1K       914ns ± 1%   1381ns ± 2%  +51.04%  (p=0.000 n=10+8)
Int/Ascent-10K     6.76µs ± 4%   9.60µs ± 1%  +41.95%  (p=0.000 n=10+9)
Int/Ascent-100K    59.3µs ± 1%   85.9µs ± 1%  +44.98%  (p=0.000 n=10+10)
Int/Ascent-1M       581µs ± 0%    881µs ± 0%  +51.61%  (p=0.000 n=10+9)
Int/Mixed-1K       18.8µs ± 1%   17.1µs ± 2%   -8.82%  (p=0.000 n=10+9)
Int/Mixed-10K       239µs ± 0%    217µs ± 0%   -9.05%  (p=0.000 n=10+9)
Int/Mixed-100K     2.91ms ± 0%   2.61ms ± 0%  -10.52%  (p=0.000 n=10+9)
Int/Mixed-1M       35.0ms ± 0%   31.0ms ± 0%  -11.28%  (p=0.000 n=10+9)
Float/1K           47.3µs ± 1%   39.9µs ± 1%  -15.67%  (p=0.000 n=10+10)
Float/10K           617µs ± 0%    517µs ± 0%  -16.10%  (p=0.000 n=9+10)
Float/100K         7.61ms ± 0%   6.42ms ± 0%  -15.70%  (p=0.000 n=10+10)
Float/1M           90.5ms ± 0%   76.4ms ± 0%  -15.64%  (p=0.000 n=10+9)
Str/1K              104µs ± 0%     93µs ± 0%  -10.39%  (p=0.000 n=10+10)
Str/10K            1.41ms ± 0%   1.27ms ± 0%   -9.46%  (p=0.000 n=10+10)
Str/100K           18.4ms ± 0%   16.9ms ± 0%   -7.80%  (p=0.000 n=7+10)
Str/1M              258ms ± 1%    251ms ± 1%   -2.74%  (p=0.000 n=10+9)
Struct/1K           104µs ± 0%     71µs ± 1%  -31.82%  (p=0.000 n=10+10)
Struct/10K         1.37ms ± 0%   1.11ms ± 1%  -19.29%  (p=0.000 n=10+9)
Struct/100K        17.1ms ± 0%   13.4ms ± 0%  -21.45%  (p=0.000 n=9+10)
Struct/1M           208ms ± 0%    159ms ± 0%  -23.36%  (p=0.000 n=10+10)
Stable/1K           189µs ± 0%     82µs ± 1%  -56.57%  (p=0.000 n=10+10)
Stable/10K         2.84ms ± 0%   1.36ms ± 1%  -52.01%  (p=0.000 n=10+10)
Stable/100K        40.3ms ± 0%   16.5ms ± 1%  -59.06%  (p=0.000 n=10+10)
Stable/1M           578ms ± 1%    207ms ± 2%  -64.24%  (p=0.000 n=10+10)
Pointer/1K         68.6µs ± 1%   62.5µs ± 1%   -8.84%  (p=0.000 n=10+10)
Pointer/10K         961µs ± 0%    884µs ± 0%   -8.03%  (p=0.000 n=10+10)
Pointer/100K       13.3ms ± 1%   12.4ms ± 1%   -6.33%  (p=0.000 n=10+10)
Pointer/1M          205ms ± 3%    203ms ± 2%     ~     (p=0.315 n=10+10)
```
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       33.1µs ± 0%  15.9µs ± 1%  -51.87%  (p=0.000 n=9+10)
Int/Small-10K       579µs ± 0%   263µs ± 0%  -54.53%  (p=0.000 n=10+10)
Int/Small-100K     8.30ms ± 0%  3.75ms ± 0%  -54.79%  (p=0.000 n=10+10)
Int/Small-1M        108ms ± 0%    49ms ± 0%  -54.76%  (p=0.000 n=9+9)
Int/Random-1K      80.8µs ± 0%  36.7µs ± 1%  -54.61%  (p=0.000 n=8+10)
Int/Random-10K     1.06ms ± 0%  0.47ms ± 0%  -55.28%  (p=0.000 n=10+9)
Int/Random-100K    13.1ms ± 0%   5.8ms ± 0%  -55.22%  (p=0.000 n=10+10)
Int/Random-1M       155ms ± 0%    70ms ± 0%  -55.24%  (p=0.000 n=7+10)
Int/Constant-1K    6.15µs ± 6%  1.36µs ± 1%  -77.86%  (p=0.000 n=10+10)
Int/Constant-10K   46.8µs ± 1%  10.3µs ± 5%  -78.01%  (p=0.000 n=10+10)
Int/Constant-100K   451µs ± 0%    89µs ± 1%  -80.18%  (p=0.000 n=10+10)
Int/Constant-1M    4.50ms ± 0%  0.89ms ± 1%  -80.30%  (p=0.000 n=10+10)
Int/Descent-1K     27.0µs ± 2%   2.2µs ± 2%  -91.75%  (p=0.000 n=10+10)
Int/Descent-10K     317µs ± 0%    16µs ± 3%  -94.96%  (p=0.000 n=10+10)
Int/Descent-100K   3.87ms ± 0%  0.15ms ± 1%  -96.20%  (p=0.000 n=10+10)
Int/Descent-1M     46.2ms ± 0%   1.5ms ± 0%  -96.83%  (p=0.000 n=10+10)
Int/Ascent-1K      25.3µs ± 1%   1.4µs ± 2%  -94.55%  (p=0.000 n=10+8)
Int/Ascent-10K      302µs ± 0%    10µs ± 1%  -96.82%  (p=0.000 n=10+9)
Int/Ascent-100K    3.73ms ± 0%  0.09ms ± 1%  -97.70%  (p=0.000 n=10+10)
Int/Ascent-1M      45.0ms ± 0%   0.9ms ± 0%  -98.04%  (p=0.000 n=9+9)
Int/Mixed-1K       42.9µs ± 2%  17.1µs ± 2%  -60.12%  (p=0.000 n=10+9)
Int/Mixed-10K       565µs ± 0%   217µs ± 0%  -61.62%  (p=0.000 n=10+9)
Int/Mixed-100K     6.96ms ± 0%  2.61ms ± 0%  -62.54%  (p=0.000 n=10+9)
Int/Mixed-1M       83.4ms ± 0%  31.0ms ± 0%  -62.81%  (p=0.000 n=10+9)
Float/1K           92.9µs ± 1%  39.9µs ± 1%  -57.01%  (p=0.000 n=10+10)
Float/10K          1.22ms ± 0%  0.52ms ± 0%  -57.52%  (p=0.000 n=10+10)
Float/100K         15.1ms ± 0%   6.4ms ± 0%  -57.56%  (p=0.000 n=10+10)
Float/1M            181ms ± 0%    76ms ± 0%  -57.69%  (p=0.000 n=10+9)
Str/1K              129µs ± 0%    93µs ± 0%  -27.59%  (p=0.000 n=10+10)
Str/10K            1.73ms ± 0%  1.27ms ± 0%  -26.58%  (p=0.000 n=9+10)
Str/100K           22.7ms ± 0%  16.9ms ± 0%  -25.27%  (p=0.000 n=9+10)
Str/1M              322ms ± 0%   251ms ± 1%  -21.86%  (p=0.000 n=7+9)
Struct/1K           133µs ± 1%    71µs ± 1%  -46.62%  (p=0.000 n=10+10)
Struct/10K         1.71ms ± 0%  1.11ms ± 1%  -35.11%  (p=0.000 n=10+9)
Struct/100K        21.1ms ± 0%  13.4ms ± 0%  -36.22%  (p=0.000 n=9+10)
Struct/1M           253ms ± 0%   159ms ± 0%  -37.02%  (p=0.000 n=10+10)
Stable/1K           469µs ± 1%    82µs ± 1%  -82.52%  (p=0.000 n=10+10)
Stable/10K         7.83ms ± 0%  1.36ms ± 1%  -82.61%  (p=0.000 n=10+10)
Stable/100K         125ms ± 0%    16ms ± 1%  -86.81%  (p=0.000 n=10+10)
Stable/1M           1.81s ± 0%   0.21s ± 2%  -88.54%  (p=0.000 n=10+10)
Pointer/1K         89.4µs ± 1%  62.5µs ± 1%  -30.09%  (p=0.000 n=10+10)
Pointer/10K        1.24ms ± 0%  0.88ms ± 0%  -28.68%  (p=0.000 n=9+10)
Pointer/100K       17.0ms ± 0%  12.4ms ± 1%  -26.77%  (p=0.000 n=10+10)
Pointer/1M          255ms ± 2%   203ms ± 2%  -20.33%  (p=0.000 n=10+10)
```