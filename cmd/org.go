/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"wolfcli/controller"
	"wolfcli/global"

	"github.com/spf13/cobra"
)

// orgCmd represents the org command
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "组织的增、删、改、查",
	Long:  `组织的增、删、改、查`,
}

var orgSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "查询组织",
	Long:  `查询组织`,
	Run: func(cmd *cobra.Command, args []string) {
		token := controller.GetToken()
		controller.GetOrgs(token, global.OrgId)
	},
}

var orgAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加组织",
	Long:  `添加组织`,
	Run: func(cmd *cobra.Command, args []string) {
		token := controller.GetToken()
		controller.AddOrg(token, global.OrgName, "0000")
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(orgSearchCmd)
	orgCmd.AddCommand(orgAddCmd)
	orgSearchCmd.Flags().StringVar(&global.OrgId, "orgid", "0000", "组织id")
	orgAddCmd.Flags().StringVar(&global.OrgName, "orgname", "", "组织名称")
}
