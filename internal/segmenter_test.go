package topfew

import (
	"fmt"
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
	s := segment{4176, 4951, file}
	kf := newKeyFinder([]uint{7})
	ch := make(chan segmentResult)
	f := filters{nil, nil, nil}
	go readSegment(&s, &f, kf, ch)

	segres := <-ch
	if segres.err != nil {
		t.Fatalf("got error from segment reader %v", segres.err)
	}
	counter := newCounter(10)
	counter.merge(segres.segCounter)

	res := counter.getTop()
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

func TestReadSegmentFiltering(t *testing.T) {
	input := "foo\nbar 2\n"
	args := []string{"-f", "2", "-g", "foo"}
	c, err := Configure(args)
	if err != nil {
		t.Error("config!")
	}

	tmpName := fmt.Sprintf("/tmp/tf-%d", os.Getpid())
	tmpfile, err := os.Create(tmpName)
	if err != nil {
		t.Fatal("can't make tmpfile: " + err.Error())
	}
	defer func() { _ = os.Remove(tmpName) }()
	_, _ = fmt.Fprint(tmpfile, input)
	_ = tmpfile.Close()
	counter := newCounter(10)
	err = readFileInSegments(tmpName, &c.filter, counter, newKeyFinder(c.fields), 1)
	if err != nil {
		t.Error("Run? " + err.Error())
	}
	if len(counter.counts) != 0 {
		t.Error("nothing should have matched")
	}
}

// TestVeryLongLines is specifically aimed at the readSegment() code that deals with the
// ErrBufferFull condition, had to create lines 80k long to execute that, so rather than clutter
// up the filesystem with this junk, we create them synthetically
func TestVeryLongLines(t *testing.T) {
	tmpName := fmt.Sprintf("/tmp/tf-%d", os.Getpid())
	tmpfile, err := os.Create(tmpName)
	if err != nil {
		t.Fatal("can't make tmpfile: " + err.Error())
	}
	defer func() { _ = os.Remove(tmpName) }()
	a80k := strings.Repeat("a", 80*1000)
	c3 := "ccc"
	b30k := strings.Repeat("b", 30*1000)
	for i := 0; i < 2; i++ {
		_, _ = fmt.Fprintln(tmpfile, b30k)
	}
	for i := 0; i < 5; i++ {
		_, _ = fmt.Fprintln(tmpfile, a80k)
	}
	for i := 0; i < 3; i++ {
		_, _ = fmt.Fprintln(tmpfile, c3)
	}
	_ = tmpfile.Close()
	counter := newCounter(10)
	err = readFileInSegments(tmpName, &filters{}, counter, newKeyFinder(nil), 1)
	if err != nil {
		t.Fatal("Failed to read long-lines file")
	}
	assertKeyCountsEqual(t,
		[]*keyCount{
			{a80k, pv(5)},
			{c3, pv(3)},
			{b30k, pv(2)}},
		counter.getTop())
}
