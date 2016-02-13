package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/integralist/gophant-renderer/local"
)

type network []struct {
	Host     string
	HostPort string
}

type spurious struct {
	Sqs    network `json:"spurious-sqs"`
	S3     network `json:"spurious-s3"`
	Dynamo network `json:"spurious-dynamo"`
}

var (
	queueName          string
	regionName         string
	bucketName         string
	sqsEndpointName    string
	s3EndpointName     string
	dynamoEndpointName string
	err                error
	spur               spurious
)

var region = flag.String("region", "eu-west-1", "Name of AWS region to create the resources within")
var queue = flag.String("queue", "producer", "Name of SQS message queue to be created")
var bucket = flag.String("bucket", "gophant", "Name of S3 storage bucket to be created")
var sqsEndpoint = flag.String("sqs-endpoint", "", "Fake AWS SQS endpoint")
var s3Endpoint = flag.String("s3-endpoint", "", "Fake AWS S3 endpoint")
var dynamoEndpoint = flag.String("dynamo-endpoint", "", "Fake AWS DynamoDB endpoint")

func init() {
	flag.Parse()

	queueName = *queue
	regionName = *region
	bucketName = *bucket
	sqsEndpointName = *sqsEndpoint
	s3EndpointName = *s3Endpoint
	dynamoEndpointName = *dynamoEndpoint

	cmdOut, err := getSpuriousNetworkDetails()
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running 'spurious ports --json' command: ", err)
		os.Exit(1)
	}

	if sqsEndpointName == "" || s3EndpointName == "" || dynamoEndpointName == "" {
		json.Unmarshal(cmdOut, &spur)
		setMissingSpuriousNetworkDetails()
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go local.NewSqs(&wg, regionName, sqsEndpointName, queueName).Setup()
	go local.NewS3(&wg, regionName, s3EndpointName, bucketName).Setup()
	go local.NewDynamo(&wg, regionName, dynamoEndpointName).Setup()

	wg.Wait()
	fmt.Println("\n\nFinished creating local resources")
}

func getSpuriousNetworkDetails() ([]byte, error) {
	var cmdOut []byte

	cmdName := "spurious"
	cmdArgs := []string{"ports", "--json"}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return nil, err
	}

	return cmdOut, nil
}

func setMissingSpuriousNetworkDetails() {
	if sqsEndpointName == "" {
		sqsEndpointName = spur.Sqs[0].Host + ":" + spur.Sqs[0].HostPort
	}

	if s3EndpointName == "" {
		s3EndpointName = spur.S3[0].Host + ":" + spur.S3[0].HostPort
	}

	if dynamoEndpointName == "" {
		dynamoEndpointName = spur.Dynamo[0].Host + ":" + spur.Dynamo[0].HostPort
	}
}
