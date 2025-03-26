// Copyright (c) 2011 Mikkel Krautz <mikkel@krautz.dk>
// The use of this source code is goverened by a BSD-style
// license that can be found in the LICENSE-file.

// Package xar provides for reading and writing XAR archives.
package xar

import (
	"bytes"
	"compress/bzip2"
	"compress/zlib"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	ErrBadMagic      = errors.New("xar: bad magic")
	ErrBadVersion    = errors.New("xar: bad version")
	ErrBadHeaderSize = errors.New("xar: bad header size")

	ErrNoTOCChecksum        = errors.New("xar: no TOC checksum info in TOC")
	ErrChecksumUnsupported  = errors.New("xar: unsupported checksum type")
	ErrChecksumTypeMismatch = errors.New("xar: header and toc checksum type mismatch")
	ErrChecksumMismatch     = errors.New("xar: checksum mismatch")

	ErrNoCertificates             = errors.New("xar: no certificates stored in xar")
	ErrCertificateTypeMismatch    = errors.New("xar: certificate type and public key type mismatch")
	ErrCertificateTypeUnsupported = errors.New("xar: unsupported certificate type")

	ErrFileNoData              = errors.New("xar: file has no data")
	ErrFileEncodingUnsupported = errors.New("xar: unsupported file encoding")
)

const xarVersion = 1
const xarHeaderMagic = 0x78617221 // 'xar!'
const xarHeaderSize = 28

type xarHeader struct {
	magic         uint32
	size          uint16
	version       uint16
	toc_len_zlib  uint64
	toc_len_plain uint64
	checksum_kind uint32
}

const (
	xarChecksumKindNone = iota
	xarChecksumKindSHA1
	xarChecksumKindMD5
)

type FileType int

const (
	FileTypeFile FileType = iota
	FileTypeDirectory
	FileTypeSymlink
	FileTypeFifo
	FileTypeCharDevice
	FileTypeBlockDevice
	FileTypeSocket
)

type FileChecksumKind int

const (
	FileChecksumKindSHA1 FileChecksumKind = iota
	FileChecksumKindMD5
)

type FileInfo struct {
	DeviceNo uint64
	Mode     uint32
	Inode    uint64
	Uid      int
	User     string
	Gid      int
	Group    string
	Atime    int64
	Mtime    int64
	Ctime    int64
}

type FileChecksum struct {
	Kind FileChecksumKind
	Sum  []byte
}

type File struct {
	Type FileType
	Info FileInfo
	Id   uint64
	Name string

	EncodingMimetype   string
	CompressedChecksum FileChecksum
	ExtractedChecksum  FileChecksum
	// The size of the archived file (the size of the file after decompressing)
	Size int64

	offset int64
	length int64
	heap   io.ReaderAt
}

type Reader struct {
	File map[uint64]*File

	Certificates          []*x509.Certificate
	SignatureCreationTime int64
	SignatureError        error

	xar        io.ReaderAt
	root       *xmlXar
	size       int64
	heapOffset int64
}

// OpenReader will open the XAR file specified by name and return a Reader.
func OpenReader(name string) (*Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return NewReader(f, info.Size())
}

// NewReader returns a new reader reading from r, which is assumed to have the given size in bytes.
func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	xr := &Reader{
		File: make(map[uint64]*File),
		xar:  r,
		size: size,
	}

	hdr := make([]byte, xarHeaderSize)
	_, err := xr.xar.ReadAt(hdr, 0)
	if err != nil {
		return nil, err
	}

	xh := &xarHeader{}
	xh.magic = binary.BigEndian.Uint32(hdr[0:4])
	xh.size = binary.BigEndian.Uint16(hdr[4:6])
	xh.version = binary.BigEndian.Uint16(hdr[6:8])
	xh.toc_len_zlib = binary.BigEndian.Uint64(hdr[8:16])
	xh.toc_len_plain = binary.BigEndian.Uint64(hdr[16:24])
	xh.checksum_kind = binary.BigEndian.Uint32(hdr[24:28])

	if xh.magic != xarHeaderMagic {
		return nil, ErrBadMagic
	}

	if xh.version != xarVersion {
		return nil, ErrBadVersion
	}

	if xh.size != xarHeaderSize {
		return nil, ErrBadHeaderSize
	}

	ztoc := make([]byte, xh.toc_len_zlib)
	_, err = xr.xar.ReadAt(ztoc, xarHeaderSize)
	if err != nil {
		return nil, err
	}

	br := bytes.NewBuffer(ztoc)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return nil, err
	}
	// dat, err := io.ReadAll(zr)
	// if err != nil {
	// 	return nil, err
	// }
	// os.WriteFile("toc.xml", dat, 0644)

	xr.root = &xmlXar{}
	decoder := xml.NewDecoder(zr)
	decoder.Strict = false
	err = decoder.Decode(xr.root)
	if err != nil {
		return nil, err
	}

	xr.heapOffset = xarHeaderSize + int64(xh.toc_len_zlib)

	if xr.root.Toc.Checksum == nil {
		return nil, ErrNoTOCChecksum
	}

	// Check whether the XAR checksum matches
	storedsum := make([]byte, xr.root.Toc.Checksum.Size)
	_, err = io.ReadFull(io.NewSectionReader(xr.xar, xr.heapOffset+xr.root.Toc.Checksum.Offset, xr.root.Toc.Checksum.Size), storedsum)
	if err != nil {
		return nil, err
	}

	var hasher hash.Hash
	switch xh.checksum_kind {
	case xarChecksumKindNone:
		return nil, ErrChecksumUnsupported
	case xarChecksumKindSHA1:
		if xr.root.Toc.Checksum.Style != "sha1" {
			return nil, ErrChecksumTypeMismatch
		}
		hasher = sha1.New()
	case xarChecksumKindMD5:
		if xr.root.Toc.Checksum.Style != "md5" {
			return nil, ErrChecksumTypeMismatch
		}
		hasher = md5.New()
	default:
		return nil, ErrChecksumUnsupported
	}

	hasher.Write(ztoc)
	calcedsum := hasher.Sum(nil)

	if !bytes.Equal(calcedsum, storedsum) {
		return nil, ErrChecksumMismatch
	}

	// Ignore error. The method automatically sets xr.SignatureError with
	// the returned error.
	_ = xr.readAndVerifySignature(xr.root, xh.checksum_kind, calcedsum)

	// Add files to Reader
	for _, xmlFile := range xr.root.Toc.File {
		err := xr.readXmlFileTree(xmlFile, "")
		if err != nil {
			return nil, err
		}
	}

	return xr, nil
}

