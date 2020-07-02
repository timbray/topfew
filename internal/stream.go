package topfew

import (
	"bufio"
	"io"
)

func FromStream(ioReadr io.Reader, kf *KeyFinder, size uint) ([]*KeyCount, error) {
	counter := NewCounter(size)
	reader := bufio.NewReader(ioReadr)
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
		counter.Add(string(keyBytes))
	}
	return counter.GetTop(), nil
}
