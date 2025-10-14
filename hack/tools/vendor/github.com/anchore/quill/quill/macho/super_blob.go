package macho

import (
	"unsafe"

	"github.com/anchore/quill/internal/log"
)

// Definition From: https://github.com/Apple-FOSS-Mirror/Security/blob/5bcad85836c8bbb383f660aaf25b555a805a48e4/OSX/sec/Security/Tool/codesign.c#L53-L89

type SuperBlob struct {
	SuperBlobHeader
	Index []BlobIndex // (count) entries
	Blobs []Blob      // payload
	Pad   []byte
}

type SuperBlobHeader struct {
	Magic  Magic  // magic number
	Length uint32 // total length of SuperBlob
	Count  uint32 // number of index entries following
}

func NewSuperBlob(magic Magic) SuperBlob {
	return SuperBlob{
		SuperBlobHeader: SuperBlobHeader{
			Magic: magic,
		},
	}
}

func (s *SuperBlob) Add(t SlotType, b *Blob) {
	if b == nil {
		return
	}
	index := BlobIndex{
		Type: t,
		// Note: offset can only be set after all blobs are added
	}
	s.Index = append(s.Index, index)
	s.Blobs = append(s.Blobs, *b)
	s.Count++
	s.Length += b.Length + uint32(unsafe.Sizeof(index))
}

func (s *SuperBlob) Finalize(paddingTarget int) {
	// find the currentOffset of the first blob (header size + size of index * number of indexes)
	currentOffset := uint32(unsafe.Sizeof(s.SuperBlobHeader)) + uint32(unsafe.Sizeof(BlobIndex{}))*uint32(len(s.Index))

	// update each blob index with the currentOffset to the start of each blob relative to the start of the super blob header
	for idx := range s.Index {
		s.Index[idx].Offset = currentOffset
		currentOffset += s.Blobs[idx].Length
	}

	// add extra few pages of 0s (wanted by the codesign tool for validation)

	// why is there a padding target? There are 4 pages worth of padding we put in at the end of the superblob, and the
	// CD is sensitive to the size value of the __LINKEDIT segment. Additionally, fetching the timestamp when adding
	// a new unsigned attribute could change the superblob size. We can't recalculate since the CD hashes is finalized,
	// but we also can't claim that the superblob size is larger than it is (or else validation will fail with an
	// unexpected EOF). So, we subtract a page size to ensure that the superblob size is always smaller than what
	// is claimed on the first signing pass, even if the timestamp increases the superblob size on the second pass.

	paddingSize := PageSize * 4
	lengthWithPadding := int(s.Length) + paddingSize
	lenBeforeCorrection := s.Length

	var padCorrection int
	if paddingTarget > 0 {
		padCorrection = paddingTarget - lengthWithPadding
	}

	s.Pad = make([]byte, (PageSize*4)+padCorrection)
	s.Length += uint32(len(s.Pad))

	log.WithFields("bytes", s.Length, "correction", padCorrection, "target", paddingTarget, "bytes-before-correction", lenBeforeCorrection).Trace("superblob size")
}
