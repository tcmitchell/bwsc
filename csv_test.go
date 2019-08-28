package main

import (
	"testing"
)

func TestSliceEqual(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g"}
	if sliceEqual(data[:2], data[:3]) {
		t.Errorf("slices not equal")
	}
	if !sliceEqual(data[1:3], data[1:3]) {
		t.Errorf("slices equal")
	}
}

func TestMergeCSV(t *testing.T) {
	src := [][]string{
		{"08/01/2019", "123456789", "87654321", "1267720", "130"},
		{"08/02/2019", "123456789", "87654321", "1267860", "140"},
		{"08/03/2019", "123456789", "87654321", "1268000", "150"},
		{"08/04/2019", "123456789", "87654321", "1268140", "80"},
	}
	dest := [][]string{
		{"08/04/2019", "123456789", "87654321", "1267720", "160"},
		{"08/05/2019", "123456789", "87654321", "1267860", "170"},
		{"08/06/2019", "123456789", "87654321", "1268000", "180"},
	}

	actual := mergeCSV(src, dest)
	expected := append(src[:3], dest[0], dest[1], dest[2])
	if len(actual) != len(expected) {
		t.Errorf("Expected length %d, got length %d", len(expected), len(actual))
		return
	}
	for i := range actual {
		if len(actual[i]) != len(expected[i]) {
			t.Errorf("Element %d len expected %d, got %d", i, len(expected[i]), len(actual[i]))
			return
		}
		for j := range actual[i] {
			if actual[i][j] != expected[i][j] {
				t.Errorf("Expected %s at %d, %d, got %s",
					expected[i][j], i, j, actual[i][j])
			}
		}
	}
}
