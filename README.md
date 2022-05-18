# Fast generic sort for slices in golang

## API for builtin types
```go
func BinarySearch[E constraints.Ordered](list []E, x E) (int, bool)
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

func (od *Order[E]) BinarySearch(list []E, x E) (int, bool)
func (od *Order[E]) IsSorted(list []E) bool
func (od *Order[E]) Sort(list []E)
func (od *Order[E]) SortStable(list []E)
func (od *Order[E]) SortWithOption(list []E, stable, inplace bool)
```

## Func API for custom types
```go
func BinarySearchFunc[E any](list []E, x E, less func(a, b E) bool) (int, bool)
func IsSortedFunc[E any](list []E, less func(a, b E) bool) bool
func SortFunc[E any](list []E, less func(a, b E) bool)
func SortStableFunc[E any](list []E, less func(a, b E) bool)
```

## [Benchmark](https://gist.github.com/PeterRK/625e8fad081267d00e5f9e9f7a8e2084) Result on Xeon-8372C
This algorithm runs fast in many cases, but pdqsort is too fast for sorted list. Usually, sorted list is handled well enough, won't be the bottleneck. We should pay more attention to general cases. 

### Compared to generic sort in golang.org/x/exp/slices
```
name               exp time/op  new time/op  delta
Int/Small-1K       24.6µs ± 1%  22.3µs ± 1%   -9.50%  (p=0.000 n=9+10)
Int/Small-10K       286µs ± 0%   264µs ± 0%   -7.46%  (p=0.000 n=10+10)
Int/Small-100K     3.40ms ± 0%  3.15ms ± 0%   -7.27%  (p=0.000 n=10+10)
Int/Small-1M       40.3ms ± 0%  36.8ms ± 0%   -8.58%  (p=0.000 n=10+10)
Int/Random-1K      43.7µs ± 1%  36.1µs ± 1%  -17.31%  (p=0.000 n=10+10)
Int/Random-10K      563µs ± 0%   467µs ± 0%  -16.94%  (p=0.000 n=10+9)
Int/Random-100K    6.94ms ± 0%  5.80ms ± 0%  -16.46%  (p=0.000 n=10+10)
Int/Random-1M      82.4ms ± 0%  69.2ms ± 0%  -15.94%  (p=0.000 n=10+9)
Int/Constant-1K     856ns ± 1%   941ns ± 1%   +9.88%  (p=0.000 n=10+10)
Int/Constant-10K   6.71µs ± 3%  6.89µs ± 3%   +2.70%  (p=0.002 n=9+10)
Int/Constant-100K  59.7µs ± 1%  59.3µs ± 2%     ~     (p=0.105 n=10+10)
Int/Constant-1M     591µs ± 1%   578µs ± 0%   -2.15%  (p=0.000 n=10+10)
Int/Ascent-1K       859ns ± 0%   933ns ± 1%   +8.59%  (p=0.000 n=10+10)
Int/Ascent-10K     6.57µs ± 3%  6.48µs ± 3%   -1.47%  (p=0.043 n=10+10)
Int/Ascent-100K    59.6µs ± 1%  59.1µs ± 1%   -0.95%  (p=0.005 n=10+10)
Int/Ascent-1M       590µs ± 1%   578µs ± 0%   -1.98%  (p=0.000 n=10+10)
Int/Descent-1K     1.39µs ± 1%  1.49µs ± 1%   +7.60%  (p=0.000 n=10+9)
Int/Descent-10K    10.9µs ± 2%  10.9µs ± 3%     ~     (p=0.481 n=10+10)
Int/Descent-100K   96.4µs ± 0%  96.3µs ± 0%     ~     (p=0.393 n=10+10)
Int/Descent-1M      971µs ± 2%   971µs ± 1%     ~     (p=0.912 n=10+10)
Int/Mixed-1K       17.4µs ± 2%  16.3µs ± 1%   -6.04%  (p=0.000 n=10+10)
Int/Mixed-10K       193µs ± 0%   184µs ± 1%   -4.77%  (p=0.000 n=10+10)
Int/Mixed-100K     2.24ms ± 0%  2.10ms ± 0%   -6.32%  (p=0.000 n=10+10)
Int/Mixed-1M       26.2ms ± 0%  24.4ms ± 0%   -6.52%  (p=0.000 n=10+10)
Hybrid/5%          3.51ms ± 0%  3.06ms ± 1%  -13.00%  (p=0.000 n=10+10)
Hybrid/10%         6.29ms ± 0%  5.40ms ± 1%  -14.22%  (p=0.000 n=10+10)
Hybrid/20%         11.8ms ± 0%  10.1ms ± 0%  -15.16%  (p=0.000 n=9+10)
Hybrid/30%         17.4ms ± 0%  14.6ms ± 0%  -16.18%  (p=0.000 n=9+10)
Hybrid/50%         28.6ms ± 0%  23.8ms ± 0%  -16.57%  (p=0.000 n=10+9)
Float/1K           47.5µs ± 1%  39.7µs ± 1%  -16.44%  (p=0.000 n=9+10)
Float/10K           617µs ± 0%   516µs ± 0%  -16.38%  (p=0.000 n=9+9)
Float/100K         7.62ms ± 0%  6.41ms ± 0%  -15.94%  (p=0.000 n=9+10)
Float/1M           90.7ms ± 0%  76.6ms ± 0%  -15.55%  (p=0.000 n=10+9)
Str/1K              105µs ± 0%    96µs ± 0%   -8.35%  (p=0.000 n=10+10)
Str/10K            1.36ms ± 0%  1.26ms ± 0%   -7.42%  (p=0.000 n=10+10)
Str/100K           17.5ms ± 0%  16.3ms ± 0%   -6.52%  (p=0.000 n=10+10)
Str/1M              221ms ± 1%   211ms ± 1%   -4.48%  (p=0.000 n=10+10)
Struct/1K           102µs ± 0%    70µs ± 0%  -31.79%  (p=0.000 n=10+10)
Struct/10K         1.35ms ± 0%  0.96ms ± 0%  -29.01%  (p=0.000 n=9+10)
Struct/100K        16.9ms ± 0%  13.1ms ± 0%  -22.63%  (p=0.000 n=10+9)
Struct/1M           205ms ± 0%   155ms ± 0%  -24.25%  (p=0.000 n=10+10)
Stable/1K           189µs ± 0%    81µs ± 0%  -57.25%  (p=0.000 n=10+10)
Stable/10K         2.84ms ± 0%  1.10ms ± 0%  -61.20%  (p=0.000 n=10+10)
Stable/100K        40.5ms ± 0%  16.0ms ± 1%  -60.54%  (p=0.000 n=10+10)
Stable/1M           577ms ± 1%   195ms ± 1%  -66.18%  (p=0.000 n=10+10)
Pointer/1K         68.7µs ± 0%  61.8µs ± 1%  -10.06%  (p=0.000 n=10+10)
Pointer/10K         946µs ± 0%   854µs ± 0%   -9.66%  (p=0.000 n=10+8)
Pointer/100K       12.4ms ± 1%  11.6ms ± 1%   -6.37%  (p=0.000 n=10+10)
Pointer/1M          181ms ± 2%   155ms ± 1%  -13.88%  (p=0.000 n=10+8)
```
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       48.8µs ± 1%  22.3µs ± 1%  -54.40%  (p=0.000 n=10+10)
Int/Small-10K       583µs ± 0%   264µs ± 0%  -54.67%  (p=0.000 n=9+10)
Int/Small-100K     6.94ms ± 0%  3.15ms ± 0%  -54.57%  (p=0.000 n=10+10)
Int/Small-1M       81.4ms ± 0%  36.8ms ± 0%  -54.74%  (p=0.000 n=10+10)
Int/Random-1K      79.5µs ± 0%  36.1µs ± 1%  -54.57%  (p=0.000 n=10+10)
Int/Random-10K     1.04ms ± 0%  0.47ms ± 0%  -55.13%  (p=0.000 n=10+9)
Int/Random-100K    12.9ms ± 0%   5.8ms ± 0%  -55.11%  (p=0.000 n=10+10)
Int/Random-1M       154ms ± 0%    69ms ± 0%  -55.13%  (p=0.000 n=10+9)
Int/Constant-1K    6.17µs ± 3%  0.94µs ± 1%  -84.75%  (p=0.000 n=10+10)
Int/Constant-10K   47.7µs ± 1%   6.9µs ± 3%  -85.56%  (p=0.000 n=10+10)
Int/Constant-100K   460µs ± 0%    59µs ± 2%  -87.10%  (p=0.000 n=9+10)
Int/Constant-1M    4.59ms ± 0%  0.58ms ± 0%  -87.41%  (p=0.000 n=10+10)
Int/Ascent-1K      25.3µs ± 2%   0.9µs ± 1%  -96.31%  (p=0.000 n=10+10)
Int/Ascent-10K      301µs ± 0%     6µs ± 3%  -97.85%  (p=0.000 n=10+10)
Int/Ascent-100K    3.72ms ± 0%  0.06ms ± 1%  -98.41%  (p=0.000 n=10+10)
Int/Ascent-1M      44.8ms ± 0%   0.6ms ± 0%  -98.71%  (p=0.000 n=9+10)
Int/Descent-1K     26.6µs ± 1%   1.5µs ± 1%  -94.37%  (p=0.000 n=7+9)
Int/Descent-10K     316µs ± 1%    11µs ± 3%  -96.54%  (p=0.000 n=10+10)
Int/Descent-100K   3.86ms ± 0%  0.10ms ± 0%  -97.51%  (p=0.000 n=10+10)
Int/Descent-1M     45.9ms ± 0%   1.0ms ± 1%  -97.88%  (p=0.000 n=10+10)
Int/Mixed-1K       43.8µs ± 1%  16.3µs ± 1%  -62.78%  (p=0.000 n=9+10)
Int/Mixed-10K       513µs ± 0%   184µs ± 1%  -64.23%  (p=0.000 n=10+10)
Int/Mixed-100K     5.97ms ± 0%  2.10ms ± 0%  -64.77%  (p=0.000 n=10+10)
Int/Mixed-1M       68.0ms ± 0%  24.4ms ± 0%  -64.02%  (p=0.000 n=9+10)
Hybrid/5%          26.6ms ± 0%   3.1ms ± 1%  -88.49%  (p=0.000 n=10+10)
Hybrid/10%         30.5ms ± 0%   5.4ms ± 1%  -82.29%  (p=0.000 n=10+10)
Hybrid/20%         38.8ms ± 0%  10.1ms ± 0%  -74.09%  (p=0.000 n=10+10)
Hybrid/30%         47.2ms ± 0%  14.6ms ± 0%  -69.05%  (p=0.000 n=10+10)
Hybrid/50%         63.4ms ± 0%  23.8ms ± 0%  -62.44%  (p=0.000 n=10+9)
Float/1K           90.9µs ± 0%  39.7µs ± 1%  -56.37%  (p=0.000 n=10+10)
Float/10K          1.19ms ± 0%  0.52ms ± 0%  -56.68%  (p=0.000 n=10+9)
Float/100K         14.8ms ± 0%   6.4ms ± 0%  -56.65%  (p=0.000 n=9+10)
Float/1M            177ms ± 0%    77ms ± 0%  -56.65%  (p=0.000 n=10+9)
Str/1K              129µs ± 0%    96µs ± 0%  -25.35%  (p=0.000 n=10+10)
Str/10K            1.69ms ± 0%  1.26ms ± 0%  -25.07%  (p=0.000 n=9+10)
Str/100K           21.6ms ± 0%  16.3ms ± 0%  -24.46%  (p=0.000 n=10+10)
Str/1M              274ms ± 1%   211ms ± 1%  -22.98%  (p=0.000 n=10+10)
Struct/1K           127µs ± 1%    70µs ± 0%  -45.09%  (p=0.000 n=10+10)
Struct/10K         1.65ms ± 0%  0.96ms ± 0%  -41.90%  (p=0.000 n=9+10)
Struct/100K        20.5ms ± 0%  13.1ms ± 0%  -36.40%  (p=0.000 n=10+9)
Struct/1M           247ms ± 0%   155ms ± 0%  -37.22%  (p=0.000 n=9+10)
Stable/1K           464µs ± 1%    81µs ± 0%  -82.62%  (p=0.000 n=10+10)
Stable/10K         7.75ms ± 0%  1.10ms ± 0%  -85.77%  (p=0.000 n=9+10)
Stable/100K         124ms ± 1%    16ms ± 1%  -87.10%  (p=0.000 n=10+10)
Stable/1M           1.78s ± 0%   0.20s ± 1%  -89.05%  (p=0.000 n=10+10)
Pointer/1K         88.8µs ± 0%  61.8µs ± 1%  -30.46%  (p=0.000 n=10+10)
Pointer/10K        1.21ms ± 0%  0.85ms ± 0%  -29.44%  (p=0.000 n=9+8)
Pointer/100K       15.7ms ± 1%  11.6ms ± 1%  -26.19%  (p=0.000 n=10+10)
Pointer/1M          214ms ± 0%   155ms ± 1%  -27.18%  (p=0.000 n=8+8)
```