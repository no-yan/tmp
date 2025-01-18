package main

type progress struct {
	n int
}

func (p *progress) Write(b []byte) (int, error) {
	p.n += len(b)
	return len(b), nil
}

func (p *progress) Show() int {
	return p.n
}
