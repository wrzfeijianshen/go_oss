package ossclient

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AlyOssBucket struct {
	bucket *oss.Bucket
	Opt    AlyOssConfigOption // 主要用到了 Bucket Endpoint
}

type AlyOssObj struct {
	Opt       AlyOssConfigOption
	Newtype   int // 以什么方式创建Clinet
	client    *oss.Client
	BucketMap map[string]*AlyOssBucket
}

type AlyOssBucketJsonData struct {
	Code  int         `json:"code"`  //错误代码
	Count int         `json:"count"` // 数据数量
	Msg   string      `json:"msg"`   //输出信息
	Data  interface{} `json:"data"`  //数据
}

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

// 初始化
func (c *AlyOssObj) Init() bool {
	if c.client != nil {
		return true
	}

	if c.Opt.Endpoint == "" {
		fmt.Println("Endpoint nil")
		return false
	}

	if c.Opt.Accesskey == "" {
		fmt.Println("Accesskey nil")
		return false
	}

	if c.Opt.Secretkey == "" {
		fmt.Println("Secretkey nil")
		return false
	}

	var err error
	switch c.Newtype {
	case AlyOssDomain:
		{
			c.client, err = oss.New(c.Opt.Endpoint, c.Opt.Accesskey, c.Opt.Secretkey)
		}
	case AlyOssDiyDomain:
		{
			c.client, err = oss.New(c.Opt.Endpoint, c.Opt.Accesskey, c.Opt.Secretkey, oss.UseCname(true))
		}
	case AlyOssSts:
		{
			c.client, err = oss.New(c.Opt.Endpoint, c.Opt.Accesskey, c.Opt.Secretkey, oss.SecurityToken(c.Opt.SecurityToken))
		}
	}

	if err != nil {
		handleError(err)
		return false
	}

	c.BucketMap = make(map[string]*AlyOssBucket)

	fmt.Println("oss Init")
	return true
}

// 存储空间是否存在
func (c *AlyOssObj) IsBucketExist(bucketName string) bool {
	if !c.Init() {
		fmt.Println("Init no")
		return false
	}
	isExist, err := c.client.IsBucketExist(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
		return false
	}
	fmt.Println("IsBucketExist result : ", isExist)
	return isExist
}

// 创建存储空间,根目录下的 Bucket
func (c *AlyOssObj) CreateBucket(bucketName string) {
	if !c.Init() {
		fmt.Println("Init no")
		return
	}

	if c.IsBucketExist(bucketName) {
		return
	}

	err := c.client.CreateBucket(bucketName)
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println("CreateBucket ok ", bucketName)
}

