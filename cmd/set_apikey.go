package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func set_api_key() *cobra.Command {
	return &cobra.Command{
		Use:   "apikey",
		Short: "Set you Gemini Api key",
		Long: `Set you Gemini Api key
Example: 
	gemmie apikey <your_api_key>
`,
		Run: func(cmd *cobra.Command, args []string) {
			api_key := args[0]
			if api_key == "" {
				fmt.Println("Your key is empty.")
				return
			}
			saved_opts.ApiKey = api_key
			err := saved_opts.SaveToFile(config_path)
			if err != nil {
				log.Fatal("Cannot save to file: ", err)
			}
			fmt.Println("Api key has been saved!")
		},
	}
}
