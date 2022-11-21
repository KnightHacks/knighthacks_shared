package utils

import "fmt"

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, len(m))
	i := 0
	for k := range m {
		r[i] = k
		i++
	}
	return r
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, len(m))
	i := 0
	for _, v := range m {
		r[i] = v
		i++
	}
	return r
}

func MapToStrings[M ~map[K]V, K comparable, V any](m M, delimiter string) (keys string, values string) {
	last := len(m) - 1
	for k, v := range m {
		// if last element
		if last == 0 {
			// set delimiter to be an empty string so that nothing is formatted into the strings
			delimiter = ""
		}
		keys += fmt.Sprintf("%s%s", k, delimiter)
		values += fmt.Sprintf("%s%s", v, delimiter)
		last--
	}
	return keys, values
}
