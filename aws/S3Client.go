package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "strings"
)

// S3ClientInterface is an interface that our S3 Client will adhere to.
type S3ClientInterface interface {
    UploadFile(bucket, key, filePath string) error
    DownloadFile(bucket, key, downloadPath string) error
}

// S3Client is our actual implementation of S3ClientInterface
type S3Client struct {
    uploader * s3manager.Uploader
    downloader * s3manager.Downloader
}

// MyS3Client is a mock implementation of S3ClientInterface for testing
type MyS3Client struct {}

// NewS3Client creates a new S3 client
func NewS3Client(sess * session.Session) S3ClientInterface {
    return &S3Client {
        uploader: s3manager.NewUploader(sess),
        downloader: s3manager.NewDownloader(sess),
    }
}

// UploadFile will upload the file to S3 and return an error if it fails
func(c * S3Client) UploadFile(bucket, key, filePath string) error {
    file, err: = os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file %q, %v", filePath, err)
    }

    _, err = c.uploader.Upload( & s3manager.UploadInput {
        Bucket: aws.String(bucket),
        Key: aws.String(key),
        Body: file,
        // Enable server side encryption
        ServerSideEncryption: aws.String("AES256"),
    })

    if err != nil {
        return fmt.Errorf("failed to upload file, %v", err)
    }

    return nil
}

// DownloadFile will download the file from S3 and return an error if it fails
func(c * S3Client) DownloadFile(bucket, key, downloadPath string) error {
    file, err: = os.Create(downloadPath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = c.downloader.Download(file, & s3.GetObjectInput {
        Bucket: aws.String(bucket),
        Key: aws.String(key),
    })

    if err != nil {
        return err
    }

    return nil
}
