# slices
generic sort for slices in golang

## basic API
```go
func BinarySearch[E constraints.Ordered](list []E, x E) int
func IsSorted[E constraints.Ordered](list []E) bool
func Sort[E constraints.Ordered](list []E)
func SortStable[E constraints.Ordered](list []E)
```

## fast API
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

## func API
```go
func BinarySearchFunc[E any](list []E, x E, less func(a, b E) bool) int
func IsSortedFunc[E any](list []E, less func(a, b E) bool) bool
func SortFunc[E any](list []E, less func(a, b E) bool)
func SortStableFunc[E any](list []E, less func(a, b E) bool)
```
