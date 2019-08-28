package main

// import (
// 	"fmt"
// )

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func mergeCSV(data1 [][]string, data2 [][]string) [][]string {
	result := make([][]string, 0)
	for i := range data1 {
		result = append(result, data1[i])
	}
	for i := range data2 {
		inserted := false
		for j := range result {
			if sliceEqual(result[j][0:3], data2[i][0:3]) {
				result[j] = data2[i]
				inserted = true
				break
			}
		}
		if !inserted {
			result = append(result, data2[i])
		}
	}
	return result
}
