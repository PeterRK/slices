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

## [Benchmark](https://gist.github.com/PeterRK/625e8fad081267d00e5f9e9f7a8e2084) Result 
Compared to generic sort in golang.org/x/exp/slices

### On Xeon-8374C (X86-64)
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

### On Yitian-710 (ARM64)
```
name               exp time/op  new time/op  delta
Int/Small-1K       19.7µs ± 0%  16.7µs ± 0%  -15.48%  (p=0.008 n=5+5)
Int/Small-10K       231µs ± 0%   199µs ± 0%  -13.89%  (p=0.008 n=5+5)
Int/Small-100K     2.75ms ± 0%  2.37ms ± 0%  -13.71%  (p=0.008 n=5+5)
Int/Small-1M       32.5ms ± 0%  27.7ms ± 0%  -14.69%  (p=0.008 n=5+5)
Int/Random-1K      35.1µs ± 0%  26.6µs ± 0%  -24.21%  (p=0.008 n=5+5)
Int/Random-10K      455µs ± 0%   350µs ± 0%  -23.18%  (p=0.008 n=5+5)
Int/Random-100K    5.62ms ± 0%  4.36ms ± 0%  -22.30%  (p=0.008 n=5+5)
Int/Random-1M      66.7ms ± 0%  52.3ms ± 0%  -21.67%  (p=0.008 n=5+5)
Int/Constant-1K     908ns ± 2%   626ns ± 3%  -31.08%  (p=0.008 n=5+5)
Int/Constant-10K   7.71µs ± 1%  5.39µs ± 0%  -30.07%  (p=0.008 n=5+5)
Int/Constant-100K  76.1µs ± 0%  52.6µs ± 0%  -30.89%  (p=0.008 n=5+5)
Int/Constant-1M     759µs ± 0%   527µs ± 0%  -30.52%  (p=0.008 n=5+5)
Int/Ascent-1K       915ns ± 2%   616ns ± 0%  -32.65%  (p=0.016 n=5+4)
Int/Ascent-10K     7.71µs ± 0%  5.39µs ± 0%  -30.17%  (p=0.008 n=5+5)
Int/Ascent-100K    76.2µs ± 1%  52.7µs ± 0%  -30.78%  (p=0.008 n=5+5)
Int/Ascent-1M       758µs ± 0%   527µs ± 0%  -30.52%  (p=0.008 n=5+5)
Int/Descent-1K     1.27µs ± 1%  1.00µs ± 0%  -21.61%  (p=0.008 n=5+5)
Int/Descent-10K    11.4µs ± 0%   9.1µs ± 1%  -20.41%  (p=0.008 n=5+5)
Int/Descent-100K    113µs ± 0%    90µs ± 1%  -19.85%  (p=0.008 n=5+5)
Int/Descent-1M     1.15ms ± 0%  0.91ms ± 0%  -20.93%  (p=0.008 n=5+5)
Int/Mixed-1K       14.5µs ± 0%  11.6µs ± 0%  -20.13%  (p=0.008 n=5+5)
Int/Mixed-10K       170µs ± 0%   144µs ± 0%  -15.30%  (p=0.008 n=5+5)
Int/Mixed-100K     1.98ms ± 0%  1.70ms ± 0%  -14.41%  (p=0.008 n=5+5)
Int/Mixed-1M       22.8ms ± 0%  19.6ms ± 0%  -14.01%  (p=0.008 n=5+5)
Hybrid/5%          3.20ms ± 0%  2.45ms ± 0%  -23.44%  (p=0.008 n=5+5)
Hybrid/10%         5.43ms ± 0%  4.17ms ± 0%  -23.20%  (p=0.016 n=5+4)
Hybrid/20%         9.89ms ± 0%  7.60ms ± 0%  -23.15%  (p=0.008 n=5+5)
Hybrid/30%         14.4ms ± 0%  11.0ms ± 0%  -23.15%  (p=0.008 n=5+5)
Hybrid/50%         23.3ms ± 0%  17.9ms ± 0%  -23.15%  (p=0.008 n=5+5)
Float/1K           41.1µs ± 1%  32.5µs ± 0%  -20.92%  (p=0.008 n=5+5)
Float/10K           534µs ± 0%   430µs ± 0%  -19.45%  (p=0.008 n=5+5)
Float/100K         6.59ms ± 0%  5.38ms ± 0%  -18.27%  (p=0.008 n=5+5)
Float/1M           78.2ms ± 0%  64.5ms ± 0%  -17.52%  (p=0.008 n=5+5)
Str/1K             89.6µs ± 0%  86.2µs ± 0%   -3.88%  (p=0.008 n=5+5)
Str/10K            1.14ms ± 0%  1.08ms ± 0%   -4.86%  (p=0.008 n=5+5)
Str/100K           15.0ms ± 0%  14.5ms ± 0%   -3.08%  (p=0.008 n=5+5)
Str/1M              218ms ± 1%   210ms ± 1%   -3.77%  (p=0.008 n=5+5)
Struct/1K           121µs ± 0%    62µs ± 0%  -48.39%  (p=0.008 n=5+5)
Struct/10K         1.62ms ± 0%  1.25ms ± 1%  -22.75%  (p=0.008 n=5+5)
Struct/100K        20.3ms ± 0%  14.6ms ± 2%  -28.03%  (p=0.008 n=5+5)
Struct/1M           244ms ± 0%   155ms ± 0%  -36.53%  (p=0.008 n=5+5)
Stable/1K           235µs ± 0%    74µs ± 0%  -68.46%  (p=0.008 n=5+5)
Stable/10K         3.67ms ± 0%  1.44ms ± 1%  -60.75%  (p=0.008 n=5+5)
Stable/100K        53.5ms ± 0%  16.4ms ± 1%  -69.44%  (p=0.008 n=5+5)
Stable/1M           722ms ± 1%   183ms ± 0%  -74.67%  (p=0.008 n=5+5)
Pointer/1K         61.7µs ± 0%  54.2µs ± 0%  -12.24%  (p=0.008 n=5+5)
Pointer/10K         866µs ± 0%   772µs ± 0%  -10.83%  (p=0.008 n=5+5)
Pointer/100K       12.4ms ± 1%  11.4ms ± 0%   -8.36%  (p=0.008 n=5+5)
Pointer/1M          162ms ± 0%   156ms ± 1%   -3.44%  (p=0.008 n=5+5)
```