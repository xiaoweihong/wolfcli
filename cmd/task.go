package cmd

import (
	"fmt"
	"wolfcli/controller"

	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "任务的增、删、改、查",
	Long:  `任务的增、删、改、查`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("task called")
	},
}

var taskSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "任务状态保存",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("任务状态保存")
		token := controller.GetToken()
		controller.TaskStatusSave(token)
	},
}

var taskResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "任务状态恢复",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("任务状态恢复")
		token := controller.GetToken()
		controller.TaskResume(token)
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskSaveCmd)
	taskCmd.AddCommand(taskResumeCmd)
}
