package main

import (
	"fmt"
	"math"
	"math/bits"
)

func Combine(n, m int) ([][]int, error) {
	result := make([][]int, 0)
	if n > math.MaxInt {
		// copied from math/const.go
		// https://cs.opensource.google/go/go/+/refs/tags/go1.24.0:src/math/const.go;l=40
		intSize := 32 << (^uint(0) >> 63) // 32 or 64
		return nil, fmt.Errorf("%d is too large; max: %d", n, intSize)
	}
	nbits := uint(1) << n

	for i := range nbits {
		if bits.OnesCount(i) != m {
			continue
		}

		picks := make([]int, 0, m)
		for j := 0; j < n; j++ {
			if i&(1<<j) > 0 {
				picks = append(picks, j)
			}
		}

		result = append(result, picks)
	}

	return result, nil
}
