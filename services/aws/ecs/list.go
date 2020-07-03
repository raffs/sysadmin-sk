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
package ecs

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
)

// ListOptions defines the options used on the `aws ecs list` command
type listOptions struct {

	// cluster name to list services information
	clusterName string `type:"string" required:"true"`

	// filter out specific services or label information
	filter string `type:"string" required:"false"`

	// the format to output the string
	format string `type:"string" required:"false"`

	// filter out specific services or label information
	listTaskInstances bool `type:"string" required:"false"`

	// Define which AWS region to connect to the service
	awsRegion string `type:"string" required:"false"`

	// Define the AWS API endpoint. Usually this is use for lower-level API call
	// and for testing and/or mocking.
	awsEndpoint string `type:"string" required:"false"`

	// Define the AWS profile
	awsProfile string `type:"string" required:"false"`
}

// ecsClient Return a AWS ECS client with an open session.
func ecsClient(options *listOptions) (*ecs.ECS, error) {
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

	return ecs.New(session), nil
}

// listServices
func listServices(options *listOptions) error {
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

	client, err := ecsClient(options)
	if err != nil {
		return err
	}

	listServicesInput := &ecs.ListServicesInput{
		Cluster:    aws.String(options.clusterName),
		MaxResults: aws.Int64(100),
	}

	// loop until there's no more page
	for serviceList, err := client.ListServices(listServicesInput); ; {
		if err != nil {
			return err
		}

		lenServices := len(serviceList.ServiceArns)
		for i := 0; i < lenServices; i += 10 {
			upperBound := i + 10
			if upperBound >= lenServices {
				upperBound = i + (lenServices - i)
			}

			describeInput := &ecs.DescribeServicesInput{
				Cluster:  aws.String(options.clusterName),
				Services: serviceList.ServiceArns[i:upperBound],
			}

			servicesDescription, err := client.DescribeServices(describeInput)
			if err != nil {
				return err
			}

			for _, service := range servicesDescription.Services {
				template.Execute(tabw, service)
				tabw.Flush()
			}
		}

		if serviceList.NextToken == nil {
			break
		}

		listServicesInput.NextToken = serviceList.NextToken
	}

	if err != nil {
		return err
	}

	fmt.Println("")
	return nil
}

func validateArgs(options *listOptions, args []string) error {
	if len(args) != 1 {
		return errors.New("Invalid number of arguments for aws-ecs list command. Use --help for details")
	}

	if options.filter != "" {
		fmt.Println("WARN: the filter is not yet use, ignoring the filters")
	}

	return nil
}

// ListCommand Return the aws-sqs command in cobra format. Essentially, we should keep the
// logic short and move the heavy logic to another place.
// The following command will provide the ability to move messages from one queue to another.
func ListCommand() *cobra.Command {
	var options listOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List information about ECS service and Task Definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := validateArgs(&options, args)
			if err != nil {
				return err
			}

			options.clusterName = args[0]
			return listServices(&options)
		},
	}

	// default output template
	defaultTempl := "\n{{.SchedulingStrategy}}\t{{.Status}}\t{{.RunningCount}}/{{.DesiredCount}}\t{{.ServiceName}}"

	cmd.PersistentFlags().StringVarP(&options.awsRegion, "aws-region", "r", "", "define AWS region.")
	cmd.PersistentFlags().StringVarP(&options.awsProfile, "aws-profile", "p", "", "define AWS profile")
	cmd.PersistentFlags().StringVarP(&options.awsEndpoint, "aws-endpoint", "e", "", "Define the AWS API endpoint (usually for low-level and testing")
	cmd.PersistentFlags().StringVarP(&options.filter, "filter", "f", "", "Filter out service")
	cmd.PersistentFlags().StringVarP(&options.format, "format", "", defaultTempl, "Display the format")
	cmd.PersistentFlags().BoolVarP(&options.listTaskInstances, "list-instance", "", false, "Filter out service")
	return cmd
}
