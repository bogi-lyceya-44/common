package set

type Set[T comparable] map[T]struct{}

func New[T comparable](needles ...T) Set[T] {
	result := make(Set[T], len(needles))

	for _, x := range needles {
		result[x] = struct{}{}
	}

	return result
}

func NewEmptyWithCapacity[T comparable](capacity int) Set[T] {
	return make(Set[T], capacity)
}

func (s Set[T]) Add(needles ...T) {
	for _, x := range needles {
		s[x] = struct{}{}
	}
}

func (s Set[T]) Contains(needle T) bool {
	_, ok := s[needle]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Slice() []T {
	result := make([]T, 0, len(s))

	for x := range s {
		result = append(result, x)
	}

	return result
}

func (s Set[T]) Equal(another Set[T]) bool {
	if s.Len() != another.Len() {
		return false
	}

	for x := range s {
		if !another.Contains(x) {
			return false
		}
	}

	return true
}

func (s Set[T]) Substitute(another Set[T]) Set[T] {
	nonExcluded := make([]T, 0, len(s))

	for x := range s {
		if !another.Contains(x) {
			nonExcluded = append(nonExcluded, x)
		}
	}

	return New(nonExcluded...)
}
