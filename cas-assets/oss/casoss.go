package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	// "log"
	// "reflect"
)

var (
	Bucket     *oss.Bucket
	ImgOptions = []oss.Option{
		// oss.Expires(expires),
		oss.ObjectACL(oss.ACLPublicRead),
		// oss.Meta("MyProp", "MyPropVal"),
		oss.CacheControl("max-age=60"),
	}
	ImgPrefix = "assets/img/"
	Region    = "zhangjiakou"
	RegionUrl = "oss-cn-zhangjiakou.aliyuncs.com"
)

func InitBucket(regionEndpoint, key, secret, bucketName string) {
	client, err := oss.New(regionEndpoint, key, secret)
	if err != nil {
		panic(err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	Bucket = bucket
}