func (r *Reader) Subdoc() xmlSubdoc {
	return r.root.Subdoc
}

func (r *Reader) TOC() xmlToc {
	return r.root.Toc
}

// Reads signature information from the xmlXar element into
// the Reader. Also attempts to verify any signatures found.
func (r *Reader) readAndVerifySignature(root *xmlXar, checksumKind uint32, checksum []byte) (err error) {
	defer func() {
		r.SignatureError = err
	}()

	// Check if there's a signature ...
	r.SignatureCreationTime = root.Toc.SignatureCreationTime
	if root.Toc.Signature != nil {
		if len(root.Toc.Signature.Certificates) == 0 {
			return ErrNoCertificates
		}

		signature := make([]byte, root.Toc.Signature.Size)
		_, err = r.xar.ReadAt(signature, r.heapOffset+root.Toc.Signature.Offset)
		if err != nil {
			return err
		}

		// Read certificates
		for i := 0; i < len(root.Toc.Signature.Certificates); i++ {
			cb64 := []byte(strings.Replace(root.Toc.Signature.Certificates[i], "\n", "", -1))
			cder := make([]byte, base64.StdEncoding.DecodedLen(len(cb64)))
			ndec, err := base64.StdEncoding.Decode(cder, cb64)
			if err != nil {
				return err
			}

			cert, err := x509.ParseCertificate(cder[0:ndec])
			if err != nil {
				return err
			}

			r.Certificates = append(r.Certificates, cert)
		}

		// Verify validity of chain
		for i := 1; i < len(r.Certificates); i++ {
			if err := r.Certificates[i-1].CheckSignatureFrom(r.Certificates[i]); err != nil {
				return err
			}
		}

		var sighash crypto.Hash
		switch checksumKind {
		case xarChecksumKindNone:
			return ErrChecksumUnsupported
		case xarChecksumKindSHA1:
			sighash = crypto.SHA1
		case xarChecksumKindMD5:
			sighash = crypto.MD5
		}

		if root.Toc.Signature.Style == "RSA" {
			pubkey, ok := r.Certificates[0].PublicKey.(*rsa.PublicKey)
			if !ok {
				return ErrCertificateTypeMismatch
			}
			err = rsa.VerifyPKCS1v15(pubkey, sighash, checksum, signature)
			if err != nil {
				return err
			}
		} else {
			return ErrCertificateTypeUnsupported
		}
	}

	return nil
}

// This is a convenience method that returns true if the opened XAR archive
// has a signature. Internally, it checks whether the SignatureCreationTime
// field of the Reader is > 0.
func (r *Reader) HasSignature() bool {
	return r.SignatureCreationTime > 0
}

// This is a convenience method that returns true of the signature if the
// opened XAR archive was successfully verified.
//
// For a signature to be valid, it must have been signed by the leaf certificate
// in the certificate chain of the archive.
//
// If there is more than one certificate in the chain, each certificate must come
// before the one that has issued it. This is verified by checking whether the
// signature of each certificate can be verified against the public key of the
// certificate following it.
//
// The Reader does not do anything to check whether the leaf certificate and/or
// any intermediate certificates are trusted. It is up to users of this package
// to determine whether they wish to trust a given certificate chain.
// If an archive has a signature, the certificate chain of the archive can be
// accessed through the Certificates field of the Reader.
//
// Internally, this method checks whether the SignatureError field is non-nil,
// and whether the SignatureCreationTime is > 0.
//
// If the signature is not valid, and the XAR file has a signature, the
// SignatureError field of the Reader can be used to determine a possible
// cause.
func (r *Reader) ValidSignature() bool {
	return r.SignatureCreationTime > 0 && r.SignatureError == nil
}

