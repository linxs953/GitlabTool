package cmd

import (
	"gitlab/automation/common"

	"github.com/spf13/cobra"
	"github.com/rs/zerolog/log"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "set config",
	Long: `Set config for current user. For example:
mt config set [key] [value].`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) ==  0 {
			log.Print("Args can not be nil")
			return
		}
		if len(args) == 1 {
			key := args[0]
			if key == "init" {
				common.Init()
				log.Print("Init Successfully")
				return
			}else {
				log.Print("lack [value] option")
				return
			}

		}
		key,value := args[0],args[1]
		if key == "" || value == "" {
			log.Print("value of arg can not be nil")
			return
		}
		 common.WriteConfig(key,value,"config.json")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
