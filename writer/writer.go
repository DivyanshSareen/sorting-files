package writer

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
)

func WriteChunk(ch <-chan []int) {
	err := os.MkdirAll("intermediate", 0755)
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

		fileName := randomFileName("intermediate/sorted", "txt")
		err := os.WriteFile(fileName, buf.Bytes(), 0644)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("saved file", fileName)
		}
	}
}

func randomFileName(prefix, extention string) string {
	for {
		name := fmt.Sprintf("%s-%d.%s", prefix, rand.Int63(), extention)
		if _, err := os.Stat(name); os.IsNotExist(err) {
			return name
		} else {
			fmt.Println("Dupe file name, generating new file name")
		}
	}
}
