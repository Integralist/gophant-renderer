package main

import (
	"flag"
	"fmt"

	"github.com/integralist/gophant-renderer/flags"
	"github.com/integralist/gophant-renderer/local"
)

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
	fmt.Println(flags.SqsEndpointName, flags.S3EndpointName, flags.DynamoEndpointName)
	// Values will be Spurious if run without the `-production true` flag
	// Otherwise they'll be empty values so that actual AWS endpoints will be utilised
}
