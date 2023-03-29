package active_storage

import (
	"time"

	"github.com/go-web-kits/dbx"
)

var Config Configuration

type Configuration struct {
	Service       Service
	URLExpire     time.Duration
	UploadTimeout time.Duration
}

func Migrate() error {
	return dbx.Conn().AutoMigrate(&Blob{}, &Attachment{}).Error
}
