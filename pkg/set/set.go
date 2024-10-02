package set

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) Remove(element T) {
	delete(s, element)
}

func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

func (s Set[T]) Size() int {
	return len(s)
}
