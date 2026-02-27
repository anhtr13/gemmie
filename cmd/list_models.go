package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"google.golang.org/genai"
)

func list_models() *cobra.Command {
	return &cobra.Command{
		Use:   "list-models",
		Short: "List all available models",
		Long: `List all available models
Example: 
	gemmie list-models
`,
		Run: func(cmd *cobra.Command, args []string) {
			if client == nil {
				log.Fatal("Api key not found, cannot create client!")
			}

			models, err := client.Models.List(ctx, &genai.ListModelsConfig{})
			if err != nil {
				log.Fatal("Cannot list models: ", err)
			}

			termWidth, _, err := term.GetSize(0)
			if err != nil {
				termWidth = 80
			}
			fmt.Print(string(markdown.Render("**Models that support generateContent:**", termWidth, 0)))

			for _, m := range models.Items {
				if slices.Contains(m.SupportedActions, "generateContent") {
					fmt.Println(strings.TrimPrefix(m.Name, "models/"))
				}
			}

			fmt.Println()
			fmt.Print(string(markdown.Render("**Models that support embedContent:**", termWidth, 0)))
			for _, m := range models.Items {
				if slices.Contains(m.SupportedActions, "embedContent") {
					fmt.Println(strings.TrimPrefix(m.Name, "models/"))
				}
			}
		},
	}
}
