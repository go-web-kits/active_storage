package active_storage

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/go-web-kits/dbx"
)

type ActiveStorageBlob struct {
	ID       uint   `json:"id" db:"id" gorm:"primary_key; index"`
	Key      string `json:"key" db:"key" gorm:"not null; unique_index:index_active_storage_blobs_on_key"`
	Filename string `json:"filename" db:"filename" gorm:"not null"`
	// ContentType  string    `json:"content_type" db:"content_type"` TODO
	// Metadata     string    `json:"metadata" db:"metadata"` TODO
	ByteSize     uint      `json:"byte_size" db:"byte_size" gorm:"not null"`
	Checksum     string    `json:"checksum" db:"checksum" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UploadStatus string    `json:"upload_status" db:"upload_status"`
	CnSyncStatus string    `json:"cn_sync_status" db:"cn_sync_status"`
	UsSyncStatus string    `json:"us_sync_status" db:"us_sync_status"`

	ActiveStorageAttachments []Attachment `json:"-" gorm:"ForeignKey:BlobID"`
}

type Blob = ActiveStorageBlob

func FindBlobBy(condition interface{}, preload bool) dbx.Result {
	opt := dbx.Opt{}
	if preload {
		opt.Preload = "ActiveStorageAttachments"
	}
	return dbx.Find(&Blob{}, condition, opt)
}

func CreateBlobAfterUpload(fileHeader *multipart.FileHeader, opts ...dbx.Opt) dbx.Result {
	file, _ := fileHeader.Open()
	bs, _ := ioutil.ReadAll(file)
	m := md5.Sum(bs)
	blob := Blob{
		Key:      generateBase36Key(),
		Filename: fileHeader.Filename,
		ByteSize: uint(fileHeader.Size),
		Checksum: base64.StdEncoding.EncodeToString(m[:]),
	}

	if err := blob.Upload(bytes.NewReader(bs)); err != nil {
		return dbx.Result{Data: blob, Err: err, Tx: dbx.Conn(opts...)}
	}
	return dbx.Create(&blob, opts...)
}

func CreateBlobAfterUploadByOpened(file *OpenedFileHeader, opts ...dbx.Opt) dbx.Result {
	_md5, _ := hex.DecodeString(file.MD5)
	blob := Blob{
		Key:      generateBase36Key(),
		Filename: file.Filename,
		ByteSize: uint(file.Size),
		Checksum: base64.StdEncoding.EncodeToString(_md5),
	}

	if err := blob.Upload(bytes.NewReader(file.Content)); err != nil {
		return dbx.Result{Data: blob, Err: err, Tx: dbx.Conn(opts...)}
	}
	return dbx.Create(&blob, opts...)
}

func CreateBlob(filename string, byteSize uint, md5 string, opts ...dbx.Opt) dbx.Result {
	_md5, _ := hex.DecodeString(md5)
	return dbx.Create(&Blob{
		Key:      generateBase36Key(),
		Filename: filename,
		ByteSize: byteSize,
		Checksum: base64.StdEncoding.EncodeToString(_md5)}, opts...)
}

// ==============
// Blob's Methods
// ==============

func (blob Blob) SignedId() uint {
	return blob.ID
}

func (blob Blob) Size() float64 {
	size, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(blob.ByteSize)/1024.0/1024.0), 64)
	return size
}

func (blob Blob) MD5() interface{} {
	if blob.Checksum == "" {
		return nil
	}
	bs, _ := base64.StdEncoding.DecodeString(blob.Checksum)
	return hex.EncodeToString(bs)
}

func (blob *Blob) DirectUploadInfo() (interface{}, error) {
	return Config.Service.DirectUploadInfo(blob)
}

func (blob *Blob) Upload(file io.Reader, timeout ...time.Duration) error {
	t := Config.UploadTimeout
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return Config.Service.Upload(blob, file, blob.Checksum, t)
}
func (blob Blob) URLWithHeader() (string, map[string]interface{}, error) {
	return Config.Service.URLWithHeader(blob, Config.URLExpire)
}

func (blob Blob) URL() (string, error) {
	return Config.Service.URL(blob, Config.URLExpire)
}

func (blob Blob) Download() ([]byte, error) {
	return Config.Service.Download(blob)
}

func (blob Blob) Delete() error {
	return Config.Service.Delete(&blob)
}

func (blob Blob) Sync() error {
	return Config.Service.Sync(&blob)
}
