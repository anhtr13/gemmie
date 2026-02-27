package model

import (
	"context"
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
	"golang.org/x/term"
	"google.golang.org/genai"

	"github.com/anhtr13/gemmie/internal/ui"
)

const ROLLBACK_MODEL = "gemini-2.5-flash"

type Model struct {
	client      *genai.Client
	config      *genai.GenerateContentConfig
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
		stream_mode: stream,
	}
}

func (m *Model) stream_answer(ctx context.Context, prompt string) error {
	stream := m.client.Models.GenerateContentStream(ctx, m.model_name, genai.Text(prompt), m.config)
	for chunk, err := range stream {
		if err != nil {
			return err
		}
		part := chunk.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)
	}
	return nil
}

func (m *Model) render_markdown_answer(ctx context.Context, prompt string) error {
	loader := ui.NewLoader([]string{"[    ]", "[=   ]", "[==  ]", "[=== ]", "[====]", "[ ===]", "[  ==]", "[   =]"}, 200)
	loader.Start()

	res, err := m.client.Models.GenerateContent(ctx, m.model_name, genai.Text(prompt), m.config)
	if err != nil {
		return err
	}

	content := res.Text()

	termWidth, _, err := term.GetSize(0)
	if err != nil {
		termWidth = 80
	}

	result := markdown.Render(content, termWidth, 0)
	loader.Stop()
	fmt.Println()
	fmt.Println(string(result))
	return nil
}

func (m *Model) GenAnswer(ctx context.Context, prompt string) error {
	prompt = fmt.Sprintf("%s\nResponse in %s", prompt, m.language)
	termWidth, _, err := term.GetSize(0)
	if err != nil {
		termWidth = 80
	}
	fmt.Print(string(markdown.Render(fmt.Sprintf("**Model**: %s\n", m.model_name), termWidth, 0)))

	fmt.Println()

	if m.stream_mode {
		return m.stream_answer(ctx, prompt)
	}
	return m.render_markdown_answer(ctx, prompt)
}

// Roll back to `ROLLBACK_MODEL` if other model failed
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
