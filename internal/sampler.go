package topfew

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Sample prints out what amounts to a debugging feed, showing how the filtering and keyrewriting are working.
func Sample(ioReader io.Reader, filters *Filters, kf *KeyFinder) error {
	reader := bufio.NewReader(ioReader)
	for {
		record, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if filters.FilterRecord(record) {
			fmt.Print("   ACCEPT: " + string(record))
		} else {
			fmt.Print("   REJECT: " + string(record))
			continue
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			return err
		}

		filtered := filters.FilterField(keyBytes)
		if bytes.Equal(keyBytes, filtered) {
			fmt.Printf("KEY AS IS: %s\n", string(filtered))
		} else {
			fmt.Printf("   KEY IN: %s\n", string(keyBytes))
			fmt.Printf(" FILTERED: %s\n", string(filtered))
		}
	}
	return nil
}
