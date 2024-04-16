package topfew

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// sample prints out what amounts to a debugging feed, showing how the filtering and keyrewriting are working.
func sample(ioReader io.Reader, filters *filters, kf *keyFinder) error {
	reader := bufio.NewReader(ioReader)
	for {
		record, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if filters.filterRecord(record) {
			fmt.Print("   ACCEPT: " + string(record))
		} else {
			fmt.Print("   REJECT: " + string(record))
			continue
		}
		keyBytes, err := kf.getKey(record)
		if err != nil {
			return err
		}

		filtered := filters.filterField(keyBytes)
		if bytes.Equal(keyBytes, filtered) {
			fmt.Printf("KEY AS IS: %s\n", string(filtered))
		} else {
			fmt.Printf("   KEY IN: %s\n", string(keyBytes))
			fmt.Printf(" FILTERED: %s\n", string(filtered))
		}
	}
	return nil
}
