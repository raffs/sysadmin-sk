/*
 * This file is part of the Sysadmin Sidekick Toolkit (Sysadmin-SK) (https://github.com/raffs/sysadmin-sk).
 * Copyright (c) 2019 Rafael Oliveira Silva
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */
package sqs

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

// MoveMessagesOptions defines all the configuration options for `aws sqs move` command
type moveMessageOptions struct {

	// Define the Source Queue Name that will be move the message from.
	SourceQueueName string `type:"string" required:"true"`

	// Define the Source Queue Name that will be move the message from.
	SourceQueueURL string `type:"string" required:"true"`

	// Define the target queue where the message will be moved to.
	TargetQueueName string `type:"string" required:"true"`

	// Define the target queue where the message will be moved to.
	TargetQueueURL string `type:"string" required:"true"`

	// Define the maximum number of messages to be processed at a time.
	BatchSize int64 `type:"int64" required:"false"`

	// TODO: Proper documentation of this here, please :)
	WaitTimeSeconds int64 `type:"int64" required:"false"`

	// Define the message visibility when ingesting the message to the target queue.
	VisibilityTimeout int64 `type:"int64" required:"false"`

	// Whether to delete a message from the source queue. default: false
	KeepMessageOnSourceQueue bool `type:"string" required:"false"`

	// Define the AWS Region to connect to. This essentially will be converted
	// to an URL with the region name.
	AwsRegion string `type:"string" required:"false"`

	// In case you want to overwrite the underlying endpoint for testing or
	// other kind black Sorcery.
	AwsEndpoint string `type:"string" required:"false"`

	// Define the AWS profile
	AwsProfile string `type:"string" required:"false"`

	// Pointer to a ReceiptHandle
	ReceiptHandlers map[string]string `type:"map[string]*string" required:"false"`
}

/**
 * Given a queue name return the URL. Just a wrapper because we need to use twice.
 */
func getQueueURL(client *sqs.SQS, queueName *string) *sqs.GetQueueUrlOutput {
	queue, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(*queueName),
	})

	if err != nil {
		panic(err)
	}

	return queue
}

/**
 * Given the queue URL return the Queue attributes which include queue type, ARN and
 * more important the approximate number of messages at the moment.
 */
func getQueueAttributes(client *sqs.SQS, queueURL *string) *sqs.GetQueueAttributesOutput {
	queue, err := client.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(*queueURL),
		AttributeNames: aws.StringSlice([]string{"All"}),
	})

	if err != nil {
		panic(err)
	}

	return queue
}

// sqsClient create and returns a sqs client object
func sqsClient(options *moveMessageOptions) (*sqs.SQS, error) {
	sessionOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		// aws configuration
		Config: aws.Config{
			Region:   aws.String(options.AwsRegion),
			Endpoint: aws.String(options.AwsEndpoint),
		},
	}

	session, err := session.NewSessionWithOptions(sessionOpts)
	if err != nil {
		fmt.Println("Unable to create new session with AWS with error: ", err.Error())
		return nil, errors.New("Unable to initialize AWS session")
	}

	return sqs.New(session), nil
}

// sendBatchMessages sends message to the target SQS queue in batch mode
func sendBatchMessages(client *sqs.SQS,
	options *moveMessageOptions,
	receiveResponse *sqs.ReceiveMessageOutput) (*sqs.SendMessageBatchOutput, error) {
	var sendBatchMessages []*sqs.SendMessageBatchRequestEntry

	// append each received message to the send and delete buffer
	for _, message := range receiveResponse.Messages {
		mRequest := sqs.SendMessageBatchRequestEntry{
			MessageAttributes: message.MessageAttributes,
			MessageBody:       message.Body,
			Id:                message.MessageId,
		}

		// keep a map between message ID and ReceiptHandle to be used when deleting the message
		options.ReceiptHandlers[*message.MessageId] = *message.ReceiptHandle

		// append message to the SendMessage array to be send in batch
		sendBatchMessages = append(sendBatchMessages, &mRequest)
	}

	batchSendMessagesInput := &sqs.SendMessageBatchInput{
		QueueUrl: &options.TargetQueueURL,
		Entries:  sendBatchMessages,
	}

	sendResponse, err := client.SendMessageBatch(batchSendMessagesInput)
	if err != nil {
		fmt.Println("Failed to send the message to target queue in batch mode")
		fmt.Println("We should abort this, as a sense something is wrong")
		fmt.Println("API returned: ", err.Error())
		return nil, errors.New("Failed to send messages to target queue")
	}

	fmt.Printf(".") // print a . (dot) for each send OP
	return sendResponse, nil
}

