package active_storage

import (
	"mime/multipart"
	"time"

	"github.com/go-web-kits/dbx"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type ActiveStorageAttachment struct {
	ID        uint      `json:"id" db:"id" gorm:"primary_key; index:index_active_storage_attachments_on_blob_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"not null"`
	OwnerType string    `json:"owner_type" db:"owner_type" gorm:"not null; unique_index:index_active_storage_attachments_uniqueness2"`
	OwnerID   uint      `json:"owner_id" db:"owner_id" gorm:"not null; unique_index:index_active_storage_attachments_uniqueness2"`
	// TODO
	RecordType string `json:"record_type" db:"record_type" gorm:"not null"`
	RecordID   uint   `json:"record_id" db:"record_id" gorm:"not null"`
	Name       string `json:"name" db:"name" gorm:"not null; unique_index:index_active_storage_attachments_uniqueness2"`
	BlobID     uint   `json:"-" gorm:"unique_index:index_active_storage_attachments_uniqueness2"`

	Blob Blob `json:"-" gorm:"ForeignKey:BlobID"`
}

type Attachment = ActiveStorageAttachment

// Attach(1234, Attachment{owner_id: 1, owner_type: "firmwares"}, db.With{Tx: &tx})
// Attach(blob, Attachment{owner_id: 1, owner_type: "firmwares"})
// Attach(params.File, Attachment{owner_id: 1, owner_type: "firmwares"})
func Attach(attachable interface{}, attachmentBase Attachment, opts ...dbx.Opt) dbx.Result {
	attachmentBase.Name = "file"
	attachmentBase.RecordType = attachmentBase.OwnerType
	attachmentBase.OwnerType = strcase.ToSnake(inflection.Plural(attachmentBase.OwnerType))
	attachmentBase.RecordID = attachmentBase.OwnerID

	switch x := attachable.(type) {
	case int, uint: // Blob id
		result := FindBlobBy(dbx.EQ{"id": x}, false)
		if result.Err != nil {
			return result
		}
		attachmentBase.BlobID = result.Data.(*Blob).ID
	case Attachment:
		attachmentBase = x
	case Blob:
		attachmentBase.BlobID = x.ID
	case *Blob:
		attachmentBase.BlobID = x.ID
	case *multipart.FileHeader:
		result := CreateBlobAfterUpload(x, opts...)
		result.Tx = nil
		if result.Err != nil {
			return result
		}
		attachmentBase.BlobID = result.Data.(*Blob).ID
	case *OpenedFileHeader:
		result := CreateBlobAfterUploadByOpened(x, opts...)
		result.Tx = nil
		if result.Err != nil {
			return result
		}
		attachmentBase.BlobID = result.Data.(*Blob).ID
	default:
		panic("unexpected attachable")
	}

	return dbx.FirstOrCreate(&attachmentBase, attachmentBase, opts...)
}
