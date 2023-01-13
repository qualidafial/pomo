package set

type Set[T comparable] map[T]struct{}

func New[T comparable]() Set[T] {
	return make(Set[T])
}

func Of[T comparable](ts ...T) Set[T] {
	s := New[T]()
	for _, t := range ts {
		s.Add(t)
	}
	return s
}

func (s Set[T]) Add(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Remove(t T) {
	delete(s, t)
}

func (s Set[T]) Toggle(t T) {
	if s.Contains(t) {
		s.Remove(t)
	} else {
		s.Add(t)
	}
}

func (s Set[T]) Contains(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Slice() []T {
	ts := make([]T, 0, len(s))
	for t := range s {
		ts = append(ts, t)
	}
	return ts
}

func (s Set[T]) Clone() Set[T] {
	clone := New[T]()
	for t := range s {
		clone.Add(t)
	}
	return clone
}
