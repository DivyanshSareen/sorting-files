package main

import (
	"bufio"
	"fmt"
	"os"
	"sorter/chunks"
	"sorter/sorter"
	"sorter/writer"
)

func main() {
	file, err := os.Open("sample.txt")
	if err != nil {
		panic(err)
	}

	ch := chunks.GenerateChunks(file, 64)

	writerCh := sorter.SortChunks(ch)

	writer.WriteChunk(writerCh)

	scanners := []bufio.Scanner{}

	files, err := os.ReadDir("./intermediate")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filename := fmt.Sprintf("./intermediate/%s", file)
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(chunks.OnComma)
		scanners = append(scanners, *scanner)
	}

}
