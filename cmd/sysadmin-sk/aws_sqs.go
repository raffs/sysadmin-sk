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
package main

import (
    "errors"
    "github.com/spf13/cobra"
    sqsLibrary "github.com/raffs/sysadmin-sk/services/aws/sqs"
)

/**
 * Aws Global Configuration shared between SQS sub-commands
 */
type AwsGlobalConfig struct {
    // Define the log level for the execution
    LogLevel          int

    // Define which AWS region to connect to the service
    AwsRegion         string

    // Define the AWS API endpoint. Usually this is use for lower-level API call
    // and for testing and/or mocking.
    AwsEndpoint       string

    // Define the AWS profile
    AwsProfile        string
}

var awsConfig AwsGlobalConfig

/**
 * Return the aws-sqs command in cobra format. Essentially, we should keep the
 * logic short and move the heavy logic to another place.
 *
 * The following command will provide the ability to move messages from one queue to another.
 */
func SqsMoveCommand() *cobra.Command {
    var maxNumberOfMessages      int64
    var waitTimeSeconds          int64
    var visibilityTimeout        int64
    var filterString             string
    var filterAttributes         string
    var KeepMessageOnSourceQueue bool

    cmd := &cobra.Command{
        Use: "move",
        Short: "Move all or part of the messages from on SQS to another",
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(args) != 2 {
                return errors.New("Invalid number of arguments for aws-sqs move command. Use --help for details")
            }

            moveOptions := &sqsLibrary.MoveMessageOptions{
               SourceQueueName: args[0],
               TargetQueueName: args[1],
               MaxNumberOfMessages: maxNumberOfMessages,
               WaitTimeSeconds: waitTimeSeconds,
               VisibilityTimeout: visibilityTimeout,
               FilterString: filterString,
               FilterAttributes: filterAttributes,
               KeepMessageOnSourceQueue: KeepMessageOnSourceQueue,
               AwsRegion: awsConfig.AwsRegion,
               AwsEndpoint: awsConfig.AwsEndpoint,
               AwsProfile: awsConfig.AwsProfile,
            }

            // way down we go
            return sqsLibrary.MoveMessages(moveOptions)
        },
    }

    cmd.PersistentFlags().Int64VarP(&maxNumberOfMessages, "max-number-of-messages", "m", 10, "How many messages at a time")
    cmd.PersistentFlags().Int64VarP(&waitTimeSeconds, "wait-time-seconds", "w", 0, "Wait until receive the message")
    cmd.PersistentFlags().Int64VarP(&visibilityTimeout, "visibility-timeout", "t", 10, "Message the visibility")
    cmd.PersistentFlags().BoolVarP(&KeepMessageOnSourceQueue, "keep-message-on-source-queue", "k", false, "Whether to keep the message from source queue")
    cmd.PersistentFlags().StringVar(&filterString, "filter-string", "", "Regex to filter out the message")
    cmd.PersistentFlags().StringVar(&filterAttributes, "filter-attribute", "", "Map key=value to use when filter messages")
    cmd.PersistentFlags().IntVarP(&awsConfig.LogLevel, "log-level", "l", 0, "define the log level when running the script")
    cmd.PersistentFlags().StringVarP(&awsConfig.AwsRegion, "aws-region", "r", "", "define AWS region.")
    cmd.PersistentFlags().StringVarP(&awsConfig.AwsProfile, "aws-profile", "p", "", "define AWS profile")
    cmd.PersistentFlags().StringVarP(&awsConfig.AwsEndpoint, "aws-endpoint", "e", "", "Define the AWS API endpoint (usually for low-level and testing")

    return cmd
}

/**
 * Return the SQS main command from sysadmin sidekick tool
 */
func NewSqsCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "aws-sqs",
        Short: "Provides script for handling AWS SQS operational tasks",
    }

    cmd.ResetFlags()
    cmd.AddCommand(SqsMoveCommand())

    return cmd
}