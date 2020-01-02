package sqs

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