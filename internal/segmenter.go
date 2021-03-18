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
	end     int64
	file    *os.File
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
	for base := int64(0); base < fileSize; base += segSize {
		end := base + segSize
		if end >= fileSize {
			end = fileSize
		}
		segment, err := newSegment(fname, base, end)
		if err != nil {
			return err
		}
		if segment != nil {
			segments = append(segments, segment)
		} else {
			// we hit a segment with no newlines. I think this can only happen at the end of file?
			// if it's the only segment, we have a one-segment file
			if len(segments) == 0 {
				segment, _ = newSegment(fname, 0, fileSize)
				segments = append(segments, segment)
			} else {
				// otherwise paste it on the end of the previous segment
				segments[len(segments) - 1].end = fileSize
			}
		}
	}

	// Fire 'em off, wait for them to report back
	ch := make(chan bool) // have fiddled with buffer sizes to no effect
	for _, segment := range segments {
		go segment.readAll(counter, kf, ch)
	}
	for done := 0; done < len(segments); done++ {
		ok := <- ch
		if !ok {
			return errors.New("botched return from segment")
		}
	}
	return nil
}

func newSegment(fname string, start int64, end int64) (*Segment, error) {
	// get the file ready to go
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	info, _ := file.Stat()
	fileSize := info.Size()
	reader := bufio.NewReader(file)

	// a segment must start at the character after a newline
	var startAt int64
	if start == 0 {
		startAt = 0
	} else {
		// we're going to back up one because if the byte before 'start' is \n then we've hit the jackpot
		target := start - 1
		offset, err := file.Seek(target, 0)
		if err != nil {
			return nil, err
		}
		if offset != target {
			return nil, errors.New(fmt.Sprintf("tried to seek to %d, went to %d", start, offset))
		}
		tillNL, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		startAt = target + int64(len(tillNL))

		// if we read the whole thing, that means there were no newlines - don't return a segment
		if startAt == fileSize {
			return nil, nil
		}
	}

	var endAt int64
	if end >= fileSize {
		endAt = fileSize
	} else {
		endAt = end
	}

	segment := Segment{startAt, endAt, file }
	return &segment, nil
}

func firstLineStartAfter(file *os.File, reader *bufio.Reader, start int64, max int64) (int64, error) {
	if start == 0 {
		return 0, nil
	}
	if start >= max {
		return 0, errors.New(fmt.Sprintf("seek to %d but filesize %d", start, max))
	}

	// we're going to back up one because if the byte before 'start' is \n then we've hit the jackpot
	target := start - 1
	offset, err := file.Seek(target, 0)
	if err != nil {
		return 0, err
	}

	if offset != target {
		return 0, errors.New(fmt.Sprintf("tried to seek to %d, went to %d", start, offset))
	}
	tillNL, err := reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}

	return target + int64(len(tillNL)), nil
}

const BUFFERSIZE = 65536

// we've already opened the file and seeked to the right place
func (s *Segment) readAll(counter *Counter, kf *KeyFinder, report chan bool) {

	// noinspection ALL
	defer s.file.Close()

	reader := bufio.NewReader(s.file)
	current := s.start
	var keys [][]byte
	inBuf := 0
	for current < s.end {
		record, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			// not sure what to do here
			_, _ = fmt.Fprintf(os.Stderr, "Can't read segment: %s\n", err.Error())
			report <- false
			return
		}
		keyBytes, err := kf.GetKey(record)
		if err != nil {
			// bypass
			_, _ = fmt.Fprintf(os.Stderr, "Can't extract key from %s\n", string(record))
			continue
		}
		keys = append(keys, keyBytes)
		inBuf += len(record)

		if inBuf > BUFFERSIZE {
			counter.ConcurrentAddKeys(keys)
			inBuf = 0
			keys = nil
		}

		// counter.ConcurrentAdd(keyBytes)
		current += int64(len(record))
	}
	if inBuf > 0 {
		counter.ConcurrentAddKeys(keys)
	}

	report <- true
}