// deleteBatchMessages deletes message to the target SQS queue in batch mode
func deleteBatchMessages(client *sqs.SQS,
	options *moveMessageOptions,
	sendResponse *sqs.SendMessageBatchOutput) (int64, error) {

	var err error
	var deletedMsgs int64
	var deleteBatchMessages []*sqs.DeleteMessageBatchRequestEntry

	// append all successfully messages to be deleted.
	for _, message := range sendResponse.Successful {
		m := &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(*message.Id),
			ReceiptHandle: aws.String(options.ReceiptHandlers[*message.Id]),
		}

		deleteBatchMessages = append(deleteBatchMessages, m)
	}

	batchDeleteMessagesInput := &sqs.DeleteMessageBatchInput{
		QueueUrl: &options.SourceQueueURL,
		Entries:  deleteBatchMessages,
	}

	deleteResult, err := client.DeleteMessageBatch(batchDeleteMessagesInput)
	if err != nil {
		fmt.Println("Failed to send the message to target queue in batch mode")
		fmt.Println("We should abort this, as a sense something is wrong")

		return 0, errors.New("Failed to delete messages after sending to the target queue")
	}

	deletedMsgs = int64(len(deleteResult.Successful))
	fmt.Printf(".") // print a . (dot) for each send OP

	return deletedMsgs, err
}

