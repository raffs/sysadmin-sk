/**
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
	"github.com/spf13/cobra"

	ecsLibrary "github.com/raffs/sysadmin-sk/services/aws/ecs"
)

// NewAwsEcsCommand return the SQS main command from sysadmin sidekick tool
func NewAwsEcsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws-ecs",
		Short: "Provides features for working with AWS ECS",
	}

	cmd.ResetFlags()
	cmd.AddCommand(ecsLibrary.ListCommand())
	return cmd
}
