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
	kf := NewKeyFinder([]uint{7})
	ch := make(chan segmentResult)
	f := Filters{nil, nil, nil}
	go readAll(&s, &f, kf, ch)

	segres := <-ch
	if segres.err != nil {
		t.Fatalf("got error from segment reader %v", segres.err)
	}
	counter := NewCounter(10)
	counter.merge(segres.counters)

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
	a := strings.Repeat("a", 5000)
	b := strings.Repeat("b", 8900)
	assertKeyCountsEqual(t, []*KeyCount{{"cc", pv(5)}, {a, pv(3)}, {b, pv(2)}}, counter.GetTop())
}
