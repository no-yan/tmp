package main_test

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	main "github.com/no-yan/tmp/permutation"
)

func TestCombine(t *testing.T) {
	tests := []struct {
		name string
		n    int
		m    int
		want [][]int
	}{
		{"1C1", 1, 1, [][]int{{0}}},
		{"2C1", 2, 1, [][]int{{0}, {1}}},
		{"0C0", 0, 0, [][]int{{}}},
		{
			"5C2", 5, 2,
			[][]int{
				{0, 1},
				{0, 2},
				{0, 3},
				{0, 4},
				{1, 2},
				{1, 3},
				{1, 4},
				{2, 3},
				{2, 4},
				{3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := main.Combine(tt.n, tt.m)
			if err != nil {
				t.Fatalf("Combine() failed unecpectedly; %v", err)
			}
			slices.SortFunc(got, func(i, j []int) int {
				n := len(i)
				if len(j) < n {
					n = len(j)
				}

				for idx := range n {
					if i[idx] == j[idx] {
						continue
					}
					return i[idx] - j[idx]
				}
				return 0
			})

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Combine() mismatch; diff\n%s", diff)
			}
		})
	}
}
