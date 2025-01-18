package main

import "os"

const defaultURL = "https://httpbin.org/get"

func parse() string {
	if len(os.Args) < 2 {
		return defaultURL
	}
	return os.Args[1]
}
