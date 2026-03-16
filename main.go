package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sorter/chunks"
	"sorter/sorter"
	"sorter/writer"
)

type token struct {
	val int
	idx int
}

type minHeap []token

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x any) { *h = append(*h, x.(token)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	file, err := os.Open("large_sample.txt")
	if err != nil {
		panic(err)
	}

	ch := chunks.GenerateChunks(file, 64 * 1024)

	writerCh := sorter.SortChunks(ch)

	writer.WriteChunk(writerCh)

	scanners := []*bufio.Scanner{}
	scannerFiles := []*os.File{}

	files, err := os.ReadDir("./intermediate")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filename := fmt.Sprintf("./intermediate/%s", file.Name())
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(chunks.OnComma)
		scanners = append(scanners, scanner)
		scannerFiles = append(scannerFiles, f)
	}

	h := &minHeap{}
	heap.Init(h)
	for i, scanner := range scanners {
		if scanner.Scan() {
			num, _ := chunks.ByteToInt(scanner.Bytes())
			heap.Push(h, token{val: num, idx: i})
		}
	}

	f, _ := os.Create("output.txt")
	w := bufio.NewWriterSize(f, 64 * 1024)
	defer f.Close()
	defer w.Flush()
	isFirst := true
	for h.Len() > 0 {
		ele := heap.Pop(h).(token)
		if !isFirst {
			w.WriteByte(',')
		}
		fmt.Fprintf(w, "%d", ele.val)
		isFirst = false
		if scanners[ele.idx].Scan() {
			num, _ := chunks.ByteToInt(scanners[ele.idx].Bytes())
			heap.Push(h, token{val: num, idx: ele.idx})
		}
	}

	for _, f := range scannerFiles {
		f.Close()
	}
}
