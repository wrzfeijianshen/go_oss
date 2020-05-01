package ossclient

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func Test() {

	// 初始化配置文件
	viper.SetConfigName("app.conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
		return
	}

	osstype := viper.GetInt("oss.osstype")
	accessKey := viper.GetString("alyoss.accessKey")
	secretKey := viper.GetString("alyoss.secretKey")
	bucket := viper.GetString("alyoss.bucket")
	endpoint := viper.GetString("alyoss.endpoint")

	fmt.Println("oss.osstype: ", osstype)
	fmt.Println("alyoss.accessKey: ", accessKey)
	fmt.Println("alyoss.secretKey: ", secretKey)
	fmt.Println("alyoss.bucket: ", bucket)
	fmt.Println("alyoss.endpoint: ", endpoint)

	// 创建oss对象
	c := NewClient(osstype)
	ossObj, err := c.NewAlyOssObj(accessKey, secretKey, bucket, endpoint, "", AlyOssDomain)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ossObj)

	newbucket := "feijianshen3"
	ossObj.CreateBucket(newbucket)
	newbucket = "feijianshen4"
	ossObj.CreateBucket(newbucket)

	ossObj.ListBuckets()

	ossObj.DeleteBucket("feijianshen4")

	newbucket = "feijianshen3"

	oss_bucket := ossObj.AddBucket(newbucket)
	fmt.Println(oss_bucket.Opt.Bucket)
	if oss_bucket == nil {
		fmt.Println("bucket null")
	} else {

		// 上传文件前,先保存到本地,

		// oss_bucket.PutObject("/img/", "aaa")  // err 不能以/开头
		urlpath := oss_bucket.PutObject("img/", "aaa", ACLPrivate)
		fmt.Println(urlpath)
		urlpath = oss_bucket.PutObject("img2", "aaa", ACLPrivate)
		fmt.Println(urlpath)

		urlpath = oss_bucket.PutObjectFromFile("img/a.png", "./bin/a.png", ACLPrivate)

		fmt.Println(urlpath)
		urlpath = oss_bucket.PutObjectFromFile("img/a2.png", "./bin/a.png", ACLPublicRead)
		fmt.Println(urlpath)

		// oss_bucket.ListObjects()

		urlpath = oss_bucket.PutObject("img3", "aaa", ACLPublicRead)

		oss_bucket.DeleteObject("img3")
	}

}
