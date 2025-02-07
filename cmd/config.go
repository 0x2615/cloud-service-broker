// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cloudfoundry/cloud-service-broker/utils"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Show system configuration",
		Long: `
The GCP Service Broker can be configured using both environment variables and
configuration files.

It accepts configuration files in YAML, JSON, TOML and Java properties formats.
You can specify a configuration file to read using the --config argument.

You can also specify configurations via environment variables.
The environment variables take the form GSB_<property> where property is the
same name as in the config file transformed to be in upper case and with all
dots replaced with underscores. For example:

    GSB_DB_USER_NAME == db.user.name == {"db":{"user":{"name":""}}}

Some older environment variables don't follow this format are aliased so either
format will work.

Precedence is in the order:

  environment vars > configuration > defaults

You can show the known configuration values using:

  ./cloud-service-broker config show
`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show the config",
		Long:  `Show the current configuration settings.`,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrettyPrintOrExit(viper.AllSettings())
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "write",
		Short: "Write configuration to a file",
		Long: `Write configuration to a file in a specified format. Valid extensions are:

 * .json
 * .yml
 * .toml
 * .properties

You can combine this command with the --config flag to translate configurations:

  GSB_DB_PASSWORD=pass cloud-service-broker --config in.json config write out.toml

out.toml:

  [api]
    port = "3340"

  [db]
    name = "servicebroker"
    password = "pass"
    port = "3306"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return viper.WriteConfigAs(args[0])
		},
	})
}
