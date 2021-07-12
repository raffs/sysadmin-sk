/**
 * This file is part of the Sysadmin Sidekick Toolkit (Sysadmin-SK) (https://github.com/raffs/sysadmin-sk).
 * Copyright (c) 2019 Rafael Oliveira Silva
 * Copyright (c) 2021 Mainak Dhar
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
	"github.com/spf13/cobra"

	ecrLibrary "github.com/raffs/sysadmin-sk/services/aws/ecr"
)

// NewAwsEcsCommand return the SQS main command from sysadmin sidekick tool
func NewAwsEcrCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws-ecr",
		Short: "Provides features for working with AWS ECR",
	}

	cmd.ResetFlags()
	cmd.AddCommand(ecrLibrary.ListImagesCommand())
	return cmd
}

