package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func list_options() *cobra.Command {
	return &cobra.Command{
		Use:   "list_opts",
		Short: "List all default options",
		Long: `List all default options
Example:
	gemmie list_opts
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Model name: %s\n", saved_opts.ModelName)
			fmt.Printf("Language: %s\n", saved_opts.Language)
			fmt.Printf("Temperature: %.2f\n", saved_opts.Temperature)
			fmt.Printf("MaxOutputTokens: %d\n", saved_opts.MaxOutputToken)
			fmt.Printf("TopP: %.2f\n", saved_opts.TopP)
			fmt.Printf("TopK: %.2f\n", saved_opts.TopK)
			fmt.Printf("Chat mode: %t\n", saved_opts.ChatMode)
			fmt.Printf("Stream mode: %t\n", saved_opts.StreamMode)
		},
	}
}
