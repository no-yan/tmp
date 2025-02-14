package main

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			got, err := Combine(tt.n, tt.m)
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

func TestPermutate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		n    int
		m    int
		want [][]int
	}{
		{"0P0", 0, 0, [][]int{{}}},
		{"1P0", 1, 0, [][]int{{}}},
		{"1P1", 1, 1, [][]int{{0}}},
		{"2P2", 2, 2, [][]int{{0, 1}, {1, 0}}},

		{
			"5P3", 5, 3,
			[][]int{{0, 1, 2}, {0, 1, 3}, {0, 1, 4}, {0, 2, 1}, {0, 2, 3}, {0, 2, 4}, {0, 3, 1}, {0, 3, 2}, {0, 3, 4}, {0, 4, 1}, {0, 4, 2}, {0, 4, 3}, {1, 0, 2}, {1, 0, 3}, {1, 0, 4}, {1, 2, 0}, {1, 2, 3}, {1, 2, 4}, {1, 3, 0}, {1, 3, 2}, {1, 3, 4}, {1, 4, 0}, {1, 4, 2}, {1, 4, 3}, {2, 0, 1}, {2, 0, 3}, {2, 0, 4}, {2, 1, 0}, {2, 1, 3}, {2, 1, 4}, {2, 3, 0}, {2, 3, 1}, {2, 3, 4}, {2, 4, 0}, {2, 4, 1}, {2, 4, 3}, {3, 0, 1}, {3, 0, 2}, {3, 0, 4}, {3, 1, 0}, {3, 1, 2}, {3, 1, 4}, {3, 2, 0}, {3, 2, 1}, {3, 2, 4}, {3, 4, 0}, {3, 4, 1}, {3, 4, 2}, {4, 0, 1}, {4, 0, 2}, {4, 0, 3}, {4, 1, 0}, {4, 1, 2}, {4, 1, 3}, {4, 2, 0}, {4, 2, 1}, {4, 2, 3}, {4, 3, 0}, {4, 3, 1}, {4, 3, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Permutate(tt.n, tt.m)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Permutate() mismatch; diff:\n%s", diff)
			}
		})
	}
}
