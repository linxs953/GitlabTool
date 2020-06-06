
package cmd

import (
	"gitlab/automation/common"
	"strconv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New a issue by  cli",
	Long: `New a issue by cli. For example:
mt new [issue-title] [issue-template-filename] [assignee-name].`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Print("Args List Length Nil")
			return
		}
		config,err := common.ReadConfig("config","json")
		if err != nil {
			log.Error().Err(err).Msg("Read Config Error")
			return
		}
		token := config.TOKEN
		if token == "" {
			log.Print("Get Config Error")
			return
		}
		filename,issueTitle,assignee := args[1],args[0],args[2]
		assignID := common.GetAssigneeID(config,assignee)
		projectID := common.GetProjectID(config)
		description := common.GetDescription(config,filename)
		if description == "" {
			log.Print("Description is nil")
			return
		}
		if assignID == -1  {
			log.Print("Get assignee userid error")
			return
		}
		if projectID == -1 {
			log.Print("Get projectid error")
			return
		}
		common.NewIssue(token,strconv.FormatInt(projectID,10),issueTitle,description,strconv.FormatInt(assignID,10),config)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
