package model

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	app_dir := fmt.Sprintf("%s/%s", home, ".gemmie")
	err = os.MkdirAll(app_dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	file_path := fmt.Sprintf("%s/%s", app_dir, "options.json")
	return file_path, nil
}

type ModelOptions struct {
	StreamMode     bool    `json:"stream_mode"`
	MaxOutputToken int32   `json:"max_output_token"`
	Temperature    float32 `json:"temperature"`
	TopK           float32 `json:"top_k"`
	TopP           float32 `json:"top_p"`
	ModelName      string  `json:"model_name"`
	Language       string  `json:"language"`
	ApiKey         string  `json:"api_key"`
}

func BaseOptions() *ModelOptions {
	return &ModelOptions{
		ModelName:      ROLLBACK_MODEL,
		Language:       "English",
		Temperature:    1,
		TopP:           0.95,
		TopK:           40,
		MaxOutputToken: 8192,
		StreamMode:     false,
	}
}

func (opt *ModelOptions) LoadFromFile(file_path string) error {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, opt)
}

func (opt *ModelOptions) SaveToFile(file_path string) error {
	data, err := json.MarshalIndent(opt, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(file_path, data, 0644)
}
