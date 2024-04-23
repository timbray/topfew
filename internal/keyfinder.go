package topfew

// Extract a Key from a record based on a list of keys. If the list is empty, the Key is the whole record.
//  Otherwise, there's a list of fields. They are extracted, joined with spaces, and that's the Key

// First implementation was regexp based but Golang regexps are slow.  So we'll use a hand-built state machine that
//  only cares whether each byte encodes space-or-tab or not.

import (
	"errors"
	"regexp"
)

// NER is the error message returned when the input has fewer fields than the keyFinder is configured for.
const NER = "not enough bytes in record"

// keyFinder extracts a Key based on the specified fields from a record. fields is a slice of small integers
// representing field numbers; 1-based on the command line, 0-based here.
type keyFinder struct {
	fields    []uint
	key       []byte
	separator *regexp.Regexp
}

// newKeyFinder creates a new Key finder with the supplied field numbers, the input should be 1 based.
// keyFinder is not thread-safe, you should clone it for each goroutine that uses it.
func newKeyFinder(keys []uint, separator *regexp.Regexp) *keyFinder {
	kf := keyFinder{
		key: make([]byte, 0, 128),
	}
	for _, knum := range keys {
		kf.fields = append(kf.fields, knum-1)
	}
	kf.separator = separator
	return &kf
}

// clone returns a new keyFinder with the same configuration. Each goroutine should use its own
// keyFinder instance.
func (kf *keyFinder) clone() *keyFinder {
	return &keyFinder{
		fields:    kf.fields,
		key:       make([]byte, 0, 128),
		separator: kf.separator,
	}
}

// getKey extracts a key from the supplied record. This is applied to every record,
// so efficiency matters.
func (kf *keyFinder) getKey(record []byte) ([]byte, error) {
	// if there are no Key-finders just return the record, minus any trailing newlines
	if len(kf.fields) == 0 && kf.separator == nil {
		if record[len(record)-1] == '\n' {
			record = record[0 : len(record)-1]
		}
		return record, nil
	}
	var err error
	if kf.separator == nil {
		kf.key = kf.key[:0]
		field := 0
		index := 0
		first := true

		// for each field in the Key
		for _, keyField := range kf.fields {
			// bypass fields before the one we want
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

			// attach desired field to Key
			kf.key, index, err = gather(kf.key, record, index)
			if err != nil {
				return nil, err
			}

			field++
		}
	} else {
		kf.key = kf.key[:0]
		allFields := kf.separator.Split(string(record), -1)
		for i, field := range kf.fields {
			if int(field) >= len(allFields) {
				return nil, errors.New(NER)
			}
			if i > 0 {
				kf.key = append(kf.key, ' ')
			}
			kf.key = append(kf.key, []byte(allFields[field])...)
		}
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

	// copy Key bytes
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
