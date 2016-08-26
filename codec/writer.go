package codec

import (
	"io"
)

type Writer struct {
	W   io.Writer
	wb  [MaxVarintLen64]byte
	err error
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{W: w}
}

func (writer *Writer) Reset(w io.Writer) {
	writer.W = w
	writer.err = nil
}

func (writer *Writer) Error() error {
	return writer.err
}

func (writer *Writer) Write(b []byte) (n int, err error) {
	n, err = writer.W.Write(b)
	writer.err = err
	return
}

func (writer *Writer) WritePacket(b []byte, spliter HeadSpliter) {
	if writer.err != nil {
		return
	}
	spliter.Write(writer, b)
}

func (writer *Writer) WriteBytes(b []byte) {
	writer.Write(b)
}

func (writer *Writer) WriteString(s string) {
	writer.WriteBytes([]byte(s))
}

func (writer *Writer) WriteUvarint(v uint64) {
	writer.Write(writer.wb[:PutUvarint(writer.wb[:], v)])
}

func (writer *Writer) WriteVarint(v int64) {
	writer.Write(writer.wb[:PutVarint(writer.wb[:], v)])
}

func (writer *Writer) WriteUint8(v uint8) {
	writer.wb[0] = v
	writer.Write(writer.wb[:1])
}

func (writer *Writer) WriteUint16BE(v uint16) {
	PutUint16BE(writer.wb[:2], v)
	writer.Write(writer.wb[:2])
}

func (writer *Writer) WriteUint16LE(v uint16) {
	PutUint16LE(writer.wb[:2], v)
	writer.Write(writer.wb[:2])
}

func (writer *Writer) WriteUint24BE(v uint32) {
	PutUint24BE(writer.wb[:3], v)
	writer.Write(writer.wb[:3])
}

func (writer *Writer) WriteUint24LE(v uint32) {
	PutUint24LE(writer.wb[:3], v)
	writer.Write(writer.wb[:3])

}

func (writer *Writer) WriteUint32BE(v uint32) {
	PutUint32BE(writer.wb[:4], v)
	writer.Write(writer.wb[:4])
}

func (writer *Writer) WriteUint32LE(v uint32) {
	PutUint32LE(writer.wb[:4], v)
	writer.Write(writer.wb[:4])
}

func (writer *Writer) WriteUint40BE(v uint64) {
	PutUint40BE(writer.wb[:5], v)
	writer.Write(writer.wb[:5])
}

func (writer *Writer) WriteUint40LE(v uint64) {
	PutUint40LE(writer.wb[:5], v)
	writer.Write(writer.wb[:5])
}

func (writer *Writer) WriteUint48BE(v uint64) {
	PutUint48BE(writer.wb[:6], v)
	writer.Write(writer.wb[:6])
}

func (writer *Writer) WriteUint48LE(v uint64) {
	PutUint48LE(writer.wb[:6], v)
	writer.Write(writer.wb[:6])
}

func (writer *Writer) WriteUint56BE(v uint64) {
	PutUint56BE(writer.wb[:7], v)
	writer.Write(writer.wb[:7])
}

func (writer *Writer) WriteUint56LE(v uint64) {
	PutUint56LE(writer.wb[:7], v)
	writer.Write(writer.wb[:7])
}

func (writer *Writer) WriteUint64BE(v uint64) {
	PutUint64BE(writer.wb[:8], v)
	writer.Write(writer.wb[:8])
}

func (writer *Writer) WriteUint64LE(v uint64) {
	PutUint64LE(writer.wb[:8], v)
	writer.Write(writer.wb[:8])
}

func (writer *Writer) WriteFloat32BE(v float32) {
	PutFloat32BE(writer.wb[:4], v)
	writer.Write(writer.wb[:4])
}

func (writer *Writer) WriteFloat32LE(v float32) {
	PutFloat32LE(writer.wb[:4], v)
	writer.Write(writer.wb[:4])
}

func (writer *Writer) WriteFloat64BE(v float64) {
	PutFloat64BE(writer.wb[:8], v)
	writer.Write(writer.wb[:8])
}

func (writer *Writer) WriteFloat64LE(v float64) {
	PutFloat64LE(writer.wb[:8], v)
	writer.Write(writer.wb[:8])
}

func (writer *Writer) WriteInt8(v int8)     { writer.WriteUint8(uint8(v)) }
func (writer *Writer) WriteInt16BE(v int16) { writer.WriteUint16BE(uint16(v)) }
func (writer *Writer) WriteInt16LE(v int16) { writer.WriteUint16LE(uint16(v)) }
func (writer *Writer) WriteInt24BE(v int32) { writer.WriteUint24BE(uint32(v)) }
func (writer *Writer) WriteInt24LE(v int32) { writer.WriteUint24LE(uint32(v)) }
func (writer *Writer) WriteInt32BE(v int32) { writer.WriteUint32BE(uint32(v)) }
func (writer *Writer) WriteInt32LE(v int32) { writer.WriteUint32LE(uint32(v)) }
func (writer *Writer) WriteInt40BE(v int64) { writer.WriteUint40BE(uint64(v)) }
func (writer *Writer) WriteInt40LE(v int64) { writer.WriteUint40LE(uint64(v)) }
func (writer *Writer) WriteInt48BE(v int64) { writer.WriteUint48BE(uint64(v)) }
func (writer *Writer) WriteInt48LE(v int64) { writer.WriteUint48LE(uint64(v)) }
func (writer *Writer) WriteInt56BE(v int64) { writer.WriteUint56BE(uint64(v)) }
func (writer *Writer) WriteInt56LE(v int64) { writer.WriteUint56LE(uint64(v)) }
func (writer *Writer) WriteInt64BE(v int64) { writer.WriteUint64BE(uint64(v)) }
func (writer *Writer) WriteInt64LE(v int64) { writer.WriteUint64LE(uint64(v)) }
func (writer *Writer) WriteIntBE(v int)     { writer.WriteUint64BE(uint64(v)) }
func (writer *Writer) WriteIntLE(v int)     { writer.WriteUint64LE(uint64(v)) }
func (writer *Writer) WriteUintBE(v uint)   { writer.WriteUint64BE(uint64(v)) }
func (writer *Writer) WriteUintLE(v uint)   { writer.WriteUint64LE(uint64(v)) }
