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
	tasks     Tasks
}

func NewConfig(outputDir string, workers uint, timeout time.Duration, tasks Tasks) *Config {
	return &Config{
		outputDir: outputDir,
		workers:   workers,
		timeout:   timeout,
		tasks:     tasks,
	}
}

func NewConfigFromFlags() *Config {
	outputDir := flag.String("output-dir", defaultOutputDir, "output directory")
	workers := flag.Uint("workers", defaultWorkers, "number of worker goroutines")
	timeout := flag.Duration("request-timeout", defaultTimeout, "timeout per request")

	flag.Parse()
	urls := flag.Args()
	tasks := NewTasks(urls...)

	return NewConfig(*outputDir, *workers, *timeout, tasks)
}
