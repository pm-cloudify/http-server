package acs3

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	} else {
		log.Println("Successfully created session.")
	}
}

func GetBuckets() ([]*s3.Bucket, error) {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return result.Buckets, nil
}

func ListObjects(bucketName string) ([]*s3.Object, error) {
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func GetObject(bucketName, objectKey, filePath string) error {
	downloader := s3manager.NewDownloaderWithClient(svc)
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})

	return err
}

// TODO: add a version which gets a stream of data and sends it
func UploadObject(bucketName, objectKey, filePath string) error {
	uploader := s3manager.NewUploaderWithClient(svc)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	return err
}

func DeleteObject(bucketName, objectKey string) error {
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	return err
}
