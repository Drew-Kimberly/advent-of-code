package list

func Map[T any, M any](original []T, fn func(T) M) []M {
	mapped := make([]M, len(original))
	for i, val := range original {
		mapped[i] = fn(val)
	}
	return mapped
}