// 列举存储空间
func (c *AlyOssObj) ListBuckets() {
	if !c.Init() {
		fmt.Println("Init no")
		return
	}

	marker := ""
	for {
		lsRes, err := c.client.ListBuckets(oss.Marker(marker))
		if err != nil {
			handleError(err)
		}

		// 默认情况下一次返回100条记录。
		for _, bucket := range lsRes.Buckets {
			fmt.Println("Bucket: ", bucket.Name)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
}

// 打开 bucket空间
func (c *AlyOssObj) AddBucket(bucketName string) *AlyOssBucket {
	isExist, err := c.client.IsBucketExist(bucketName)
	if err != nil {
		handleError(err)
		return nil
	}
	if isExist {
		bucket, err := c.client.Bucket(bucketName)

		if err != nil {
			handleError(err)
			return nil
		}
		ossbucket := new(AlyOssBucket)
		ossbucket.bucket = bucket
		ossbucket.Opt = c.Opt
		ossbucket.Opt.Bucket = bucketName

		c.BucketMap[bucketName] = ossbucket
		return ossbucket
	}
	return nil
}

func (c *AlyOssObj) DeleteBucket(bucketName string) {
	isExist, err := c.client.IsBucketExist(bucketName)
	if err != nil {
		handleError(err)
		return
	}
	if isExist {
		err := c.client.DeleteBucket(bucketName)
		if err != nil {
			handleError(err)
			return
		}
	}
	return
}

// 移除map
func (c *AlyOssObj) DeleteBucketObj(bucketName string) {
	delete(c.BucketMap, bucketName)
}

func (c *AlyOssBucket) GetUrl(bucket, dndpoint, objectName string) string {
	putUrl := "http://" + bucket + "." + dndpoint + "/" + objectName
	return putUrl
}

// 上传字符串,末尾带/表示创建文件夹 如 /img/
func (c *AlyOssBucket) PutObject(objectName, fpath, aclPublicRead string) string {
	var putUrl string
	isExist, err := c.bucket.IsObjectExist(objectName)
	if err != nil {
		handleError(err)
		return putUrl
	}
	if isExist {
		fmt.Println("Exist:", isExist)
		if aclPublicRead != ACLPrivate {
			putUrl = c.GetUrl(c.Opt.Bucket, c.Opt.Endpoint, objectName)
		}
		return putUrl
	}

	// 指定访问权限为公共读，缺省为继承bucket的权限。
	objectAcl := oss.ObjectACL(oss.ACLType(aclPublicRead))

	// 上传字符串。带/表示是创建文件夹，并带权限
	// err = bucket.PutObject(objectName, strings.NewReader(fpath), storageType, objectAcl)
	err = c.bucket.PutObject(objectName, strings.NewReader(fpath), objectAcl)
	if err != nil {
		handleError(err)
		return putUrl
	}

	// 权限不是私有 返回地址

	if aclPublicRead != ACLPrivate {
		// publish-read文件的话，直接拼接路径。
		// 如果想获取临时URL的话，可以用ossutil工具sign命令生成
		// url就是http://$bucket.$endpoit/$object
		// http://xxxxxx.oss-cn-zhangjiakou.aliyuncs.com/img/xx.png
		putUrl = c.GetUrl(c.Opt.Bucket, c.Opt.Endpoint, objectName)
		fmt.Println("上传的文件url： ", putUrl)
	}
	return putUrl
}

// 上传文件
func (c *AlyOssBucket) PutObjectFromFile(objectName, localFileName, aclPublicRead string) string {
	var putUrl string

	isExist, err := c.bucket.IsObjectExist(objectName)
	if err != nil {
		handleError(err)
	}
	if isExist {
		fmt.Println("Exist:", isExist)
		if aclPublicRead != ACLPrivate {
			putUrl = c.GetUrl(c.Opt.Bucket, c.Opt.Endpoint, objectName)
		}
		return putUrl
	}

	// 由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt
	objectAcl := oss.ObjectACL(oss.ACLType(aclPublicRead))

	err = c.bucket.PutObjectFromFile(objectName, localFileName, objectAcl)
	if err != nil {
		handleError(err)
		return putUrl
	}

	// 通过参数设置权限或者如下代码进行设置权限
	// err = c.bucket.SetObjectACL(file, oss.ACLPublicRead)
	// if err != nil {
	// 	handleError(err)
	// }

	// 权限不是私有 返回地址

	if aclPublicRead != ACLPrivate {
		// url就是http://$bucket.$endpoit/$object
		putUrl = c.GetUrl(c.Opt.Bucket, c.Opt.Endpoint, objectName)
		fmt.Println("上传的文件url： ", putUrl)
	}
	return putUrl
}

func (c *AlyOssBucket) ListObjects() {

	// 列举所有文件。
	marker := ""
	lsRes, err := c.bucket.ListObjects(oss.Marker(marker))
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(lsRes)

	s := AlyOssBucketJsonData{Data: lsRes.Objects}
	fmt.Println(s)
}

// 删除文件
func (c *AlyOssBucket) DeleteObject(objectName string) {

	err := c.bucket.DeleteObject(objectName)
	if err != nil {
		handleError(err)
	}

}
