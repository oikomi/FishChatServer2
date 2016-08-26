package codec

import (
	"encoding/binary"
	"io"
	"math"
)

const (
	MaxVarintLen16 = binary.MaxVarintLen16
	MaxVarintLen32 = binary.MaxVarintLen32
	MaxVarintLen64 = binary.MaxVarintLen64
)

func GetUint16LE(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func PutUint16LE(b []byte, v uint16) {
	binary.LittleEndian.PutUint16(b, v)
}

func GetUint16BE(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func PutUint16BE(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

func GetUint24LE(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func PutUint24LE(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func GetUint24BE(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func PutUint24BE(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func GetUint32LE(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func PutUint32LE(b []byte, v uint32) {
	binary.LittleEndian.PutUint32(b, v)
}

func GetUint32BE(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func PutUint32BE(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v)
}

func GetUint40LE(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32
}

func PutUint40LE(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
}

func GetUint40BE(b []byte) uint64 {
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 | uint64(b[0])<<32
}

func PutUint40BE(b []byte, v uint64) {
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
}

func GetUint48LE(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40
}

func PutUint48LE(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
}

func GetUint48BE(b []byte) uint64 {
	return uint64(b[5]) | uint64(b[4])<<8 | uint64(b[3])<<16 | uint64(b[2])<<24 | uint64(b[1])<<32 | uint64(b[0])<<40
}

func PutUint48BE(b []byte, v uint64) {
	b[0] = byte(v >> 40)
	b[1] = byte(v >> 32)
	b[2] = byte(v >> 24)
	b[3] = byte(v >> 16)
	b[4] = byte(v >> 8)
	b[5] = byte(v)
}

func GetUint56LE(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48
}

func PutUint56LE(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
}

func GetUint56BE(b []byte) uint64 {
	return uint64(b[6]) | uint64(b[5])<<8 | uint64(b[4])<<16 | uint64(b[3])<<24 | uint64(b[2])<<32 | uint64(b[1])<<40 | uint64(b[0])<<48
}

func PutUint56BE(b []byte, v uint64) {
	b[0] = byte(v >> 48)
	b[1] = byte(v >> 40)
	b[2] = byte(v >> 32)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 16)
	b[5] = byte(v >> 8)
	b[6] = byte(v)
}

func GetUint64LE(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func PutUint64LE(b []byte, v uint64) {
	binary.LittleEndian.PutUint64(b, v)
}

func GetUint64BE(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func PutUint64BE(b []byte, v uint64) {
	binary.BigEndian.PutUint64(b, v)
}

func GetFloat32BE(b []byte) float32 {
	return math.Float32frombits(GetUint32BE(b))
}

func PutFloat32BE(b []byte, v float32) {
	PutUint32BE(b, math.Float32bits(v))
}

func GetFloat32LE(b []byte) float32 {
	return math.Float32frombits(GetUint32LE(b))
}

func PutFloat32LE(b []byte, v float32) {
	PutUint32LE(b, math.Float32bits(v))
}

func GetFloat64BE(b []byte) float64 {
	return math.Float64frombits(GetUint64BE(b))
}

func PutFloat64BE(b []byte, v float64) {
	PutUint64BE(b, math.Float64bits(v))
}

func GetFloat64LE(b []byte) float64 {
	return math.Float64frombits(GetUint64LE(b))
}

func PutFloat64LE(b []byte, v float64) {
	PutUint64LE(b, math.Float64bits(v))
}

func UvarintSize(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

func VarintSize(x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return UvarintSize(ux)
}

func GetUvarint(b []byte) (uint64, int) {
	return binary.Uvarint(b)
}

func PutUvarint(b []byte, v uint64) int {
	return binary.PutUvarint(b, v)
}

func GetVarint(b []byte) (int64, int) {
	return binary.Varint(b)
}

func PutVarint(b []byte, v int64) int {
	return binary.PutVarint(b, v)
}

func ReadUvarint(r io.ByteReader) (uint64, error) {
	return binary.ReadUvarint(r)
}

func ReadVarint(r io.ByteReader) (int64, error) {
	return binary.ReadVarint(r)
}
