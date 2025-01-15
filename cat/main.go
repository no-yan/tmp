package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("sample.txt")
	if err != nil {
		return
	}

	w := bufio.NewWriter(os.Stdout)
	if _, err := io.Copy(w, f); err != nil {
		log.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		return
	}
}
