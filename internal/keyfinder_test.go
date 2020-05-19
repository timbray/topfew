package topfew

import (
	"testing"
)

func TestKeyFinder_GetKey(t *testing.T) {
	var records = []string{
		"a x c",
		"a b c",
		"a b c d e",
	}
	var kf, kf2 *KeyFinder
	kf = NewKeyFinder(nil)
	kf2 = NewKeyFinder([]uint{})
	for _, record := range records {
		r, err := kf.GetKey(record)
		if err != nil || r != record {
			t.Errorf("bad result on nil for %s", record)
		}
		r, err = kf2.GetKey(record)
		if err != nil || r != record {
			t.Errorf("bad result on empty for %s", record)
		}
	}

	kf = NewKeyFinder([]uint{1, 3})
	for _, record := range records {
		r, err := kf.GetKey(record)
		if err != nil || r != "a c" {
			t.Errorf("wanted a c from %s, got %s", record, r)
		}
	}

	kf = NewKeyFinder([]uint{1, 4})
	tooShorts := []string{"a", "a b", "a b c"}
	for _, tooShort := range tooShorts {
		_, err := kf.GetKey(tooShort)
		if err == nil {
			t.Errorf("no error on %s", tooShort)
		}
	}
	r, err := kf.GetKey("a b c d")
	if err != nil || r != "a d" {
		t.Error("border condition")
	}
}
