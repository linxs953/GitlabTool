package cmd

import (
	"fmt"
	"gitlab/automation/common"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "sw",
	Short: "Change Current Group or Project",
	Long: `Change Current Group or Project. For example:
mt sw group [groupname]
mt sw project [projectname]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Print("Args List length can not be zero")
			return
		}
		optionType,opsName := args[0],args[1]
		switch  {
		case strings.ToUpper(optionType) == "GROUP":change("GROUPS",opsName)
		case strings.ToUpper(optionType) == "PROJECT":change("PROJECT",opsName)
		default:log.Printf("Can not support arg %s",optionType);return
		}
	},
}

func change(opsType string,ops string ) {
	common.WriteConfig(strings.ToUpper(opsType),ops,"config.json")
	if strings.ToUpper(opsType) == "GROUPS" {
		fmt.Printf("Switch to group %s\n",ops)
	}else if strings.ToUpper(opsType) == "PROJECT" {
		fmt.Printf("Switch to project %s\n",ops)
	}else {
		log.Printf("Can not Support key %s",strings.ToUpper(opsType))
		return
	}
}


func init() {
	rootCmd.AddCommand(switchCmd)
}
