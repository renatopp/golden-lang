package iter

func Each[T any](items []T, fn func(T)) {
	for _, item := range items {
		fn(item)
	}
}

func Map[T, R any](items []T, fn func(T) R) []R {
	var result []R
	for _, item := range items {
		result = append(result, fn(item))
	}
	return result
}

func Filter[T any](items []T, fn func(T) bool) []T {
	var result []T
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}
