package compress

import (
	"bytes"
	"floo/util"
	"io"
)

// Writer implements a compression writer.
type Writer struct {
	// Underlying raw, uncompressed data stream.
	rawW io.Writer

	// Buffers data into maxChunkSize chunks.
	chunkBuf *bytes.Buffer

	// Index with records which contain chunk offsets.
	index []record

	// Accumulator representing uncompressed offset.
	rawOff int64

	// Accumulator representing compressed offset.
	zipOff int64

	// Holds trailer data.
	trailer *trailer

	// Holds algorithm interface.
	algo Algorithm

	// Type of the algorithm
	algoType AlgorithmType

	// Becomes true after the first write.
	headerWritten bool
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if err := w.writeHeaderIfNeeded(); err != nil {
		return 0, nil
	}
	written := len(p)
	// compress only maxChunkSize equal chunks
	for {
		n, _ := w.chunkBuf.Write(p[:util.Min(len(p), maxChunkSize)])

		if w.chunkBuf.Len() < maxChunkSize {
			break
		}

		if err := w.flushBuffer(w.chunkBuf.Next(maxChunkSize)); err != nil {
			return 0, err
		}
		p = p[n:]
	}
	return written, nil
}

func (w *Writer) addRecordToIndex() {
	w.index = append(w.index, record{
		rawOff: w.rawOff,
		zipOff: w.zipOff,
	})
}

// compress and flush
func (w *Writer) flushBuffer(data []byte) error {
	if len(data) < 0 {
		return nil
	}
	// add record with start offset of the current chunk
	w.addRecordToIndex()

	// compress and flush current chunk
	encData, err := w.algo.Encode(data)
	if err != nil {
		return err
	}
	n, err := w.rawW.Write(encData)
	if err != nil {
		return err
	}

	// update offset for the current chunk
	// the compressed data offset is updated in the background
	// using a SizeAccumulator in combination with a MultiWriter
	w.rawOff += int64(len(data))
	w.zipOff += int64(n)
	return nil
}

func (w *Writer) writeHeaderIfNeeded() error {
	if w.headerWritten {
		return nil
	}

	if _, err := w.rawW.Write(makeHeader(w.algoType, currentVersion)); err != nil {
		return err
	}

	w.headerWritten = true
	w.zipOff += headerSize
	return nil
}

// ReadFrom implements io.ReaderFrom
func (w *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	read := 0
	buf := [maxChunkSize]byte{}

	if err := w.writeHeaderIfNeeded(); err != nil {
		return 0, err
	}

	for {
		n, rerr := r.Read(buf[:])
		read += n
		if rerr != nil && rerr != io.EOF {
			return int64(read), rerr
		}

		werr := w.flushBuffer(buf[:n])
		if werr != nil && werr != io.EOF {
			return int64(read), werr
		}
		if werr == io.EOF || rerr == io.EOF {
			return int64(read), nil
		}
	}
}

// NewWriter returns a Writer with compression support.
func NewWriter(w io.Writer, algoType AlgorithmType) (*Writer, error) {
	algo, err := AlgorithmFromType(algoType)
	if err != nil {
		return nil, err
	}
	return &Writer{
		rawW:     w,
		algo:     algo,
		algoType: algoType,
		chunkBuf: &bytes.Buffer{},
		trailer:  &trailer{},
	}, nil
}

// Close cleans up all internal resources
func (w *Writer) Close() error {
	if err := w.writeHeaderIfNeeded(); err != nil {
		return err
	}

	// write remaining bytes left in buffer and update index.
	if err := w.flushBuffer(w.chunkBuf.Bytes()); err != nil {
		return err
	}
	w.addRecordToIndex()

	// handle trailer of uncompressed file
	w.trailer.indexSize = uint64(indexChunkSize * len(w.index))
	indexBuf := make([]byte, w.trailer.indexSize)
	indexBufStartOff := indexBuf
	for _, record := range w.index {
		record.marshal(indexBuf)
		indexBuf = indexBuf[indexChunkSize:]
	}
	// write index
	if n, err := w.rawW.Write(indexBufStartOff); err != nil || uint64(n) != w.trailer.indexSize {
		return err
	}

	// write trailer buffer (chunksize, indexSize)
	// at the end of file and close the stream.
	trailerBuf := make([]byte, trailerSize)
	w.trailer.marshal(trailerBuf)

	if _, err := w.rawW.Write(trailerBuf); err != nil {
		return err
	}
	return nil
}
