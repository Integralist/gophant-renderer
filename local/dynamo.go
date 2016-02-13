package local

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Dynamo is the container for our DynamoDB bootstrap data
type Dynamo struct {
	Wait     *sync.WaitGroup
	Region   string
	Endpoint string
}

// Setup will bootstrap a new DynamoDB Table
func (s *Dynamo) Setup() {
	defer s.Wait.Done()

	svc := dynamodb.New(
		session.New(),
		&aws.Config{
			Region:     aws.String(s.Region),
			Endpoint:   aws.String(s.Endpoint),
			DisableSSL: aws.Bool(true),
		},
	)

	seqResp, err := svc.CreateTable(sequencerTableInput())
	handleError("DynamoDB", err)

	lookupResp, err := svc.CreateTable(lookupTableInput())
	handleError("DynamoDB", err)

	fmt.Printf("Sequence Table:\n%v\n\nLookup Table:\n%v", seqResp, lookupResp)
}

// NewDynamo creates a new instance of our Dynamo data container
func NewDynamo(wg *sync.WaitGroup, region, endpoint string) *Dynamo {
	return &Dynamo{
		wg, region, endpoint,
	}
}

func handleError(service string, err error) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			log.Printf("%s Error: %v, %v", service, awsErr.Code(), awsErr.Message())

			// Prints out full error message, including original error if there was one.
			log.Printf("%s Error: %v", service, awsErr.Error())

			// Get original error
			if origErr := awsErr.OrigErr(); origErr != nil {
				// operate on original error.
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

func sequencerTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("sequencer"),
	}
}

func lookupTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("component_key"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("batch_version"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("component_key"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("batch_version"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("lookup"),
	}
}