// MoveMessages Given a moveMessageOptions struct with the proper source and target queue
// along with additional options for fine control migration. And sync and/or move
// all or the partially (see filters options) from source queue to target queue.
func MoveMessages(options *moveMessageOptions) error {
	client, err := sqsClient(options)
	if err != nil {
		return err
	}

	// get Queue's url and related attributes
	sourceQueue := getQueueURL(client, &options.SourceQueueName)
	targetQueue := getQueueURL(client, &options.TargetQueueName)

	options.SourceQueueURL = *sourceQueue.QueueUrl
	options.TargetQueueURL = *targetQueue.QueueUrl

	sourceQueueAttr := getQueueAttributes(client, sourceQueue.QueueUrl)
	targetQueueAttr := getQueueAttributes(client, targetQueue.QueueUrl)

	sourceNumMessages, err := strconv.Atoi(*sourceQueueAttr.Attributes["ApproximateNumberOfMessages"])
	if err != nil {
		fmt.Println("Failed when trying to convert messages from string to integer")
		return errors.New("Failed to retrieve information from source queue")
	}

	targetNumMessages, err := strconv.Atoi(*targetQueueAttr.Attributes["ApproximateNumberOfMessages"])
	if err != nil {
		fmt.Println("Failed when trying to convert messages from string to integer")
		return errors.New("Failed to retrieve information from target queue")
	}

	// if there's no message, our job is done here, let's pack it and go home
	if sourceNumMessages <= 0 {
		fmt.Println(fmt.Sprintf("No messages in Queue: '%s'", *sourceQueue.QueueUrl))
		fmt.Println("No actions to be done here partner")
		return nil
	}

	// Displaying summary of queues
	fmt.Printf("Source Queue '%s' contains %d of messages\n", options.SourceQueueName, sourceNumMessages)
	fmt.Printf("Target Queue '%s' contains %d of messages\n", options.TargetQueueName, targetNumMessages)
	fmt.Printf("Number of the messages to be processed at a time: %d\n", options.BatchSize)
	fmt.Printf("\nStarting migrating, these could take a while ")

	messageInOptions := &sqs.ReceiveMessageInput{
		QueueUrl:              sourceQueue.QueueUrl,
		MaxNumberOfMessages:   aws.Int64(options.BatchSize),
		WaitTimeSeconds:       aws.Int64(options.WaitTimeSeconds),
		VisibilityTimeout:     aws.Int64(options.VisibilityTimeout),
		MessageAttributeNames: []*string{aws.String(sqs.QueueAttributeNameAll)},
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
	}

	// loop over all the message until we are done.
	for {
		receiveResponse, err := client.ReceiveMessage(messageInOptions)
		if err != nil {
			return errors.New("Failed to receive message from source queue")
		}

		if len(receiveResponse.Messages) <= 0 {
			break /* no messages receive, no actions to be done */
		}

		sendResponse, err := sendBatchMessages(client, options, receiveResponse)
		if err != nil {
			return err
		}

		// Delete successfully migrated message from source queue
		if !options.KeepMessageOnSourceQueue && len(sendResponse.Successful) > 0 {

			// if sendBatch does not return any successful message, no message do be deleted
			if len(sendResponse.Successful) <= 0 {
				continue
			}

			// delete messages
			_, err := deleteBatchMessages(client, options, sendResponse)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("\n\n+ Summary:\n")
	fmt.Printf("Successfully sync/move all the messages, my done job is done here partner!\n")

	return nil // return null because, there is no error to be return.
}

// validatedArgs
func validateArgs(options *moveMessageOptions, args []string) error {

	if len(args) != 2 {
		return errors.New("Invalid number of arguments for aws-sqs move command. Use --help for details")
	}

	if options.BatchSize < 0 || options.BatchSize > 10 {
		return errors.New("Invalid number for batch size, The 'batch size' needs to be between 1 and 10")
	}

	if options.WaitTimeSeconds < 0 || options.WaitTimeSeconds > 20 {
		return errors.New("Invalid 'Wait Time Seconds', The 'wait time seconds' needs to be between 0 and 20")
	}

	if options.VisibilityTimeout < 0 || options.VisibilityTimeout > 43200 {
		return errors.New("The 'visibility timeout' cannot be negative, needs to between 0 and 12 hours (43200 seconds)")
	}

	return nil
}

// MoveCommand Return the aws-sqs command in cobra format. Essentially, we should keep the
// logic short and move the heavy logic to another place.
// The following command will provide the ability to move messages from one queue to another
func MoveCommand() *cobra.Command {
	var options moveMessageOptions
	options.ReceiptHandlers = make(map[string]string)

	cmd := &cobra.Command{
		Use:   "move",
		Short: "Move all or part of the messages from on SQS to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := validateArgs(&options, args)
			if err != nil {
				return err
			}

			options.SourceQueueName = args[0]
			options.TargetQueueName = args[1]
			return MoveMessages(&options)
		},
	}

	cmd.PersistentFlags().Int64VarP(&options.BatchSize, "batch-size", "b", 10, "How many messages at a time")
	cmd.PersistentFlags().Int64VarP(&options.WaitTimeSeconds, "wait-time-seconds", "w", 0, "Wait until receive the message")
	cmd.PersistentFlags().Int64VarP(&options.VisibilityTimeout, "visibility-timeout", "t", 10, "Message the visibility")
	cmd.PersistentFlags().BoolVarP(&options.KeepMessageOnSourceQueue, "keep-message-on-source-queue", "k", false, "Whether to keep the message from source queue")
	cmd.PersistentFlags().StringVarP(&options.AwsRegion, "aws-region", "r", "", "define AWS region.")
	cmd.PersistentFlags().StringVarP(&options.AwsProfile, "aws-profile", "p", "", "define AWS profile")
	cmd.PersistentFlags().StringVarP(&options.AwsEndpoint, "aws-endpoint", "e", "", "Define the AWS API endpoint (usually for low-level and testing")

	return cmd
}
