package model

import (
	"context"
	"fmt"
	"iter"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"golang.org/x/term"
	"google.golang.org/genai"

	"github.com/anhtr13/gemmie/internal/ui"
)

const (
	APP_DIR        = ".gemmie"
	ROLLBACK_MODEL = "gemini-2.5-flash"
)

type Model struct {
	client      *genai.Client
	config      *genai.GenerateContentConfig
	chat_mode   bool
	stream_mode bool
	model_name  string
	language    string
}

func NewModel(
	client *genai.Client,
	model_name string,
	language string,
	temperature float32,
	max_output_token int32,
	top_p float32,
	top_k float32,
	chat bool,
	stream bool,
) *Model {
	config := &genai.GenerateContentConfig{
		Temperature:     &temperature,
		TopP:            &top_p,
		TopK:            &top_k,
		MaxOutputTokens: max_output_token,
	}
	return &Model{
		client:      client,
		config:      config,
		language:    language,
		model_name:  model_name,
		chat_mode:   chat,
		stream_mode: stream,
	}
}

func (m *Model) stream_answer(stream iter.Seq2[*genai.GenerateContentResponse, error]) (string, error) {
	answer := []string{}
	for chunk, err := range stream {
		if err != nil {
			return "", err
		}
		part := chunk.Candidates[0].Content.Parts[0]
		text := part.Text
		fmt.Print(text)
		answer = append(answer, text)
	}
	return strings.Join(answer, ""), nil
}

func (m *Model) render_markdown_answer(res *genai.GenerateContentResponse) (string, error) {
	answer := res.Text()

	termWidth, _, err := term.GetSize(0)
	if err != nil {
		termWidth = 80
	}

	result := markdown.Render(answer, termWidth, 0)
	fmt.Println(string(result))
	return answer, nil
}

func (m *Model) GenAnswer(ctx context.Context, prompt string) error {
	termWidth, _, err := term.GetSize(0)
	if err != nil {
		termWidth = 80
	}
	fmt.Print(string(markdown.Render(fmt.Sprintf("**Model**: %s\n", m.model_name), termWidth, 0)))
	fmt.Println()

	prompt = fmt.Sprintf("%s?\nResponse in %s.", prompt, m.language)

	if m.chat_mode {
		chat_history_path, err := GetChatPath()
		if err != nil {
			return err
		}
		chat_history := ChatHistory{}
		err = chat_history.LoadFromFile(chat_history_path)
		if err != nil {
			chat_history.SaveToFile(chat_history_path)
		}

		history := []*genai.Content{}
		for _, c := range chat_history {
			history = append(history, genai.NewContentFromText(c.Text, c.Role))
		}

		chat, err := m.client.Chats.Create(ctx, m.model_name, m.config, history)
		if err != nil {
			return err
		}

		if m.stream_mode {
			stream := chat.SendMessageStream(ctx, genai.Part{Text: prompt})
			answer, err := m.stream_answer(stream)
			if err == nil {
				chat_history = append(chat_history, Message{Text: prompt, Role: genai.RoleUser})
				chat_history = append(chat_history, Message{Text: answer, Role: genai.RoleModel})
				chat_history.SaveToFile(chat_history_path)
			}
			return err
		}

		res, err := get_generate_content_response_with_loader(func() (*genai.GenerateContentResponse, error) {
			return chat.SendMessage(ctx, genai.Part{Text: prompt})
		})
		if err != nil {
			return err
		}

		answer, err := m.render_markdown_answer(res)
		if err == nil {
			chat_history = append(chat_history, Message{Text: prompt, Role: genai.RoleUser})
			chat_history = append(chat_history, Message{Text: answer, Role: genai.RoleModel})
			chat_history.SaveToFile(chat_history_path)
		}
		return err
	}

	if m.stream_mode {
		stream := m.client.Models.GenerateContentStream(ctx, m.model_name, genai.Text(prompt), m.config)
		_, err = m.stream_answer(stream)
		return err
	}

	res, err := get_generate_content_response_with_loader(func() (*genai.GenerateContentResponse, error) {
		return m.client.Models.GenerateContent(ctx, m.model_name, genai.Text(prompt), m.config)
	})
	if err != nil {
		return err
	}

	_, err = m.render_markdown_answer(res)
	return err
}

// Roll back to `ROLLBACK_MODEL` if current model failed
func (m *Model) GenAnswerRollBack(ctx context.Context, prompt string) error {
	if m.model_name == ROLLBACK_MODEL {
		return m.GenAnswer(ctx, prompt)
	}

	err := m.GenAnswer(ctx, prompt)
	if err != nil {
		fmt.Println(err.Error())
		m.model_name = ROLLBACK_MODEL
		return m.GenAnswer(ctx, prompt)
	}
	return nil
}

func get_generate_content_response_with_loader(
	f func() (*genai.GenerateContentResponse, error),
) (*genai.GenerateContentResponse, error) {
	loader := ui.NewLoader([]string{"[    ]", "[=   ]", "[==  ]", "[=== ]", "[====]", "[ ===]", "[  ==]", "[   =]"}, 150)
	loader.Start()
	defer loader.Stop()
	return f()
}
