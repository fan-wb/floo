package mio

import (
	"bytes"
	"floo/catfs/mio/compress"
	"floo/util/testutil"
	"fmt"
	"github.com/stretchr/testify/require"
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

func TestLimitedStream(t *testing.T) {
	t.Parallel()

	testData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	r := bytes.NewReader(testData)

	// Fake a stream without all the encryption/compression fuzz.
	stream := struct {
		io.Reader
		io.Seeker
		io.Closer
		io.WriterTo
	}{
		Reader:   r,
		Seeker:   r,
		WriterTo: r,
		Closer:   io.NopCloser(r),
	}

	for idx := 0; idx <= 10; idx++ {
		// Seek back to beginning:
		_, err := stream.Seek(0, io.SeekStart)
		require.Nil(t, err)

		smallStream := LimitStream(stream, uint64(idx))
		data, err := io.ReadAll(smallStream)
		require.Nil(t, err)
		require.Equal(t, testData[:idx], data)
	}

	var err error

	// Reset and do some special torturing:
	_, err = stream.Seek(0, io.SeekStart)
	require.Nil(t, err)

	limitStream := LimitStream(stream, 5)

	n, err := limitStream.Seek(4, io.SeekStart)
	require.Nil(t, err)
	require.Equal(t, int64(4), n)

	n, err = limitStream.Seek(6, io.SeekStart)
	require.True(t, err == io.EOF)

	n, err = limitStream.Seek(-5, io.SeekEnd)
	require.Nil(t, err)
	require.Equal(t, int64(0), n)

	_, err = limitStream.Seek(-6, io.SeekEnd)
	require.True(t, err == io.EOF)

	_, err = stream.Seek(0, io.SeekStart)
	require.Nil(t, err)

	limitStream = LimitStream(stream, 5)

	buf := &bytes.Buffer{}
	n, err = limitStream.WriteTo(buf)
	require.Nil(t, err)
	require.Equal(t, n, int64(10))
	require.Equal(t, buf.Bytes(), testData[:5])

	buf.Reset()
	_, err = stream.Seek(0, io.SeekStart)
	require.Nil(t, err)
	limitStream = LimitStream(stream, 11)

	n, err = limitStream.WriteTo(buf)
	require.Nil(t, err)
	require.Equal(t, n, int64(10))
	require.Equal(t, buf.Bytes(), testData)
}
