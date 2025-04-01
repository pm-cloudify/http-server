package acs3

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pm-cloudify/http-server/internal/config"
)

var sess *session.Session
var svc *s3.S3

func InitConnection() {
	var err error

	sess, err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.LoadedEnv.AC_AccessKey, config.LoadedEnv.AC_SecretKey, ""),
	})

	svc = s3.New(sess, &aws.Config{
		Region:   aws.String(config.LoadedEnv.AC_S3_Region),
		Endpoint: aws.String(config.LoadedEnv.AC_S3_Endpoint),
	})

	if err != nil {
		log.Panicln(err.Error())
	}
}

// TODO: create a bucket
func CreateBucket(bucket_name string) error {
	return nil
}

// TODO: get buckets
func GetBuckets() {}