func xmlFileToFileInfo(xmlFile *xmlFile) (fi FileInfo, err error) {
	var t time.Time
	if xmlFile.Ctime != "" {
		t, err = time.Parse(time.RFC3339, xmlFile.Ctime)
		if err != nil {
			return
		}
		fi.Ctime = t.Unix()
	}

	if xmlFile.Mtime != "" {
		t, err = time.Parse(time.RFC3339, xmlFile.Mtime)
		if err != nil {
			return
		}
		fi.Mtime = t.Unix()
	}

	if xmlFile.Atime != "" {
		t, err = time.Parse(time.RFC3339, xmlFile.Atime)
		if err != nil {
			return
		}
		fi.Atime = t.Unix()
	}

	fi.Group = xmlFile.Group
	fi.Gid = xmlFile.Gid

	fi.User = xmlFile.User
	fi.Uid = xmlFile.Uid

	fi.Mode = xmlFile.Mode

	fi.Inode = xmlFile.Inode
	fi.DeviceNo = xmlFile.DeviceNo

	return
}

// Convert a xmlFileChecksum to a FileChecksum.
func fileChecksumFromXml(f *FileChecksum, x *xmlFileChecksum) (err error) {
	f.Sum, err = hex.DecodeString(x.Digest)
	if err != nil {
		return
	}

	switch strings.ToUpper(x.Style) {
	case "MD5":
		f.Kind = FileChecksumKindMD5
	case "SHA1":
		f.Kind = FileChecksumKindSHA1
	default:
		return ErrChecksumUnsupported
	}

	return nil
}

// Create a new SectionReader that is limited to reading from the file's heap
func (r *Reader) newHeapReader() *io.SectionReader {
	return io.NewSectionReader(r.xar, r.heapOffset, r.size-r.heapOffset)
}

// Reads the file tree from a parse XAR TOC into the Reader.
func (r *Reader) readXmlFileTree(xmlFile *xmlFile, dir string) (err error) {
	xf := &File{}
	xf.heap = r.newHeapReader()

	if xmlFile.Type == "file" {
		xf.Type = FileTypeFile
	} else if xmlFile.Type == "directory" {
		xf.Type = FileTypeDirectory
	} else {
		return
	}

	xf.Id, err = strconv.ParseUint(xmlFile.Id, 10, 0)
	if err != nil {
		return
	}

	xf.Name = path.Join(dir, xmlFile.Name)

	xf.Info, err = xmlFileToFileInfo(xmlFile)
	if err != nil {
		return
	}

	if xf.Type == FileTypeFile && xmlFile.Data == nil {
		err = ErrFileNoData
		return
	}
	if xf.Type == FileTypeFile {
		xf.EncodingMimetype = xmlFile.Data.Encoding.Style
		xf.Size = xmlFile.Data.Size
		xf.length = xmlFile.Data.Length
		xf.offset = xmlFile.Data.Offset

		err = fileChecksumFromXml(&xf.CompressedChecksum, &xmlFile.Data.ArchivedChecksum)
		if err != nil {
			return
		}

		err = fileChecksumFromXml(&xf.ExtractedChecksum, &xmlFile.Data.ExtractedChecksum)
		if err != nil {
			return
		}
	}

	r.File[xf.Id] = xf

	if xf.Type == FileTypeDirectory {
		for _, subXmlFile := range xmlFile.File {
			err = r.readXmlFileTree(subXmlFile, xf.Name)
			if err != nil {
				return
			}
		}
	}

	return
}

// Open returns a ReadCloser that provides access to the file's
// uncompressed content.
func (f *File) Open() (rc io.ReadCloser, err error) {
	r := io.NewSectionReader(f.heap, f.offset, f.length)
	switch f.EncodingMimetype {
	case "application/octet-stream":
		rc = ioutil.NopCloser(r)
	case "application/x-gzip":
		rc, err = zlib.NewReader(r)
	case "application/x-bzip2":
		rc = ioutil.NopCloser(bzip2.NewReader(r))
	default:
		err = ErrFileEncodingUnsupported
	}

	return rc, err
}

// OpenRaw returns a ReadCloser that provides access to the file's
// raw content. The encoding of the raw content is specified in
// the File's EncodingMimetype field.
func (f *File) OpenRaw() (rc io.ReadCloser, err error) {
	rc = ioutil.NopCloser(io.NewSectionReader(f.heap, f.offset, f.length))
	return
}

// Verify that the compressed content of the File in the
// archive matches the stored checksum.
func (f *File) VerifyChecksum() bool {
	// Non-files are implicitly OK, since all metadata
	// is stored in the TOC.
	if f.Type != FileTypeFile {
		return true
	}

	var hasher hash.Hash
	switch f.CompressedChecksum.Kind {
	case FileChecksumKindSHA1:
		hasher = sha1.New()
	case FileChecksumKindMD5:
		hasher = md5.New()
	default:
		return false
	}

	io.Copy(hasher, io.NewSectionReader(f.heap, f.offset, f.length))
	sum := hasher.Sum(nil)
	return bytes.Equal(sum, f.CompressedChecksum.Sum)
}
