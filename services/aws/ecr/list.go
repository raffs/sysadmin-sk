package ecr

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecr"
	"html/template"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

// ListOptions defines the options used on the `aws ecs list` command
type listOptions struct {

	repositoryName string `type:"string" required:"true"`

	registryId string `type:"string" required:"false"`

	tagStatus string `type:"string" required:"true"`

	// the format to output the string
	format string `type:"string" required:"false"`

	// Define which AWS region to connect to the service
	awsRegion string `type:"string" required:"false"`

	// Define the AWS API endpoint. Usually this is use for lower-level API call
	// and for testing and/or mocking.
	awsEndpoint string `type:"string" required:"false"`

	// Define the AWS profile
	awsProfile string `type:"string" required:"false"`
}

// ecsClient Return a AWS ECS client with an open session.
func ecrClient(options *listOptions) (*ecr.ECR, error) {
	sessionOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		// aws configuration
		Config: aws.Config{
			Region:   aws.String(options.awsRegion),
			Endpoint: aws.String(options.awsEndpoint),
		},
	}

	session, err := session.NewSessionWithOptions(sessionOpts)
	if err != nil {
		fmt.Println("Unable to create new session with AWS with error: ", err.Error())
		return nil, errors.New("Unable to initialize AWS session")
	}

	return ecr.New(session), nil
}

// listServices
func listImages(options *listOptions) error {
	tabw := new(tabwriter.Writer)
	tabw.Init(os.Stdout, 20, 1, 3, ' ', 0)
	defer tabw.Flush()

	template, err := template.New("").Parse(options.format)
	if err != nil {
		fmt.Println("Not able to parse template: ", options.format)
		return errors.New("Failed to parse the above template")
	}

	columns := strings.Replace(options.format, "{{", "", -1)
	columns = strings.Replace(columns, "}}", "", -1)
	columns = strings.Replace(columns, ".", "", -1)
	columns = strings.Replace(columns, "\n", "", -1)
	fmt.Fprintf(tabw, "%v", columns)

	client, err := ecrClient(options)
	if err != nil {
		return err
	}


	listImagesFilter := &ecr.ListImagesFilter{TagStatus: &options.tagStatus}


	listImagesInput := &ecr.ListImagesInput{
		Filter:    listImagesFilter,
		MaxResults: aws.Int64(100),
		RepositoryName: aws.String(options.repositoryName),
	}

	if options.repositoryName != "" && options.registryId != "" {
		listImagesInput = &ecr.ListImagesInput{
			Filter:    listImagesFilter,
			MaxResults: aws.Int64(100),
			RepositoryName: aws.String(options.repositoryName),
			RegistryId: aws.String(options.registryId),
		}
	}




	// loop until there's no more page
	for imageList, err := client.ListImages(listImagesInput); ; {
		if err != nil {
			return err
		}

		lenImages := len(imageList.ImageIds)
		for i := 0; i < lenImages; i += 10 {
			upperBound := i + 10
			if upperBound >= lenImages {
				upperBound = i + (lenImages - i)
			}

			describeInput := &ecr.DescribeImagesInput{
				ImageIds: imageList.ImageIds[i:upperBound],
			}

			imageDescription, err := client.DescribeImages(describeInput)
			if err != nil {
				return err
			}

			for _, image := range imageDescription.ImageDetails {
				template.Execute(tabw, image)
				tabw.Flush()
			}
		}

		if imageList.NextToken == nil {
			break
		}

		imageList.NextToken = imageList.NextToken
	}

	if err != nil {
		return err
	}

	fmt.Println("")
	return nil
}

func validateArgs(options *listOptions, args []string) error {
	if len(args) != 1 {
		return errors.New("Invalid number of arguments for aws-ecr list images command. Use --help for details")
	}

	return nil
}

// ListCommand Return the aws-sqs command in cobra format. Essentially, we should keep the
// logic short and move the heavy logic to another place.
// The following command will provide the ability to move messages from one queue to another.
func ListImagesCommand() *cobra.Command {
	var options listOptions

	cmd := &cobra.Command{
		Use:   "listImages",
		Short: "List information about ECR images in a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := validateArgs(&options, args)
			if err != nil {
				return err
			}

			return listImages(&options)
		},
	}

	// default output template
	defaultTempl := "\n{{.ImageId}}\t{{.Status}}\t{{.RunningCount}}/{{.DesiredCount}}\t{{.ServiceName}}"

	cmd.PersistentFlags().StringVarP(&options.awsRegion, "aws-region", "r", "", "define AWS region.")
	cmd.PersistentFlags().StringVarP(&options.awsProfile, "aws-profile", "p", "", "define AWS profile")
	cmd.PersistentFlags().StringVarP(&options.awsEndpoint, "aws-endpoint", "e", "", "Define the AWS API endpoint (usually for low-level and testing")
	cmd.PersistentFlags().StringVarP(&options.repositoryName, "repository-name", "f", "", "Name of the ECR repository to scan")
	cmd.PersistentFlags().StringVarP(&options.format, "format", "", defaultTempl, "Display the format")
	cmd.PersistentFlags().StringVarP(&options.registryId, "registry-id", "", "", "The AWS ECR registry to use")
	cmd.PersistentFlags().StringVarP(&options.tagStatus, "tag-status", "", "", "The tag status with which to filter your ListImages results. You can filter results on the following vars TAGGED/UNTAGGED/ANY")
	return cmd
}
