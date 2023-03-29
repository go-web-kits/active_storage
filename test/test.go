package test

import (
	"time"

	. "github.com/go-web-kits/active_storage"
)

type Post struct {
	ID        uint       `json:"id" db:"id" gorm:"primary_key;index"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"` // `sql:"index"`
	Picture   Attachment `json:"-" gorm:"polymorphic:Owner"`
}

var Models = []interface{}{&Post{}, &Attachment{}, &Blob{}}
