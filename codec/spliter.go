package codec

import (
	"io"
)

type HeadSpliter struct {
	ReadHead  func(r *Reader) int
	WriteHead func(w *Writer, l int)
}

func (s HeadSpliter) Read(r *Reader) []byte {
	n := s.ReadHead(r)
	if r.Error() != nil {
		return nil
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil
	}
	return b
}

func (s HeadSpliter) Write(w *Writer, b []byte) {
	s.WriteHead(w, len(b))
	if w.Error() != nil {
		return
	}
	w.Write(b)
}

func (s HeadSpliter) Limit(r *Reader) *io.LimitedReader {
	n := s.ReadHead(r)
	return &io.LimitedReader{r, int64(n)}
}

var (
	SplitByUint16BE = HeadSpliter{
		ReadHead:  func(r *Reader) int { return int(r.ReadUint16BE()) },
		WriteHead: func(w *Writer, l int) { w.WriteUint16BE(uint16(l)) },
	}
	SplitByUint16LE = HeadSpliter{
		ReadHead:  func(r *Reader) int { return int(r.ReadUint16LE()) },
		WriteHead: func(w *Writer, l int) { w.WriteUint16LE(uint16(l)) },
	}
)
