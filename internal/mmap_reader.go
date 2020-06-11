package topfew

import (
	"github.com/edsrzf/mmap-go"
	"io"
	"os"
)

type MemReader struct {
	mapper mmap.MMap
	offset int64
	max    int64
}

func NewMmap(fname string) (*MemReader, error) {
	var m mmap.MMap
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	m, err = mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &MemReader{m, 0, fi.Size()}, nil
}

func (m *MemReader) Read(p []byte) (n int, err error) {
	if m.offset >= m.max {
		return 0, io.EOF
	}
	var i int
	for i = 0; i < len(p) && m.offset < m.max; i++ {
		p[i] = m.mapper[m.offset]
		m.offset++
	}
	return i, nil
}
