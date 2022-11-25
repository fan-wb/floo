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
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"errors"
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

// HeaderInfo represents a parsed header.
type HeaderInfo struct {
	// Version of the file format. Currently, it's always 1.
	Version uint16

	// Cipher type used in the file.
	CipherBit Flags

	// KeyLen is the number of bytes in the encryption key.
	KeyLen uint32

	// BlockLen is the max. number of bytes in a block.
	// The last block may be smaller.
	BlockLen uint32

	// Flags control the encryption algorithm and other things.
	Flags Flags
}

var (
	// ErrSmallHeader is returned when the header is too small to parse.
	// Usually happens when trying to decrypt a raw stream.
	ErrSmallHeader = errors.New("header is too small")

	// ErrBadMagic is returned when the stream does not start with the magic number.
	// Usually happens when trying to decrypt a raw or compressed stream.
	ErrBadMagic = errors.New("magic number missing")

	// ErrBadFlags means that you passed an invalid flags combination
	// or the stream was modified to have wrong flags.
	ErrBadFlags = errors.New("inconsistent header flags")

	// ErrBadHeaderMAC means that the header is not what the writer originally
	// put into the stream. Usually means somebody or something changed it.
	ErrBadHeaderMAC = errors.New("header mac differs from expected")
)

// ParseHeader parses the header of the format file
// returns the flags, key and block length.
func ParseHeader(header, key []byte) (*HeaderInfo, error) {
	if len(header) < len(MagicNumber) {
		return nil, ErrSmallHeader
	}

	if bytes.Compare(header[:len(MagicNumber)], MagicNumber) != 0 {
		return nil, ErrBadMagic
	}

	if len(header) < headerSize {
		return nil, ErrSmallHeader
	}

	flags := Flags(binary.LittleEndian.Uint32(header[8:12]))
	keyLen := binary.LittleEndian.Uint32(header[12:16])
	blockLen := binary.LittleEndian.Uint32(header[16:20])

	cipherBit, err := cipherTypeBitFromFlags(flags)
	if err != nil {
		return nil, err
	}

	// check the header MAC
	headerMac := hmac.New(sha3.New224, key)
	if _, err := headerMac.Write(header[:headerSize-macSize]); err != nil {
		return nil, err
	}

	shortHeaderMac := headerMac.Sum(nil)[:macSize]
	storedMac := header[headerSize-macSize : headerSize]
	if !hmac.Equal(shortHeaderMac, storedMac) {
		return nil, ErrBadHeaderMAC
	}

	return &HeaderInfo{
		Version:   version,
		CipherBit: cipherBit,
		KeyLen:    keyLen,
		BlockLen:  blockLen,
		Flags:     flags,
	}, nil
}

func cipherTypeBitFromFlags(flags Flags) (Flags, error) {
	var cipherBit Flags
	var bits = []Flags{
		FlagEncryptAES256GCM,
		FlagEncryptChaCha20,
	}

	for _, bit := range bits {
		if flags&bit == 0 {
			continue
		}

		if cipherBit != 0 {
			// only one bit at the same time allowed.
			return 0, ErrBadFlags
		}

		cipherBit = bit
	}

	if cipherBit == 0 {
		// no algorithm set: also error out.
		return 0, ErrBadFlags
	}

	return cipherBit, nil
}