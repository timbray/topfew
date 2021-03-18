package topfew

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestReadSegment(t *testing.T) {

	fname := "../test/data/access-1k"
	file, err := os.Open(fname)
	if err != nil {
		t.Error("Can't open file")
		return
	}
	// noinspection ALL
	defer file.Close()

	info, _ := file.Stat()
	size := info.Size()

	// find all the line breaks
	reader := bufio.NewReader(file)
	var lineBreaks []int64
	newlineAt := int64(0)
	for true {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error("read botch: " + err.Error())
			return
		}
		newlineAt += int64(len(line))
		lineBreaks = append(lineBreaks, newlineAt)
	}
	// make GoLand quit whining
	if len(lineBreaks) == 0 {
		return
	}

	// we'll go through the file, pick 100 random points and for each one pick 100 random sizes, and verify that
	//  the segment in each case begins with the first char after the next newline after start and ends with the
	//  first character after the next newline after end
	lastStart := lineBreaks[len(lineBreaks) - 2]
	delta := size / 200
	for start := int64(0); start < lastStart; start += delta {
		end := start + (delta * (start % 3))
		if end > size {
			end = size
		}
		seg, err := newSegment(fname, start, end)
		if err != nil {
			t.Errorf("newSeg failed start %d end %d (size %d)", start, end, size)
		}
		if seg == nil {
			t.Errorf("Nil seg for start %d", start)
			return
		}
		_  = seg.file.Close()

		want := lineStartAfter(start, lineBreaks)
		if want != seg.start {
			t.Errorf("start %d, seg.start %d, wanted %d", start, seg.start, want)
		}


	}
}

func lineStartAfter(offset int64, lineBreaks []int64) int64 {
	if offset == 0 {
		return 0
	}
	var i int
	for i = 0; lineBreaks[i] < offset; i++ {
		// no-op
	}
	return lineBreaks[i]
}

func TestReadAll(t *testing.T) {
	file, err := os.Open("../test/data/access-1k")
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
	s := Segment{4176, 4951, file }
	counter := NewCounter(10)
	kf := NewKeyFinder([]uint{7})
	ch := make(chan bool)
	go s.readAll(counter, kf, ch)

	done := <-ch
	if !done {
		t.Error("Didn't get report back")
	}

	res := counter.GetTop()
	var want = map[string]bool{
		"/ongoing/picInfo.xml?o=https://old.tbray.org/ongoing/":true,
		"/ongoing/in-feed.xml":true,
		"/ongoing/When/202x/2020/04/29/Leaving-Amazon":true,
		"/ongoing/picInfo.xml?o=https://old.tbray.org/ongoing/When/202x/2020/04/29/Leaving-Amazon":true,
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
		if kc.Count != 1 {
			t.Errorf("Bogus count, should be 1: %d", kc.Count)
		}
	}
	if len(want) != 0 {
		t.Errorf("Remaining %d", len(want))
	}
}

