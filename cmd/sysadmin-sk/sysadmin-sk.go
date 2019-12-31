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
    "github.com/spf13/cobra"
    "github.com/lithammer/dedent"
)

func main() {

    cmd := &cobra.Command{
        Use: "sysadmin-sk",
        Short: "You sysadmin sidekick toolkit",
        SilenceUsage: true,
        Long: dedent.Dedent(`
            SysAdmin Sidekick Toolkit (aka: sysadmin-sk) is a set of feature to help
            Day to Day sysadmin operations. Mostly cloud due to their nature of exposing
            Infrastructure as Service though an API
        `),
    }

    cmd.ResetFlags()

    // registering command to the global command handler
    cmd.AddCommand(NewVersionCmd())
    cmd.AddCommand(NewSqsCmd())
    cmd.Execute()
}