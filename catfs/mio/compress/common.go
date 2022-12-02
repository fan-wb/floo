package compress

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// AlgorithmType user defined type to store the algorithm type.
type AlgorithmType byte

const (
	// AlgoNone represents a `uncompressed` algorithm.
	AlgoNone = iota

	// AlgoSnappy represents the snappy compression algorithm:
	// https://en.wikipedia.org/wiki/Snappy_(software)
	AlgoSnappy

	//AlgoLZ4 represents the lz4 compression algorithm:
	// https://en.wikipedia.org/wiki/LZ4_(compression_algorithm)
	AlgoLZ4
)

// IsValid returns true if `at` is a valid algorithm type.
func (at AlgorithmType) IsValid() bool {
	switch at {
	case AlgoNone, AlgoSnappy, AlgoLZ4:
		return true
	}

	return false
}

// record structure represents an offset mapping {uncompressed offset, compressedOffset}.
// A chunk of maxChunkSize is defined by two records. The size of a specific
// record can be determinate by a simple substitution of two record offsets.
type record struct {
	rawOff int64
	zipOff int64
}

// trailer holds basic information about the compressed file.
type trailer struct {
	chunksize uint32
	indexSize uint64
}

const (
	maxChunkSize   = 64 * 1024
	indexChunkSize = 16
	trailerSize    = 12
	headerSize     = 12
	currentVersion = 1
)

var (
	// ErrBadIndex is returned on invalid compression index.
	ErrBadIndex = errors.New("broken compression index")

	// ErrHeaderTooSmall is returned if the header is less than 10 bytes.
	// It usually indicates a broken file or a non-compressed file.
	ErrHeaderTooSmall = errors.New("header is less than 10 bytes")

	// ErrBadMagicNumber is returned if the first 8 bytes of the stream is not
	// the expected "elchwald".
	ErrBadMagicNumber = errors.New("bad magic number in compressed stream")

	// ErrBadAlgorithm is returned when the algorithm was either not present
	// or it had an invalid value
	ErrBadAlgorithm = errors.New("invalid algorithm")

	// ErrUnsupportedVersion is returned when we don't have a reader that
	// understands that format.
	ErrUnsupportedVersion = errors.New("version of this format is not supported")
)

func (t *trailer) marshal(buf []byte) {
	binary.LittleEndian.PutUint32(buf[0:4], t.chunksize)
	binary.LittleEndian.PutUint64(buf[4:12], t.indexSize)
}

func (t *trailer) unmarshal(buf []byte) {
	t.chunksize = binary.LittleEndian.Uint32(buf[0:4])
	t.indexSize = binary.LittleEndian.Uint64(buf[4:12])
}

func (rc *record) marshal(buf []byte) {
	binary.LittleEndian.PutUint64(buf[0:8], uint64(rc.rawOff))
	binary.LittleEndian.PutUint64(buf[8:16], uint64(rc.zipOff))
}

func (rc *record) unmarshal(buf []byte) {
	rc.rawOff = int64(binary.LittleEndian.Uint64(buf[0:8]))
	rc.zipOff = int64(binary.LittleEndian.Uint64(buf[8:16]))
}

type header struct {
	algo    AlgorithmType
	version uint16
}

func makeHeader(algo AlgorithmType, version byte) []byte {
	algoField := make([]byte, 2)
	binary.LittleEndian.PutUint16(algoField, uint16(algo))

	versionField := make([]byte, 2)
	binary.LittleEndian.PutUint16(versionField, uint16(version))

	suffix := append(versionField, algoField...)
	return append([]byte("engorgio"), suffix...)
}

func readHeader(bHeader []byte) (*header, error) {
	if len(bHeader) < 10 {
		return nil, ErrHeaderTooSmall
	}
	if !bytes.Equal(bHeader[:8], []byte("engorgio")) {
		return nil, ErrBadMagicNumber
	}

	// this version only understands itself currently
	version := binary.LittleEndian.Uint16(bHeader[8:10])
	if version != currentVersion {
		return nil, ErrUnsupportedVersion
	}

	if len(bHeader) < 12 {
		return nil, ErrBadAlgorithm
	}

	algo := AlgorithmType(binary.LittleEndian.Uint16(bHeader[10:12]))
	if !algo.IsValid() {
		return nil, ErrBadAlgorithm
	}

	return &header{
		algo:    algo,
		version: version,
	}, nil
}
