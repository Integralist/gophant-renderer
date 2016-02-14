package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	svc := sqs.New(
		session.New(),
		&aws.Config{
			Region:     aws.String(flags.RegionName),
			Endpoint:   aws.String(flags.SqsEndpointName),
			DisableSSL: aws.Bool(true),
		},
	)

	params := &sqs.ReceiveMessageInput{
		QueueUrl:        aws.String(flags.QueueName),
		WaitTimeSeconds: aws.Int64(5),
	}

	run(svc, params)
}

func run(svc *sqs.SQS, params *sqs.ReceiveMessageInput) {
	for {
		resp, err := svc.ReceiveMessage(params)
		if err != nil {
			fmt.Println("There was an issue receiving a message: ", err)
		}

		renderIntoTemplate()

		if len(resp.Messages) > 0 {
			fmt.Printf("\nDisplay received message:\n\n%+v", resp)
		}

		time.Sleep(5 * time.Second) // message needs to be explicitly deleted still
	}
}

func renderIntoTemplate() {
	templ, err := ioutil.ReadFile("./template.tpl")
	if err != nil {
		fmt.Println("Unable to load the message template")
	}

	setupTemplate := template.Must(
		template.New("component").Parse(string(templ)),
	)

	type dataSource struct {
		RegionName     string
		CouncilWebSite string
	}

	ds := dataSource{
		"region name",
		"council website",
	}

	if err := setupTemplate.Execute(os.Stdout, ds); err != nil {
		log.Fatal(err)
	}
}
