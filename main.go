package main

import (
	"os"
	"sorter/chunks"
	"sorter/merger"
	"sorter/sorter"
	"sorter/writer"
)

func main() {
	file, err := os.Open("large_sample.txt")
	if err != nil {
		panic(err)
	}

	ch := chunks.GenerateChunks(file, 64*1024)

	writerCh := sorter.SortChunks(ch)

	writer.WriteChunk(writerCh)

	merger.MergeFiles()
}
