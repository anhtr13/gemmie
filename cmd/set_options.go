package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func set_options() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Set you default app options",
		Long: `Set your default app options
Example: 
	gemmie config --lang=Vietnamese
`,
		Run: func(cmd *cobra.Command, args []string) {
			changed := false
			if saved_opts.Language != user_opts.Language {
				saved_opts.Language = user_opts.Language
				changed = true
			}
			if saved_opts.ModelName != user_opts.ModelName {
				saved_opts.ModelName = user_opts.ModelName
				changed = true
			}
			if saved_opts.Temperature != user_opts.Temperature {
				saved_opts.Temperature = user_opts.Temperature
				changed = true
			}
			if saved_opts.MaxOutputToken != user_opts.MaxOutputToken {
				saved_opts.MaxOutputToken = user_opts.MaxOutputToken
				changed = true
			}
			if saved_opts.TopP != user_opts.TopP {
				saved_opts.TopP = user_opts.TopP
				changed = true
			}
			if saved_opts.TopK != user_opts.TopK {
				saved_opts.TopK = user_opts.TopK
				changed = true
			}
			if saved_opts.ChatMode != user_opts.ChatMode {
				saved_opts.ChatMode = user_opts.ChatMode
				changed = true
			}
			if saved_opts.StreamMode != user_opts.StreamMode {
				saved_opts.StreamMode = user_opts.StreamMode
				changed = true
			}
			if !changed {
				fmt.Println("Nothing changed!")
				return
			}

			err := saved_opts.SaveToFile(config_path)
			if err != nil {
				log.Fatal("Cannot save to file: ", err)
			}
			fmt.Println("Options have been saved!")
		},
	}
}
