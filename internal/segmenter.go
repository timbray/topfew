package topfew

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
)

type Segment struct {
	start int64
	end   int64
	file  *os.File
}

/**
 * The idea is that we break the file up into segments, one for each core that Golang thinks the machine
 *  has, and read them in parallel
 */
func ReadFileInSegments(fname string, counter *Counter, kf *KeyFinder) error {

	// find file size
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := info.Size()
	_ = file.Close()

	cores := runtime.NumCPU()
	segSize := fileSize / int64(cores)

	var segments []*Segment
	base := int64(0)
	for base < fileSize {

		segment, err := newSegment(fname, base, base+segSize)
		if err != nil {
			return err
		}
		segments = append(segments, segment)
		base = segment.end
	}

	// Fire 'em off, wait for them to report back
	ch := make(chan bool) // have fiddled with buffer sizes to no effect
	for _, segment := range segments {
		go readAll(segment, counter, kf, ch)
	}
	for done := 0; done < len(segments); done++ {
		ok := <-ch
		if !ok {
			return errors.New("botched return from segment")
		}
	}
	return nil
}

// the start value is guaranteed to be at file start or after newline
func newSegment(fname string, start int64, end int64) (*Segment, error) {
	// get the file ready to go
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	info, _ := file.Stat()
	fileSize := info.Size()
	reader := bufio.NewReader(file)

	var offset int64
	if end >= fileSize {
		end = fileSize
	} else {

		offset, err = file.Seek(end, 0)
		if err != nil {
			return nil, err
		}
		if offset != end {
			return nil, errors.New(fmt.Sprintf("tried to seek to %d, went to %d", end, offset))
		}
		tillNL, err := reader.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return nil, err
		}
		end += int64(len(tillNL))
	}
	offset, err = file.Seek(start, 0)
	if err != nil {
		return nil, err
	}
	if offset != start {
		return nil, errors.New(fmt.Sprintf("tried to seek to %d, went to %d", start, offset))
	}
	return &Segment{start, end, file}, nil
}

//const BUFFERSIZE = 65536
const BUFFERSIZE = 131072

// we've already opened the file and seeked to the right place
func readAll(s *Segment, counter *Counter, kf *KeyFinder, report chan bool) {

	// noinspection ALL
	defer s.file.Close()

	reader := bufio.NewReader(s.file)
	current := s.start
	var keys [][]byte
	inBuf := 0
	for current < s.end {
		record, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			// not sure what to do here
			_, _ = fmt.Fprintf(os.Stderr, "Can't read segment: %s\n", err.Error())
			report <- false
			return
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			// bypass
			_, _ = fmt.Fprintf(os.Stderr, "Can't extract key from %s\n", string(record))
		} else {
			keys = append(keys, keyBytes)
			inBuf += len(record)

			if inBuf > BUFFERSIZE {
				counter.ConcurrentAddKeys(keys)
				inBuf = 0
				keys = nil
			}
		}

		current += int64(len(record))
	}
	if inBuf > 0 {
		counter.ConcurrentAddKeys(keys)
	}

	report <- true
}
