package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/integralist/gophant-renderer/flags"
	"github.com/integralist/gophant-renderer/local"
)

var err error

func init() {
	flag.Parse()

	flags.QueueName = *flags.Queue
	flags.RegionName = *flags.Region
	flags.BucketName = *flags.Bucket
	flags.SqsEndpointName = *flags.SqsEndpoint
	flags.S3EndpointName = *flags.S3Endpoint
	flags.DynamoEndpointName = *flags.DynamoEndpoint

	local.CheckFlags(flags.Production, &flags.SqsEndpointName, &flags.S3EndpointName, &flags.DynamoEndpointName)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go local.NewSqs(&wg, flags.RegionName, flags.SqsEndpointName, flags.QueueName).Setup()
	go local.NewS3(&wg, flags.RegionName, flags.S3EndpointName, flags.BucketName).Setup()
	go local.NewDynamo(&wg, flags.RegionName, flags.DynamoEndpointName).Setup()

	wg.Wait()

	fmt.Println("\n\nFinished creating local resources")
}
