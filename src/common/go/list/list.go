package list

func Map[T any, M any](original []T, fn func(val T, i int) M) []M {
	mapped := make([]M, len(original))
	for i, val := range original {
		mapped[i] = fn(val, i)
	}
	return mapped
}

func Reduce[T any, R any](original []T, fn func(acc R, nextVal T, i int) R, initial R) R {
	reduced := initial
	for i, val := range original {
		reduced = fn(reduced, val, i)
	}
	return reduced
}

func Filter[T any](original []T, predicate func(val T, i int) bool) []T {
	var filtered []T
	for i, val := range original {
		if predicate(val, i) {
			filtered = append(filtered, val)
		}
	}

	return filtered
}
