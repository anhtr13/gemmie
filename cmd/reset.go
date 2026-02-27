package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/anhtr13/gemmie/internal/model"
)

func reset() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset default options",
		Long: `Reset default options
Example: 
	gemmie reset
`,
		Run: func(cmd *cobra.Command, args []string) {
			default_opts := model.BaseOptions()
			default_opts.ApiKey = saved_opts.ApiKey
			err := default_opts.SaveToFile(config_path)
			if err != nil {
				log.Fatal("Cannot save to file: ", err)
			}
			fmt.Println("Options have been reset")
		},
	}
}
