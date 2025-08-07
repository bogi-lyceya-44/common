package utils

// Map maps the array to another array
// with elements that were converted by provided function.
//
// Example: Map([1, 2, 3], f(x) -> x*2) -> [2, 4, 6]
func Map[T, U any](arr []T, f func(x T) U) []U {
	result := make([]U, 0, len(arr))

	for _, x := range arr {
		result = append(result, f(x))
	}

	return result
}

// MapWithError maps the array to another array
// with elements that were converted by provided function.
// If the provided function returns an error,
// MapWithError returns nil array and this error as a result.
//
// Example: Map([1, 2, 3], f(x) -> x*2) -> [2, 4, 6]
func MapWithError[T, U any](arr []T, f func(x T) (U, error)) ([]U, error) {
	result := make([]U, 0, len(arr))

	for _, x := range arr {
		mapped, err := f(x)
		if err != nil {
			return nil, err
		}

		result = append(result, mapped)
	}

	return result, nil
}

// GroupBy groups the array by result of provided function.
// The value by key contains ALL elements of the array
// where the key is equal to the result of the provided function.
//
// Example: GroupBySingle([1, 2, 3, 1, 2], f(x) -> x) -> map[1:[1, 1], 2:[2, 2], 3:[3]]
func GroupBy[T any, U comparable](arr []T, f func(x T) U) map[U][]T {
	result := make(map[U][]T, len(arr))

	for _, x := range arr {
		result[f(x)] = append(result[f(x)], x)
	}

	return result
}

// GroupBySingle groups the array by result of provided function.
// The value by key contains the LAST element of the array
// where the key is equal to the result of the provided function.
//
// Example: GroupBySingle([1, 2, 3], f(x) -> x) -> map[1:1, 2:2, 3:3]
func GroupBySingle[T any, U comparable](arr []T, f func(x T) U) map[U]T {
	result := make(map[U]T, len(arr))

	for _, x := range arr {
		result[f(x)] = x
	}

	return result
}

// Split splits array by provided predicate in such way
// that the first resulting array contains elements where
// predicate(x) == true and the second resulting array
// contains elements where predicate(x) == false.
//
// Example: Split([1, 2, 3], isOdd) -> [[1, 3], [2]]
func Split[T any](arr []T, f func(x T) bool) [][]T {
	result := make([][]T, 2)

	for _, x := range arr {
		part := 0

		if !f(x) {
			part = 1
		}

		result[part] = append(result[part], x)
	}

	return result
}

// Filter filters array by provided predicate.
//
// Example: Filter([1, 2, 3], isOdd) -> [1, 3]
func Filter[T any](arr []T, f func(x T) bool) []T {
	result := make([]T, 0, len(arr))

	for _, x := range arr {
		if f(x) {
			result = append(result, x)
		}
	}

	return result
}

// Batch batches array into small arrays of max length batchSize.
//
// Example: Batch([1, 2, 3, 4, 5, 6, 7], 3) -> [[1, 2, 3], [4, 5, 6], [7]]
func Batch[T any](arr []T, batchSize int) [][]T {
	if batchSize <= 0 {
		// fail fast
		panic("batch size <= 0")
	}

	result := make([][]T, 0, (len(arr)+batchSize-1)/batchSize)
	i := 0

	for i < len(arr) {
		l, r := i, min(i+batchSize, len(arr))
		result = append(result, arr[l:r:r])
		i = r
	}

	return result
}

// Flatten flattens two-dimensional array into one-dimensional.
//
// Example: Flatten([[1, 2, 3], [4, 5], [6]]) -> [1, 2, 3, 4, 5, 6]
func Flatten[T any](arr [][]T) []T {
	totalLength := 0

	for _, x := range arr {
		totalLength += len(x)
	}

	result := make([]T, 0, totalLength)

	for _, x := range arr {
		result = append(result, x...)
	}

	return result
}

// Unique removes duplicates from the provided array.
//
// Example: Unique([1, 2, 3, 3, 2, 1]) -> [1, 2, 3]
func Unique[T comparable](arr []T) []T {
	result := make([]T, 0, len(arr))
	exist := make(map[T]struct{}, len(arr))

	for _, x := range arr {
		if _, ok := exist[x]; ok {
			continue
		}

		exist[x] = struct{}{}
		result = append(result, x)
	}

	return result
}

// Keys returns a slice of keys retrieved from provided map.
//
// Example: Keys(map[1:1, 2:2, 3:4]) -> [1, 2, 3]
func Keys[T comparable, U any](items map[T]U) []T {
	result := make([]T, 0, len(items))

	for k := range items {
		result = append(result, k)
	}

	return result
}

// Values returns a slice of values retrieved from provided map.
//
// Example: Values(map[1:1, 2:2, 3:4]) -> [1, 2, 4]
func Values[T comparable, U any](items map[T]U) []U {
	result := make([]U, 0, len(items))

	for _, v := range items {
		result = append(result, v)
	}

	return result
}

// Zip "zips" two slices into one slice of pairs,
// where one pair has elements from slices under the same index.
// Extra elements will be discarded.
//
// Example: Zip([1, 2, 3, 4], [5, 6]) -> [{1, 5}, {2, 6}]
func Zip[T, U any](l []T, r []U) []Pair[T, U] {
	totalLength := min(len(l), len(r))

	result := make([]Pair[T, U], 0, totalLength)

	for i := range totalLength {
		result = append(
			result,
			Pair[T, U]{
				First:  l[i],
				Second: r[i],
			},
		)
	}

	return result
}

// Reduce accumulates some value over the passed array.
// JS analogue: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/reduce
//
// Example: Reduce([1, 2, 3, 4], f(accumulator, x) -> accumulator + x, 0) -> 10
func Reduce[T, U any](
	arr []T,
	reducer func(U, T) U,
	initialValue U,
) U {
	accumulator := initialValue

	for _, item := range arr {
		accumulator = reducer(accumulator, item)
	}

	return accumulator
}
