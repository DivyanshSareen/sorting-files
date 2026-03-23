package chunks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// FixedSizeChunker splits a file into fixed-size integer chunks.
type FixedSizeChunker struct {
	chunkSize int
}

func NewFixedSizeChunker(chunkSize int) *FixedSizeChunker {
	return &FixedSizeChunker{chunkSize: chunkSize}
}

func (c *FixedSizeChunker) GenerateChunks(file *os.File) (<-chan []int, <-chan error) {
	ch := make(chan []int, 3)
	errCh := make(chan error, 1)
	go c.biteChunk(file, ch, errCh)
	return ch, errCh
}

func (c *FixedSizeChunker) biteChunk(file *os.File, ch chan<- []int, errCh chan<- error) {
	defer close(ch)
	defer close(errCh)
	scanner := bufio.NewScanner(file)
	scanner.Split(OnComma)
	chunk := make([]int, 0, c.chunkSize)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		digit, err := ByteToInt(scanner.Bytes())
		if err != nil {
			errCh <- fmt.Errorf("parsing token: %w", err)
			return
		}
		chunk = append(chunk, digit)
		if len(chunk) == c.chunkSize {
			ch <- chunk
			chunk = make([]int, 0, c.chunkSize)
		}
	}
	if err := scanner.Err(); err != nil {
		errCh <- fmt.Errorf("scanning file: %w", err)
		return
	}
	if len(chunk) > 0 {
		ch <- chunk
	}
}

func ByteToInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("empty byte slice")
	}
	n := 0
	for _, c := range b {
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func OnComma(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, ','); i >= 0 {
		return i + 1, cleanNumber(data[:i]), nil
	}
	if atEOF {
		return len(data), cleanNumber(data), nil
	}
	return 0, nil, nil
}

func cleanNumber(data []byte) []byte {
	idx := 0
	for i := 0; i < len(data); i++ {
		if data[i] >= '0' && data[i] <= '9' {
			data[idx] = data[i]
			idx++
		}
	}
	return data[:idx]
}
