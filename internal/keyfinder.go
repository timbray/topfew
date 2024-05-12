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
// representing field numbers; 1-based on the command line, 0-based here. The key field is used to store the
// key as it is built up fromm the record's fields; it is truncated at the beginning of each call.
// The idea is to reuse the same storage for each record and minimize allocation and garbage collection. It
// does mean that the contents of the field are only valid until you call getKey again, and also that
// the keyFinder type is not thread-safe
type keyFinder struct {
	fields       []uint
	key          []byte
	separator    *regexp.Regexp
	quotedFields bool
}

// newKeyFinder creates a new Key finder with the supplied field numbers, the input should be 1 based.
// keyFinder is not thread-safe, you should clone it for each goroutine that uses it.
func newKeyFinder(keys []uint, separator *regexp.Regexp, quotedFields bool) *keyFinder {
	kf := keyFinder{
		key: make([]byte, 0, 128),
	}
	for _, knum := range keys {
		kf.fields = append(kf.fields, knum-1)
	}
	kf.separator = separator
	kf.quotedFields = quotedFields
	return &kf
}

// clone returns a new keyFinder with the same configuration. Each goroutine should use its own
// keyFinder instance.
func (kf *keyFinder) clone() *keyFinder {
	return &keyFinder{
		fields:       kf.fields,
		key:          make([]byte, 0, 128),
		separator:    kf.separator,
		quotedFields: kf.quotedFields,
	}
}

// getKey extracts a key from the supplied record. This is applied to every record,
// so efficiency matters.
func (kf *keyFinder) getKey(record []byte) ([]byte, error) {
	// chomp
	if record[len(record)-1] == '\n' {
		record = record[:len(record)-1]
	}
	// if there are no Key-finders the key is the record
	if len(kf.fields) == 0 {
		return record, nil
	}
	var err error
	kf.key = kf.key[:0]
	if kf.separator == nil {
		// no regex provided, we're doing space-separation
		if kf.quotedFields {
			// if we're doing apache httpd style access_log files, with some "-quoted fields
			field := 0
			index := 0
			first := true

			// for each field in the key
			for _, keyField := range kf.fields {
				// bypass fields before the one we want
				for field < int(keyField) {
					index, err = passQuoted(record, index)
					if err != nil {
						return nil, err
					}
					// in the special case where we might have just passed a quoted fields, we will
					// advance index past the closing quote
					if index < len(record) && record[index] == '"' {
						index++
					}
					field++
				}

				// join(' ', kf)
				if first {
					first = false
				} else {
					kf.key = append(kf.key, ' ')
				}

				kf.key, index, err = gatherQuoted(kf.key, record, index)
				if err != nil {
					return nil, err
				}
				// in the special case where we might have just passed a quoted fields, we will
				// advance index past the closing quote
				if index < len(record) && record[index] == '"' {
					index++
				}
				field++
			}
		} else {
			// basic space-separation
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
		}
	} else {
		// regex separator provided, less code but probably slower
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

// gather pulls in the bytes from a desired field, and leaves index positioned at the first white-space
// character following the field, or at the end of the record, i.e. len(record)
func gather(key []byte, record []byte, index int) ([]byte, int, error) {
	// eat leading space - if we're already at the end of the record, the loop is a no-op
	for index < len(record) && (record[index] == ' ' || record[index] == '\t') {
		index++
	}
	if index == len(record) {
		return nil, 0, errors.New(NER)
	}

	// copy Key bytes
	startAt := index
	for index < len(record) && record[index] != ' ' && record[index] != '\t' {
		index++
	}
	key = append(key, record[startAt:index]...)
	return key, index, nil
}

// same semantics as gather, but respects quoted fields that might create spaces. Leaves the index
// value pointing at the closing quote
func gatherQuoted(key []byte, record []byte, index int) ([]byte, int, error) {
	// eat leading space
	for index < len(record) && (record[index] == ' ' || record[index] == '\t') {
		index++
	}
	if index >= len(record) {
		return nil, 0, errors.New(NER)
	}

	if record[index] == '"' {
		index++
		startAt := index
		for index < len(record) && record[index] != '"' {
			index++
		}
		key = append(key, record[startAt:index]...)
		// if we hit end-of-record before the closing quote, that's an error
		if index == len(record) {
			return nil, 0, errors.New(NER)
		}
	} else {
		startAt := index
		for index < len(record) && record[index] != ' ' && record[index] != '\t' {
			index++
		}
		key = append(key, record[startAt:index]...)
	}
	return key, index, nil
}

// pass moves the index variable past any white space and a space-separated field,
// leaving index pointing at the first white-space character after the field or
// at the end of record, i.e. == len(record)
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

// same semantics as pass, but for quoted fields. Leaves the index value pointing at the
// closing "
func passQuoted(record []byte, index int) (int, error) {
	// eat leading space
	for index < len(record) && (record[index] == ' ' || record[index] == '\t') {
		index++
	}
	if index == len(record) {
		return 0, errors.New(NER)
	}
	if record[index] == '"' {
		index++
		for index < len(record) && record[index] != '"' {
			index++
		}
		// if we hit end of record before the closing quote, that's a bug
		if index >= len(record) {
			return 0, errors.New(NER)
		}
	} else {
		for index < len(record) && record[index] != ' ' && record[index] != '\t' {
			index++
		}
	}
	return index, nil
}
