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

## [Benchmark](https://gist.github.com/PeterRK/625e8fad081267d00e5f9e9f7a8e2084) Result on Xeon-8374C
### Compared to generic sort in golang.org/x/exp/slices
```
name               exp time/op  new time/op  delta
Int/Small-1K       26.8µs ± 1%  23.8µs ± 2%  -10.90%  (p=0.008 n=5+5)
Int/Small-10K       304µs ± 0%   249µs ± 0%  -18.16%  (p=0.008 n=5+5)
Int/Small-100K     3.61ms ± 0%  2.87ms ± 0%  -20.55%  (p=0.008 n=5+5)
Int/Small-1M       42.8ms ± 0%  35.5ms ± 0%  -16.96%  (p=0.008 n=5+5)
Int/Random-1K      46.3µs ± 0%  38.5µs ± 1%  -16.79%  (p=0.008 n=5+5)
Int/Random-10K      595µs ± 0%   465µs ± 0%  -21.89%  (p=0.008 n=5+5)
Int/Random-100K    7.33ms ± 0%  5.61ms ± 0%  -23.44%  (p=0.008 n=5+5)
Int/Random-1M      87.2ms ± 0%  66.0ms ± 0%  -24.31%  (p=0.008 n=5+5)
Int/Constant-1K    1.06µs ± 1%  0.86µs ± 1%  -18.99%  (p=0.008 n=5+5)
Int/Constant-10K   8.12µs ± 0%  6.86µs ± 5%  -15.44%  (p=0.016 n=4+5)
Int/Constant-100K  64.4µs ± 2%  53.0µs ± 2%  -17.64%  (p=0.008 n=5+5)
Int/Constant-1M     621µs ± 0%   518µs ± 0%  -16.61%  (p=0.008 n=5+5)
Int/Ascent-1K      1.05µs ± 1%  0.86µs ± 1%  -18.11%  (p=0.008 n=5+5)
Int/Ascent-10K     7.53µs ± 4%  6.17µs ± 3%  -18.11%  (p=0.008 n=5+5)
Int/Ascent-100K    64.3µs ± 2%  53.2µs ± 0%  -17.25%  (p=0.016 n=5+4)
Int/Ascent-1M       621µs ± 0%   519µs ± 0%  -16.43%  (p=0.008 n=5+5)
Int/Descent-1K     1.58µs ± 1%  1.41µs ± 2%  -10.68%  (p=0.008 n=5+5)
Int/Descent-10K    11.7µs ± 5%  10.4µs ± 6%  -10.91%  (p=0.008 n=5+5)
Int/Descent-100K    100µs ± 1%    91µs ± 2%   -9.55%  (p=0.008 n=5+5)
Int/Descent-1M     1.01ms ± 0%  0.91ms ± 0%  -10.25%  (p=0.008 n=5+5)
Int/Mixed-1K       18.6µs ± 3%  16.9µs ± 2%   -9.10%  (p=0.008 n=5+5)
Int/Mixed-10K       204µs ± 1%   229µs ± 0%  +12.20%  (p=0.008 n=5+5)
Int/Mixed-100K     2.37ms ± 0%  2.53ms ± 0%   +6.55%  (p=0.008 n=5+5)
Int/Mixed-1M       27.6ms ± 0%  27.9ms ± 0%   +0.94%  (p=0.008 n=5+5)
Hybrid/5%          3.70ms ± 0%  2.94ms ± 0%  -20.52%  (p=0.008 n=5+5)
Hybrid/10%         6.64ms ± 0%  5.23ms ± 0%  -21.17%  (p=0.008 n=5+5)
Hybrid/20%         12.5ms ± 0%   9.8ms ± 0%  -21.60%  (p=0.008 n=5+5)
Hybrid/30%         18.4ms ± 0%  14.5ms ± 0%  -21.61%  (p=0.008 n=5+5)
Hybrid/50%         30.2ms ± 0%  23.6ms ± 0%  -21.80%  (p=0.008 n=5+5)
Float/1K           50.6µs ± 1%  42.2µs ± 1%  -16.47%  (p=0.008 n=5+5)
Float/10K           654µs ± 0%   499µs ± 0%  -23.65%  (p=0.008 n=5+5)
Float/100K         8.08ms ± 0%  5.92ms ± 0%  -26.82%  (p=0.008 n=5+5)
Float/1M           96.1ms ± 0%  68.5ms ± 0%  -28.70%  (p=0.008 n=5+5)
Str/1K              112µs ± 1%   101µs ± 0%   -9.56%  (p=0.008 n=5+5)
Str/10K            1.44ms ± 0%  1.32ms ± 0%   -8.21%  (p=0.008 n=5+5)
Str/100K           18.5ms ± 0%  17.1ms ± 0%   -7.20%  (p=0.008 n=5+5)
Str/1M              232ms ± 0%   221ms ± 1%   -4.73%  (p=0.008 n=5+5)
Struct/1K           108µs ± 1%    74µs ± 0%  -31.01%  (p=0.008 n=5+5)
Struct/10K         1.41ms ± 0%  1.02ms ± 0%  -27.81%  (p=0.008 n=5+5)
Struct/100K        17.6ms ± 0%  14.7ms ± 0%  -16.33%  (p=0.008 n=5+5)
Struct/1M           213ms ± 0%   165ms ± 0%  -22.44%  (p=0.008 n=5+5)
Stable/1K           200µs ± 0%    86µs ± 1%  -56.93%  (p=0.008 n=5+5)
Stable/10K         3.02ms ± 0%  1.17ms ± 0%  -61.27%  (p=0.008 n=5+5)
Stable/100K        43.0ms ± 0%  17.2ms ± 0%  -60.00%  (p=0.008 n=5+5)
Stable/1M           615ms ± 0%   219ms ± 1%  -64.31%  (p=0.008 n=5+5)
Pointer/1K         72.7µs ± 1%  66.3µs ± 1%   -8.74%  (p=0.008 n=5+5)
Pointer/10K         989µs ± 0%   911µs ± 0%   -7.93%  (p=0.008 n=5+5)
Pointer/100K       12.7ms ± 0%  12.4ms ± 0%   -2.67%  (p=0.016 n=5+4)
Pointer/1M          183ms ± 0%   175ms ± 0%   -4.68%  (p=0.016 n=5+4)
```
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       49.0µs ± 0%  23.8µs ± 2%  -51.38%  (p=0.008 n=5+5)
Int/Small-10K       582µs ± 0%   249µs ± 0%  -57.22%  (p=0.008 n=5+5)
Int/Small-100K     6.93ms ± 0%  2.87ms ± 0%  -58.54%  (p=0.008 n=5+5)
Int/Small-1M       81.6ms ± 0%  35.5ms ± 0%  -56.46%  (p=0.008 n=5+5)
Int/Random-1K      86.1µs ± 0%  38.5µs ± 1%  -55.23%  (p=0.008 n=5+5)
Int/Random-10K     1.11ms ± 0%  0.47ms ± 0%  -58.04%  (p=0.008 n=5+5)
Int/Random-100K    13.7ms ± 0%   5.6ms ± 0%  -58.87%  (p=0.008 n=5+5)
Int/Random-1M       163ms ± 0%    66ms ± 0%  -59.45%  (p=0.008 n=5+5)
Int/Constant-1K    4.20µs ± 1%  0.86µs ± 1%  -79.60%  (p=0.008 n=5+5)
Int/Constant-10K   28.4µs ± 3%   6.9µs ± 5%  -75.82%  (p=0.008 n=5+5)
Int/Constant-100K   248µs ± 1%    53µs ± 2%  -78.60%  (p=0.008 n=5+5)
Int/Constant-1M    2.45ms ± 0%  0.52ms ± 0%  -78.88%  (p=0.008 n=5+5)
Int/Ascent-1K      4.20µs ± 2%  0.86µs ± 1%  -79.58%  (p=0.008 n=5+5)
Int/Ascent-10K     27.4µs ± 2%   6.2µs ± 3%  -77.46%  (p=0.008 n=5+5)
Int/Ascent-100K     247µs ± 1%    53µs ± 0%  -78.45%  (p=0.016 n=5+4)
Int/Ascent-1M      2.45ms ± 0%  0.52ms ± 0%  -78.84%  (p=0.008 n=5+5)
Int/Descent-1K     6.03µs ± 5%  1.41µs ± 2%  -76.62%  (p=0.008 n=5+5)
Int/Descent-10K    39.1µs ± 2%  10.4µs ± 6%  -73.28%  (p=0.008 n=5+5)
Int/Descent-100K    361µs ± 1%    91µs ± 2%  -74.90%  (p=0.008 n=5+5)
Int/Descent-1M     3.59ms ± 0%  0.91ms ± 0%  -74.79%  (p=0.008 n=5+5)
Int/Mixed-1K       40.6µs ± 1%  16.9µs ± 2%  -58.31%  (p=0.008 n=5+5)
Int/Mixed-10K       444µs ± 0%   229µs ± 0%  -48.43%  (p=0.008 n=5+5)
Int/Mixed-100K     5.13ms ± 0%  2.53ms ± 0%  -50.78%  (p=0.008 n=5+5)
Int/Mixed-1M       57.5ms ± 0%  27.9ms ± 0%  -51.56%  (p=0.008 n=5+5)
Hybrid/5%          8.29ms ± 0%  2.94ms ± 0%  -64.50%  (p=0.008 n=5+5)
Hybrid/10%         13.7ms ± 0%   5.2ms ± 0%  -61.74%  (p=0.008 n=5+5)
Hybrid/20%         24.5ms ± 0%   9.8ms ± 0%  -59.84%  (p=0.008 n=5+5)
Hybrid/30%         35.3ms ± 0%  14.5ms ± 0%  -59.07%  (p=0.008 n=5+5)
Hybrid/50%         56.9ms ± 0%  23.6ms ± 0%  -58.48%  (p=0.008 n=5+5)
Float/1K            100µs ± 1%    42µs ± 1%  -57.61%  (p=0.008 n=5+5)
Float/10K          1.29ms ± 0%  0.50ms ± 0%  -61.36%  (p=0.008 n=5+5)
Float/100K         15.9ms ± 0%   5.9ms ± 0%  -62.87%  (p=0.008 n=5+5)
Float/1M            190ms ± 0%    69ms ± 0%  -63.87%  (p=0.008 n=5+5)
Str/1K              139µs ± 0%   101µs ± 0%  -27.40%  (p=0.008 n=5+5)
Str/10K            1.80ms ± 0%  1.32ms ± 0%  -26.49%  (p=0.008 n=5+5)
Str/100K           22.9ms ± 0%  17.1ms ± 0%  -25.19%  (p=0.008 n=5+5)
Str/1M              289ms ± 1%   221ms ± 1%  -23.43%  (p=0.008 n=5+5)
Struct/1K           141µs ± 1%    74µs ± 0%  -47.35%  (p=0.008 n=5+5)
Struct/10K         1.78ms ± 0%  1.02ms ± 0%  -42.81%  (p=0.008 n=5+5)
Struct/100K        21.9ms ± 0%  14.7ms ± 0%  -32.83%  (p=0.008 n=5+5)
Struct/1M           262ms ± 0%   165ms ± 0%  -36.87%  (p=0.008 n=5+5)
Stable/1K           497µs ± 0%    86µs ± 1%  -82.67%  (p=0.008 n=5+5)
Stable/10K         8.35ms ± 0%  1.17ms ± 0%  -85.99%  (p=0.008 n=5+5)
Stable/100K         133ms ± 0%    17ms ± 0%  -87.05%  (p=0.008 n=5+5)
Stable/1M           1.92s ± 0%   0.22s ± 1%  -88.57%  (p=0.008 n=5+5)
Pointer/1K         95.4µs ± 1%  66.3µs ± 1%  -30.47%  (p=0.008 n=5+5)
Pointer/10K        1.27ms ± 0%  0.91ms ± 0%  -28.32%  (p=0.008 n=5+5)
Pointer/100K       16.6ms ± 2%  12.4ms ± 0%  -25.35%  (p=0.016 n=5+4)
Pointer/1M          225ms ± 2%   175ms ± 0%  -22.41%  (p=0.016 n=5+4)
```
