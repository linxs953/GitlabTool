package cmd

import (
	"gitlab/automation/common"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Print("Args List Can not be emptyed")
			return
		}
		option := args[0]
		switch  {
			case strings.ToUpper(option) == "PROJECTS":listProjects()
			case strings.ToUpper(option) == "GROUPS":listGroups()
			default:log.Printf("Can not support arg %s",option)
		}
	},
}

func listProjects(){
	config,err := common.ReadConfig("config","json")
	if err != nil {
		log.Error().Err(err).Msg("Read Config error")
		return
	}
	currentGroup := config.GROUPS
	if currentGroup == "" {
		log.Print("Get config key nil")
		return
	}
	common.ListAllProjects(config,currentGroup)
}


func listGroups(){
	config,err := common.ReadConfig("config","json")
	if err != nil {
		log.Error().Err(err).Msg("Read Config error")
		return
	}
	common.ListAllGroups(config)

}



func init() {
	rootCmd.AddCommand(listCmd)
}
