package macho

import (
	"fmt"
	"unsafe"

	"github.com/go-restruct/restruct"
)

// Definition From: https://github.com/Apple-FOSS-Mirror/Security/blob/5bcad85836c8bbb383f660aaf25b555a805a48e4/OSX/sec/Security/Tool/codesign.c#L53-L89

type SlotType uint32

type Blob struct {
	BlobHeader
	Payload []byte
}

type BlobHeader struct {
	Magic  Magic  // magic number
	Length uint32 // total length of blob
}

func NewBlob(m Magic, p []byte) Blob {
	return Blob{
		BlobHeader: BlobHeader{
			Magic:  m,
			Length: uint32(len(p) + int(unsafe.Sizeof(Blob{}.Magic)) + int(unsafe.Sizeof(Blob{}.Length))),
		},
		Payload: p,
	}
}

func (b Blob) Pack() ([]byte, error) {
	by, err := restruct.Pack(SigningOrder, &b)
	if err != nil {
		return nil, fmt.Errorf("unable to pack blob (%x): %w", b.Magic, err)
	}
	return by, err
}
