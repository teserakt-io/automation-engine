// Copyright 2020 Teserakt AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"github.com/spf13/cobra"

	"github.com/teserakt-io/automation-engine/internal/cli"
)

// Command defines a cli Command
type Command interface {
	CobraCmd() *cobra.Command
}

type rootCommandFlags struct {
	Endpoint string
	Cert     string
}

type rootCommand struct {
	cobraCmd *cobra.Command
	flags    rootCommandFlags
}

var _ Command = &rootCommand{}

// NewRootCommand creates and configure a new cli root command
func NewRootCommand(c2aeClientFactory cli.APIClientFactory, version string) Command {
	rootCmd := &rootCommand{}

	listCmd := NewListCommand(c2aeClientFactory)
	createCmd := NewCreateCommand(c2aeClientFactory)
	addTriggerCmd := NewAddTriggerCommand(c2aeClientFactory)
	addTargetCmd := NewAddTargetCommand(c2aeClientFactory)
	showCmd := NewShowCommand(c2aeClientFactory)
	deleteCmd := NewDeleteCommand(c2aeClientFactory)

	completionCmd := NewCompletionCommand(rootCmd)

	cobraCmd := &cobra.Command{
		Use:                    "c2ae-cli",
		BashCompletionFunction: completionCmd.GenerateCustomCompletionFuncs(),
		Version:                version,
		SilenceUsage:           true,
		SilenceErrors:          true,
	}

	cobraCmd.PersistentFlags().StringVarP(
		&rootCmd.flags.Endpoint,
		cli.EndpointFlag,
		"e",
		"127.0.0.1:5556", "url to the c2ae grpc api",
	)

	cobraCmd.PersistentFlags().StringVarP(
		&rootCmd.flags.Cert,
		cli.CertFlag,
		"c",
		"configs/c2ae-cert.pem", "path to the c2ae grpc api certificate",
	)

	cobraCmd.AddCommand(
		listCmd.CobraCmd(),
		createCmd.CobraCmd(),
		addTriggerCmd.CobraCmd(),
		addTargetCmd.CobraCmd(),
		showCmd.CobraCmd(),
		deleteCmd.CobraCmd(),

		// Autocompletion script generation command
		completionCmd.CobraCmd(),
	)

	cobraCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	rootCmd.cobraCmd = cobraCmd

	return rootCmd
}

func (c *rootCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}
