package topfew

// Extract a key from a record based on a list of keys. If the list is empty, the key is the whole record.
//  Otherwise there's a list of fields. They are extracted, joined with spaces, and that's the key

// First implementation was regexp based but Golang regexps are slow.  So we'll use a hand-built state machine that
//  only cares whether each byte encodes space-or-tab or not.

import (
	"errors"
)

const NER = "not enough bytes in record"

// just a slice of small integers representing field numbers; 1-based on the command line, 0-based here
type KeyFinder []uint

func NewKeyFinder(keys []uint) *KeyFinder {
	if keys == nil {
		return nil
	}

	var kf KeyFinder
	for _, knum := range keys {
		kf = append(kf, knum-1)
	}
	return &kf
}

// This is applied to every record, so efficiency matters
func (kf *KeyFinder) GetKey(record []byte) ([]byte, error) {
	var err error
	key := make([]byte, 0, 100)

	// if there are no keyfinders just return the record, minus any trailing newlines
	if kf == nil || len(*kf) == 0 {
		if record[len(record)-1] == '\n' {
			record = record[0 : len(record)-1]
		}
		return record, nil
	}

	field := 0
	index := 0
	first := true

	// for each field in the key
	for _, keyField := range *kf {

		// bypass fields before the one we want
		for field < int(keyField) {
			index, err = pass(record, index)
			if err != nil {
				return nil, err
			}
			field++
		}

		// join(' ', keyfields)
		if first {
			first = false
		} else {
			key = append(key, ' ')
		}

		// attach desired field to key
		key, index, err = gather(key, record, index)
		if err != nil {
			return nil, err
		}

		field++
	}
	return key, err
}

// pull in the bytes from a desired field
func gather(key []byte, record []byte, index int) ([]byte, int, error) {

	// eat leading space
	for index < len(record) && (record[index] == ' ' || record[index] == '\t') {
		index++
	}
	if index == len(record) {
		return nil, 0, errors.New(NER)
	}

	// copy key bytes
	for index < len(record) && record[index] != ' ' && record[index] != '\t' && record[index] != '\n' {
		key = append(key, record[index])
		index++
	}
	return key, index, nil
}

func pass(record []byte, index int) (int, error) {
	// eat leading space
	for index < len(record) && (record[index] == ' ' || record[index] == '\t') {
		index++
	}
	if index == len(record) {
		return 0, errors.New(NER)
	}
	for index < len(record) && record[index] != ' ' && record[index] != '\t' {
		index++
	}
	return index, nil
}
