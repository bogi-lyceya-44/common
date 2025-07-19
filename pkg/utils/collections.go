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
// Example: [[1, 2, 3], [4, 5], [6]] -> [1, 2, 3, 4, 5, 6]
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
