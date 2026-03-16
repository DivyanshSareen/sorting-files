package sorter

import "slices"

func SortChunks(ch <-chan []int) <-chan []int {
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
