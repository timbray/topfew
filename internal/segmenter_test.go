package topfew

import (
	"os"
	"strings"
	"testing"
)

func TestReadAll(t *testing.T) {
	file, err := os.Open("../test/data/small")
	if err != nil {
		t.Error("Can't open file")
		return
	}
	//noinspection ALL
	defer file.Close()

	offs, err := file.Seek(4176, 0)
	if offs != 4176 || err != nil {
		t.Error("OUCH")
	}
	s := Segment{4176, 4951, file}
	counter := NewCounter(10)
	kf := NewKeyFinder([]uint{7})
	ch := make(chan bool)
	f := Filters{nil, nil, nil}
	go readAll(&s, &f, counter, kf, ch)

	done := <-ch
	if !done {
		t.Error("Didn't get report back")
	}

	res := counter.GetTop()
	var want = map[string]bool{
		"/ongoing/picInfo.xml?o=https://old.tbray.org/ongoing/": true,
		"/ongoing/in-feed.xml":                         true,
		"/ongoing/When/202x/2020/04/29/Leaving-Amazon": true,
		"/ongoing/picInfo.xml?o=https://old.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon": true,
	}
	if len(res) != 4 {
		t.Errorf("len(res) should be 4 is %d", len(res))
	}

	for _, kc := range res {
		_, ok := want[kc.Key]
		if !ok {
			t.Error("Missing: " + kc.Key)
		} else {
			delete(want, kc.Key)
		}
		if *kc.Count != 1 {
			t.Errorf("Bogus count, should be 1: %d", kc.Count)
		}
	}
	if len(want) != 0 {
		t.Errorf("Remaining %d", len(want))
	}
}

func TestReadAllLongLine(t *testing.T) {
	counter := NewCounter(10)
	err := ReadFileInSegments("../test/data/long_lines", &Filters{}, counter, nil, 1)
	if err != nil {
		t.Fatalf("Failed to process file %v", err)
	}
	res := counter.GetTop()
	a := strings.Repeat("a", 5000)
	b := strings.Repeat("b", 8900)
	n5 := uint64(5)
	n3 := uint64(3)
	n2 := uint64(2)
	if len(res) != 3 {
		t.Errorf("Expecting 3 results, but got %d", len(res))
	} else {
		for i, exp := range []KeyCount{{"cc", &n5}, {a, &n3}, {b, &n2}} {
			if exp.Key != res[i].Key {
				t.Errorf("Unexpected key %v at index %d, expecting %v", res[i].Key, i, exp.Key)
			}
			if *exp.Count != *res[i].Count {
				t.Errorf("Unexpected count of %d at index %d, expecting %d", res[i].Count, i, exp.Count)
			}
		}
	}
}
