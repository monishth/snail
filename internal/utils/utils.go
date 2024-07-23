package utils

func CopyMapArray[K comparable, V any](dst, src map[K][]V) {
	for k, vv := range src {
		dst[k] = append(dst[k], vv...)
	}
}
