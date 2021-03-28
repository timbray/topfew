package topfew

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func Sample(ioReader io.Reader, filters *Filters, kf *KeyFinder) error {

	reader := bufio.NewReader(ioReader)
	for true {
		record, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if filters.FilterRecord(record) {
			fmt.Print("   ACCEPT: " + string(record))
		} else {
			fmt.Print("   REJECT: " + string(record))
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			return err
		}

		filtered := filters.FilterField(keyBytes)
		if filtered == nil {
			fmt.Printf("  REJECT: %s\n", string(filtered))
		} else if bytes.Equal(keyBytes, filtered) {
			fmt.Printf("KEY AS IS: %s\n", string(filtered))
		} else {
			fmt.Printf("   KEY IN: %s\n", string(keyBytes))
			fmt.Printf(" FILTERED: %s\n", string(filtered))
		}
	}
	return nil
}
