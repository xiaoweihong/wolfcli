/*
Copyright © 2020 xiaoweihong

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
	"net"
	"os"
	"time"
	"wolfcli/global"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wolfcli",
	Short: "a",
	Long:  `战狼小工具`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) {
	//	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IPVarP(&global.IP, "ip", "i", net.IPv4(127, 0, 0, 1), "战狼服务器ip")
	rootCmd.PersistentFlags().IPVar(&global.DBIP, "dip", net.IPv4(127, 0, 0, 1), "战狼服务器数据库ip")
	rootCmd.PersistentFlags().IPVar(&global.PicIP, "pip", net.IPv4(127, 0, 0, 1), "图片服务器ip")
	rootCmd.PersistentFlags().StringVar(&global.PicPort, "pport", "9333", "图片服务器端口")
	rootCmd.PersistentFlags().DurationVar(&global.HttpTimeOut, "timeout", 5*time.Second, "http超时时间(单位:秒)")
}
