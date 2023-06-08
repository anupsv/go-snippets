package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
)

type S3Client struct {
	s3        *s3.S3
	bucket string
}

func NewS3Client(bucket string) *S3Client {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &S3Client{
		s3:     s3.New(sess),
		bucket: bucket,
	}
}

func (client *S3Client) UploadFile(filePath string, objectName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	_, err = client.s3.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(client.bucket),
		Key:                  aws.String(objectName),
		ACL:                  aws.String("private"),
		Body:                 aws.ReadSeekCloser(file),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
	})

	return err
}

func (client *S3Client) DownloadFile(objectName string, filePath string) error {
	downloader := s3.New(session.New(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filePath, err)
	}

	defer file.Close()

	_, err = downloader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(client.bucket),
		Key:    aws.String(objectName),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return fmt.Errorf("file %q does not exist", objectName)
			default:
				return aerr
			}
		} else {
			return err
		}
	}

	if _, err := io.Copy(file, downloader.Body); err != nil {
		return fmt.Errorf("failed to download file %q, %v", objectName, err)
	}

	return nil
}

func main() {
	client := NewS3Client("mybucket")

	if err := client.UploadFile("test.txt", "test.txt"); err != nil {
		log.Fatalf("failed to upload file, %v", err)
	}

	if err := client.DownloadFile("test.txt", "downloaded_test.txt"); err != nil {
		log.Fatalf("failed to download file, %v", err)
	}
}
