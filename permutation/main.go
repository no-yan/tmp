package main

import (
	"fmt"
)

func main() {
	for n := range 10 {
		for m := range min(n, 3) {
			fmt.Printf("Permutation(%v, %v)\n", n, m)
			res := Permutate(n, m)
			printResult(res)
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
			printResult(res)
		}
	}
}

func printResult(lines [][]int) {
	fmt.Println("Count: ", len(lines))
	for _, line := range lines {
		fmt.Println(line)
	}
}

func Permutate(n, m int) [][]int {
	result := make([][]int, 0, productRange(m, n)) // nPm  = n! / (n-r)!
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

func productRange(start, end int) uint {
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
