package main

import (
	"fmt"
	"io"
	"os"
)

func bigger() error {
	f, err := os.OpenFile("sample.txt", os.O_RDWR, 0o644)
	if err != nil {
		return err
	}

	dst, err := os.Create("large.txt")
	if err != nil {
		return err
	}
	defer dst.Close()

	for range 1000 {
		if _, err := f.Seek(0, 0); err != nil {
			return err
		}
		if _, err = io.Copy(dst, f); err != nil {
			return err
		}
	}

	fileInfo, err := dst.Stat()
	if err != nil {
		return err
	}
	sizeBytes := fileInfo.Size()
	sizeMB := float64(sizeBytes) / (1024 * 1024)
	fmt.Printf("%s: %f MB\n", fileInfo.Name(), sizeMB)

	return nil
}
