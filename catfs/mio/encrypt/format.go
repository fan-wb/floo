/*
Package encrypt implement the encryption layer of floo

file format look like this:

[HEADER][[BLOCKHEADER][PAYLOAD]...]

HEADER is 36 bytes big, contains following fields:
	-	8-byte: Magic Number (to identify non-floo files quickly)
	-	2-Byte: Format version
    -   2-Byte: Used cipher type (ChaCha20 or AES-GCM currently)
	-	4-byte: Key length in bytes
	-	4-byte: Block size (the last one may be smaller)
	-  16-byte: MAC (protect the header from forgery)

BLOCKHEADER is 8 bytes big,  contains following fields:
	-	8 byte: Nonce (derived from current block number)

PAYLOAD contains the actual encrypted data, including a MAC at the end
The size of the MAC depends on the algorithm, for poly1305 it is 16 bytes

All header metadata is encoded in little endian.

Reader/Writer are capable or reading/writing this format.  Additionally,
Reader supports efficient seeking into the encrypted data, provided the
underlying datastream supports seeking.  SEEK_END is only supported when the
number of encrypted blocks is present in the header.
*/

package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/sha3"
)

const (
	aeadCipherChaCha = iota
	aeadCipherAES
)

// Other constants:
const (
	// Size of the header mac:
	macSize = 16

	// current file format version, increment on incompatible changes.
	version = 1

	// Size of the initial header:
	headerSize = 20 + macSize

	// Chacha20 appears to be twice as fast as AES-GCM on my machine
	defaultCipherType = aeadCipherAES

	// Default maxBlockSize if not set
	defaultMaxBlockSize = 64 * 1024

	defaultDecBufferSize = defaultMaxBlockSize
	defaultEncBufferSize = defaultMaxBlockSize + 40
)

var (
	// MagicNumber contains the first 8 byte of every floo header.
	// It is the ASCII string "evanesco" -- the Vanishing Spell.
	MagicNumber = []byte{
		0x65, 0x76, 0x61, 0x6e,
		0x65, 0x73, 0x63, 0x6f,
	}
)

// KeySize of the used cipher's key in bytes.
var KeySize = chacha20poly1305.KeySize

// GenerateHeader generates a valid header for the format file
func GenerateHeader(key []byte, maxBlockSize int64, cipher uint16) []byte {
	// This is in big endian:
	header := []byte{
		// floo's magic number (8 Byte):
		0, 0, 0, 0, 0, 0, 0, 0,
		// File format version (2 Byte):
		0, 0,
		// Cipher type (2 Byte):
		0, 0,
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

	binary.LittleEndian.PutUint16(header[8:10], version)

	binary.LittleEndian.PutUint16(header[10:12], cipher)
	// Encode key size:
	binary.LittleEndian.PutUint32(header[12:16], uint32(KeySize))
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
	Cipher uint16

	// Keylen is the number of bytes in the encryption key.
	Keylen uint32

	// Blocklen is the max number of bytes in a block.
	// The last block may be smaller.
	Blocklen uint32
}

// ParseHeader parses the header of the format file
// returns the flags, key and block length.
func ParseHeader(header, key []byte) (*HeaderInfo, error) {
	if bytes.Compare(header[:len(MagicNumber)], MagicNumber) != 0 {
		return nil, fmt.Errorf("magic number in header differs")
	}

	version := binary.LittleEndian.Uint16(header[8:10])
	cipherType := binary.LittleEndian.Uint16(header[10:12])
	switch cipherType {
	case aeadCipherAES:
	case aeadCipherChaCha:
		// we support this!
	default:
		return nil, fmt.Errorf("unknown cipher type: %d", cipherType)
	}
	keylen := binary.LittleEndian.Uint32(header[12:16])
	blocklen := binary.LittleEndian.Uint32(header[16:20])

	// check the header MAC
	headerMac := hmac.New(sha3.New224, key)
	if _, err := headerMac.Write(header[:headerSize-macSize]); err != nil {
		return nil, err
	}

	shortHeaderMac := headerMac.Sum(nil)[:macSize]
	storedMac := header[headerSize-macSize : headerSize]
	if !hmac.Equal(shortHeaderMac, storedMac) {
		return nil, fmt.Errorf("header MAC differs from expected")
	}

	return &HeaderInfo{
		Version:  version,
		Cipher:   cipherType,
		Keylen:   keylen,
		Blocklen: blocklen,
	}, nil
}

/*
	Common Utilities
*/

type aeadCommon struct {
	// Nonce that form the first aead.NonceSize() bytes of the output
	nonce []byte

	// Key used for encryption/decryption
	key []byte

	// For more information, see:
	// https://en.wikipedia.org/wiki/Authenticated_encryption
	aead cipher.AEAD

	// Buffer for encrypted data (maxBlockSize + overhead)
	encBuf []byte
}

func createAEADWorker(cipherType uint16, key []byte) (cipher.AEAD, error) {
	switch cipherType {
	case aeadCipherAES:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		return cipher.NewGCM(block)
	case aeadCipherChaCha:
		return chacha20poly1305.New(key)
	}

	return nil, fmt.Errorf("no such cipher type: %d", cipherType)
}
