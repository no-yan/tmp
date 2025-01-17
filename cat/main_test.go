package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

const (
	shouldErr = true
	noErr     = false
)

func Test_cat(t *testing.T) {
	tests := map[string]struct {
		srcs    []string
		in      io.Reader
		want    string
		wantErr bool
	}{
		"stdin":     {[]string{}, strings.NewReader("test"), "test", noErr},
		"stdin_err": {[]string{}, iotest.ErrReader(fmt.Errorf("cannot read")), "", shouldErr},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := new(bytes.Buffer)
			gotErr := cat(tt.srcs, tt.in, b)

			if b.String() != tt.want {
				t.Fatalf("cat() failed: got: %s, want: %s", b.String(), tt.want)
			}
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("cat() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("cat() succeeded unexpectedly")
			}
		})
	}
}

func Test_copyFileToWriter(t *testing.T) {
	tests := map[string]struct {
		src     string
		w       io.Writer
		wantErr bool
	}{
		"testdata": {"testdata/in/sample.txt", new(bytes.Buffer), noErr},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := copyFileToWriter(tt.src, tt.w)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("copyFileToWriter() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("copyFileToWriter() succeeded unexpectedly")
			}
		})
	}
}
