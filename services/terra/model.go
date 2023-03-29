package terra

type Response struct {
	Message  string `json:"message"`
	Code     int    `json:"code"`
	InvokeId string `json:"invokeId"`
}

type STSResponse struct {
	Response
	Result STS `json:"result"`
}

type STS struct {
	// 上传类型，0: 内网上传，1: AWS外网，2: OSS外网
	UploadType      int                    `json:"uploadType,omitempty"`
	CloudName       string                 `json:"cloudName,omitempty"`
	Authorization   STSAuth                `json:"authorization,omitempty"`
	ObjectKey       string                 `json:"objectKey,omitempty"`
	Region          string                 `json:"region,omitempty"`
	CloudBucketName string                 `json:"cloudBucketName,omitempty"`
	EndPoint        string                 `json:"endPoint,omitempty"`
	FileID          string                 `json:"fileId,omitempty"`
	Headers         map[string]interface{} `json:"headers,omitempty"`

	// 内网
	UploadToken string `json:"uploadToken,omitempty"`
	UploadURL   string `json:"uploadURL,omitempty"`
}

type STSAuth struct {
	Ak           string `json:"ak"`
	Sk           string `json:"sk"`
	SessionToken string `json:"sessionToken"`
	CreateTime   int64  `json:"createTime"`
}

type DownloadResonse struct {
	Response
	Result DownloadInfo `json:"result"`
}

type DownloadInfo struct {
	URL     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
}

type CallbackData struct {
	// Type 事件类型， 20 代表上传完成， 30 代表同步完成， 31 代表同步完成(文件已存在)
	Type          int    `json:"type,omitempty"`
	AppID         string `json:"appId,omitempty"`
	BucketName    string `json:"bucketName,omitempty"`
	FileID        string `json:"fileId,omitempty"`
	FileSize      uint   `json:"fileSize,omitempty"`
	MD5           string `json:"md5,omitempty"`
	FileCreatedAt int64  `json:"fileCreatedAt,omitempty"`
	FromFileId    string `json:"fromFileId,omitempty"`
}
