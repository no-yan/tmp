package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	fileNames := parse()
	if err := cat(fileNames); err != nil {
		log.Fatal(err)
	}
}

func parse() []string {
	return os.Args[1:]
}

func cat(srcs []string) error {
	// if src is empty, read from Stdin
	if len(srcs) == 0 {
		if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
			return err
		}
		return nil
	}

	w := bufio.NewWriter(os.Stdout)
	for _, src := range srcs {
		if err := copyFileToWriter(src, w); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func copyFileToWriter(src string, w io.Writer) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return err
	}

	return nil
}
