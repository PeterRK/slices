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
Int/Small-1K       26.8µs ± 1%  24.0µs ± 1%  -10.23%  (p=0.000 n=9+10)
Int/Small-10K       304µs ± 0%   280µs ± 0%   -7.98%  (p=0.000 n=10+9)
Int/Small-100K     3.62ms ± 0%  3.34ms ± 0%   -7.85%  (p=0.000 n=10+10)
Int/Small-1M       42.8ms ± 0%  39.2ms ± 1%   -8.28%  (p=0.000 n=10+10)
Int/Random-1K      46.6µs ± 1%  39.9µs ± 3%  -14.37%  (p=0.000 n=10+10)
Int/Random-10K      597µs ± 0%   503µs ± 0%  -15.82%  (p=0.000 n=9+9)
Int/Random-100K    7.33ms ± 0%  6.24ms ± 0%  -14.89%  (p=0.000 n=10+10)
Int/Random-1M      87.2ms ± 0%  74.6ms ± 0%  -14.45%  (p=0.000 n=10+10)
Int/Constant-1K    1.05µs ± 2%  0.86µs ± 1%  -18.27%  (p=0.000 n=10+10)
Int/Constant-10K   7.89µs ± 6%  6.64µs ± 8%  -15.78%  (p=0.000 n=10+10)
Int/Constant-100K  64.1µs ± 3%  55.0µs ± 4%  -14.08%  (p=0.000 n=9+10)
Int/Constant-1M     623µs ± 1%   519µs ± 1%  -16.58%  (p=0.000 n=10+10)
Int/Ascent-1K      1.06µs ± 1%  0.86µs ± 1%  -18.69%  (p=0.000 n=10+10)
Int/Ascent-10K     7.32µs ± 4%  6.03µs ± 4%  -17.69%  (p=0.000 n=10+10)
Int/Ascent-100K    64.5µs ± 3%  53.4µs ± 5%  -17.12%  (p=0.000 n=10+10)
Int/Ascent-1M       623µs ± 0%   520µs ± 0%  -16.46%  (p=0.000 n=10+10)
Int/Descent-1K     1.57µs ± 2%  1.41µs ± 1%  -10.30%  (p=0.000 n=10+10)
Int/Descent-10K    11.5µs ± 3%  10.2µs ± 4%  -11.11%  (p=0.000 n=9+10)
Int/Descent-100K    101µs ± 2%    91µs ± 2%   -9.42%  (p=0.000 n=10+10)
Int/Descent-1M     1.01ms ± 0%  0.91ms ± 1%   -9.82%  (p=0.000 n=9+10)
Int/Mixed-1K       18.7µs ± 3%  17.0µs ± 2%   -8.86%  (p=0.000 n=10+10)
Int/Mixed-10K       203µs ± 1%   191µs ± 1%   -6.03%  (p=0.000 n=10+10)
Int/Mixed-100K     2.37ms ± 0%  2.17ms ± 0%   -8.30%  (p=0.000 n=10+10)
Int/Mixed-1M       27.7ms ± 0%  25.2ms ± 0%   -8.96%  (p=0.000 n=10+10)
Hybrid/5%          3.71ms ± 0%  3.15ms ± 0%  -15.10%  (p=0.000 n=10+10)
Hybrid/10%         6.66ms ± 0%  5.63ms ± 0%  -15.44%  (p=0.000 n=10+9)
Hybrid/20%         12.6ms ± 0%  10.6ms ± 0%  -15.43%  (p=0.000 n=9+10)
Hybrid/30%         18.5ms ± 0%  15.7ms ± 0%  -15.29%  (p=0.000 n=9+10)
Hybrid/50%         30.3ms ± 0%  25.6ms ± 0%  -15.43%  (p=0.000 n=9+10)
Float/1K           50.5µs ± 1%  42.4µs ± 1%  -16.02%  (p=0.000 n=10+9)
Float/10K           655µs ± 0%   547µs ± 0%  -16.42%  (p=0.000 n=10+10)
Float/100K         8.07ms ± 0%  6.80ms ± 0%  -15.73%  (p=0.000 n=10+10)
Float/1M           96.1ms ± 0%  81.0ms ± 0%  -15.65%  (p=0.000 n=9+10)
Str/1K              111µs ± 0%   101µs ± 1%   -9.17%  (p=0.000 n=8+10)
Str/10K            1.44ms ± 0%  1.33ms ± 0%   -7.99%  (p=0.000 n=10+10)
Str/100K           18.5ms ± 0%  17.2ms ± 0%   -7.16%  (p=0.000 n=10+10)
Str/1M              233ms ± 1%   223ms ± 1%   -4.52%  (p=0.000 n=10+10)
Struct/1K           108µs ± 1%    75µs ± 1%  -30.86%  (p=0.000 n=10+10)
Struct/10K         1.41ms ± 0%  1.02ms ± 0%  -27.80%  (p=0.000 n=10+10)
Struct/100K        17.6ms ± 0%  14.8ms ± 1%  -16.14%  (p=0.000 n=9+10)
Struct/1M           213ms ± 0%   167ms ± 1%  -21.86%  (p=0.000 n=10+10)
Stable/1K           201µs ± 1%    87µs ± 0%  -56.82%  (p=0.000 n=10+10)
Stable/10K         3.02ms ± 0%  1.18ms ± 0%  -61.09%  (p=0.000 n=10+10)
Stable/100K        43.1ms ± 0%  17.3ms ± 1%  -59.85%  (p=0.000 n=10+10)
Stable/1M           619ms ± 0%   232ms ± 3%  -62.52%  (p=0.000 n=8+10)
Pointer/1K         72.6µs ± 1%  66.2µs ± 1%   -8.83%  (p=0.000 n=9+10)
Pointer/10K         990µs ± 0%   913µs ± 0%   -7.82%  (p=0.000 n=10+10)
Pointer/100K       12.8ms ± 1%  12.2ms ± 0%   -4.24%  (p=0.000 n=10+8)
Pointer/1M          180ms ± 1%   175ms ± 1%   -2.93%  (p=0.000 n=10+9)
```
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       49.7µs ± 1%  24.0µs ± 1%  -51.65%  (p=0.000 n=10+10)
Int/Small-10K       584µs ± 0%   280µs ± 0%  -52.08%  (p=0.000 n=10+9)
Int/Small-100K     6.95ms ± 0%  3.34ms ± 0%  -52.01%  (p=0.000 n=10+10)
Int/Small-1M       81.7ms ± 0%  39.2ms ± 1%  -52.02%  (p=0.000 n=10+10)
Int/Random-1K      86.2µs ± 2%  39.9µs ± 3%  -53.69%  (p=0.000 n=10+10)
Int/Random-10K     1.11ms ± 0%  0.50ms ± 0%  -54.76%  (p=0.000 n=9+9)
Int/Random-100K    13.7ms ± 0%   6.2ms ± 0%  -54.42%  (p=0.000 n=9+10)
Int/Random-1M       163ms ± 0%    75ms ± 0%  -54.28%  (p=0.000 n=10+10)
Int/Constant-1K    4.23µs ± 1%  0.86µs ± 1%  -79.64%  (p=0.000 n=10+10)
Int/Constant-10K   28.2µs ± 3%   6.6µs ± 8%  -76.48%  (p=0.000 n=10+10)
Int/Constant-100K   248µs ± 1%    55µs ± 4%  -77.82%  (p=0.000 n=10+10)
Int/Constant-1M    2.46ms ± 0%  0.52ms ± 1%  -78.85%  (p=0.000 n=10+10)
Int/Ascent-1K      4.23µs ± 3%  0.86µs ± 1%  -79.62%  (p=0.000 n=10+10)
Int/Ascent-10K     27.3µs ± 4%   6.0µs ± 4%  -77.90%  (p=0.000 n=10+10)
Int/Ascent-100K     249µs ± 1%    53µs ± 5%  -78.50%  (p=0.000 n=10+10)
Int/Ascent-1M      2.46ms ± 0%  0.52ms ± 0%  -78.84%  (p=0.000 n=10+10)
Int/Descent-1K     5.99µs ± 3%  1.41µs ± 1%  -76.42%  (p=0.000 n=9+10)
Int/Descent-10K    38.9µs ± 4%  10.2µs ± 4%  -73.72%  (p=0.000 n=10+10)
Int/Descent-100K    362µs ± 0%    91µs ± 2%  -74.74%  (p=0.000 n=7+10)
Int/Descent-1M     3.60ms ± 0%  0.91ms ± 1%  -74.70%  (p=0.000 n=9+10)
Int/Mixed-1K       40.0µs ± 3%  17.0µs ± 2%  -57.53%  (p=0.000 n=10+10)
Int/Mixed-10K       444µs ± 0%   191µs ± 1%  -56.92%  (p=0.000 n=10+10)
Int/Mixed-100K     5.14ms ± 0%  2.17ms ± 0%  -57.74%  (p=0.000 n=10+10)
Int/Mixed-1M       57.7ms ± 0%  25.2ms ± 0%  -56.34%  (p=0.000 n=10+10)
Hybrid/5%          8.31ms ± 0%  3.15ms ± 0%  -62.13%  (p=0.000 n=9+10)
Hybrid/10%         13.8ms ± 0%   5.6ms ± 0%  -59.22%  (p=0.000 n=9+9)
Hybrid/20%         24.7ms ± 0%  10.6ms ± 0%  -56.87%  (p=0.000 n=10+10)
Hybrid/30%         35.5ms ± 0%  15.7ms ± 0%  -55.85%  (p=0.000 n=10+10)
Hybrid/50%         57.2ms ± 0%  25.6ms ± 0%  -55.21%  (p=0.000 n=10+10)
Float/1K            100µs ± 1%    42µs ± 1%  -57.45%  (p=0.000 n=10+9)
Float/10K          1.29ms ± 0%  0.55ms ± 0%  -57.72%  (p=0.000 n=10+10)
Float/100K         16.0ms ± 0%   6.8ms ± 0%  -57.50%  (p=0.000 n=10+10)
Float/1M            191ms ± 0%    81ms ± 0%  -57.46%  (p=0.000 n=10+10)
Str/1K              139µs ± 0%   101µs ± 1%  -27.33%  (p=0.000 n=9+10)
Str/10K            1.80ms ± 0%  1.33ms ± 0%  -26.42%  (p=0.000 n=10+10)
Str/100K           23.0ms ± 0%  17.2ms ± 0%  -25.17%  (p=0.000 n=10+10)
Str/1M              297ms ± 1%   223ms ± 1%  -25.10%  (p=0.000 n=10+10)
Struct/1K           141µs ± 1%    75µs ± 1%  -46.97%  (p=0.000 n=10+10)
Struct/10K         1.78ms ± 0%  1.02ms ± 0%  -42.66%  (p=0.000 n=9+10)
Struct/100K        21.9ms ± 0%  14.8ms ± 1%  -32.76%  (p=0.000 n=10+10)
Struct/1M           263ms ± 0%   167ms ± 1%  -36.71%  (p=0.000 n=10+10)
Stable/1K           493µs ± 1%    87µs ± 0%  -82.41%  (p=0.000 n=10+10)
Stable/10K         8.33ms ± 0%  1.18ms ± 0%  -85.89%  (p=0.000 n=10+10)
Stable/100K         134ms ± 0%    17ms ± 1%  -87.10%  (p=0.000 n=10+10)
Stable/1M           1.93s ± 0%   0.23s ± 3%  -87.97%  (p=0.000 n=10+10)
Pointer/1K         95.2µs ± 1%  66.2µs ± 1%  -30.47%  (p=0.000 n=10+10)
Pointer/10K        1.27ms ± 0%  0.91ms ± 0%  -28.36%  (p=0.000 n=9+10)
Pointer/100K       16.2ms ± 0%  12.2ms ± 0%  -24.79%  (p=0.000 n=10+8)
Pointer/1M          235ms ± 1%   175ms ± 1%  -25.48%  (p=0.000 n=9+9)
```
