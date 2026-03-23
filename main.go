package main

import (
	"fmt"
	"os"
	"sorter/chunks"
	"sorter/merger"
	"sorter/pipeline"
	"sorter/sorter"
	"sorter/writer"
)

func run(
	file *os.File,
	chunker pipeline.Chunker,
	sorter pipeline.Sorter,
	w pipeline.Writer,
	merger pipeline.Merger,
) error {
	ch, errCh := chunker.GenerateChunks(file)
	sorted := sorter.Sort(ch)
	if err := w.Write(sorted); err != nil {
		return fmt.Errorf("write stage: %w", err)
	}
	if err := <-errCh; err != nil {
		return fmt.Errorf("chunk stage: %w", err)
	}
	return merger.Merge()
}

func main() {
	file, err := os.Open("large_sample.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	err = run(
		file,
		chunks.NewFixedSizeChunker(64*1024),
		sorter.NewSliceSorter(),
		writer.NewFileWriter("./intermediate"),
		merger.NewKWayMerger("./intermediate", "output.txt"),
	)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
