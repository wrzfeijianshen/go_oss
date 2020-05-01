package ossclient

type OSS_Client struct {
	OssType int
}

func NewClient(ossType int) *OSS_Client {
	client := OSS_Client{
		OssType: ossType,
	}
	return &client
}

func (c *OSS_Client) NewAlyOssObj(accessKey, secretKey, bucket, endpoint string, securityToken string, newType int) (obj AlyOssObj, err error) {
	obj.Opt.Accesskey = accessKey
	obj.Opt.Secretkey = secretKey
	obj.Opt.Bucket = bucket
	obj.Opt.Endpoint = endpoint
	obj.Opt.SecurityToken = securityToken
	obj.Newtype = newType

	obj.Init()
	return obj, nil
}
