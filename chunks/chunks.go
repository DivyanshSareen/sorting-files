package chunks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func GenerateChunks(file *os.File, chunkSize int) <-chan []int {
	ch := make(chan []int, 3)
	go biteChunk(file, chunkSize, ch)
	return ch
}

func biteChunk(file *os.File, chunkSize int, ch chan<- []int) {
	defer close(ch)
	scanner := bufio.NewScanner(file)
	scanner.Split(OnComma)
	chunk := make([]int, 0, chunkSize)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}

		// digit, err := strconv.Atoi(string(scanner.Bytes()))
		digit, err := ByteToInt(scanner.Bytes())
		if err != nil {
			fmt.Println("Error converting string :(", err)
			continue
		}
		chunk = append(chunk, digit)

		if len(chunk) == chunkSize {
			ch <- chunk
			chunk = make([]int, 0)
		}
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
		return (len(data)), cleanNumber(data), nil
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
