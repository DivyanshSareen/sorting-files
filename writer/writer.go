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

func (fw *FileWriter) Write(ch <-chan []int) {
	err := os.MkdirAll(fw.dir, 0755)
	if err != nil {
		fmt.Println(err)
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
		err := os.WriteFile(fileName, buf.Bytes(), 0644)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("saved file", fileName)
		}
	}
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
