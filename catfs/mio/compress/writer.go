package compress

import (
	"bytes"
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

}
