// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/kuops/ngconf/pkg/hck"
	"github.com/spf13/cobra"
)

var hckOptions hck.Options

// hckCmd represents the hck command
var hckCmd = &cobra.Command{
	Use:                   "hck --upstream=upstream_name --domain=healthcheck_domain [--uri=uri] --hckfile=configfile",
	DisableFlagsInUseLine: true,
	Long:                  `add healthcheck to healthcheck.conf`,
	Example:               `ngconf hck  --upstream=demo --domain=demo.kuops.com --uri=/healthCheck --hckfile=healthcheck.conf`,
	Run: func(cmd *cobra.Command, args []string) {
		err := hck.AppendHck(&hckOptions)
		if err != nil {
			fmt.Printf("err is %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(hckCmd)

	hckCmd.PersistentFlags().StringVar(&hckOptions.Upstream, "upstream", "", "upstram is upstream_name same common as projectname.")
	hckCmd.PersistentFlags().StringVar(&hckOptions.Domain, "domain", "", "domain is healthcheck header Host field")
	hckCmd.PersistentFlags().StringVar(&hckOptions.HckFile, "hckfile", "", "hckfile is healthcheck.conf file path")
	hckCmd.PersistentFlags().StringVar(&hckOptions.URI, "uri", "", "hckfile is healthcheck.conf file path")
}
