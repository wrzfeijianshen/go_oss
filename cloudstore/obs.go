package cloudstore

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wrzfeijianshen/go_oss"
	"github.com/wrzfeijianshen/go_oss/obs"
)

type ObsObj struct {
	Opt       ObsConfigOption
	Newtype   int // 以什么方式创建Clinet
	client    *obs.ObsClient
	BucketMap map[string]*obs.GetObjectOutput
}

func (c *ObsObj) NewObj(accessKey, secretKey, bucket, endpoint string, securityToken string, newType int) (err error) {
	c.Opt.Accesskey = accessKey
	c.Opt.Secretkey = secretKey
	c.Opt.Bucket = bucket
	c.Opt.Endpoint = endpoint
	c.Newtype = newType

	c.Init()
	return nil
}

// 初始化
func (c *ObsObj) Init() bool {
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
	c.client, err = obs.New(c.Opt.Accesskey, c.Opt.Secretkey, c.Opt.Endpoint)

	if err != nil {
		handleError(err)
		return false
	}
	// obs.InitLog("./logs/OBS-SDK.log", 20480, 10, obs.LEVEL_INFO, false)
	c.BucketMap = make(map[string]*obs.GetObjectOutput)

	fmt.Println("oss Init")
	return true
}
func (c *ObsObj) Close() {
	c.client.Close()
}

// obsClient.Close()

func (c *ObsObj) IsBucketExist(bucketName string) bool {
	if !c.Init() {
		fmt.Println("Init no")
		return false
	}
	output, err := c.client.HeadBucket(bucketName)
	if err == nil {
		fmt.Println("code : ", output.StatusCode, " ", bucketName)

		if output.StatusCode == 404 {
			fmt.Println("Init no")
			return false
		}
	} else {
		return false
	}

	return true
}
func (c *ObsObj) CreateBucket(bucketName string) {
	if c.IsBucketExist(bucketName) {
		fmt.Println("exist", bucketName)
		return
	}
	input := &obs.CreateBucketInput{}
	input.ACL = obs.AclPublicRead
	input.Bucket = bucketName
	_, err := c.client.CreateBucket(input)
	// 判断桶是否存在，返回的结果中HTTP状态码为200表明桶存在，否则返回404表明桶不存在
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.StatusCode)
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}

	}
}

func (c *ObsObj) DeleteBucket(bucketName string) {
	if !c.Init() {
		fmt.Println("Init no")
		return
	}
	c.client.DeleteBucket(bucketName)
}

func (c *ObsObj) ListBuckets() {
	if !c.Init() {
		fmt.Println("Init no")
		return
	}
	output, err := c.client.ListBuckets(nil)
	if err == nil {
		return
		fmt.Printf("RequestId:%s\n", output.RequestId)
		fmt.Printf("Owner.ID:%s\n", output.Owner.ID)
		for index, val := range output.Buckets {
			fmt.Printf("Bucket[%d]-Name:%s,CreationDate:%s\n", index, val.Name, val.CreationDate)
		}
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func (c *ObsObj) OpenBucket(bucketName string) (err error) {
	// 此函数在华为云不存在
	return nil
	isExist := c.IsBucketExist(bucketName)
	if isExist {
		input := &obs.GetObjectInput{}
		input.Bucket = bucketName
		input.Key = "aa/"
		output, err := c.client.GetObject(input) // 打开此对象

		if err == nil {
			return nil
			fmt.Println("err : nil")
			defer output.Body.Close()
			fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
			fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
				output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
			p := make([]byte, 1024)
			var readErr error
			var readCount int
			for {
				readCount, readErr = output.Body.Read(p)
				if readCount > 0 {
					fmt.Printf("%s", p[:readCount])
				}
				if readErr != nil {
					break
				}
			}
		} else {
			fmt.Println("err :no  nil")

			if obsError, ok := err.(obs.ObsError); ok {
				fmt.Println(obsError.StatusCode)
				fmt.Println(obsError.Code)
				fmt.Println(obsError.Message)
			} else {
				fmt.Println(err)
			}
		}

		c.BucketMap[bucketName] = output

	}
	return nil
}
func (c *ObsObj) DeleteBucketObj(bucketName string) {

}
func (c *ObsObj) UploadString(bucketName, objectName string, args ...string) string {
	var aclPublicRead string
	for _, arg := range args {
		if arg == ACLPrivate || arg == ACLPublicRead || arg == ACLPublicReadWrite {
			aclPublicRead = arg
		}
	}
	input := &obs.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = objectName
	input.ACL = obs.AclType(aclPublicRead)
	_, err := c.client.PutObject(input)

	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
		return ""
	}

	var putUrl string
	if aclPublicRead != ACLPrivate {
		putUrl = c.GetUrl(bucketName, c.Opt.Endpoint, objectName)
	}
	return putUrl
}

func (c *ObsObj) GetUrl(bucket, dndpoint, objectName string) string {
	putUrl := "https://" + bucket + "." + dndpoint + "/" + objectName
	return putUrl
}

func (c *ObsObj) PutObjectFromFile(bucketName, objectName, localFileName string, args ...string) string {
	var aclPublicRead string
	for _, arg := range args {
		if arg == ACLPrivate || arg == ACLPublicRead || arg == ACLPublicReadWrite {
			aclPublicRead = arg
		}
	}
	input := &obs.PutFileInput{}
	input.Bucket = bucketName
	input.Key = objectName
	input.SourceFile = localFileName
	input.ACL = obs.AclType(aclPublicRead)
	_, err := c.client.PutFile(input)

	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
		return ""
	}
	// https://feijianshen3.obs.myhuaweicloud.com/img5/a2.png
	var putUrl string
	if aclPublicRead != ACLPrivate {
		putUrl = c.GetUrl(bucketName, c.Opt.Endpoint, objectName)
	}
	return putUrl

}
func (c *ObsObj) DeleteObject(bucketName, objectName string) {
	input := &obs.DeleteObjectInput{}
	input.Bucket = bucketName
	input.Key = objectName
	_, err := c.client.DeleteObject(input)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func (c *ObsObj) Download(bucketName, object string, savePath string) (err error) {
	input := &obs.GetObjectInput{}
	input.Bucket = bucketName
	input.Key = object

	output := &obs.GetObjectOutput{}
	output, err = c.client.GetObject(input)
	if err != nil {
		return
	}
	defer output.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(output.Body)
	if err != nil {
		return
	}

	return ioutil.WriteFile(savePath, b, os.ModePerm)

}
func (c *ObsObj) GetInfo(bucketName, object string) (info go_oss.File, err error) {
	input := &obs.GetObjectMetadataInput{
		Bucket: bucketName,
		Key:    object,
	}
	output := &obs.GetObjectMetadataOutput{}
	output, err = c.client.GetObjectMetadata(input)
	if err != nil {
		return
	}
	info = go_oss.File{
		Name:    object,
		Size:    output.ContentLength,
		IsDir:   output.ContentLength == 0,
		ModTime: output.LastModified,
	}
	fmt.Println(info)

	return

}
func (c *ObsObj) ListObjects(bucketName string) (files []go_oss.File, err error) {
	input := &obs.ListObjectsInput{}
	input.Bucket = bucketName
	output := &obs.ListObjectsOutput{}
	output, err = c.client.ListObjects(input)
	if err != nil {
		return
	}

	for _, item := range output.Contents {
		files = append(files, go_oss.File{
			ModTime: item.LastModified,
			Name:    item.Key,
			Size:    item.Size,
			IsDir:   item.Size == 0,
		})
		fmt.Println(item.Key)
	}
	return

}
