package main

import "os"

func GetOrDefault[K comparable, V any](mp map[K]V, key K, defaultValue V) V {
	val, ok := mp[key]
	if !ok {
		val = defaultValue
	}
	return val
}

func GetEnvOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = defaultValue
	}
	return value
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
