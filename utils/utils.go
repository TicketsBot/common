package utils

func Ptr[T any](v T) *T {
	return &v
}

func ValueOrZero[T any](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}

func Keys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
