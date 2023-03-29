package minio

import (
	"strings"

	"github.com/go-web-kits/active_storage"
	"github.com/go-web-kits/dbx"
)

// https://docs.minio.io/docs/minio-bucket-notification-guide#webhooks

type ParamsMinIOCallback struct {
	EventName string `json:"EventName"`
	Key       string `json:"key"`
}

func UploadedCallback(params ParamsMinIOCallback) (*active_storage.Blob, error) {
	keyPathes := strings.Split(params.Key, "/")
	// 第一个 path 是 bucket name
	key := strings.Join(keyPathes[1:], "/")

	var blob active_storage.Blob
	err := dbx.FindBy(&blob, dbx.EQ{"key": key}).Err
	if err != nil {
		return nil, err
	}

	if params.EventName == "s3:ObjectCreated:Put" {
		err = dbx.UpdateBy(&blob, dbx.H{"upload_status": "upload success"}).Err
		if err != nil {
			return nil, err
		}
	}

	return &blob, nil
}

/* callback data example
{
        "EventName":"s3:ObjectCreated:Put",
        "Key":"test/1zqf605vclhogjnvsji1wrsz/w11/aaf.txt",
        "Records":[
            {
                "eventVersion":"2.0",
                "eventSource":"minio:s3",
                "awsRegion":"",
                "eventTime":"2020-03-17T04:18:27Z",
                "eventName":"s3:ObjectCreated:Put",
                "userIdentity":{
                    "principalId":"minio"
                },
                "requestParameters":{
                    "accessKey":"minio",
                    "region":"",
                    "sourceIPAddress":"127.0.0.1"
                },
                "responseElements":{
                    "content-length":"0",
                    "x-amz-request-id":"15FCFC6B648234D8",
                    "x-minio-deployment-id":"7dc77c88-ba9d-4035-8ff6-72cbef0008b2",
                    "x-minio-origin-endpoint":"http://192.168.125.120:9000"
                },
                "s3":{
                    "s3SchemaVersion":"1.0",
                    "configurationId":"Config",
                    "bucket":{
                        "name":"test",
                        "ownerIdentity":{
                            "principalId":"minio"
                        },
                        "arn":"arn:aws:s3:::test"
                    },
                    "object":{
                        "key":"1zqf605vclhogjnvsji1wrsz%2Fw11%2Faaf.txt",
                        "size":3380,
                        "eTag":"c8e9883eca7c30177ad252cde5434e72-1",
                        "contentType":"multipart/form-data; boundary=--------------------------709471177751309366066290",
                        "userMetadata":{
                            "cache-control":"no-cache",
                            "content-type":"multipart/form-data; boundary=--------------------------709471177751309366066290"
                        },
                        "versionId":"1",
                        "sequencer":"15FCFC6B649EF820"
                    }
                },
                "source":{
                    "host":"127.0.0.1",
                    "port":"",
                    "userAgent":"PostmanRuntime/7.6.0"
                }
            }
        ],
    }
*/
