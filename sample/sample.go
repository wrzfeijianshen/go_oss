package sample

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/wrzfeijianshen/go_oss"
	"github.com/wrzfeijianshen/go_oss/cloudstore"
)

func OssTest() {

	// 初始化配置文件
	viper.SetConfigName("app.conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
		return
	}

	osstype := viper.GetInt("oss.osstype")
	ossaccessKey := viper.GetString("alyoss.accessKey")
	osssecretKey := viper.GetString("alyoss.secretKey")
	ossbucket := viper.GetString("alyoss.bucket")
	ossendpoint := viper.GetString("alyoss.endpoint")

	fmt.Println("oss.osstype: ", osstype)
	fmt.Println("alyoss.accessKey: ", ossaccessKey)
	fmt.Println("alyoss.secretKey: ", osssecretKey)
	fmt.Println("alyoss.bucket: ", ossbucket)
	fmt.Println("alyoss.endpoint: ", ossendpoint)

	fmt.Println("----------------Test end-------------------")

	// 初始化阿里云sdk
	ossClientObj := new(cloudstore.AlyOssObj)
	err = ossClientObj.NewAlyOssObj(ossaccessKey, osssecretKey, ossbucket, ossendpoint, "", cloudstore.AlyOssDomain)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ossClientObj)

	var ossClient go_oss.CloudStore
	ossClient = ossClientObj

	newbucket := "feijianshen3"
	ossClient.CreateBucket(newbucket)
	newbucket = "feijianshen4"
	ossClient.CreateBucket(newbucket)

	ossClient.ListBuckets()

	ossClient.DeleteBucket("feijianshen4")
	fmt.Println("list delete feijianshen4 : ")
	ossClient.ListBuckets()

	newbucket = "feijianshen3"

	err = ossClient.OpenBucket(newbucket)
	if err != nil {
		fmt.Println("bucket null")
	} else {

		// 上传文件前,先保存到本地,

		// oss_bucket.PutObject("/img/", "aaa")  // err 不能以/开头
		urlpath := ossClient.UploadString(newbucket, "img3/", "", cloudstore.ACLPrivate)
		fmt.Println(urlpath)

		urlpath = ossClient.UploadString(newbucket, "img4/", "", cloudstore.ACLPublicRead)
		fmt.Println(urlpath)
		urlpath = ossClient.UploadString(newbucket, "img5/", "aaaa", cloudstore.ACLPublicRead)
		fmt.Println(urlpath)

		urlpath = ossClient.PutObjectFromFile(newbucket, "img5/a.png", "./bin/a.png", cloudstore.ACLPrivate)

		fmt.Println(urlpath)
		urlpath = ossClient.PutObjectFromFile(newbucket, "img5/a2.png", "./bin/a.png", cloudstore.ACLPublicRead)
		fmt.Println(urlpath)

		urlpath = ossClient.UploadString(newbucket, "img6/", "", cloudstore.ACLPrivate)
		//urlpath = ossClient.PutObjectFromFile(newbucket, "img6/a2.png", "./bin/a.png", cloudstore.ACLPublicRead)
		// 有文件删除不了文件夹
		ossClient.DeleteObject(newbucket, "img6/")

		ossClient.Download(newbucket, "img5/a.png", "./bin/down_a.png")
		ossClient.GetInfo(newbucket, "img5/a.png")

		ossClient.ListObjects(newbucket)
	}

}

func ObsTest() {
	// 初始化配置文件
	viper.SetConfigName("app.conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
		return
	}

	osstype := viper.GetInt("oss.osstype")
	accessKey := viper.GetString("obs.accessKey")
	secretKey := viper.GetString("obs.secretKey")
	bucket := viper.GetString("obs.bucket")
	endpoint := viper.GetString("obs.endpoint")

	fmt.Println("oss.osstype: ", osstype)
	fmt.Println("obs.accessKey: ", accessKey)
	fmt.Println("obs.secretKey: ", secretKey)
	fmt.Println("obs.bucket: ", bucket)
	fmt.Println("obs.endpoint: ", endpoint)

	fmt.Println("----------------Test end-------------------")

	// 初始化华为云存储
	ossClientObj := new(cloudstore.ObsObj)
	err = ossClientObj.NewObj(accessKey, secretKey, bucket, endpoint, "", cloudstore.AlyOssDomain)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ossClientObj)

	var ossClient go_oss.CloudStore
	ossClient = ossClientObj

	newbucket := "feijianshen3"
	ossClient.CreateBucket(newbucket)
	newbucket = "feijianshen10"
	bret := ossClient.IsBucketExist(newbucket)
	if !bret {
		fmt.Println("exist no ", newbucket)
	}
	newbucket = "feijianshen4"
	ossClient.CreateBucket(newbucket)

	ossClient.ListBuckets()

	ossClient.DeleteBucket("feijianshen4")
	fmt.Println("list delete feijianshen4 : ")
	ossClient.ListBuckets()

	newbucket = "feijianshen3"

	// 	// 上传文件前,先保存到本地,

	// oss_bucket.PutObject("/img/", "aaa")  // err 不能以/开头
	urlpath := ossClient.UploadString(newbucket, "img3/", "", cloudstore.ACLPrivate)
	fmt.Println(urlpath)
	urlpath = ossClient.UploadString(newbucket, "img34", "", cloudstore.ACLPrivate)
	fmt.Println(urlpath)
	urlpath = ossClient.UploadString(newbucket, "img345", "", cloudstore.ACLPublicRead)
	fmt.Println(urlpath)
	urlpath = ossClient.PutObjectFromFile(newbucket, "img5/a.png", "./bin/a.png", cloudstore.ACLPrivate)

	fmt.Println(urlpath)
	urlpath = ossClient.PutObjectFromFile(newbucket, "img5/a2.png", "./bin/a.png", cloudstore.ACLPublicRead)
	fmt.Println(urlpath)

	urlpath = ossClient.UploadString(newbucket, "img6/", "", cloudstore.ACLPrivate)
	urlpath = ossClient.PutObjectFromFile(newbucket, "img6/a2.png", "./bin/a.png", cloudstore.ACLPublicRead)
	// 有文件也删除文件夹
	ossClient.DeleteObject(newbucket, "img6/")

	ossClient.Download(newbucket, "img5/a.png", "./bin/down_obs.png")
	ossClient.GetInfo(newbucket, "img5/a.png")

	ossClient.ListObjects(newbucket)
}
