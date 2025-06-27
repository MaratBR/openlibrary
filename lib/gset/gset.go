package gset

type Set[T comparable] struct {
	seen map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		seen: make(map[T]struct{}),
	}
}

func (s *Set[T]) Clear() {
	s.seen = make(map[T]struct{})
}

func (s Set[T]) GetMap() map[T]struct{} {
	return s.seen
}

func (s *Set[T]) Add(v T) {
	if _, ok := s.seen[v]; !ok {
		s.seen[v] = struct{}{}
	}
}

func (s *Set[T]) Arr() []T {
	arr := make([]T, 0, len(s.seen))
	for k := range s.seen {
		arr = append(arr, k)
	}
	return arr
}

func (s *Set[T]) Len() int {
	return len(s.seen)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.seen[v]
	return ok
}

func (s *Set[T]) Remove(v T) {
	delete(s.seen, v)
}
