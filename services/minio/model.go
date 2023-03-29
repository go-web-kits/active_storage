package minio

import "time"

// type STSResponse struct {
// 	Result STS `json:"result"`
// }
//
// type STS struct {
// 	// 上传类型，0: 内网上传，1: AWS外网，2: OSS外网
// 	UploadType      int                    `json:"uploadType,omitempty"`
// 	CloudName       string                 `json:"cloudName,omitempty"`
// 	Authorization   STSAuth                `json:"authorization,omitempty"`
// 	ObjectKey       string                 `json:"objectKey,omitempty"`
// 	Region          string                 `json:"region,omitempty"`
// 	CloudBucketName string                 `json:"cloudBucketName,omitempty"`
// 	EndPoint        string                 `json:"endPoint,omitempty"`
// 	FileID          string                 `json:"fileId,omitempty"`
// 	Headers         map[string]interface{} `json:"headers,omitempty"`
//
// 	// 内网
// 	UploadToken string `json:"uploadToken,omitempty"`
// 	UploadURL   string `json:"uploadURL,omitempty"`
// }
//
// type STSAuth struct {
// 	Ak           string `json:"ak"`
// 	Sk           string `json:"sk"`
// 	SessionToken string `json:"sessionToken"`
// 	CreateTime   int64  `json:"createTime"`
// }

type PresignedUploadInfo struct {
	Filename   string
	URL        string
	Expiry     time.Duration // seconds
	CreateTime int64
}
