package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileSaver struct {
	dir string
}

func NewFileSaver(dir string) *FileSaver {
	return &FileSaver{dir: dir}
}

func (fs FileSaver) Save(r io.Reader, url string) (int64, error) {
	err := os.Mkdir(outDir, 0o644)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	fName := fs.createFileName(url)
	path := filepath.Join(outDir, fName)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(f, r)
}

func (fs FileSaver) createFileName(url string) string {
	b := sha256.Sum256([]byte(url))
	hex := hex.EncodeToString(b[:])
	return fmt.Sprintf("%s.log", hex)
}
