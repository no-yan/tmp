package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	fmt.Printf("tmp file: %s\n", f.Name())

	p := &progress{}
	w := io.MultiWriter(os.Stdout, f, p)

	if err := fetch(w, p); err != nil {
		log.Fatal(err)
	}
}

func fetch(w io.Writer, p *progress) error {
	t := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	defer func() {
		t.Stop()
		close(done)
	}()

	url := parse()
	res, err := request(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-t.C:
				bytesRead := p.Show()
				if res.ContentLength > 0 {
					percentage := float64(bytesRead) / float64(res.ContentLength) * 100
					fmt.Printf("%s: Downloaded %d bytes (%.2f%%)\n", t, bytesRead, percentage)
				} else {
					fmt.Printf("%s: Downloaded %d bytes\n", t, bytesRead)
				}
			}
		}
	}()

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return err
	}

	t.Stop()
	done <- true

	return nil
}

func request(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("response failed with status code : %d and\nbody: %s", res.StatusCode, body)
	}

	return res, err
}
