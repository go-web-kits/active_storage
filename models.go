package active_storage

import (
	"io"
	"mime/multipart"
	"time"
)

type OpenedFileHeader struct {
	*multipart.FileHeader
	Content []byte
	MD5     string
}

type Service interface {
	Upload(blob *Blob, file io.Reader, checksum string, timeout ...time.Duration) error
	Download(blob Blob) ([]byte, error)
	Delete(blob *Blob) error
	URLWithHeader(blob Blob, expire ...time.Duration) (string, map[string]interface{}, error)
	URL(blob Blob, expire ...time.Duration) (string, error)
	DirectUploadInfo(blob *Blob) (interface{}, error)
	Sync(blob *Blob) error
}
