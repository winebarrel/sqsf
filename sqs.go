package sqsf

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	maxNumberOfMessages = 10
	waitTimeSeconds     = 20
	interval            = 1 * time.Second
)

type SqsfOpts struct {
	QueueName         string
	Decode            bool
	Delete            bool
	VisibilityTimeout int32
	Limit             int
}

type Client struct {
	*SqsfOpts
	sqs      *sqs.Client
	QueueUrl string
}

func NewClient(ctx context.Context, opts *SqsfOpts) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sqs.NewFromConfig(cfg)
	queueUrl, err := getQueueUrl(ctx, client, opts.QueueName)

	if err != nil {
		return nil, fmt.Errorf("failed to get queue URL: %w", err)
	}

	sqs := &Client{
		SqsfOpts: opts,
		sqs:      client,
		QueueUrl: queueUrl,
	}

	return sqs, nil
}

func getQueueUrl(ctx context.Context, client *sqs.Client, queueName string) (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	output, err := client.GetQueueUrl(ctx, input)

	if err != nil {
		return "", err
	}

	return aws.ToString(output.QueueUrl), nil
}

func (client *Client) Follow(ctx context.Context) error {
	maxNum := maxNumberOfMessages
	counter := 0

	if 0 < client.Limit && client.Limit < maxNumberOfMessages {
		maxNum = client.Limit
	}

	for {
		if client.Limit > 0 {
			if counter >= client.Limit {
				return nil
			}

			counter++
		}

		messages, err := client.receiveMessage(ctx, maxNum)

		if err != nil {
			return fmt.Errorf("failed to receive message: %w", err)
		}

		for _, m := range messages {
			j, err := marshalMessage(m, client.Decode)

			if err != nil {
				return fmt.Errorf("failed to marshal message: %w", err)
			}

			fmt.Println(string(j))
		}

		if client.Delete && len(messages) > 0 {
			err := client.deleteMessages(ctx, messages)

			if err != nil {
				return fmt.Errorf("failed to delete messages: %w", err)
			}
		}

		time.Sleep(interval)
	}
}

func (client *Client) receiveMessage(ctx context.Context, maxNum int) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(client.QueueUrl),
		MaxNumberOfMessages: int32(maxNum),
		WaitTimeSeconds:     waitTimeSeconds,
		VisibilityTimeout:   client.VisibilityTimeout,
	}

	output, err := client.sqs.ReceiveMessage(ctx, input)

	if err != nil {
		return nil, err
	}

	return output.Messages, nil
}

func (client *Client) deleteMessages(ctx context.Context, messages []types.Message) error {
	input := &sqs.DeleteMessageBatchInput{
		Entries:  make([]types.DeleteMessageBatchRequestEntry, 0, len(messages)),
		QueueUrl: aws.String(client.QueueUrl),
	}

	for _, m := range messages {
		input.Entries = append(input.Entries, types.DeleteMessageBatchRequestEntry{
			Id:            m.MessageId,
			ReceiptHandle: m.ReceiptHandle,
		})
	}

	_, err := client.sqs.DeleteMessageBatch(ctx, input)

	if err != nil {
		return err
	}

	return nil
}
