package topfew

// Extract a key from a record based on a list of keys. If the list is empty, the key is the whole record.
//  Otherwise there's a list of fields. They are extracted, joined with spaces, and that's the key

// First implementation was regexp based but Golang regexps are slow.  So we'll use a hand-built state machine that
//  only cares whether each byte encodes space-or-tab or not.

import (
	"errors"
)

const NER = "not enough records"

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

func (kf *KeyFinder) GetKey(record []byte) (key []byte, err error) {
	if kf == nil || len(*kf) == 0 {
		return record, nil
	}

	field := 0
	index := 0
	first := true
	for _, keyField := range *kf {
		for field < int(keyField) {
			index, err = pass(record, index)
			if err != nil {
				return
			}
			field++
		}
		if first {
			first = false
		} else {
			key = append(key, ' ')
		}
		key, index, err = gather(key, record, index)

		if err != nil {
			return
		}
		field++
	}
	return
}

func gather(key []byte, record []byte, index int) ([]byte, int, error) {

	// eat leading space
	for index < len(record) && (record[index] == ' ' || record[index] == '\t') {
		index++
	}
	if index == len(record) {
		return nil, 0, errors.New(NER)
	}
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
