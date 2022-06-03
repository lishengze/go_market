package utils

import "strings"

func UniqueInt64(list []int64) []int64 {
	keys := make(map[int64]bool)
	var result []int64
	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

func SeparateSymbol(symbol string) (string, string) {
	symbols := strings.Split(symbol, "_")
	if len(symbols) > 1 {
		return symbols[0], symbols[1]
	}

	return "", ""
}
