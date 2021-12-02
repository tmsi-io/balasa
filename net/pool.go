package net

type Pool struct {
	pos int
	buf []byte
}

func (pool *Pool) Get(size int) []byte {
	b := pool.buf[pool.pos : pool.pos+size]
	pool.pos += size
	return b
}

func (pool *Pool) ReInit() {
	pool.pos = 0
}

func NewPool(maxSize int) *Pool {
	return &Pool{
		buf: make([]byte, maxSize),
	}
}
