package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/anhtr13/gemmie/internal/model"
)

func clear_chat() *cobra.Command {
	return &cobra.Command{
		Use:   "clear_history",
		Short: "Clear your chat history",
		Long: `Clear your chat history
Example: 
	gemmie clear_history
`,
		Run: func(cmd *cobra.Command, args []string) {
			chat_history_path, err := model.GetChatPath()
			if err != nil {
				log.Fatal("Cannot get chat history: ", err)
			}

			chat_history := model.ChatHistory{}
			err = chat_history.SaveToFile(chat_history_path)
			if err != nil {
				log.Fatal("Cannot save chat history: ", err)
			}
			fmt.Println("Chat history has been cleaned up.")
		},
	}
}
