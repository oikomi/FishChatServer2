package codec

import (
	"io"
)

var zero [MaxVarintLen64]byte

type Reader struct {
	R   io.Reader
	buf [MaxVarintLen64]byte
	err error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{R: r}
}

func (reader *Reader) Reset(r io.Reader) {
	reader.R = r
	reader.err = nil
}

func (reader *Reader) Error() error {
	return reader.err
}

func (reader *Reader) Read(b []byte) (n int, err error) {
	if reader.err == nil {
		n, err = reader.R.Read(b)
		reader.err = err
	}
	return
}

func (reader *Reader) ReadPacket(spliter HeadSpliter) (b []byte) {
	if reader.err != nil {
		return nil
	}
	b = spliter.Read(reader)
	return
}

func (reader *Reader) ReadBytes(n int) (b []byte) {
	b = make([]byte, n)
	_, reader.err = io.ReadFull(reader.R, b)
	return b
}

func (reader *Reader) ReadString(n int) string {
	return string(reader.ReadBytes(n))
}

func (reader *Reader) ReadUvarint() (v uint64) {
	if reader.err == nil {
		v, reader.err = ReadUvarint(reader)
	}
	return
}

func (reader *Reader) ReadVarint() (v int64) {
	if reader.err == nil {
		v, reader.err = ReadVarint(reader)
	}
	return
}

func (reader *Reader) seek(n int) (b []byte) {
	if reader.err == nil {
		b = reader.buf[:n]
		_, reader.err = io.ReadFull(reader.R, b)
		if reader.err == nil {
			return
		}
	}
	return zero[:n]
}

func (reader *Reader) ReadByte() (byte, error) {
	if byteReader, ok := reader.R.(io.ByteReader); ok {
		return byteReader.ReadByte()
	}
	return byte(reader.seek(1)[0]), reader.err
}

func (reader *Reader) ReadUint8() (v uint8) {
	return uint8(reader.seek(1)[0])
}

func (reader *Reader) ReadUint16BE() uint16 {
	return GetUint16BE(reader.seek(2))
}

func (reader *Reader) ReadUint16LE() uint16 {
	return GetUint16LE(reader.seek(2))
}

func (reader *Reader) ReadUint24BE() uint32 {
	return GetUint24BE(reader.seek(3))
}

func (reader *Reader) ReadUint24LE() uint32 {
	return GetUint24LE(reader.seek(3))
}

func (reader *Reader) ReadUint32BE() uint32 {
	return GetUint32BE(reader.seek(4))
}

func (reader *Reader) ReadUint32LE() uint32 {
	return GetUint32LE(reader.seek(4))
}

func (reader *Reader) ReadUint40BE() uint64 {
	return GetUint40BE(reader.seek(5))
}

func (reader *Reader) ReadUint40LE() uint64 {
	return GetUint40LE(reader.seek(5))
}

func (reader *Reader) ReadUint48BE() uint64 {
	return GetUint48BE(reader.seek(6))
}

func (reader *Reader) ReadUint48LE() uint64 {
	return GetUint48LE(reader.seek(6))
}

func (reader *Reader) ReadUint56BE() uint64 {
	return GetUint56BE(reader.seek(7))
}

func (reader *Reader) ReadUint56LE() uint64 {
	return GetUint56LE(reader.seek(7))
}

func (reader *Reader) ReadUint64BE() uint64 {
	return GetUint64BE(reader.seek(8))
}

func (reader *Reader) ReadUint64LE() uint64 {
	return GetUint64LE(reader.seek(8))
}

func (reader *Reader) ReadFloat32BE() float32 {
	return GetFloat32BE(reader.seek(4))
}

func (reader *Reader) ReadFloat32LE() float32 {
	return GetFloat32LE(reader.seek(4))
}

func (reader *Reader) ReadFloat64BE() float64 {
	return GetFloat64BE(reader.seek(8))
}

func (reader *Reader) ReadFloat64LE() float64 {
	return GetFloat64LE(reader.seek(8))
}

func (reader *Reader) ReadInt8() int8     { return int8(reader.ReadUint8()) }
func (reader *Reader) ReadInt16BE() int16 { return int16(reader.ReadUint16BE()) }
func (reader *Reader) ReadInt16LE() int16 { return int16(reader.ReadUint16LE()) }
func (reader *Reader) ReadInt24BE() int32 { return int32(reader.ReadUint24BE()) }
func (reader *Reader) ReadInt24LE() int32 { return int32(reader.ReadUint24LE()) }
func (reader *Reader) ReadInt32BE() int32 { return int32(reader.ReadUint32BE()) }
func (reader *Reader) ReadInt32LE() int32 { return int32(reader.ReadUint32LE()) }
func (reader *Reader) ReadInt40BE() int64 { return int64(reader.ReadUint40BE()) }
func (reader *Reader) ReadInt40LE() int64 { return int64(reader.ReadUint40LE()) }
func (reader *Reader) ReadInt48BE() int64 { return int64(reader.ReadUint48BE()) }
func (reader *Reader) ReadInt48LE() int64 { return int64(reader.ReadUint48LE()) }
func (reader *Reader) ReadInt56BE() int64 { return int64(reader.ReadUint56BE()) }
func (reader *Reader) ReadInt56LE() int64 { return int64(reader.ReadUint56LE()) }
func (reader *Reader) ReadInt64BE() int64 { return int64(reader.ReadUint64BE()) }
func (reader *Reader) ReadInt64LE() int64 { return int64(reader.ReadUint64LE()) }
func (reader *Reader) ReadIntBE() int     { return int(reader.ReadUint64BE()) }
func (reader *Reader) ReadIntLE() int     { return int(reader.ReadUint64LE()) }
func (reader *Reader) ReadUintBE() uint   { return uint(reader.ReadUint64BE()) }
func (reader *Reader) ReadUintLE() uint   { return uint(reader.ReadUint64LE()) }
