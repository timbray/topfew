package topfew

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// fromStream reads a stream and hands each line to the top-occurrence counter. Currently only used on stdin.
func fromStream(ioReader io.Reader, filters *filters, kf *keyFinder, size int) ([]*keyCount, error) {
	counter := newCounter(size)
	reader := bufio.NewReader(ioReader)
	for {
		record, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		if !filters.filterRecord(record) {
			continue
		}
		keyBytes, err := kf.getKey(record)
		if err != nil {
			// bypass
			_, _ = fmt.Fprintf(os.Stderr, "Can't extract Key from %s\n", string(record))
			continue
		}
		keyBytes = filters.filterField(keyBytes)

		counter.add(keyBytes)
	}
	return counter.getTop(), nil
}
