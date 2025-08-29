package macho

import (
	"bytes"
	"hash"
	"io"
)

const (
	HashTypeNohash          HashType = 0
	HashTypeSha1            HashType = 1
	HashTypeSha256          HashType = 2
	HashTypeSha256Truncated HashType = 3
	HashTypeSha384          HashType = 4
	HashTypeSha512          HashType = 5
)

type HashType uint8

func hashChunks(hasher hash.Hash, chunkSize int, data []byte) (hashes [][]byte, err error) {
	var dataSize = len(data)
	var dataReader = bytes.NewReader(data)
	var buf = make([]byte, chunkSize)

loop:
	for idx := 0; idx < dataSize; {
		bufferLen, err := io.ReadFull(dataReader, buf)
		switch err {
		case nil, io.ErrUnexpectedEOF:
			break
		case io.EOF:
			break loop
		default:
			return nil, err
		}

		if idx+bufferLen > dataSize {
			bufferLen = dataSize - idx
		}
		idx += bufferLen

		hasher.Reset()
		hasher.Write(buf[:bufferLen])
		sum := hasher.Sum(nil)

		hashes = append(hashes, sum)
	}
	return hashes, nil
}
