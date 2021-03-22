package topfew

import (
	"bufio"
	"io"
)

func FromStream(ioReader io.Reader, kf *KeyFinder, size uint) ([]*KeyCount, error) {
	counter := NewCounter(size)
	reader := bufio.NewReader(ioReader)
	for true {
		record, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			return nil, err
		}
		counter.Add(keyBytes)
	}
	return counter.GetTop(), nil
}
