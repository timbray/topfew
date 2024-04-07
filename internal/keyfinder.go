package topfew

// Extract a key from a record based on a list of keys. If the list is empty, the key is the whole record.
//  Otherwise, there's a list of Fields. They are extracted, joined with spaces, and that's the key

// First implementation was regexp based but Golang regexps are slow.  So we'll use a hand-built state machine that
//  only cares whether each byte encodes space-or-tab or not.

import (
	"errors"
)

// NER is the error message returned when the input has fewer Fields than the KeyFinder is configured for.
const NER = "not enough bytes in record"

// KeyFinder extracts a key based on the specified Fields from a record. Fields is a slice of small integers
// representing field numbers; 1-based on the command line, 0-based here.
type KeyFinder struct {
	fields []uint
	key    []byte
}

// NewKeyFinder creates a new key finder with the supplied field numbers, the input should be 1 based.
// KeyFinder is not thread-safe, you should Clone it for each goroutine that uses it.
func NewKeyFinder(keys []uint) *KeyFinder {
	kf := KeyFinder{
		key: make([]byte, 0, 128),
	}
	for _, knum := range keys {
		kf.fields = append(kf.fields, knum-1)
	}
	return &kf
}

// Clone returns a new KeyFinder with the same configuration. Each goroutine should use its own
// KeyFinder instance.
func (kf *KeyFinder) Clone() *KeyFinder {
	return &KeyFinder{
		fields: kf.fields,
		key:    make([]byte, 0, 128),
	}
}

// GetKey extracts a key from the supplied record. This is applied to every record,
// so efficiency matters.
func (kf *KeyFinder) GetKey(record []byte) ([]byte, error) {
	// if there are no key-finders just return the record, minus any trailing newlines
	if len(kf.fields) == 0 {
		if record[len(record)-1] == '\n' {
			record = record[0 : len(record)-1]
		}
		return record, nil
	}

	var err error
	kf.key = kf.key[:0]
	field := 0
	index := 0
	first := true

	// for each field in the key
	for _, keyField := range kf.fields {
		// bypass Fields before the one we want
		for field < int(keyField) {
			index, err = pass(record, index)
			if err != nil {
				return nil, err
			}
			field++
		}

		// join(' ', kf)
		if first {
			first = false
		} else {
			kf.key = append(kf.key, ' ')
		}

		// attach desired field to key
		kf.key, index, err = gather(kf.key, record, index)
		if err != nil {
			return nil, err
		}

		field++
	}
	return kf.key, err
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
