package main

import (
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("sample.txt")
	if err != nil {
		return
	}

	if _, err := io.Copy(os.Stdout, f); err != nil {
		log.Fatal(err)
	}
}
