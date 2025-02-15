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

func CombineLex(n, m int) [][]int {
	result := make([][]int, 0, binomial(n, m))
	comb := make([]int, m)

	var dfs func(int)
	dfs = func(depth int) {
		if depth == m {
			cp := make([]int, m)
			copy(cp, comb)
			result = append(result, cp)
			return
		}

		minI := 0
		if depth > 0 {
			minI = comb[depth-1] + 1
		}
		for i := minI; i < n; i++ {
			comb[depth] = i
			dfs(depth + 1)
		}
		comb[depth] = -1
	}
	dfs(0)
	return result
}

func binomial(n, k int) uint {
	if k > n {
		return 0
	}

	if k == 0 || k == n {
		return 1
	}

	// 計算量を抑えるために nCk = nC(n-k) を利用
	if k > n-k {
		k = n - k
	}

	res := uint(1)
	for i := 0; i < k; i++ {
		// 累積的にかけて割る方法で大きな中間値を抑える
		res *= uint(n-i) / uint(i+1)
		// n-1/i+1 が割り切れる保証:
		// 1. nCk は整数である
		// 2. また、nC0 は1である
		// 3. i回目のループでresはnCiに等しい

		// 3(alt). nCk+1 = nCk * (n-1) / (i+1) の等式が成り立つ
	}

	return res
}
