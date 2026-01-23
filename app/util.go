package main

func GetOrDefault[K comparable, V any](mp map[K]V, key K, defaultValue V) V {
	val, ok := mp[key]
	if !ok {
		val = defaultValue
	}
	return val
}
