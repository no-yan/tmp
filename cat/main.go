package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	fileNames := parse()
	if err := cat(fileNames, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func parse() []string {
	return os.Args[1:]
}

func cat(srcs []string, in io.Reader, out io.Writer) error {
	// if src is empty, read from Stdin
	if len(srcs) == 0 {
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		return nil
	}

	w := bufio.NewWriter(out)
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
