package tests

import (
	"sort"
	"testing"
)

func TestSortSlice(t *testing.T) {
	a := [...]int{5, 4, 3, 2, 1}

	// Check for a subslice.
	s1 := a[1:4]
	sort.Slice(s1, func(i, j int) bool { return s1[i] < s1[j] })
	if a != [...]int{5, 2, 3, 4, 1} {
		t.Fatal("not equal")
	}

	// Check a slice of the whole array.
	s2 := a[:]
	sort.Slice(s2, func(i, j int) bool { return s2[i] < s2[j] })
	if a != [...]int{1, 2, 3, 4, 5} {
		t.Fatal("not equal")
	}

	// Try using a slice with cap.
	a2 := [...]int{6, 5, 4, 3, 2, 1}
	s3 := a2[1:4:4]
	sort.Slice(s3, func(i, j int) bool { return s3[i] < s3[j] })
	if a2 != [...]int{6, 3, 4, 5, 2, 1} {
		t.Fatal("not equal")
	}
}
