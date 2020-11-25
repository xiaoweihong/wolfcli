package cmd

import (
	"github.com/spf13/cobra"
	"wolfcli/controller"
	"wolfcli/global"
)

// seaweedfsCmd represents the seaweedfs command
var seaweedfsCmd = &cobra.Command{
	Use:   "seaweedfs",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		day := cmd.Flag("day").Value.String()
		controller.GetVolumeInfo(day)
	},
}

// seaweedfsDel represents the seaweedfsDel command
var seaweedfsDel = &cobra.Command{
	Use:   "delete",
	Short: "根据id删除volume",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		id := cmd.Flag("id").Value.String()
		controller.DeleteVolumeById(id)
	},
}

// seaweedfsDel represents the seaweedfsDel command
var seaweedfsDelMore = &cobra.Command{
	Use:   "deleteByTime",
	Short: "根据时间(天)删除volume",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//time := cmd.Flag("time").Value.String()
		controller.DeleteVolumeByTime()
	},
}

func init() {
	rootCmd.AddCommand(seaweedfsCmd)
	seaweedfsCmd.AddCommand(seaweedfsDel)
	seaweedfsCmd.AddCommand(seaweedfsDelMore)
	seaweedfsCmd.PersistentFlags().StringVar(&global.Day, "day", "180", "需要删除的图片的天数")

	seaweedfsDel.Flags().Int("id", 0, "volumeID")
}
