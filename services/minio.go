package services

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-web-kits/active_storage"
	"github.com/go-web-kits/active_storage/services/minio"
	"github.com/go-web-kits/dbx"
	"github.com/pkg/errors"
)

type ASMinIO struct{}

func (s ASMinIO) Upload(blob *active_storage.Blob, file io.Reader, checksum string, timeout ...time.Duration) error {
	// filename := blob.Key + "/" + blob.Filename
	// fileID, err := terra.DirectUpload(file, filename, filename, int64(blob.ByteSize), timeout...)
	// blob.Key = fileID
	// TODO using minio's API
	return nil
}

func (s ASMinIO) Download(blob active_storage.Blob) ([]byte, error) {
	url, err := minio.RequestCachedDownloadURL(blob.Key)
	if err != nil {
		return nil, errors.Wrapf(err, "Download get url")
	}

	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Download Do. url: %v", url)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Download read body")
	}

	return body, nil
}

func (s ASMinIO) Delete(blob *active_storage.Blob) error {
	//
	return nil
}

func (s ASMinIO) URLWithHeader(blob active_storage.Blob, expire ...time.Duration) (string, map[string]interface{}, error) {
	if blob.Key == "" {
		return "", map[string]interface{}{}, nil
	}
	url, err := minio.RequestCachedDownloadURL(blob.Key)
	return url, map[string]interface{}{}, err
}

func (s ASMinIO) URL(blob active_storage.Blob, expire ...time.Duration) (string, error) {
	if blob.Key == "" {
		return "", nil
	}
	url, err := minio.RequestCachedDownloadURL(blob.Key)
	return url, err
}

// returns minio.PresignedUploadInfo
func (s ASMinIO) DirectUploadInfo(blob *active_storage.Blob) (interface{}, error) {
	filename := blob.Key + "/" + blob.Filename
	uploadInfo, err := minio.RequestPresignedUploadInfo(filename)
	if err != nil {
		return nil, err
	}
	blob.Key = filename

	return uploadInfo, dbx.UpdateBy(blob, dbx.H{"key": blob.Key}).Err
}

func (s ASMinIO) Sync(blob *active_storage.Blob) error {
	return nil
}
