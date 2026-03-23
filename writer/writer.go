package writer

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
)

// FileWriter writes sorted chunks to intermediate files.
type FileWriter struct {
	dir string
}

func NewFileWriter(dir string) *FileWriter {
	return &FileWriter{dir: dir}
}

func (fw *FileWriter) Write(ch <-chan []int) error {
	if err := os.MkdirAll(fw.dir, 0755); err != nil {
		return fmt.Errorf("creating intermediate dir: %w", err)
	}
	for chunk := range ch {
		buf := new(bytes.Buffer)
		for i, num := range chunk {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%d", num)
		}
		fileName := fw.randomFileName()
		if err := os.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing chunk to %s: %w", fileName, err)
		}
		fmt.Println("saved file", fileName)
	}
	return nil
}

func (fw *FileWriter) randomFileName() string {
	for {
		name := fmt.Sprintf("%s/sorted-%d.txt", fw.dir, rand.Int63())
		if _, err := os.Stat(name); os.IsNotExist(err) {
			return name
		}
		fmt.Println("Dupe file name, generating new file name")
	}
}
