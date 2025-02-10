package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type FileSaver struct {
	dir  string
	once *sync.Once
	err  error
	fs   FileSystem
}

type FileSystem interface {
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	MkdirAll(path string, perm fs.FileMode) error
	IsExist(err error) bool
}

type osfs struct{}

func NewOSFS() osfs {
	return osfs{}
}

func (o osfs) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (o osfs) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o osfs) IsExist(err error) bool {
	return os.IsExist(err)
}

func NewFileSaver(dir string, fs FileSystem) *FileSaver {
	return &FileSaver{dir: dir, once: &sync.Once{}, err: nil, fs: fs}
}

func (fs FileSaver) Save(r io.Reader, url string) (int64, error) {
	err := fs.ensureDir()
	if err != nil {
		return 0, err
	}

	fName := fs.createFileName(url)
	path := filepath.Join(fs.dir, fName)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
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

func (fs *FileSaver) ensureDir() error {
	fs.once.Do(func() {
		err := os.MkdirAll(fs.dir, 0o755)
		if err != nil && !os.IsExist(err) {
			fs.err = err
			return
		}
	})

	return fs.err
}
