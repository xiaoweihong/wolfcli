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
	"fmt"
	"github.com/spf13/cobra"
	"wolfcli/controller"
	"wolfcli/global"
)

// sensorCmd represents the sensor command
var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "设备的增、删、改、查",
	Long:  `设备的增、删、改、查`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	token := controller.GetToken()
	//	zap.L().Info(token)
	//},
}

var sensorSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "查询设备",
	Long:  `查询设备`,
	Run: func(cmd *cobra.Command, args []string) {
		token := controller.GetToken()
		controller.GetSensors(token, global.OrgId)
	},
}

var sensorAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加设备",
	Long:  `添加设备`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("添加设备")
	},
}

func init() {
	rootCmd.AddCommand(sensorCmd)
	sensorCmd.AddCommand(sensorSearchCmd)
	sensorCmd.AddCommand(sensorAddCmd)

	sensorSearchCmd.Flags().StringVar(&global.OrgId, "orgid", "0000", "组织id")

}
