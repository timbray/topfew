package topfew

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// FromStream reads a stream and hands each line to the top-occurrence counter. Currently only used on stdin.
func FromStream(ioReader io.Reader, filters *Filters, kf *KeyFinder, size int) ([]*KeyCount, error) {
	counter := NewCounter(size)
	reader := bufio.NewReader(ioReader)
	for {
		record, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		if !filters.FilterRecord(record) {
			continue
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			// bypass
			_, _ = fmt.Fprintf(os.Stderr, "Can't extract key from %s\n", string(record))
			continue
		}
		keyBytes = filters.FilterField(keyBytes)

		counter.Add(keyBytes)
	}
	return counter.GetTop(), nil
}
