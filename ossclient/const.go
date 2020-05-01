package ossclient

const (
	ACLPrivate         = "private"           // 私有权限
	ACLPublicRead      = "public-read"       // 公共读
	ACLPublicReadWrite = "public-read-write" // 公共读写
)

const (
	AlyOssType    = 0 // 阿里云oss
	QiniuyOssType = 1 // 七牛云
	UpyOssType    = 2 // 又拍云
	TxOssType     = 3 // 腾讯云
)

const (
	AlyOssDomain    = 0 // oss 域名
	AlyOssDiyDomain = 1 // oss 自定义域名
	AlyOssSts       = 2 // sts新建client
)

// 阿里云对象存储配置
type AlyOssConfigOption struct {
	Bucket    string // 命令存储空间
	Accesskey string // 用户 key id
	Secretkey string // 用户 key value
	// https://helpcdn.aliyun.com/document_detail/31837.html
	Endpoint      string // 地域节点 oss-cn-zhangjiakou.aliyuncs.com 不带http://
	SecurityToken string
}

// 又拍云配置
type UpyunConfigOption struct {
	Bucket   string
	Operator string
	Password string
	Domain   string
}

//腾讯云对象存储配置
type CosConfigOption struct {
	Bucket    string
	APPID     string
	Region    string
	SecretID  string
	SecretKey string
}
