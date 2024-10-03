package set

type Set[T comparable] map[T]struct{}

func New[T comparable]() Set[T] {
	return Set[T]{}
}

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

func (s Set[T]) Clear() {
    for k := range s {
        delete(s, k)
    }
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) Empty() bool {
	return len(s) == 0
}

// Experimenting...
func (s Set[T]) Iter() <-chan T {
    output := make(chan T)
    go func() {
        defer close(output)
        for v := range s {
            output <- v
        }
    }()
    return output
}
