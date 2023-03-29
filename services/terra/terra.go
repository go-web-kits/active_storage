package terra

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/go-web-kits/cache"
	// "github.com/go-web-kits/cloudbus"
	"github.com/go-web-kits/utils"
	"github.com/pkg/errors"
)

var (
	CloudBus        *cloudbus.CloudBus
	Service         string
	AppID           string
	Secret          string
	Bucket          string
	Region          string
	UrlExpire       string // minutes
	ExtranetBuckets []string
)

func RequestSTS(fileName string) (*STS, error) {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	params := url.Values{
		"app_id":      {AppID},
		"bucket_name": {Bucket},
		"file_name":   {fileName},
		"timestamp":   {timestamp},
		"sign":        {signature(AppID + Bucket + fileName + timestamp)},
	}

	response, err := CloudBus.POSTRequestForm(Service+".Upload_STSInit.post", params)
	if err != nil {
		return nil, errors.Wrapf(err, "RequestSTS")
	}

	rsp := STSResponse{}
	if err = parse([]byte(response), &rsp); err != nil {
		return nil, errors.Wrapf(err, "RequestSTS parse. bucket: %v. file name: %v", Bucket, fileName)
	}

	return &(rsp.Result), nil
}

func RequestCachedURL(fileId string) (DownloadInfo, error) {
	var info DownloadInfo
	exp, _ := strconv.ParseInt(UrlExpire, 10, 64)
	_, err := cache.Fetch("terra/"+fileId, cache.Opt{
		Default: func() interface{} {
			info, err := RequestURL(fileId)
			if err != nil {
				return err
			}
			return info
		},
		ExpiresIn: time.Duration(exp-1) * time.Minute,
		To:        &info,
	})

	return info, err
}

func RequestURL(fileId string) (DownloadInfo, error) {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	exp, _ := strconv.Atoi(UrlExpire)
	params := url.Values{
		"app_id":      {AppID},
		"bucket_name": {Bucket},
		"file_id":     {fileId},
		"url_expire":  {fmt.Sprint(exp * 60000)},
		"timestamp":   {timestamp},
		"sign":        {signature(AppID + Bucket + fileId + timestamp + UrlExpire)},
	}

	response, err := CloudBus.POSTRequestForm(Service+".Get_Url.post", params)
	if err != nil {
		return DownloadInfo{}, errors.Wrapf(err, "RequestCachedURL. file id: %v. bucket: %v", fileId, Bucket)
	}

	rsp := DownloadResonse{}
	if err = parse([]byte(response), &rsp); err != nil {
		return DownloadInfo{}, errors.Wrapf(err, "RequestCachedURL parse. file id: %v. bucket: %v", fileId, Bucket)
	}

	return rsp.Result, nil
}

func RequestSync(fromID, toBucket string, overwrite bool) error {
	timestamp, o := strconv.FormatInt(time.Now().Unix()*1000, 10), "0"
	if overwrite {
		o = "1"
	}
	params := url.Values{
		"app_id":           {AppID},
		"from_files_id":    {fromID},
		"from_bucket_name": {Bucket},
		"to_bucket_name":   {toBucket},
		"overwrite":        {o},
		"timestamp":        {timestamp},
		"sign":             {signature(AppID + fromID + Bucket + toBucket + o + timestamp)},
	}

	response, err := CloudBus.POSTRequestForm(Service+".Sync_Files.post", params)
	if err != nil {
		return errors.Wrapf(err, "RequestSync. from id: %v. from bucket: %v, to bucket: %v", fromID, Bucket, toBucket)
	}

	rsp := Response{}
	if err = parse([]byte(response), &rsp); err != nil {
		return errors.Wrapf(err,
			"RequestSync parse. from id: %v. from bucket: %v, to bucket: %v", fromID, Bucket, toBucket)
	}

	return nil
}

func DirectUpload(file io.Reader, objectKey, fileName string, fileSize int64, timeout ...time.Duration) (fileId string, err error) {
	result, err := RequestSTS(objectKey)
	if err != nil {
		return "", errors.Wrapf(err, "DirectUpload")
	}

	switch result.UploadType {
	case 0: // 内网
		err = UploadToInternal(result.UploadURL, file, fileName, result.UploadToken, fileSize, timeout...)
		if err != nil {
			return "", err
		}
		return result.FileID, nil
	case 1: // AWS外网
	case 2: // OSS外网
	}

	return "", nil
}

func UploadToInternal(uploadUrl string, file io.Reader, fileName, uploadToken string, fileSize int64, timeout ...time.Duration) error {
	rsp, err := utils.PostByForm(uploadUrl+"/upload/form", fileName, nil, map[string]string{"token": uploadToken}, file, fileSize, timeout...)
	if err != nil {
		return errors.Wrapf(err, "UploadToInternal post form")
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrapf(err, "UploadToInternal read body")
	}

	if err = parse(body, &Response{}); err != nil {
		return errors.Wrapf(err, "UploadToInternal parse")
	}

	return nil
}

// ======

func signature(content string) string {
	s := sha256.New()
	s.Write([]byte(content + Secret))
	return hex.EncodeToString(s.Sum(nil))
}

func parse(body []byte, out interface{}) error {
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		return errors.Wrapf(err, "parse unmarshal body: %v", string(body))
	}

	if response.Code != 0 {
		return fmt.Errorf("parse return error. invoke id: %v. msg: %v", response.InvokeId, response.Message)
	}

	return json.Unmarshal(body, out)
}
