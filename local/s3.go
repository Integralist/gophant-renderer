package local

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 is the container for our S3 bootstrap data
type S3 struct {
	Wait     *sync.WaitGroup
	Region   string
	Endpoint string
	Bucket   string
}

// Setup will bootstrap a new S3 bucket
func (s *S3) Setup() {
	defer s.Wait.Done()

	svc := s3.New(
		session.New(),
		&aws.Config{
			Region:           aws.String(s.Region),
			Endpoint:         aws.String(s.Endpoint),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		},
	)

	params := &s3.CreateBucketInput{
		Bucket: aws.String(s.Bucket),
	}

	_, err := svc.CreateBucket(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			log.Println("S3 Error:", awsErr.Code(), awsErr.Message())

			// Prints out full error message, including original error if there was one.
			log.Println("S3 Error:", awsErr.Error())

			// Get original error
			if origErr := awsErr.OrigErr(); origErr != nil {
				// operate on original error.
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

// NewS3 creates a new instance of our S3 data container
func NewS3(wg *sync.WaitGroup, region, endpoint, bucket string) *S3 {
	return &S3{
		wg, region, endpoint, bucket,
	}
}
