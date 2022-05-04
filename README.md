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

### Compared to generic pdqsort in golang.org/x/exp/slices
```
name               exp time/op  new time/op  delta
Int/Small-1K       17.5µs ± 2%   18.0µs ± 2%   +2.87%  (p=0.000 n=10+10)
Int/Small-10K       305µs ± 0%    304µs ± 1%   -0.40%  (p=0.010 n=9+10)
Int/Small-100K     4.42ms ± 0%   4.36ms ± 0%   -1.35%  (p=0.000 n=9+10)
Int/Small-1M       58.1ms ± 0%   56.7ms ± 0%   -2.39%  (p=0.000 n=9+9)
Int/Random-1K      46.8µs ± 2%   39.7µs ± 1%  -15.23%  (p=0.000 n=10+10)
Int/Random-10K      601µs ± 0%    520µs ± 0%  -13.37%  (p=0.000 n=10+10)
Int/Random-100K    7.39ms ± 1%   6.52ms ± 0%  -11.69%  (p=0.000 n=10+10)
Int/Random-1M      87.7ms ± 0%   78.4ms ± 1%  -10.56%  (p=0.000 n=10+10)
Int/Constant-1K    1.04µs ± 1%   1.60µs ± 1%  +52.78%  (p=0.000 n=10+10)
Int/Constant-10K   7.80µs ± 3%  12.55µs ± 2%  +60.95%  (p=0.000 n=10+10)
Int/Constant-100K  63.9µs ± 2%  106.1µs ± 5%  +66.13%  (p=0.000 n=10+10)
Int/Constant-1M     621µs ± 1%   1033µs ± 4%  +66.22%  (p=0.000 n=10+10)
Int/Descent-1K     1.64µs ± 2%   3.04µs ± 3%  +85.05%  (p=0.000 n=10+10)
Int/Descent-10K    11.4µs ± 4%   20.8µs ± 5%  +82.11%  (p=0.000 n=10+10)
Int/Descent-100K    102µs ± 2%    186µs ± 1%  +82.69%  (p=0.000 n=10+10)
Int/Descent-1M     1.03ms ± 1%   1.87ms ± 1%  +82.01%  (p=0.000 n=8+10)
Int/Ascent-1K      1.05µs ± 1%   1.67µs ± 1%  +59.54%  (p=0.000 n=10+10)
Int/Ascent-10K     7.05µs ± 2%  11.18µs ± 3%  +58.47%  (p=0.000 n=9+10)
Int/Ascent-100K    63.4µs ± 1%   94.1µs ± 1%  +48.44%  (p=0.000 n=10+10)
Int/Ascent-1M       621µs ± 1%    948µs ± 2%  +52.81%  (p=0.000 n=10+10)
Int/Mixed-1K       19.9µs ± 2%   19.4µs ± 1%   -2.29%  (p=0.000 n=10+10)
Int/Mixed-10K       250µs ± 1%    237µs ± 1%   -5.07%  (p=0.000 n=9+10)
Int/Mixed-100K     3.05ms ± 0%   2.86ms ± 0%   -6.27%  (p=0.000 n=10+10)
Int/Mixed-1M       36.4ms ± 0%   34.5ms ± 0%   -5.20%  (p=0.000 n=10+10)
Float/1K           50.5µs ± 1%   42.4µs ± 1%  -16.03%  (p=0.000 n=10+10)
Float/10K           653µs ± 0%    562µs ± 0%  -14.04%  (p=0.000 n=10+10)
Float/100K         8.06ms ± 0%   7.07ms ± 0%  -12.22%  (p=0.000 n=10+10)
Float/1M           95.8ms ± 0%   85.3ms ± 0%  -10.94%  (p=0.000 n=9+10)
Str/1K              111µs ± 0%     99µs ± 0%  -10.48%  (p=0.000 n=9+10)
Str/10K            1.49ms ± 0%   1.35ms ± 0%   -9.26%  (p=0.000 n=10+10)
Str/100K           19.5ms ± 0%   18.0ms ± 1%   -7.67%  (p=0.000 n=10+10)
Str/1M              257ms ± 1%    249ms ± 3%   -2.96%  (p=0.000 n=10+10)
Struct/1K           110µs ± 0%     75µs ± 1%  -31.87%  (p=0.000 n=10+10)
Struct/10K         1.45ms ± 1%   1.16ms ± 0%  -20.26%  (p=0.000 n=10+9)
Struct/100K        18.1ms ± 0%   14.0ms ± 0%  -22.79%  (p=0.000 n=10+10)
Struct/1M           218ms ± 1%    165ms ± 0%  -24.16%  (p=0.000 n=10+9)
Stable/1K           200µs ± 1%     86µs ± 0%  -56.89%  (p=0.000 n=9+10)
Stable/10K         3.01ms ± 1%   1.39ms ± 0%  -53.99%  (p=0.000 n=10+10)
Stable/100K        42.8ms ± 1%   16.8ms ± 0%  -60.71%  (p=0.000 n=10+10)
Stable/1M           608ms ± 1%    209ms ± 2%  -65.68%  (p=0.000 n=10+10)
Pointer/1K         73.0µs ± 1%   66.2µs ± 1%   -9.30%  (p=0.000 n=10+10)
Pointer/10K        1.02ms ± 0%   0.94ms ± 0%   -7.90%  (p=0.000 n=9+10)
Pointer/100K       14.1ms ± 1%   13.2ms ± 1%   -6.22%  (p=0.000 n=10+10)
Pointer/1M          192ms ± 1%    195ms ± 5%     ~     (p=0.315 n=8+10)
```
### Compared to non-generic pdqsort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       34.5µs ± 1%  18.0µs ± 2%  -47.98%  (p=0.000 n=9+10)
Int/Small-10K       583µs ± 0%   304µs ± 1%  -47.92%  (p=0.000 n=10+10)
Int/Small-100K     8.40ms ± 0%  4.36ms ± 0%  -48.04%  (p=0.000 n=9+10)
Int/Small-1M        110ms ± 0%    57ms ± 0%  -48.68%  (p=0.000 n=10+9)
Int/Random-1K      86.3µs ± 0%  39.7µs ± 1%  -54.02%  (p=0.000 n=10+10)
Int/Random-10K     1.10ms ± 0%  0.52ms ± 0%  -52.76%  (p=0.000 n=9+10)
Int/Random-100K    13.6ms ± 0%   6.5ms ± 0%  -51.93%  (p=0.000 n=9+10)
Int/Random-1M       162ms ± 0%    78ms ± 1%  -51.60%  (p=0.000 n=10+10)
Int/Constant-1K    5.71µs ± 0%  1.60µs ± 1%  -72.04%  (p=0.000 n=10+10)
Int/Constant-10K   33.0µs ± 0%  12.5µs ± 2%  -62.01%  (p=0.000 n=10+10)
Int/Constant-100K   252µs ± 0%   106µs ± 5%  -57.96%  (p=0.000 n=9+10)
Int/Constant-1M    2.47ms ± 1%  1.03ms ± 4%  -58.22%  (p=0.000 n=10+10)
Int/Descent-1K     7.90µs ± 0%  3.04µs ± 3%  -61.47%  (p=0.000 n=9+10)
Int/Descent-10K    42.0µs ± 1%  20.8µs ± 5%  -50.48%  (p=0.000 n=10+10)
Int/Descent-100K    366µs ± 0%   186µs ± 1%  -49.24%  (p=0.000 n=10+10)
Int/Descent-1M     3.62ms ± 1%  1.87ms ± 1%  -48.29%  (p=0.000 n=10+10)
Int/Ascent-1K      5.71µs ± 0%  1.67µs ± 1%  -70.70%  (p=0.000 n=8+10)
Int/Ascent-10K     32.5µs ± 1%  11.2µs ± 3%  -65.63%  (p=0.000 n=10+10)
Int/Ascent-100K     252µs ± 0%    94µs ± 1%  -62.72%  (p=0.000 n=10+10)
Int/Ascent-1M      2.47ms ± 1%  0.95ms ± 2%  -61.66%  (p=0.000 n=10+10)
Int/Mixed-1K       44.9µs ± 1%  19.4µs ± 1%  -56.74%  (p=0.000 n=10+10)
Int/Mixed-10K       532µs ± 0%   237µs ± 1%  -55.33%  (p=0.000 n=10+10)
Int/Mixed-100K     6.41ms ± 0%  2.86ms ± 0%  -55.38%  (p=0.000 n=9+10)
Int/Mixed-1M       76.1ms ± 0%  34.5ms ± 0%  -54.63%  (p=0.000 n=10+10)
Float/1K           99.4µs ± 0%  42.4µs ± 1%  -57.32%  (p=0.000 n=10+10)
Float/10K          1.28ms ± 0%  0.56ms ± 0%  -56.16%  (p=0.000 n=10+10)
Float/100K         15.8ms ± 0%   7.1ms ± 0%  -55.29%  (p=0.000 n=10+10)
Float/1M            188ms ± 0%    85ms ± 0%  -54.72%  (p=0.000 n=8+10)
Str/1K              139µs ± 0%    99µs ± 0%  -28.70%  (p=0.000 n=10+10)
Str/10K            1.83ms ± 0%  1.35ms ± 0%  -26.14%  (p=0.000 n=10+10)
Str/100K           23.9ms ± 0%  18.0ms ± 1%  -24.71%  (p=0.000 n=10+10)
Str/1M              325ms ± 1%   249ms ± 3%  -23.33%  (p=0.000 n=9+10)
Struct/1K           143µs ± 0%    75µs ± 1%  -47.36%  (p=0.000 n=10+10)
Struct/10K         1.79ms ± 0%  1.16ms ± 0%  -35.35%  (p=0.000 n=9+9)
Struct/100K        21.9ms ± 0%  14.0ms ± 0%  -36.14%  (p=0.000 n=10+10)
Struct/1M           262ms ± 0%   165ms ± 0%  -36.72%  (p=0.000 n=9+9)
Stable/1K           500µs ± 1%    86µs ± 0%  -82.74%  (p=0.000 n=10+10)
Stable/10K         8.32ms ± 1%  1.39ms ± 0%  -83.34%  (p=0.000 n=10+10)
Stable/100K         132ms ± 0%    17ms ± 0%  -87.27%  (p=0.000 n=10+10)
Stable/1M           1.91s ± 0%   0.21s ± 2%  -89.06%  (p=0.000 n=10+10)
Pointer/1K         95.4µs ± 0%  66.2µs ± 1%  -30.57%  (p=0.000 n=10+10)
Pointer/10K        1.29ms ± 0%  0.94ms ± 0%  -27.39%  (p=0.000 n=10+10)
Pointer/100K       17.6ms ± 0%  13.2ms ± 1%  -24.87%  (p=0.000 n=9+10)
Pointer/1M          266ms ± 7%   195ms ± 5%  -26.59%  (p=0.000 n=10+10)
```

