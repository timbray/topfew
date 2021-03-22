package topfew

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
)

type FileAccess struct {
	name    string
	size    int64
	file    *os.File
	reader  *bufio.Reader
	offset  int64
	mapper  mmap.MMap
	bufsize int
}

func NewFileAccessViaFilesystem(filename string) (*FileAccess, error) {
	var fa FileAccess

	fa.name = filename
	var err error
	fa.file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	info, _ := fa.file.Stat()
	fa.size = info.Size()
	fa.reader = bufio.NewReader(fa.file)
	return &fa, nil
}

func NewFileAccessViaMmap(filename string) (*FileAccess, error) {
	var fa FileAccess

	fa.bufsize = 80
	fa.name = filename
	var err error
	fa.file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	info, _ := fa.file.Stat()
	fa.size = info.Size()
	fa.mapper, err = mmap.Map(fa.file, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return &fa, nil
}

func (f *FileAccess) SetOffset(target int64) error {
	if target >= f.size {
		return errors.New(fmt.Sprintf("tried to seek to %d but file size %d", target, f.size))
	}
	if f.mapper != nil {
		f.offset = target
	} else {
		reached, err := f.file.Seek(target, 0)
		if err != nil {
			return err
		}
		if reached != target {
			return errors.New(fmt.Sprintf("tried to seek to %d, reached %d", target, reached))
		}
	}
	return nil
}

func (f *FileAccess) ReadLine() ([]byte, error) {
	if f.mapper != nil {
		buf := make([]byte, 0, f.bufsize)
		bufIndex := 0
		for f.offset < f.size {
			c := f.mapper[f.offset]
			f.offset++
			if bufIndex == f.bufsize {
				buf = expand(buf, f)
			}
			buf[bufIndex] = c
			bufIndex++
			if c == '\n' {
				break
			}
		}
		return buf[0:bufIndex], nil
	} else {
		bytes, err := f.reader.ReadBytes('\n')
		return bytes, err
	}
}

func expand(buf []byte, fa *FileAccess) []byte {
	newBuf := make([]byte, 0, fa.bufsize*2)
	for i := 0; i < fa.bufsize; i++ {
		newBuf[i] = buf[i]
	}
	fa.bufsize *= 2
	return newBuf
}
