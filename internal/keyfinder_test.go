package topfew

import (
	"bytes"
	"testing"
)

func TestKeyFinder(t *testing.T) {
	var records = []string{
		"a x c",
		"a b c",
		"a b c d e",
	}
	var kf, kf2 *KeyFinder

	kf = NewKeyFinder(nil)
	kf2 = NewKeyFinder([]uint{})

	for _, recordString := range records {
		record := []byte(recordString)
		r, err := kf.GetKey(record)
		if (err != nil) || !bytes.Equal(r, record) {
			t.Errorf("bad result on nil for %s", record)
		}
		r, err = kf2.GetKey(record)
		if (err != nil) || !bytes.Equal(r, record) {
			t.Errorf("bad result on empty for %s", record)
		}
	}

	singles := []string { "x", "b", "b" }
	kf = NewKeyFinder([]uint{2})
	for i, record := range(records) {
		k, err := kf.GetKey([]byte(record))
		if err != nil {
			t.Error("KF fail on: " + record)
		} else {
			if string(k) != singles[i] {
				t.Errorf("got '%s' wanted '%s'", string(k), singles[i])
			}
		}

	}

	kf = NewKeyFinder([]uint{1, 3})
	for _, recordstring := range records {
		record := []byte(recordstring)
		r, err := kf.GetKey(record)
		if err != nil || string(r) != "a c" {
			t.Errorf("wanted a c from %s, got %s", record, r)
		}
	}

	kf = NewKeyFinder([]uint{1, 4})
	tooShorts := []string{"a", "a b", "a b c"}
	for _, tooShortString := range tooShorts {
		tooShort := []byte(tooShortString)
		_, err := kf.GetKey(tooShort)
		if err == nil {
			t.Errorf("no error on %s", tooShort)
		}
	}
	r, err := kf.GetKey([]byte("a b c d"))
	if err != nil || string(r) != "a d" {
		t.Error("border condition")
	}
}
