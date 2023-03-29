package services

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-web-kits/active_storage"
	"github.com/go-web-kits/active_storage/services/terra"
	"github.com/go-web-kits/dbx"
	"github.com/pkg/errors"
)

type ASTerra struct{}

func (s ASTerra) Upload(blob *active_storage.Blob, file io.Reader, checksum string, timeout ...time.Duration) error {
	filename := blob.Key + "/" + blob.Filename
	fileID, err := terra.DirectUpload(file, filename, filename, int64(blob.ByteSize), timeout...)
	blob.Key = fileID
	return err
}

func (s ASTerra) Download(blob active_storage.Blob) ([]byte, error) {
	info, err := terra.RequestCachedURL(blob.Key)
	if err != nil {
		return nil, errors.Wrapf(err, "Download get url")
	}

	req, _ := http.NewRequest("GET", info.URL, nil)
	for k, v := range info.Headers {
		req.Header.Set(k, v.(string))
	}
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Download Do. url: %v", info)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Download read body")
	}

	return body, nil
}

func (s ASTerra) Delete(blob *active_storage.Blob) error {
	//
	return nil
}

func (s ASTerra) URLWithHeader(blob active_storage.Blob, expire ...time.Duration) (string, map[string]interface{}, error) {
	if blob.Key == "" {
		return "", map[string]interface{}{}, nil
	}
	info, err := terra.RequestCachedURL(blob.Key)
	return info.URL, info.Headers, err
}

func (s ASTerra) URL(blob active_storage.Blob, expire ...time.Duration) (string, error) {
	if blob.Key == "" {
		return "", nil
	}
	info, err := terra.RequestCachedURL(blob.Key)
	return info.URL, err
}

func (s ASTerra) DirectUploadInfo(blob *active_storage.Blob) (interface{}, error) {
	filename := blob.Key + "/" + blob.Filename
	sts, err := terra.RequestSTS(filename)
	if err != nil {
		return nil, err
	}
	blob.Key = sts.FileID

	return sts, dbx.UpdateBy(blob, dbx.H{"key": blob.Key}).Err
}

func (s ASTerra) Sync(blob *active_storage.Blob) error {
	for _, bk := range terra.ExtranetBuckets {
		err := terra.RequestSync(blob.Key, bk, true)
		if err != nil {
			return err
		}
	}
	return nil
}
