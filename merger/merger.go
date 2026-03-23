package merger

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sorter/chunks"
)

type token struct {
	val int
	idx int
}

type minHeap []token

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(token)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// KWayMerger merges sorted intermediate files using a k-way heap merge.
type KWayMerger struct {
	intermediateDir string
	outputFile      string
}

func NewKWayMerger(intermediateDir, outputFile string) *KWayMerger {
	return &KWayMerger{intermediateDir: intermediateDir, outputFile: outputFile}
}

func (m *KWayMerger) Merge() error {
	files, err := os.ReadDir(m.intermediateDir)
	if err != nil {
		return fmt.Errorf("reading intermediate dir: %w", err)
	}

	var scanners []*bufio.Scanner
	var openFiles []*os.File

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", m.intermediateDir, file.Name())
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening %s: %w", path, err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(chunks.OnComma)
		scanners = append(scanners, scanner)
		openFiles = append(openFiles, f)
	}

	defer func() {
		for _, f := range openFiles {
			f.Close()
		}
	}()

	return m.kWayMerge(scanners)
}

func (m *KWayMerger) kWayMerge(scanners []*bufio.Scanner) error {
	h := &minHeap{}
	heap.Init(h)
	for i, scanner := range scanners {
		if scanner.Scan() {
			num, _ := chunks.ByteToInt(scanner.Bytes())
			heap.Push(h, token{val: num, idx: i})
		}
	}

	f, err := os.Create(m.outputFile)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriterSize(f, 64*1024)
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
	return nil
}
