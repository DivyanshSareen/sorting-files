package sorter

import "slices"

// SliceSorter sorts each chunk using the standard slice sort.
type SliceSorter struct{}

func NewSliceSorter() *SliceSorter {
	return &SliceSorter{}
}

func (s *SliceSorter) Sort(ch <-chan []int) <-chan []int {
	sortedChannel := make(chan []int, 1)
	go func() {
		defer close(sortedChannel)
		for chunk := range ch {
			slices.Sort(chunk)
			sortedChannel <- chunk
		}
	}()
	return sortedChannel
}
