package topfew

import (
	"bufio"
	"errors"
	"io"
)

func TopFewFromStream(ioReader io.Reader, kf *KeyFinder, size uint) ([]*KeyCount, error) {
	scanner := bufio.NewScanner(ioReader)
	counter := NewCounter(size)
	for scanner.Scan() {
		key, err := kf.GetKey(scanner.Text())
		// for now, we'll skip records we can't make a key for. Later, a dead-letter-file would be nice
		if err != nil {
			continue
		}
		counter.Add(key)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("I/O error: " + err.Error())
	}

	return counter.GetTop(), nil
}
