package helpers

import (
	"errors"
	"log"
	"os"
)

func ReadFileChunk(path string, offset, limit int32) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("unable to open file %v", path)
		log.Println(err)
		return nil, err
	}

	stat, err := f.Stat()

	if err != nil {
		log.Println("unable to get file stat")
		return nil, err
	}

	if int32(stat.Size()) <= offset {
		log.Println("the file size is smaller than or equal to the offset")
		return nil, errors.New("the file size is smaller than or equal to the offset")
	}

	if (int32(stat.Size()) - offset) < limit {
		limit = (int32(stat.Size()) - offset)
	}
	content := make([]byte, limit)

	if _, err := f.ReadAt(content, int64(offset)); err != nil {
		log.Println("unable to read file chunck: ", err)
		return nil, err
	}

	return content, nil
}
