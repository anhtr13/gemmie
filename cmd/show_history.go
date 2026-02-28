package cmd

import (
	"fmt"
	"log"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/anhtr13/gemmie/internal/model"
)

const MAX_MESSAGE_LENGTH = 96

func format_message(msg string) string {
	text := strings.Join(strings.Split(msg, "\n"), " ")
	if len(text) <= MAX_MESSAGE_LENGTH {
		return text
	}
	return fmt.Sprintf("%s...", text[:MAX_MESSAGE_LENGTH])
}

func chat_history() *cobra.Command {
	return &cobra.Command{
		Use:   "show_history",
		Short: "Show your chat history",
		Long: `Show your chat history
Example: 
	gemmie show_history
`,
		Run: func(cmd *cobra.Command, args []string) {
			chat_history_path, err := model.GetChatPath()
			if err != nil {
				log.Fatal("Cannot get chat history: ", err)
			}
			chat_history := model.ChatHistory{}
			err = chat_history.LoadFromFile(chat_history_path)
			if err != nil {
				chat_history.SaveToFile(chat_history_path)
				fmt.Println("Something wrong with history file, it has been cleaned up.")
				return
			}
			termWidth, _, err := term.GetSize(0)
			if err != nil {
				termWidth = 80
			}

			for _, message := range chat_history {
				role := markdown.Render(fmt.Sprintf("**%s:**", message.Role), termWidth, 0)
				for role[len(role)-1] == '\n' {
					role = role[:len(role)-1]
				}
				fmt.Printf("%s %s\n", string(role), format_message(message.Text))
				if message.Role == "model" {
					fmt.Println()
				}
			}
		},
	}
}
