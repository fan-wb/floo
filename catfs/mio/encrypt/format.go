/*
Package encrypt implement the encryption layer of floo

file format look like this:

[HEADER][[BLOCKHEADER][PAYLOAD]...]

HEADER is 36 bytes big, contains following fields:
	-	8-byte: Magic Number (to identify non-floo files quickly)
	-	4-byte: Flags (describe the stream)
	-	2-byte: Key length in bytes
	-	2-byte: Reserved
	-	4-byte: Block size (the last one may be smaller)
	-  16-byte: MAC (protect the header from forgery)

BLOCKHEADER is 8 bytes big,  contains following fields:
	-	8 byte: Nonce (derived from current block number)

PAYLOAD contains the actual encrypted data, including a MAC at the end
The size of the MAC depends on the algorithm, for poly1305 it is 16 bytes

*/

package encrypt

import (
	"crypto/hmac"
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

type Flags int32

// possible ciphers in Counter mode
const (
	// FlagEmpty is invalid
	FlagEmpty = Flags(0)

	// FlagEncryptAES256GCM indicates the stream was encrypted with AES256 in GCM mode.
	// This should be fast on modern CPUs.
	FlagEncryptAES256GCM = Flags(1) << iota

	// FlagEncryptChaCha20 indicates that the stream was encrypted with ChaCha20.
	// This can be a good choice if your CPU does not support the AES-NI instruction set.
	FlagEncryptChaCha20

	// reserve some flags for more encryption types.
	// no particular reason, just want to have enc-type flags to be in line.
	flagReserved1
	flagReserved2
	flagReserved3
	flagReserved4
	flagReserved5
	flagReserved6

	// FlagCompressedInside indicates that the encrypted data was also compressed.
	// This can be used to decide at runtime what streaming is needed.
	FlagCompressedInside
)

// Other constants:
const (
	// Size of the header mac:
	macSize = 16

	// current file format version, increment on incompatible changes.
	version = 1

	// Size of the initial header:
	headerSize = 20 + macSize

	// Default maxBlockSize if not set
	defaultMaxBlockSize = 64 * 1024

	defaultDecBufferSize = defaultMaxBlockSize
	defaultEncBufferSize = defaultMaxBlockSize + 40
)

var (
	// MagicNumber contains the first 8 byte of every floo header.
	// It is the ascii string "evanesco".
	MagicNumber = []byte{
		0x65, 0x76, 0x61, 0x6e,
		0x65, 0x73, 0x63, 0x6f,
	}
)

// GenerateHeader generates a valid header for the format file
func GenerateHeader(key []byte, maxBlockSize int64, flags Flags) []byte {
	// This is in big endian:
	header := []byte{
		// Magic number (8 Byte):
		0, 0, 0, 0, 0, 0, 0, 0,
		// Flags (4 byte):
		0, 0, 0, 0,
		// Key length (4 Byte):
		0, 0, 0, 0,
		// Block length (4 Byte):
		0, 0, 0, 0,
		// MAC Header (16 Byte):
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	// Magic number
	copy(header[:len(MagicNumber)], MagicNumber)

	// Flags
	binary.LittleEndian.PutUint32(header[8:12], uint32(flags))

	// Encode key size
	binary.LittleEndian.PutUint32(header[12:16], uint32(32))

	// Encode max block size
	binary.LittleEndian.PutUint32(header[16:20], uint32(maxBlockSize))

	// Calculate a MAC of the header. This needs to be done last.
	headerMac := hmac.New(sha3.New224, key)
	if _, err := headerMac.Write(header[:(headerSize - macSize)]); err != nil {
		return nil
	}

	// Copy the MAC to the output
	shortHeaderMac := headerMac.Sum(nil)[:macSize]
	copy(header[headerSize-macSize:headerSize], shortHeaderMac)

	return header
}
