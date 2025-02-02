package main

import (
	"flag"
	"time"
)

const (
	defaultOutputDir = "out"
	defaultWorkers   = 4
	defaultTimeout   = 30 * time.Second
)

type Config struct {
	outputDir string
	workers   uint
	timeout   time.Duration
}

func NewConfig(outputDir string, workers uint, timeout time.Duration) *Config {
	return &Config{
		outputDir: outputDir,
		workers:   workers,
		timeout:   timeout,
	}
}

func NewConfigFromFlags() *Config {
	outputDir := flag.String("output-dir", defaultOutputDir, "output directory")
	workers := flag.Uint("workers", defaultWorkers, "number of worker goroutines")
	timeout := flag.Duration("request-timeout", defaultTimeout, "timeout per request")

	flag.Parse()

	return NewConfig(*outputDir, *workers, *timeout)
}
