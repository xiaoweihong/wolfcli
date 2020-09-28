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
	"wolfcli/initialize"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理图片数据",
	Long:  `清理图片数据`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initDb()
	},
	Run: func(cmd *cobra.Command, args []string) {
		controller.DeleteFaces(global.StartTime, global.EndTime)
	},
}

func initDb() {
	initialize.GormPostgreSql()
	db, err := global.Db.DB()
	if err != nil {
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVar(&global.PgUsername, "postgres", "postgres", "数据库用户名")
	cleanCmd.Flags().StringVar(&global.PgPassword, "password", "Zstvgcs@9102", "数据库密码")
	cleanCmd.Flags().StringVar(&global.DbName, "dbname", "deepface_v5", "数据库名称")
	cleanCmd.Flags().Int64Var(&global.StartTime, "start", 0, "图片删除开始时间")
	cleanCmd.Flags().Int64Var(&global.EndTime, "end", 0, "图片删除结束时间")
	cleanCmd.Flags().Float64Var(&global.ThresHold, "threshold", 0.1, "清理阈值")
	cleanCmd.Flags().IntVar(&global.ParamNum, "num", 1000, "并发数")
}
