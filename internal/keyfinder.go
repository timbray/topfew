package topfew

// Extract a key from a record based on a list of keys. If the list is empty, the key is the whole record.
//  Otherwise there's a list of fields. They are extracted, joined with spaces, and that's the key

import (
	"errors"
	"regexp"
)

type KeyFinder struct {
	keys []int
	sep  *regexp.Regexp
}

// make a new KeyFinder. If the keys argument is empty or nil, return the whole record. The field numbers
//  are 1-based.
func NewKeyFinder(keys []uint) *KeyFinder {
	kf := new(KeyFinder)
	if keys == nil {
		kf.keys = nil
	} else {
		for _, knum := range keys {
			kf.keys = append(kf.keys, int(knum-1)) // make it 0-based
		}
	}
	kf.sep = regexp.MustCompile("\\s+")
	return kf
}

// Get a key.
//  Possibly the record doesn't have enough fields, in which case the error will be set.
func (f *KeyFinder) GetKey(record string) (string, error) {
	if f.keys == nil || len(f.keys) == 0 {
		return record, nil
	}
	fields := f.sep.Split(record, f.keys[len(f.keys)-1]+2)
	if len(fields) <= f.keys[len(f.keys)-1] {
		return "", errors.New("not enough fields to make key")
	}
	key := fields[f.keys[0]]
	for field := 1; field < len(f.keys); field++ {
		key = key + " " + fields[f.keys[field]]
	}
	return key, nil
}
