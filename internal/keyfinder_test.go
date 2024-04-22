package topfew

import (
	"bytes"
	"testing"
)

func TestFieldSeparator(t *testing.T) {
	args := []string{"-p", "tt*", "-f", "2,4"}

	c, err := Configure(args)
	if err != nil {
		t.Error("Config!")
	}

	records := []string{
		"atbttctttttdtttte",
	}
	wanted := []string{
		"b d",
	}
	kf := newKeyFinder(c.fields, c.fieldSeparator)
	for i, record := range records {
		got, err := kf.getKey([]byte(record))
		if err != nil {
			t.Error("getKey: " + err.Error())
		}
		if string(got) != wanted[i] {
			t.Errorf("wanted %s got %s", wanted[i], string(got))
		}
	}
	_, err = kf.getKey([]byte("atbtc"))
	if err == nil || err.Error() != NER {
		t.Error("bad error value")
	}
}

func TestKeyFinder(t *testing.T) {
	var records = []string{
		"a x c",
		"a b c",
		"a b c d e",
	}
	var kf, kf2 *keyFinder

	kf = newKeyFinder(nil, nil)
	kf2 = newKeyFinder([]uint{}, nil)

	for _, recordString := range records {
		record := []byte(recordString)
		r, err := kf.getKey(record)
		if (err != nil) || !bytes.Equal(r, record) {
			t.Errorf("bad result on nil for %s", record)
		}
		r, err = kf2.getKey(record)
		if (err != nil) || !bytes.Equal(r, record) {
			t.Errorf("bad result on empty for %s", record)
		}
	}

	singles := []string{"x", "b", "b"}
	kf = newKeyFinder([]uint{2}, nil)
	for i, record := range records {
		k, err := kf.getKey([]byte(record))
		if err != nil {
			t.Error("KF fail on: " + record)
		} else {
			if string(k) != singles[i] {
				t.Errorf("got '%s' wanted '%s'", string(k), singles[i])
			}
		}
	}

	kf = newKeyFinder([]uint{1, 3}, nil)
	for _, recordstring := range records {
		record := []byte(recordstring)
		r, err := kf.getKey(record)
		if err != nil || string(r) != "a c" {
			t.Errorf("wanted a c from %s, got %s", record, r)
		}
	}

	kf = newKeyFinder([]uint{1, 4}, nil)
	tooShorts := []string{"a", "a b", "a b c"}
	for _, tooShortString := range tooShorts {
		tooShort := []byte(tooShortString)
		_, err := kf.getKey(tooShort)
		if err == nil {
			t.Errorf("no error on %s", tooShort)
		}
	}
	r, err := kf.getKey([]byte("a b c d"))
	if err != nil || string(r) != "a d" {
		t.Error("border condition")
	}
}
