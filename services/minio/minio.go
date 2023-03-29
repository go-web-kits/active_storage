package minio

import (
	"io"
	"net/url"
	"time"

	"github.com/go-web-kits/cache"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
)

var (
	client *minio.Client

	Endpoint          string
	Bucket            string
	AccessKeyID       string
	SecretAccessKey   string
	UseSSL            bool
	DownloadURLExpire int // minutes
	UploadURLExpire   int // minutes
	// Region          string
	// ExtranetBuckets []string
)

func Client() *minio.Client {
	if client != nil {
		return client
	}

	c, err := minio.New(Endpoint, AccessKeyID, SecretAccessKey, UseSSL)
	if err != nil {
	}
	client = c
	return client
}

// func RequestSTS(fileName string) (*STS, error) {
// 	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
// 	params := url.Values{
// 		"app_id":      {AppID},
// 		"bucket_name": {Bucket},
// 		"file_name":   {fileName},
// 		"timestamp":   {timestamp},
// 		"sign":        {signature(AppID + Bucket + fileName + timestamp)},
// 	}
//
// 	response, err := CloudBus.POSTRequestForm(Service+".Upload_STSInit.post", params)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "RequestSTS")
// 	}
//
// 	rsp := STSResponse{}
// 	if err = parse([]byte(response), &rsp); err != nil {
// 		return nil, errors.Wrapf(err, "RequestSTS parse. bucket: %v. file name: %v", Bucket, fileName)
// 	}
//
// 	return &(rsp.Result), nil
// }

func RequestPresignedUploadInfo(objectName string) (PresignedUploadInfo, error) {
	expiry := time.Duration(UploadURLExpire) * 60 * time.Second
	presignedURL, err := Client().PresignedPutObject(Bucket, objectName, expiry)
	if err != nil {
		return PresignedUploadInfo{}, err
	}
	return PresignedUploadInfo{
		Filename:   objectName,
		URL:        presignedURL.String(),
		Expiry:     expiry,
		CreateTime: time.Now().Unix(),
	}, nil
}

func RequestCachedDownloadURL(fileName string) (string, error) {
	cachedURL, err := cache.Fetch("minio/"+fileName, cache.Opt{
		Default: func() interface{} {
			url, err := RequestPresignedDownloadURL(fileName)
			if err != nil {
				return err
			}
			return url
		},
		ExpiresIn: time.Duration(DownloadURLExpire-1) * time.Minute,
	})

	return cachedURL.(string), err
}

func RequestPresignedDownloadURL(objectName string) (string, error) {
	// Set request parameters for content-disposition.
	// reqParams := url.Values{"response-content-disposition": {"attachment; filename=\"your-filename.txt\""}}
	reqParams := url.Values{}

	expiry := time.Duration(DownloadURLExpire) * 60 * time.Second
	presignedURL, err := Client().PresignedGetObject(Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", errors.Wrapf(err, "RequestPresignedDownloadURL")
	}

	return presignedURL.String(), nil
}

func RequestSync(fromID, toBucket string, overwrite bool) error {
	return nil
}

func DirectUpload(file io.Reader, objectKey, fileName string, fileSize int64, timeout ...time.Duration) (fileId string, err error) {
	// result, err := RequestSTS(objectKey)
	// if err != nil {
	// 	return "", errors.Wrapf(err, "DirectUpload")
	// }
	//
	// switch result.UploadType {
	// case 0: // 内网
	// 	err = UploadToInternal(result.UploadURL, file, fileName, result.UploadToken, fileSize, timeout...)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return result.FileID, nil
	// case 1: // AWS外网
	// case 2: // OSS外网
	// }

	return "", nil
}

// func UploadToInternal(uploadUrl string, file io.Reader, fileName, uploadToken string, fileSize int64, timeout ...time.Duration) error {
// 	rsp, err := utils.PostByForm(uploadUrl+"/upload/form", fileName, nil, map[string]string{"token": uploadToken}, file, fileSize, timeout...)
// 	if err != nil {
// 		return errors.Wrapf(err, "UploadToInternal post form")
// 	}
//
// 	body, err := ioutil.ReadAll(rsp.Body)
// 	if err != nil {
// 		return errors.Wrapf(err, "UploadToInternal read body")
// 	}
//
// 	if err = parse(body, &Response{}); err != nil {
// 		return errors.Wrapf(err, "UploadToInternal parse")
// 	}
//
// 	return nil
// }
