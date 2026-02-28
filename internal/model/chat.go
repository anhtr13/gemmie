package model

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/genai"
)

const CHAT_HISTORY = "chat_history.json"

func GetChatPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	app_dir := fmt.Sprintf("%s/%s", home, APP_DIR)
	err = os.MkdirAll(app_dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	file_path := fmt.Sprintf("%s/%s", app_dir, CHAT_HISTORY)
	return file_path, nil
}

type Message struct {
	Role genai.Role `json:"role"`
	Text string     `json:"content"`
}

type ChatHistory []Message

func (chat *ChatHistory) LoadFromFile(file_path string) error {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, chat)
}

func (chat *ChatHistory) SaveToFile(file_path string) error {
	data, err := json.MarshalIndent(chat, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(file_path, data, 0644)
}
