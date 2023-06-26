# Fast generic sort for slices in golang

This package is based on the slices package in go standard library since v1.21, but has different apis for search and sort.

## API for builtin types
```go
func BinarySearch[E constraints.Ordered](list []E, x E) (int, bool)
func IsSorted[E constraints.Ordered](list []E) bool
func Min[E cmp.Ordered](list []E) E
func Max[E cmp.Ordered](list []E) E
func Sort[E constraints.Ordered](list []E)
func SortStable[E constraints.Ordered](list []E)
```

## API for custom types
```go
type Order[E any] struct {
	Less    func(a, b E) bool
	RefLess func(a, b *E) bool
}

func (od *Order[E]) BinarySearch(list []E, x E) (int, bool)
func (od *Order[E]) IsSorted(list []E) bool
func (od *Order[E]) Min(list []E) E
func (od *Order[E]) Max(list []E) E
func (od *Order[E]) Sort(list []E)
func (od *Order[E]) SortStable(list []E)
func (od *Order[E]) SortWithOption(list []E, stable, inplace bool)
```

## Benchmark Result

### On Xeon-8374C (X86-64)
```
name               std time/op  new time/op  delta
Int/Small-1K       27.7µs ± 2%  24.1µs ± 3%  -12.92%  (p=0.008 n=5+5)
Int/Small-10K       315µs ± 0%   249µs ± 0%  -20.96%  (p=0.008 n=5+5)
Int/Small-100K     3.76ms ± 0%  2.86ms ± 0%  -24.06%  (p=0.008 n=5+5)
Int/Small-1M       44.6ms ± 0%  35.4ms ± 0%  -20.57%  (p=0.008 n=5+5)
Int/Random-1K      48.4µs ± 1%  38.8µs ± 1%  -19.69%  (p=0.008 n=5+5)
Int/Random-10K      614µs ± 1%   464µs ± 0%  -24.44%  (p=0.008 n=5+5)
Int/Random-100K    7.58ms ± 0%  5.61ms ± 0%  -26.07%  (p=0.008 n=5+5)
Int/Random-1M      90.3ms ± 1%  65.8ms ± 0%  -27.08%  (p=0.008 n=5+5)
Int/Constant-1K    1.33µs ± 1%  0.93µs ± 1%  -30.13%  (p=0.008 n=5+5)
Int/Constant-10K   11.1µs ± 4%   8.0µs ± 2%  -27.36%  (p=0.008 n=5+5)
Int/Constant-100K  98.0µs ± 5%  67.0µs ± 2%  -31.61%  (p=0.008 n=5+5)
Int/Constant-1M     932µs ± 1%   646µs ± 0%  -30.72%  (p=0.008 n=5+5)
Int/Ascent-1K      1.33µs ± 1%  0.93µs ± 1%  -30.03%  (p=0.008 n=5+5)
Int/Ascent-10K     10.8µs ± 3%   8.0µs ± 3%  -26.21%  (p=0.008 n=5+5)
Int/Ascent-100K    99.1µs ± 5%  66.2µs ± 1%  -33.20%  (p=0.008 n=5+5)
Int/Ascent-1M       926µs ± 0%   646µs ± 0%  -30.19%  (p=0.008 n=5+5)
Int/Descent-1K     1.85µs ± 1%  1.46µs ± 0%  -21.19%  (p=0.008 n=5+5)
Int/Descent-10K    14.8µs ± 3%  12.0µs ± 6%  -18.43%  (p=0.008 n=5+5)
Int/Descent-100K    132µs ± 1%   104µs ± 2%  -21.46%  (p=0.008 n=5+5)
Int/Descent-1M     1.32ms ± 0%  1.04ms ± 1%  -21.37%  (p=0.008 n=5+5)
Int/Mixed-1K       19.9µs ± 3%  17.1µs ± 3%  -13.87%  (p=0.008 n=5+5)
Int/Mixed-10K       219µs ± 1%   231µs ± 0%   +5.53%  (p=0.008 n=5+5)
Int/Mixed-100K     2.54ms ± 0%  2.55ms ± 0%     ~     (p=0.310 n=5+5)
Int/Mixed-1M       29.5ms ± 0%  28.0ms ± 0%   -5.11%  (p=0.008 n=5+5)
Hybrid/5%-4        4.09ms ± 0%  3.06ms ± 0%  -25.07%  (p=0.008 n=5+5)
Hybrid/10%-4       7.12ms ± 0%  5.33ms ± 0%  -25.04%  (p=0.008 n=5+5)
Hybrid/20%-4       13.1ms ± 0%   9.9ms ± 0%  -24.73%  (p=0.008 n=5+5)
Hybrid/30%-4       19.2ms ± 1%  14.4ms ± 0%  -25.00%  (p=0.008 n=5+5)
Hybrid/50%-4       31.2ms ± 1%  23.5ms ± 0%  -24.70%  (p=0.008 n=5+5)
Float/1K           58.2µs ± 2%  48.6µs ± 1%  -16.58%  (p=0.008 n=5+5)
Float/10K           751µs ± 0%   549µs ± 0%  -26.80%  (p=0.008 n=5+5)
Float/100K         9.28ms ± 0%  6.41ms ± 0%  -30.91%  (p=0.008 n=5+5)
Float/1M            110ms ± 0%    73ms ± 0%  -33.44%  (p=0.008 n=5+5)
Str/1K              130µs ± 0%   125µs ± 0%   -4.00%  (p=0.008 n=5+5)
Str/10K            1.69ms ± 0%  1.64ms ± 0%   -2.73%  (p=0.008 n=5+5)
Str/100K           21.6ms ± 1%  21.1ms ± 1%   -2.41%  (p=0.008 n=5+5)
Str/1M              272ms ± 1%   269ms ± 2%     ~     (p=0.095 n=5+5)
Struct/1K           104µs ± 1%    86µs ± 0%  -17.34%  (p=0.008 n=5+5)
Struct/10K         1.36ms ± 0%  1.07ms ± 0%  -20.93%  (p=0.008 n=5+5)
Struct/100K        16.8ms ± 0%  13.0ms ± 0%  -22.48%  (p=0.008 n=5+5)
Struct/1M           201ms ± 0%   152ms ± 0%  -24.21%  (p=0.008 n=5+5)
Stable/1K           183µs ± 1%    84µs ± 0%  -53.88%  (p=0.008 n=5+5)
Stable/10K         2.71ms ± 0%  1.14ms ± 0%  -58.16%  (p=0.008 n=5+5)
Stable/100K        37.9ms ± 0%  17.4ms ± 0%  -54.15%  (p=0.016 n=5+4)
Stable/1M           498ms ± 1%   175ms ± 1%  -64.74%  (p=0.008 n=5+5)
Pointer/1K         78.7µs ± 1%  66.1µs ± 1%  -15.96%  (p=0.008 n=5+5)
Pointer/10K        1.05ms ± 0%  0.90ms ± 0%  -14.87%  (p=0.008 n=5+5)
Pointer/100K       14.0ms ± 0%  12.1ms ± 0%  -13.52%  (p=0.008 n=5+5)
Pointer/1M          171ms ± 0%   159ms ± 0%   -6.82%  (p=0.008 n=5+5)
```

