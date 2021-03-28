package topfew

import (
	"bufio"
	"io"
)

func FromStream(ioReader io.Reader, filters *Filters, kf *KeyFinder, size int) ([]*KeyCount, error) {
	counter := NewCounter(size)
	reader := bufio.NewReader(ioReader)
	for true {
		record, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !filters.FilterRecord(record) {
			continue
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			return nil, err
		}

		keyBytes = filters.FilterField(keyBytes)
		if keyBytes == nil {
			continue
		}
		counter.Add(keyBytes)
	}
	return counter.GetTop(), nil
}
