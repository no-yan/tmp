package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Job struct {
	v int
}

type Result [32]byte

func worker(jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		results <- doWork(j)
	}
}

func doWork(job Job) Result {
	data := fmt.Appendf(nil, "payload-%d", job.v)
	return sha256.Sum256(data)
}

func main() {
	// Worker poolは、あらかじめ決まった数のゴルーチンを起動しておき、それらにタスクを割り当てるパターンです。
	// メリットとして、ゴルーチンの急増によるメモリ負荷を防ぐことができる。
	// そのため、大量のタスクでも安定したスループットを維持することに向いている。

	// Goの実装では、WaitGroupとchannelを使って実装されることが多い。
	// タスクキュー用のchannelと出力用のチャネルを作成し、
	// タスクキューを追加する。
	// タスクキューの追加が完了したら、close(tasks) でチャネルを閉じ
	// go func() { wg.Wait(); close(outpu); }() で終了をまつ。
	const (
		numJobs    = 5
		numWorkers = 10
	)

	wg := sync.WaitGroup{}
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	for range numWorkers {
		wg.Go(func() { worker(jobs, results) })

	}

	for i := range numJobs {
		jobs <- Job{v: i}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("%x\n", r)
	}

	fmt.Println("done")
}
