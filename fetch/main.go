package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Printf("tmp file: %s\n", f.Name())

	w := io.MultiWriter(os.Stdout, f)

	if err := fetch(w); err != nil {
		log.Fatal(err)
	}
}

func fetch(w io.Writer) error {
	url := parse()
	res, err := request(url)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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
