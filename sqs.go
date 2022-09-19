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
	DecodeBody        bool
	Delete            bool
	Limit             int
	MessageId         string
	VisibilityTimeout int32
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
				break
			}

			counter++
		}

		messages, err := client.receiveMessage(ctx, maxNum)

		if err != nil {
			return fmt.Errorf("failed to receive message: %w", err)
		}

		messagesToDelete := []types.Message{}

		for _, m := range messages {
			if client.MessageId != "" && client.MessageId != *m.MessageId {
				continue
			}

			j, err := marshalMessage(m, client.DecodeBody)

			if err != nil {
				return fmt.Errorf("failed to marshal message: %w", err)
			}

			fmt.Println(string(j))
			messagesToDelete = append(messagesToDelete, m)
		}

		if client.Delete && len(messagesToDelete) > 0 {
			err := client.deleteMessages(ctx, messagesToDelete)

			if err != nil {
				return fmt.Errorf("failed to delete messages: %w", err)
			}

			if client.MessageId != "" {
				break
			}
		}

		time.Sleep(interval)
	}

	return nil
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
