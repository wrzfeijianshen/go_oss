package cloudstore

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wrzfeijianshen/go_oss"
)

type AlyOssObj struct {
	Opt       AlyOssConfigOption
	Newtype   int // 以什么方式创建Clinet
	client    *oss.Client
	BucketMap map[string]*oss.Bucket
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

	c.BucketMap = make(map[string]*oss.Bucket)

	fmt.Println("oss Init")
	return true
}

func (c *AlyOssObj) NewAlyOssObj(accessKey, secretKey, bucket, endpoint string, securityToken string, newType int) (err error) {
	c.Opt.Accesskey = accessKey
	c.Opt.Secretkey = secretKey
	c.Opt.Bucket = bucket
	c.Opt.Endpoint = endpoint
	c.Opt.SecurityToken = securityToken
	c.Newtype = newType

	c.Init()
	return nil
}

func (c *AlyOssObj) Close() {

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
func (c *AlyOssObj) OpenBucket(bucketName string) (err error) {
	isExist, err := c.client.IsBucketExist(bucketName)
	if err != nil {
		handleError(err)
		return err
	}
	if isExist {
		bucket, err := c.client.Bucket(bucketName)

		if err != nil {
			handleError(err)
			return err
		}
		c.BucketMap[bucketName] = bucket
	}
	return nil
}

// 移除map
func (c *AlyOssObj) DeleteBucketObj(bucketName string) {
	delete(c.BucketMap, bucketName)
}

func (c *AlyOssObj) GetUrl(bucket, dndpoint, objectName string) string {
	putUrl := "http://" + bucket + "." + dndpoint + "/" + objectName
	return putUrl
}

// 上传字符串,末尾带/表示创建文件夹 如 /img/
func (c *AlyOssObj) UploadString(bucketName, tmpstr string, args ...string) string {
	if c.BucketMap[bucketName] == nil {
		return ""
	}
	var aclPublicRead string
	for _, arg := range args {
		if arg == ACLPrivate || arg == ACLPublicRead || arg == ACLPublicReadWrite {
			aclPublicRead = arg
		}
	}

	var putUrl string
	isExist, err := c.BucketMap[bucketName].IsObjectExist(tmpstr)
	if err != nil {
		handleError(err)
		return putUrl
	}
	if isExist {
		fmt.Println("Exist:", isExist)
		if aclPublicRead != ACLPrivate {
			putUrl = c.GetUrl(bucketName, c.Opt.Endpoint, tmpstr)
		}
		return putUrl
	}

	// 指定访问权限为公共读，缺省为继承bucket的权限。
	objectAcl := oss.ObjectACL(oss.ACLType(aclPublicRead))

	// 上传字符串。带/表示是创建文件夹，并带权限
	// err = bucket.PutObject(objectName, strings.NewReader(fpath), storageType, objectAcl)

	err = c.BucketMap[bucketName].PutObject(tmpstr, strings.NewReader(""), objectAcl)
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
		putUrl = c.GetUrl(bucketName, c.Opt.Endpoint, tmpstr)
		fmt.Println("上传的文件url： ", putUrl)
	}
	return putUrl
}

// 上传文件
func (c *AlyOssObj) PutObjectFromFile(bucketName, objectName, localFileName string, args ...string) string {
	var putUrl string
	if c.BucketMap[bucketName] == nil {
		return ""
	}
	var aclPublicRead string
	for _, arg := range args {
		if arg == ACLPrivate || arg == ACLPublicRead || arg == ACLPublicReadWrite {
			aclPublicRead = arg
		}
	}

	isExist, err := c.BucketMap[bucketName].IsObjectExist(objectName)
	if err != nil {
		handleError(err)
	}
	if isExist {
		fmt.Println("Exist:", isExist)
		if aclPublicRead != ACLPrivate {
			putUrl = c.GetUrl(bucketName, c.Opt.Endpoint, objectName)
		}
		return putUrl
	}

	// 由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt
	objectAcl := oss.ObjectACL(oss.ACLType(aclPublicRead))

	err = c.BucketMap[bucketName].PutObjectFromFile(objectName, localFileName, objectAcl)
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
		putUrl = c.GetUrl(bucketName, c.Opt.Endpoint, objectName)
		fmt.Println("上传的文件url： ", putUrl)
	}
	return putUrl
}

func (c *AlyOssObj) ListObjects(bucketName string) (files []go_oss.File, err error) {

	// 列举所有文件。
	marker := ""
	lsRes, err := c.BucketMap[bucketName].ListObjects(oss.Marker(marker))
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(lsRes)

	s := AlyOssBucketJsonData{Data: lsRes.Objects}
	fmt.Println(s)
	for _, object := range lsRes.Objects {
		files = append(files, go_oss.File{
			ModTime: object.LastModified,
			Name:    object.Key,
			Size:    object.Size,
			IsDir:   object.Size == 0,
			Header:  map[string]string{},
		})
	}
	return files, nil
}

// 删除文件
func (c *AlyOssObj) DeleteObject(bucketName, objectName string) {

	err := c.BucketMap[bucketName].DeleteObject(objectName)
	if err != nil {
		handleError(err)
	}

}

func (c *AlyOssObj) Download(bucketName, object string, savePath string) (err error) {
	err = c.BucketMap[bucketName].DownloadFile(object, savePath, 1048576)
	return
}

func (c *AlyOssObj) GetInfo(bucketName, object string) (info go_oss.File, err error) {
	// https://help.aliyun.com/document_detail/31859.html?spm=a2c4g.11186623.2.10.713d1592IKig7s#concept-lkf-swy-5db
	//Cache-Control	指定该 Object 被下载时的网页的缓存行为
	//Content-Disposition	指定该 Object 被下载时的名称
	//Content-Encoding	指定该 Object 被下载时的内容编码格式
	//Content-Language	指定该 Object 被下载时的内容语言编码
	//Expires	过期时间
	//Content-Length	该 Object 大小
	//Content-Type	该 Object 文件类型
	//Last-Modified	最近修改时间

	var header http.Header

	header, err = c.BucketMap[bucketName].GetObjectMeta(object)
	if err != nil {
		return
	}

	headerMap := make(map[string]string)

	for k, _ := range header {
		headerMap[k] = header.Get(k)
	}

	info.Header = headerMap
	info.Size, _ = strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	info.ModTime, _ = time.Parse(http.TimeFormat, header.Get("Last-Modified"))
	info.Name = object
	info.IsDir = false
	fmt.Println(info)
	return
}
