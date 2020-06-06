package cmd

import (
	"fmt"
	"gitlab/automation/common"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cicdCmd = &cobra.Command{
	Use:   "cicd",
	Short: "Open browser to get current project cicd",
	Long: `Open browser to get current project cicd. For example:
	mt cicd`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("cicd called")
		config,err := common.ReadConfig("config","json")
		if err != nil {
			log.Error().Err(err).Msg("Read Config error")
			return
		}
		group,project,cicd := config.GROUPS,config.PROJECT,config.CICD
		if group == "" || project == "" || cicd == "" {
			log.Print("Read Config Key error")
			return
		}
		cicdURL := fmt.Sprintf(cicd,group + "/" + project)
		execCmd := exec.Command("open",cicdURL)
		err = execCmd.Run()
		if err != nil {
			log.Error().Err(err).Msg("Run System Command Error")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(cicdCmd)
}