### On Yitian-710 (ARM64)
```
name               std time/op  new time/op  delta
Int/Small-1K       21.4µs ± 0%  17.4µs ± 0%  -18.95%  (p=0.008 n=5+5)
Int/Small-10K       253µs ± 0%   207µs ± 0%  -18.18%  (p=0.008 n=5+5)
Int/Small-100K     3.02ms ± 0%  2.47ms ± 0%  -18.42%  (p=0.008 n=5+5)
Int/Small-1M       35.7ms ± 0%  28.8ms ± 0%  -19.54%  (p=0.008 n=5+5)
Int/Random-1K      37.7µs ± 0%  27.6µs ± 0%  -26.84%  (p=0.008 n=5+5)
Int/Random-10K      492µs ± 0%   362µs ± 0%  -26.46%  (p=0.008 n=5+5)
Int/Random-100K    6.09ms ± 0%  4.51ms ± 0%  -25.96%  (p=0.008 n=5+5)
Int/Random-1M      72.5ms ± 0%  53.9ms ± 0%  -25.69%  (p=0.008 n=5+5)
Int/Constant-1K    1.08µs ± 1%  0.70µs ± 0%  -35.23%  (p=0.016 n=5+4)
Int/Constant-10K   9.60µs ± 0%  6.55µs ± 0%  -31.77%  (p=0.008 n=5+5)
Int/Constant-100K  93.1µs ± 0%  63.5µs ± 1%  -31.78%  (p=0.008 n=5+5)
Int/Constant-1M     929µs ± 0%   640µs ± 1%  -31.09%  (p=0.008 n=5+5)
Int/Ascent-1K      1.08µs ± 1%  0.71µs ± 4%  -34.36%  (p=0.008 n=5+5)
Int/Ascent-10K     9.61µs ± 0%  6.57µs ± 0%  -31.57%  (p=0.008 n=5+5)
Int/Ascent-100K    93.1µs ± 0%  63.1µs ± 0%  -32.20%  (p=0.008 n=5+5)
Int/Ascent-1M       929µs ± 0%   639µs ± 1%  -31.27%  (p=0.008 n=5+5)
Int/Descent-1K     1.45µs ± 0%  1.07µs ± 1%  -26.66%  (p=0.008 n=5+5)
Int/Descent-10K    13.2µs ± 0%   9.9µs ± 0%  -25.07%  (p=0.016 n=5+4)
Int/Descent-100K    130µs ± 0%   100µs ± 0%  -23.09%  (p=0.008 n=5+5)
Int/Descent-1M     1.30ms ± 0%  1.02ms ± 1%  -22.02%  (p=0.008 n=5+5)
Int/Mixed-1K       16.3µs ± 4%  12.1µs ± 0%  -25.65%  (p=0.008 n=5+5)
Int/Mixed-10K       187µs ± 0%   150µs ± 0%  -19.82%  (p=0.008 n=5+5)
Int/Mixed-100K     2.20ms ± 0%  1.77ms ± 0%  -19.24%  (p=0.008 n=5+5)
Int/Mixed-1M       25.3ms ± 0%  20.5ms ± 0%  -19.21%  (p=0.008 n=5+5)
Hybrid/5%-4        3.50ms ± 0%  2.59ms ± 0%  -26.14%  (p=0.008 n=5+5)
Hybrid/10%-4       5.92ms ± 0%  4.36ms ± 0%  -26.33%  (p=0.008 n=5+5)
Hybrid/20%-4       10.8ms ± 0%   7.9ms ± 0%  -26.48%  (p=0.008 n=5+5)
Hybrid/30%-4       15.6ms ± 0%  11.5ms ± 0%  -26.43%  (p=0.008 n=5+5)
Hybrid/50%-4       25.3ms ± 0%  18.6ms ± 0%  -26.49%  (p=0.008 n=5+5)
Float/1K           46.1µs ± 0%  35.6µs ± 0%  -22.75%  (p=0.008 n=5+5)
Float/10K           606µs ± 0%   467µs ± 0%  -22.85%  (p=0.008 n=5+5)
Float/100K         7.51ms ± 0%  5.79ms ± 0%  -22.83%  (p=0.008 n=5+5)
Float/1M           89.5ms ± 0%  69.0ms ± 0%  -22.84%  (p=0.008 n=5+5)
Str/1K              124µs ± 0%   119µs ± 0%   -4.40%  (p=0.008 n=5+5)
Str/10K            1.54ms ± 0%  1.48ms ± 0%   -3.70%  (p=0.008 n=5+5)
Str/100K           20.1ms ± 1%  19.6ms ± 0%   -2.37%  (p=0.008 n=5+5)
Str/1M              290ms ± 1%   286ms ± 1%   -1.10%  (p=0.032 n=5+5)
Struct/1K           100µs ± 0%    85µs ± 0%  -14.55%  (p=0.008 n=5+5)
Struct/10K         1.32ms ± 0%  1.07ms ± 0%  -19.19%  (p=0.008 n=5+5)
Struct/100K        16.4ms ± 0%  12.7ms ± 0%  -22.38%  (p=0.008 n=5+5)
Struct/1M           196ms ± 0%   141ms ± 2%  -28.06%  (p=0.008 n=5+5)
Stable/1K           186µs ± 0%    73µs ± 0%  -60.49%  (p=0.008 n=5+5)
Stable/10K         2.84ms ± 0%  1.22ms ± 0%  -57.16%  (p=0.008 n=5+5)
Stable/100K        40.5ms ± 0%  14.7ms ± 0%  -63.75%  (p=0.008 n=5+5)
Stable/1M           534ms ± 0%   163ms ± 0%  -69.43%  (p=0.016 n=5+4)
Pointer/1K         63.7µs ± 0%  55.3µs ± 0%  -13.18%  (p=0.008 n=5+5)
Pointer/10K         873µs ± 0%   775µs ± 1%  -11.23%  (p=0.008 n=5+5)
Pointer/100K       12.4ms ± 0%  11.3ms ± 0%   -9.58%  (p=0.008 n=5+5)
Pointer/1M          159ms ± 0%   147ms ± 0%   -7.59%  (p=0.008 n=5+5)
```