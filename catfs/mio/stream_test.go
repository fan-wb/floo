package mio

import (
	"bytes"
	"floo/catfs/mio/compress"
	"floo/util/testutil"
	"fmt"
	"io"
	"testing"
)

var TestKey = []byte("01234567890ABCDE01234567890ABCDE")

type wrapReader struct {
	io.Reader
	io.Seeker
	io.Closer
	io.WriterTo
}

func testWriteAndRead(t *testing.T, raw []byte, algoType compress.AlgorithmType) {
	rawBuf := &bytes.Buffer{}
	if _, err := rawBuf.Write(raw); err != nil {
		t.Errorf("buf-write failed")
		return
	}

	encStream, err := NewInStream(rawBuf, TestKey, algoType)
	if err != nil {
		t.Errorf("creating enc stream failed: %v", err)
		return
	}

	encrypted := &bytes.Buffer{}
	if _, err := io.Copy(encrypted, encStream); err != nil {
		t.Errorf("reading enc stream failed: %v", err)
		return
	}

	// Fake a close method:
	br := bytes.NewReader(encrypted.Bytes())

	r := wrapReader{
		Reader:   br,
		Seeker:   br,
		WriterTo: br,
		Closer:   io.NopCloser(nil),
	}

	decStream, err := NewOutStream(r, TestKey)
	if err != nil {
		t.Errorf("creating dec stream failed: %v", err)
		return
	}

	decrypted := &bytes.Buffer{}
	if _, err = io.Copy(decrypted, decStream); err != nil {
		t.Errorf("reading decrypted data failed: %v", err)
		return
	}

	if !bytes.Equal(decrypted.Bytes(), raw) {
		t.Errorf("raw and decrypted is not equal => BUG.")
		t.Errorf("RAW:\n  %v", raw)
		t.Errorf("DEC:\n  %v", decrypted.Bytes())
		return
	}
}

func TestWriteAndRead(t *testing.T) {
	t.Parallel()

	s64k := int64(64 * 1024)
	sizes := []int64{
		0, 1, 10, s64k, s64k - 1, s64k + 1,
		s64k * 2, s64k * 1024,
	}

	for _, size := range sizes {
		regularData := testutil.CreateDummyBuf(size)
		randomData := testutil.CreateRandomDummyBuf(size, 42)

		for algo := range compress.AlgoMap {
			prefix := fmt.Sprintf("%v-size%d-", algo, size)
			t.Run(prefix+"regular", func(t *testing.T) {
				t.Parallel()
				testWriteAndRead(t, regularData, algo)
			})
			t.Run(prefix+"random", func(t *testing.T) {
				t.Parallel()
				testWriteAndRead(t, randomData, algo)
			})
		}
	}
}
