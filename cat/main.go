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

func open(sources []string) (io.Reader, error) {
	if len(sources) == 0 {
		return os.Stdin, nil
	}

	rs := make([]io.Reader, len(sources))
	for _, s := range sources {
		f, err := os.Open(s)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		rs = append(rs, f)
	}

	return io.MultiReader(rs...), nil
}

func cat(srcs []string) error {
	r, err := open(srcs)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(os.Stdout)
	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
