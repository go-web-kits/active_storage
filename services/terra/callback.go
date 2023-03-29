package terra

import (
	"encoding/json"
	"fmt"

	"github.com/go-web-kits/active_storage"
	"github.com/go-web-kits/dbx"
)

func UploadedCallback(BizData string) (*active_storage.Blob, error) {
	var data CallbackData
	err := json.Unmarshal([]byte(BizData), &data)
	if err != nil {
		return nil, err
	}

	if data.BucketName != Bucket {
		return nil, nil
	}

	status, fid := "", ""
	var blob active_storage.Blob
	switch data.Type {
	case 20:
		status = "upload success"
		fid = data.FileID
	case 30, 31:
		status = "sync success"
		fid = data.FromFileId
	default:
		status = fmt.Sprintf("failed with %v", data.Type)
	}

	err = dbx.FindBy(&blob, dbx.EQ{"key": fid}).Err
	if err != nil {
		return nil, err
	}

	if data.Type == 20 {
		if blob.MD5() == data.MD5 {
			status += ", MD5 comparison success"
		} else {
			status += ", MD5 comparison failed"
		}
	}

	update := dbx.H{
		"upload_status": status,
	}
	if data.Type == 30 {
		update["key"] = data.FileID
		update["byte_size"] = data.FileSize
	}

	err = dbx.UpdateBy(&blob, update).Err
	if err != nil {
		return nil, err
	}

	return &blob, nil
}
