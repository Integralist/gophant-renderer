package flags

import "flag"

var (
	QueueName          string
	RegionName         string
	BucketName         string
	SqsEndpointName    string
	S3EndpointName     string
	DynamoEndpointName string
)

var Region = flag.String("region", "eu-west-1", "Name of AWS region to create the resources within")
var Queue = flag.String("queue", "producer", "Name of SQS message queue to be created")
var Bucket = flag.String("bucket", "gophant", "Name of S3 storage bucket to be created")
var SqsEndpoint = flag.String("sqs-endpoint", "", "AWS SQS endpoint")
var S3Endpoint = flag.String("s3-endpoint", "", "AWS S3 endpoint")
var DynamoEndpoint = flag.String("dynamo-endpoint", "", "AWS DynamoDB endpoint")
var Production = flag.String("production", "false", "Indicates if we're running this in production")
