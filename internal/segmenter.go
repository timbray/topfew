package topfew

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
)

// segment represents a segment of a file. Is required to begin at the start of a line, i.e. start of file or
// after a \n.
type segment struct {
	start int64
	end   int64
	file  *os.File
}

// readFileInSegments breaks the file up into multiple segments and then reads them in parallel. counter
// will be updated with the resulting occurrence counts.
func readFileInSegments(fname string, filter *filters, counter *counter, kf *keyFinder, width int) error {
	// find file size
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	info, _ := file.Stat() // can't fail I think
	fileSize := info.Size()
	_ = file.Close()

	// if user doesn't specify segment parallelism, we ask Go how many cores it thinks the CPU has and
	//  assign one segment per CPU
	var segSize int64
	if width == 0 {
		cores := runtime.NumCPU()
		segSize = fileSize / int64(cores)
	} else {
		segSize = fileSize / int64(width)
	}

	// compute segments and put them in a slice
	var segments []*segment
	base := int64(0)
	for base < fileSize {
		// each segment starts at the beginning of a line and ends after a newline (or at EOF)
		segment, err := newSegment(fname, base, base+segSize)
		if err != nil {
			return err
		}
		segments = append(segments, segment)
		base = segment.end
	}

	// Fire 'em off, wait for them to report back
	ch := make(chan segmentResult)
	for _, segment := range segments {
		go readSegment(segment, filter, kf, ch)
	}
	for done := 0; done < len(segments); done++ {
		res := <-ch
		if res.err != nil {
			return err
		}
		counter.merge(res.segCounter)
	}
	return nil
}

// the start value is guaranteed to be at file start or after newline
func newSegment(fname string, start int64, end int64) (*segment, error) {
	// All these "err != nil" tests on basic filesystem seek operations are probably superfluous and
	// drive down the test coverage

	// Get the file ready to go
	file, _ := os.Open(fname) // can't fail, we just opened this in the parent func
	info, _ := file.Stat()    // can't fail
	fileSize := info.Size()
	reader := bufio.NewReader(file)

	var err error
	var offset int64
	if end >= fileSize {
		end = fileSize
	} else {
		// seek to near where we want the end to be, then peek forward to find a line-end
		offset, err = file.Seek(end, 0)
		if err != nil {
			return nil, err
		}
		if offset != end {
			return nil, fmt.Errorf("tried to seek to %d, went to %d", end, offset)
		}
		tillNL, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		end += int64(len(tillNL))
	}

	// now seek back to the beginning of the segment to get ready for reading
	offset, err = file.Seek(start, 0)
	if err != nil {
		return nil, err
	}
	if offset != start {
		return nil, fmt.Errorf("tried to seek to %d, went to %d", start, offset)
	}
	return &segment{start, end, file}, nil
}

type segmentResult struct {
	// one of these will be set
	err        error
	segCounter segmentCounter
}

// we've already opened the file and seeked to the right place
func readSegment(s *segment, filter *filters, kf *keyFinder, reportCh chan segmentResult) {
	// noinspection ALL
	defer s.file.Close()

	reader := bufio.NewReaderSize(s.file, 16*1024)
	current := s.start
	segCounter := newSegmentCounter()
	kf = kf.clone()
	for current < s.end {
		// ReadSlice results are only valid until the next call to Read, so we need
		// to be careful about how long we hang onto the record slice. The SegmentCounter
		// is the only thing that holds onto data from record, and it has to make a copy
		// anyway when it constructs its string Key. So this is safe.
		record, err := reader.ReadSlice('\n')
		// ReadSlice returns an error if a line doesn't fit in its buffer. We
		// deal with that by switching to ReadBytes to get the remainder of the line.
		if errors.Is(err, bufio.ErrBufferFull) {
			// Copy record because ReadBytes is going to overwrite it, and it contains
			// the start of the current line.
			linestart := append([]byte(nil), record...)
			record, err = reader.ReadBytes('\n')
			record = append(linestart, record...)
		}
		// not smart enough to figure out how to test this
		if (err != nil) && !errors.Is(err, io.EOF) {
			reportCh <- segmentResult{err: fmt.Errorf("can't read segment: %w", err)}
			return
		}
		current += int64(len(record))
		if !filter.filterRecord(record) {
			continue
		}
		keyBytes, err := kf.getKey(record)
		if err != nil {
			// bypass
			_, _ = fmt.Fprintf(os.Stderr, "Can't extract Key from %s\n", string(record))
			continue
		}
		segCounter.add(filter.filterField(keyBytes))
	}
	reportCh <- segmentResult{segCounter: segCounter}
}
