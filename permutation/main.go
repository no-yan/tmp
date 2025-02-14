package main

import (
	"fmt"
	"math"
	"math/bits"
)

func main() {
	for n := range 10 {
		for m := range min(n, 3) {
			fmt.Printf("Permutation(%v, %v)\n", n, m)
			res := Permutate(n, m)
			print(res)
		}
	}

	fmt.Printf("\n\n=================================\n\n")

	for n := range 10 {
		for m := range min(n, 3) {
			fmt.Printf("Combine(%v, %v)\n", n, m)
			res, err := Combine(n, m)
			if err != nil {
				panic(err)
			}
			print(res)
		}
	}
}

func print(lines [][]int) {
	fmt.Println("Count: ", len(lines))
	for _, line := range lines {
		fmt.Println(line)
	}
}

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

func Permutate(n, m int) [][]int {
	// nPm  = n! / (n-r)!
	result := make([][]int, 0, factorial(m, n))
	perm := make([]int, m)
	used := make([]bool, n)

	var dfs func(int)
	dfs = func(depth int) {
		if depth == m {
			cp := make([]int, m)
			copy(cp, perm)
			result = append(result, cp)

			return
		}

		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}

			perm[depth] = i
			used[i] = true
			dfs(depth + 1)
			used[i] = false
		}
	}

	dfs(0)
	return result
}

func factorial(start, end int) uint {
	if start > end {
		err := fmt.Errorf("start is bigger than end; start %d, end %d", start, end)
		panic(err)
	}
	ret := uint(1)
	for i := uint(start); i <= uint(end); i++ {
		ret *= i
	}
	return ret
}
