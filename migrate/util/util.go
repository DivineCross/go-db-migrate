package util

import "slices"

func ToSorted[S ~[]E, E any](source S, compare func(a, b E) int) S {
	result := slices.Clone(source)
	slices.SortStableFunc(result, compare)
	return result
}

func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
