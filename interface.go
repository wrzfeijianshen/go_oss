package go_oss

import "time"

type File struct {
	ModTime time.Time
	Name    string
	Size    int64
	IsDir   bool
	Header  map[string]string
}

type CloudStore interface {
	IsBucketExist(bucketName string) bool                                                  // 判断存储空间是否存在
	CreateBucket(bucketName string)                                                        // 创建存储空间,根目录下的 Bucket
	DeleteBucket(bucketName string)                                                        // 删除存储空间
	ListBuckets()                                                                          // 列举存储空间
	OpenBucket(bucketName string) (err error)                                              // 打开bucket
	DeleteBucketObj(bucketName string)                                                     // 移除已经打开的bucket
	UploadString(bucketName, tmpstr string, args ...string) string                         // 上传字符串
	PutObjectFromFile(bucketName, objectName, localFileName string, args ...string) string // 上传文件
	DeleteObject(bucketName, objectName string)                                            // 删除文件
	Download(bucketName, object string, savePath string) (err error)                       // 下载文件
	GetInfo(bucketName, object string) (info File, err error)
	ListObjects(bucketName string) (files []File, err error)
	Close()
	// GetSignURL(object string, expire int64) (link string, err error)                  // 文件访问签名
	// IsExist(object string) (err error)                                                // 判断文件是否存在
	// Lists(prefix string) (files []File, err error)                                    // 文件前缀，列出文件
	// Upload(tmpFile string, saveFile string, headers ...map[string]string) (err error) // 上传文件
	// GetInfo(object string) (info File, err error)                                     // 获取指定文件信息
}
