package tus

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type Metadata map[string]string

type Upload struct {
	size   int64
	offset int64

	Fingerprint string
	Metadata    Metadata
}

// Updates the Upload information based on offset.
func (u *Upload) updateProgress(offset int64) {
	u.offset = offset
}

// Returns whether this upload is finished or not.
func (u *Upload) Finished() bool {
	return u.offset >= u.size
}

// Returns the progress in a percentage.
func (u *Upload) Progress() int64 {
	return (u.offset * 100) / u.size
}

// Returns the current upload offset.
func (u *Upload) Offset() int64 {
	return u.offset
}

// Returns the size of the upload body.
func (u *Upload) Size() int64 {
	return u.size
}

// EncodedMetadata encodes the upload metadata.
func (u *Upload) EncodedMetadata() string {
	var encoded []string

	for k, v := range u.Metadata {
		encoded = append(encoded, fmt.Sprintf("%s %s", k, b64encode(v)))
	}

	return strings.Join(encoded, ",")
}

// NewUploadFromFile creates a new Upload from an os.File.
func NewUploadFromFile(f *os.File) (*Upload, error) {
	fi, err := f.Stat()

	if err != nil {
		return nil, err
	}

	metadata := map[string]string{
		"filename": fi.Name(),
	}

	fingerprint := fmt.Sprintf("%s-%d-%s", fi.Name(), fi.Size(), fi.ModTime())

	return NewUpload(fi.Size(), metadata, fingerprint), nil
}

// NewUploadFromBytes creates a new upload from a byte array.
func NewUploadFromBytes(b []byte) *Upload {
	buffer := bytes.NewReader(b)
	return NewUpload(buffer.Size(), nil, "")
}

// NewUpload creates a new upload from an io.Reader.
func NewUpload(size int64, metadata Metadata, fingerprint string) *Upload {
	if metadata == nil {
		metadata = make(Metadata)
	}

	return &Upload{
		size: size,

		Fingerprint: fingerprint,
		Metadata:    metadata,
	}
}

func b64encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
