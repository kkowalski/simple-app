package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Config struct {
	Bucket string
	Client *s3.S3
}

func ConfigureS3Client() *S3Config {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("simple", "simple", ""),
		Endpoint:         aws.String("http://localhost:9000"),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		panic(fmt.Errorf("error creating session to aws: %w", err))
	}

	return &S3Config{Client: s3.New(newSession), Bucket: "simple"}
}

func EnsureBucket(ctx context.Context, s3Client *s3.S3, bucket string) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if !bucketExists(ctxWithTimeout, s3Client, bucket) {
		_, err := s3Client.CreateBucketWithContext(ctxWithTimeout, &s3.CreateBucketInput{Bucket: &bucket})
		if err != nil {
			panic(fmt.Errorf("error creating bucket: %w", err))
		}
	}
}

func bucketExists(ctx context.Context, s3Client *s3.S3, bucket string) bool {
	resp, err := s3Client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		panic(fmt.Errorf("couldn't list buckets: %w", err))
	}

	for _, b := range resp.Buckets {
		if *b.Name == bucket {
			return true
		}
	}
	return false
}
