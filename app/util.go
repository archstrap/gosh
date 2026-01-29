package main

import (
	"os"
	"slices"
)

func GetOrDefault[K comparable, V any](mp map[K]V, key K, defaultValue V) V {
	val, ok := mp[key]
	if !ok {
		val = defaultValue
	}
	return val
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func Do[K any](items []*K, callBack func(item *K)) {
	for _, item := range items {
		callBack(item)
	}
}

func PerformTask[K any](items []K, callBack func(item K)) {
	for _, item := range items {
		callBack(item)
	}
}

func AddItems(dest *[]string, src *[]string) {
	for _, item := range *src {
		if !slices.Contains(*dest, item) {
			*dest = append(*dest, item)
		}
	}
}
