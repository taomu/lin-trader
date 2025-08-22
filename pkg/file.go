package common

import (
	"io"
	"log"
	"os"
)

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err2 := file.Close()
		if err2 != nil {
			log.Printf("ReadFile file.Close error: %v", err2)
		}
	}(file)
	return io.ReadAll(file)
}
