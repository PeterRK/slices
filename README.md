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

## [Benchmark](https://gist.github.com/PeterRK/625e8fad081267d00e5f9e9f7a8e2084) Result on Xeon-8372C
This algorithm runs fast in many cases, but pdqsort is too fast for sorted list. Usually, sorted list is handled well enough, won't be the bottleneck. We should pay more attention to general cases. 

### Compared to generic sort in golang.org/x/exp/slices
```
name               exp time/op  new time/op   delta
Int/Small-1K       24.6µs ± 1%   22.4µs ± 0%   -9.11%  (p=0.000 n=10+8)
Int/Small-10K       286µs ± 0%    265µs ± 0%   -7.45%  (p=0.000 n=10+9)
Int/Small-100K     3.41ms ± 0%   3.16ms ± 0%   -7.10%  (p=0.000 n=10+10)
Int/Small-1M       40.3ms ± 0%   37.1ms ± 0%   -7.86%  (p=0.000 n=10+10)
Int/Random-1K      43.9µs ± 1%   36.4µs ± 0%  -17.05%  (p=0.000 n=10+10)
Int/Random-10K      565µs ± 0%    470µs ± 0%  -16.87%  (p=0.000 n=10+10)
Int/Random-100K    6.95ms ± 0%   5.82ms ± 0%  -16.31%  (p=0.000 n=8+10)
Int/Random-1M      82.6ms ± 0%   69.7ms ± 0%  -15.66%  (p=0.000 n=10+8)
Int/Constant-1K     908ns ± 1%   1362ns ± 1%  +49.97%  (p=0.000 n=10+10)
Int/Constant-10K   6.74µs ± 3%  10.29µs ± 5%  +52.76%  (p=0.000 n=10+10)
Int/Constant-100K  59.4µs ± 1%   89.5µs ± 1%  +50.75%  (p=0.000 n=10+10)
Int/Constant-1M     581µs ± 0%    887µs ± 1%  +52.53%  (p=0.000 n=10+10)
Int/Ascent-1K       914ns ± 1%   1381ns ± 2%  +51.04%  (p=0.000 n=10+8)
Int/Ascent-10K     6.76µs ± 4%   9.60µs ± 1%  +41.95%  (p=0.000 n=10+9)
Int/Ascent-100K    59.3µs ± 1%   85.9µs ± 1%  +44.98%  (p=0.000 n=10+10)
Int/Ascent-1M       581µs ± 0%    882µs ± 1%  +51.72%  (p=0.000 n=10+9)
Int/Descent-1K     1.50µs ± 1%   2.23µs ± 2%  +47.91%  (p=0.000 n=10+10)
Int/Descent-10K    11.3µs ± 2%   16.0µs ± 3%  +42.07%  (p=0.000 n=10+10)
Int/Descent-100K    102µs ± 1%    147µs ± 1%  +43.95%  (p=0.000 n=10+10)
Int/Descent-1M     1.02ms ± 0%   1.46ms ± 0%  +44.06%  (p=0.000 n=10+10)
Int/Mixed-1K       17.4µs ± 1%   16.1µs ± 1%   -7.53%  (p=0.000 n=10+10)
Int/Mixed-10K       193µs ± 0%    182µs ± 1%   -5.35%  (p=0.000 n=10+10)
Int/Mixed-100K     2.25ms ± 0%   2.09ms ± 0%   -6.97%  (p=0.000 n=9+9)
Int/Mixed-1M       26.1ms ± 0%   24.2ms ± 0%   -7.44%  (p=0.000 n=10+10)
Hybrid/5%          3.53ms ± 0%   3.49ms ± 0%   -1.11%  (p=0.000 n=10+10)
Hybrid/10%         6.31ms ± 0%   5.78ms ± 0%   -8.47%  (p=0.000 n=10+10)
Hybrid/20%         11.9ms ± 0%   10.3ms ± 0%  -13.28%  (p=0.000 n=10+10)
Hybrid/30%         17.5ms ± 0%   14.9ms ± 0%  -14.88%  (p=0.000 n=10+10)
Hybrid/50%         28.7ms ± 0%   24.1ms ± 0%  -15.96%  (p=0.000 n=9+9)
Float/1K           47.5µs ± 1%   39.9µs ± 1%  -16.09%  (p=0.000 n=10+10)
Float/10K           619µs ± 0%    519µs ± 0%  -16.21%  (p=0.000 n=9+8)
Float/100K         7.64ms ± 0%   6.45ms ± 0%  -15.59%  (p=0.000 n=10+10)
Float/1M           90.8ms ± 0%   76.7ms ± 0%  -15.60%  (p=0.000 n=8+10)
Str/1K              106µs ± 0%     96µs ± 0%   -9.16%  (p=0.000 n=10+10)
Str/10K            1.38ms ± 0%   1.26ms ± 0%   -8.30%  (p=0.000 n=10+10)
Str/100K           17.6ms ± 0%   16.3ms ± 0%   -7.52%  (p=0.000 n=10+10)
Str/1M              223ms ± 1%    210ms ± 0%   -6.17%  (p=0.000 n=10+10)
Struct/1K           103µs ± 1%     71µs ± 0%  -31.19%  (p=0.000 n=10+9)
Struct/10K         1.36ms ± 0%   0.97ms ± 1%  -28.39%  (p=0.000 n=10+10)
Struct/100K        17.0ms ± 0%   13.3ms ± 0%  -21.57%  (p=0.000 n=10+9)
Struct/1M           206ms ± 0%    159ms ± 1%  -22.87%  (p=0.000 n=10+10)
Stable/1K           189µs ± 0%     81µs ± 1%  -56.90%  (p=0.000 n=10+10)
Stable/10K         2.85ms ± 0%   1.11ms ± 0%  -61.04%  (p=0.000 n=9+10)
Stable/100K        40.5ms ± 0%   16.3ms ± 0%  -59.69%  (p=0.000 n=10+10)
Stable/1M           579ms ± 1%    203ms ± 2%  -64.97%  (p=0.000 n=9+10)
Pointer/1K         68.8µs ± 1%   62.1µs ± 0%   -9.70%  (p=0.000 n=10+8)
Pointer/10K         944µs ± 1%    863µs ± 0%   -8.56%  (p=0.000 n=10+9)
Pointer/100K       12.3ms ± 0%   11.4ms ± 0%   -7.11%  (p=0.000 n=10+10)
Pointer/1M          171ms ± 3%    163ms ± 1%   -4.52%  (p=0.000 n=10+9)
```
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       49.0µs ± 1%  22.4µs ± 0%  -54.31%  (p=0.000 n=10+8)
Int/Small-10K       583µs ± 0%   265µs ± 0%  -54.58%  (p=0.000 n=10+9)
Int/Small-100K     6.93ms ± 0%  3.16ms ± 0%  -54.35%  (p=0.000 n=9+10)
Int/Small-1M       81.1ms ± 0%  37.1ms ± 0%  -54.24%  (p=0.000 n=10+10)
Int/Random-1K      79.8µs ± 0%  36.4µs ± 0%  -54.42%  (p=0.000 n=10+10)
Int/Random-10K     1.04ms ± 0%  0.47ms ± 0%  -54.97%  (p=0.000 n=8+10)
Int/Random-100K    12.9ms ± 0%   5.8ms ± 0%  -55.03%  (p=0.000 n=10+10)
Int/Random-1M       155ms ± 0%    70ms ± 0%  -54.91%  (p=0.000 n=9+8)
Int/Constant-1K    6.15µs ± 6%  1.36µs ± 1%  -77.86%  (p=0.000 n=10+10)
Int/Constant-10K   46.8µs ± 1%  10.3µs ± 5%  -78.01%  (p=0.000 n=10+10)
Int/Constant-100K   451µs ± 0%    89µs ± 1%  -80.18%  (p=0.000 n=10+10)
Int/Constant-1M    4.50ms ± 0%  0.89ms ± 1%  -80.30%  (p=0.000 n=10+10)
Int/Ascent-1K      25.3µs ± 1%   1.4µs ± 2%  -94.55%  (p=0.000 n=10+8)
Int/Ascent-10K      302µs ± 0%    10µs ± 1%  -96.82%  (p=0.000 n=10+9)
Int/Ascent-100K    3.73ms ± 0%  0.09ms ± 1%  -97.70%  (p=0.000 n=10+10)
Int/Ascent-1M      45.0ms ± 0%   0.9ms ± 1%  -98.04%  (p=0.000 n=9+9)
Int/Descent-1K     27.0µs ± 2%   2.2µs ± 2%  -91.75%  (p=0.000 n=10+10)
Int/Descent-10K     317µs ± 0%    16µs ± 3%  -94.96%  (p=0.000 n=10+10)
Int/Descent-100K   3.87ms ± 0%  0.15ms ± 1%  -96.20%  (p=0.000 n=10+10)
Int/Descent-1M     46.2ms ± 0%   1.5ms ± 0%  -96.83%  (p=0.000 n=10+10)
Int/Mixed-1K       43.9µs ± 1%  16.1µs ± 1%  -63.33%  (p=0.000 n=10+10)
Int/Mixed-10K       515µs ± 0%   182µs ± 1%  -64.60%  (p=0.000 n=10+10)
Int/Mixed-100K     5.98ms ± 0%  2.09ms ± 0%  -65.01%  (p=0.000 n=9+9)
Int/Mixed-1M       68.0ms ± 0%  24.2ms ± 0%  -64.48%  (p=0.000 n=10+10)
Hybrid/5%          26.4ms ± 0%   3.5ms ± 0%  -86.79%  (p=0.000 n=9+10)
Hybrid/10%         30.4ms ± 0%   5.8ms ± 0%  -80.99%  (p=0.000 n=10+10)
Hybrid/20%         38.8ms ± 0%  10.3ms ± 0%  -73.33%  (p=0.000 n=9+10)
Hybrid/30%         47.1ms ± 0%  14.9ms ± 0%  -68.34%  (p=0.000 n=10+10)
Hybrid/50%         63.4ms ± 0%  24.1ms ± 0%  -62.04%  (p=0.000 n=10+9)
Float/1K           91.1µs ± 0%  39.9µs ± 1%  -56.22%  (p=0.000 n=10+10)
Float/10K          1.19ms ± 0%  0.52ms ± 0%  -56.58%  (p=0.000 n=10+8)
Float/100K         14.8ms ± 0%   6.4ms ± 0%  -56.47%  (p=0.000 n=10+10)
Float/1M            177ms ± 0%    77ms ± 0%  -56.71%  (p=0.000 n=10+10)
Str/1K              129µs ± 0%    96µs ± 0%  -25.62%  (p=0.000 n=10+10)
Str/10K            1.69ms ± 0%  1.26ms ± 0%  -25.30%  (p=0.000 n=9+10)
Str/100K           21.7ms ± 0%  16.3ms ± 0%  -24.68%  (p=0.000 n=10+10)
Str/1M              274ms ± 1%   210ms ± 0%  -23.56%  (p=0.000 n=10+10)
Struct/1K           128µs ± 1%    71µs ± 0%  -44.56%  (p=0.000 n=9+9)
Struct/10K         1.66ms ± 0%  0.97ms ± 1%  -41.32%  (p=0.000 n=10+10)
Struct/100K        20.6ms ± 0%  13.3ms ± 0%  -35.55%  (p=0.000 n=9+9)
Struct/1M           248ms ± 0%   159ms ± 1%  -35.97%  (p=0.000 n=9+10)
Stable/1K           464µs ± 1%    81µs ± 1%  -82.47%  (p=0.000 n=10+10)
Stable/10K         7.78ms ± 0%  1.11ms ± 0%  -85.74%  (p=0.000 n=9+10)
Stable/100K         124ms ± 0%    16ms ± 0%  -86.82%  (p=0.000 n=9+10)
Stable/1M           1.79s ± 0%   0.20s ± 2%  -88.69%  (p=0.000 n=10+10)
Pointer/1K         89.0µs ± 1%  62.1µs ± 0%  -30.26%  (p=0.000 n=10+8)
Pointer/10K        1.21ms ± 0%  0.86ms ± 0%  -28.84%  (p=0.000 n=9+9)
Pointer/100K       15.6ms ± 1%  11.4ms ± 0%  -26.77%  (p=0.000 n=10+10)
Pointer/1M          215ms ± 1%   163ms ± 1%  -24.26%  (p=0.000 n=10+9)
```