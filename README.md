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
### Compared to non-generic sort in stdlib
```
name               std time/op  new time/op  delta
Int/Small-1K       35.9µs ± 3%  16.0µs ± 5%  -55.29%  (p=0.000 n=10+10)
Int/Small-10K       622µs ± 0%   236µs ± 1%  -62.01%  (p=0.000 n=10+10)
Int/Small-100K     8.90ms ± 1%  3.23ms ± 0%  -63.72%  (p=0.000 n=10+10)
Int/Small-1M        115ms ± 0%    41ms ± 0%  -64.40%  (p=0.000 n=10+9)
Int/Random-1K      86.0µs ± 1%  37.1µs ± 2%  -56.82%  (p=0.000 n=10+9)
Int/Random-10K     1.12ms ± 0%  0.45ms ± 0%  -59.68%  (p=0.000 n=9+10)
Int/Random-100K    13.9ms ± 0%   5.4ms ± 0%  -61.20%  (p=0.000 n=9+9)
Int/Random-1M       165ms ± 0%    62ms ± 0%  -62.19%  (p=0.000 n=10+10)
Int/Constant-1K    7.24µs ± 3%  1.33µs ± 2%  -81.63%  (p=0.000 n=9+10)
Int/Constant-10K   51.8µs ± 1%   9.9µs ± 5%  -80.98%  (p=0.000 n=9+10)
Int/Constant-100K   489µs ± 1%    81µs ± 3%  -83.43%  (p=0.000 n=10+10)
Int/Constant-1M    4.91ms ± 1%  0.78ms ± 1%  -84.12%  (p=0.000 n=10+10)
Int/Descent-1K     29.8µs ± 2%   2.9µs ± 2%  -90.35%  (p=0.000 n=10+10)
Int/Descent-10K     337µs ± 1%    21µs ± 5%  -93.75%  (p=0.000 n=10+10)
Int/Descent-100K   4.10ms ± 0%  0.18ms ± 2%  -95.52%  (p=0.000 n=10+10)
Int/Descent-1M     48.9ms ± 0%   1.9ms ± 1%  -96.21%  (p=0.000 n=10+10)
Int/Ascent-1K      28.7µs ± 4%   1.8µs ± 2%  -93.64%  (p=0.000 n=10+10)
Int/Ascent-10K      321µs ± 0%    13µs ± 5%  -95.89%  (p=0.000 n=8+10)
Int/Ascent-100K    3.96ms ± 0%  0.11ms ± 3%  -97.22%  (p=0.000 n=10+10)
Int/Ascent-1M      47.7ms ± 0%   1.1ms ± 1%  -97.69%  (p=0.000 n=10+10)
Int/Mixed-1K       46.0µs ± 3%  21.2µs ± 3%  -53.94%  (p=0.000 n=10+10)
Int/Mixed-10K       601µs ± 0%   260µs ± 1%  -56.81%  (p=0.000 n=10+10)
Int/Mixed-100K     7.39ms ± 0%  2.96ms ± 0%  -59.97%  (p=0.000 n=10+10)
Int/Mixed-1M       88.3ms ± 0%  34.0ms ± 0%  -61.51%  (p=0.000 n=10+10)
Float/1K           96.9µs ± 1%  39.5µs ± 1%  -59.19%  (p=0.000 n=9+10)
Float/10K          1.27ms ± 0%  0.47ms ± 0%  -62.94%  (p=0.000 n=10+10)
Float/100K         15.8ms ± 0%   5.5ms ± 0%  -64.88%  (p=0.000 n=9+9)
Float/1M            188ms ± 0%    64ms ± 0%  -66.10%  (p=0.000 n=8+10)
Str/1K              138µs ± 0%    99µs ± 0%  -28.21%  (p=0.000 n=9+10)
Str/10K            1.85ms ± 0%  1.35ms ± 0%  -27.01%  (p=0.000 n=10+10)
Str/100K           24.1ms ± 0%  18.0ms ± 0%  -25.41%  (p=0.000 n=10+10)
Str/1M              312ms ± 3%   238ms ± 3%  -23.90%  (p=0.000 n=10+10)
Struct/1K           140µs ± 1%    75µs ± 1%  6.46%  (p=0.000 n=10+10)
Struct/10K         1.81ms ± 0%  1.16ms ± 0%  -35.77%  (p=0.000 n=9+10)
Struct/100K        22.3ms ± 0%  14.0ms ± 1%  -37.21%  (p=0.000 n=10+10)
Struct/1M           267ms ± 0%   167ms ± 1%  -37.51%  (p=0.000 n=10+9)
Stable/1K           495µs ± 1%    86µs ± 1%  -82.59%  (p=0.000 n=10+10)
Stable/10K         8.30ms ± 1%  1.38ms ± 1%  -83.34%  (p=0.000 n=10+10)
Stable/100K         132ms ± 0%    17ms ± 1%  -87.12%  (p=0.000 n=9+10)
Stable/1M           1.89s ± 0%   0.21s ± 1%  -89.13%  (p=0.000 n=10+10)
Pointer/1K         95.8µs ± 0%  66.5µs ± 1%  -30.66%  (p=0.000 n=9+9)
Pointer/10K        1.33ms ± 0%  0.93ms ± 0%  -29.60%  (p=0.000 n=10+9)
Pointer/100K       18.1ms ± 1%  13.1ms ± 1%  -27.80%  (p=0.000 n=10+10)
Pointer/1M          264ms ± 3%   204ms ± 3%  -22.59%  (p=0.000 n=10+10)
```
### Compared to generic sort in golang.org/x/exp/slices
```
name               exp time/op  new time/op  delta
Int/Small-1K       17.9µs ± 4%  16.0µs ± 5%  -10.62%  (p=0.000 n=10+10)
Int/Small-10K       311µs ± 0%   236µs ± 1%  -23.99%  (p=0.000 n=10+10)
Int/Small-100K     4.49ms ± 0%  3.23ms ± 0%  -28.02%  (p=0.000 n=10+10)
Int/Small-1M       58.8ms ± 0%  41.1ms ± 0%  -30.21%  (p=0.000 n=9+9)
Int/Random-1K      47.2µs ± 2%  37.1µs ± 2%  -21.25%  (p=0.000 n=10+9)
Int/Random-10K      602µs ± 0%   452µs ± 0%  -25.00%  (p=0.000 n=10+10)
Int/Random-100K    7.40ms ± 0%  5.38ms ± 0%  -27.36%  (p=0.000 n=7+9)
Int/Random-1M      88.1ms ± 1%  62.5ms ± 0%  -29.10%  (p=0.000 n=9+10)
Int/Constant-1K    1.58µs ± 3%  1.33µs ± 2%  -15.72%  (p=0.000 n=10+10)
Int/Constant-10K   11.7µs ± 4%   9.9µs ± 5%  -15.59%  (p=0.000 n=10+10)
Int/Constant-100K  94.3µs ± 3%  81.0µs ± 3%  -14.19%  (p=0.000 n=10+10)
Int/Constant-1M     967µs ± 1%   779µs ± 1%  -19.47%  (p=0.000 n=10+10)
Int/Descent-1K     9.36µs ± 4%  2.88µs ± 2%  -69.25%  (p=0.000 n=10+10)
Int/Descent-10K    93.5µs ± 3%  21.1µs ± 5%  -77.48%  (p=0.000 n=10+10)
Int/Descent-100K   1.09ms ± 0%  0.18ms ± 2%  -83.11%  (p=0.000 n=10+10)
Int/Descent-1M     13.1ms ± 0%   1.9ms ± 1%  -85.86%  (p=0.000 n=9+10)
Int/Ascent-1K      8.89µs ± 4%  1.82µs ± 2%  -79.49%  (p=0.000 n=10+10)
Int/Ascent-10K     88.4µs ± 1%  13.2µs ± 5%  -85.05%  (p=0.000 n=10+10)
Int/Ascent-100K    1.03ms ± 1%  0.11ms ± 3%  -89.38%  (p=0.000 n=10+10)
Int/Ascent-1M      12.7ms ± 0%   1.1ms ± 1%  -91.32%  (p=0.000 n=10+10)
Int/Mixed-1K       22.4µs ± 2%  21.2µs ± 3%   -5.70%  (p=0.000 n=10+10)
Int/Mixed-10K       283µs ± 1%   260µs ± 1%   -8.35%  (p=0.000 n=10+10)
Int/Mixed-100K     3.44ms ± 0%  2.96ms ± 0%  -13.87%  (p=0.000 n=10+10)
Int/Mixed-1M       40.7ms ± 0%  34.0ms ± 0%  -16.45%  (p=0.000 n=8+10)
Float/1K           50.1µs ± 2%  39.5µs ± 1%  -21.05%  (p=0.000 n=10+10)
Float/10K           650µs ± 0%   471µs ± 0%  -27.66%  (p=0.000 n=10+10)
Float/100K         8.03ms ± 0%  5.53ms ± 0%  -31.07%  (p=0.000 n=10+9)
Float/1M           95.6ms ± 0%  63.8ms ± 0%  -33.23%  (p=0.000 n=10+10)
Str/1K              112µs ± 0%    99µs ± 0%  -11.46%  (p=0.000 n=10+10)
Str/10K            1.51ms ± 0%  1.35ms ± 0%  -10.59%  (p=0.000 n=10+10)
Str/100K           19.8ms ± 0%  18.0ms ± 0%   -8.81%  (p=0.000 n=9+10)
Str/1M              255ms ± 2%   238ms ± 3%   -6.66%  (p=0.000 n=10+10)
Struct/1K           114µs ± 1%    75µs ± 1%  -34.31%  (p=0.000 n=10+10)
Struct/10K         1.52ms ± 1%  1.16ms ± 0%  -23.43%  (p=0.000 n=10+10)
Struct/100K        19.0ms ± 0%  14.0ms ± 1%  -25.99%  (p=0.000 n=10+10)
Struct/1M           227ms ± 0%   167ms ± 1%  -26.45%  (p=0.000 n=10+9)
Stable/1K           201µs ± 1%    86µs ± 1%  -57.05%  (p=0.000 n=10+10)
Stable/10K         3.02ms ± 0%  1.38ms ± 1%  -54.24%  (p=0.000 n=10+10)
Stable/100K        43.1ms ± 0%  17.0ms ± 1%  -60.65%  (p=0.000 n=10+10)
Stable/1M           607ms ± 1%   205ms ± 1%  -66.16%  (p=0.000 n=10+10)
Pointer/1K         77.7µs ± 2%  66.5µs ± 1%  -14.43%  (p=0.000 n=10+9)
Pointer/10K        1.09ms ± 0%  0.93ms ± 0%  -14.16%  (p=0.000 n=10+9)
Pointer/100K       15.1ms ± 1%  13.1ms ± 1%  -13.17%  (p=0.000 n=10+10)
Pointer/1M          219ms ± 2%   204ms ± 3%   -6.64%  (p=0.000 n=10+10)
```