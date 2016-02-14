package local

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var svc *sqs.SQS

// Sqs is the container for our SQS bootstrap data
type Sqs struct {
	Wait     *sync.WaitGroup
	Region   string
	Endpoint string
	Queue    string
}

type message struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// Setup will bootstrap a new SQS queue and send a message to it
func (s *Sqs) Setup() {
	defer s.Wait.Done()

	svc = sqs.New(
		session.New(),
		&aws.Config{
			Region:     aws.String(s.Region),
			Endpoint:   aws.String(s.Endpoint),
			DisableSSL: aws.Bool(true),
		},
	)

	resp := createQueue(s.Queue)
	sendMessage(resp.QueueUrl)
}

// NewSqs creates a new instance of our Sqs data container
func NewSqs(wg *sync.WaitGroup, region, endpoint, queue string) *Sqs {
	return &Sqs{
		wg, region, endpoint, queue,
	}
}

func createQueue(queue string) *sqs.CreateQueueOutput {
	params := &sqs.CreateQueueInput{
		QueueName: aws.String(queue),
	}

	resp, err := svc.CreateQueue(params)
	handleError("SQS", err)

	fmt.Printf("SQS Create Queue:\n%v\n\n", resp)

	return resp
}

func sendMessage(queueURL *string) {
	msg := message{exampleJSONFormat(), time.Now()}
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("JSON Marshal Error:", err)
	}

	params := &sqs.SendMessageInput{
		MessageBody: aws.String(string(b)),
		QueueUrl:    queueURL,
	}

	resp, err := svc.SendMessage(params)
	handleError("SQS", err)

	fmt.Printf("SQS Sent Message:\n%v\n\n", resp)
}

func exampleJSONFormat() string {
	return `{
		"headers": {
			"ndpSeqNo": 1
		},
		"body": {
			"results": [
				{
					"gssId": "E07000026",
					"regionName": "Foobar",
					"councilWebSite": "http://www.foobar.com/"
				}
			]
		}
	}`
}
