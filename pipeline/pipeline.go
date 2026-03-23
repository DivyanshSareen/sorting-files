package pipeline

import "os"

type Chunker interface {
	GenerateChunks(file *os.File) <-chan []int
}

type Sorter interface {
	Sort(ch <-chan []int) <-chan []int
}

type Writer interface {
	Write(ch <-chan []int)
}

type Merger interface {
	Merge() error
}
