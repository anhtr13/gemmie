package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/genai"

	"github.com/anhtr13/gemmie/internal/model"
)

var (
	root_cmd    *cobra.Command
	client      *genai.Client
	user_opts   *model.ModelOptions
	saved_opts  *model.ModelOptions
	ctx         context.Context
	config_path string
	prompt      string
)

func init() {
	ctx = context.Background()
	user_opts = &model.ModelOptions{}
	saved_opts = model.BaseOptions()

	p, err := model.GetConfigPath()
	if err != nil {
		log.Fatal("Cannot get config path: ", err)
	}
	config_path = p

	err = saved_opts.LoadFromFile(config_path)
	if err != nil {
		saved_opts.SaveToFile(config_path)
	}

	if saved_opts.ApiKey != "" {
		client, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  saved_opts.ApiKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			log.Fatal("Cannot create client: ", err)
		}
	}

	prompt = strings.Trim(prompt, " ")

	root_cmd = &cobra.Command{
		Use:   "gemmie",
		Short: "Use -p [--prompt] flag to chat with GeminiAI",
		Long: `Use -p [--prompt] flag to chat with GeminiAI
		
Example: 
	gemmie --model=gemini-3.0-pro --lang=Vietnamese --temp=2.0 --limit=6900 -p="write a story about a magic backpack."
	`,
		Run: func(cmd *cobra.Command, args []string) {
			if client == nil {
				log.Fatal("Api key not found, cannot create client!")
			}
			model := model.NewModel(
				client,
				user_opts.ModelName,
				user_opts.Language,
				user_opts.Temperature,
				user_opts.MaxOutputToken,
				user_opts.TopP,
				user_opts.TopK,
				user_opts.StreamMode,
			)

			err := model.GenAnswerRollBack(ctx, prompt)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	root_cmd.AddCommand(set_api_key())
	root_cmd.AddCommand(set_options())
	root_cmd.AddCommand(list_models())
	root_cmd.AddCommand(list_options())
	root_cmd.AddCommand(reset())

	root_cmd.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "Your prompt")
	root_cmd.PersistentFlags().
		StringVar(&user_opts.Language, "lang", saved_opts.Language, "Specify the responses language")
	root_cmd.PersistentFlags().
		StringVar(&user_opts.ModelName, "model", saved_opts.ModelName, "Specify what Gemini model to use")
	root_cmd.PersistentFlags().
		Float32Var(&user_opts.Temperature, "temp", saved_opts.Temperature, "Controls the randomness of the output. Use higher values for more creative responses, and lower values for more deterministic responses. Values can range from [0.0, 2.0].")
	root_cmd.PersistentFlags().
		Float32Var(&user_opts.TopP, "top_p", saved_opts.TopP, "Changes how the model selects tokens for output. Tokens are selected from the most to least probable until the sum of their probabilities equals the topP value.")
	root_cmd.PersistentFlags().
		Float32Var(&user_opts.TopK, "top_k", saved_opts.TopK, "Changes how the model selects tokens for output. A topK of 1 means the selected token is the most probable among all the tokens in the model's vocabulary, while a topK of 3 means that the next token is selected from among the 3 most probable using the temperature. Tokens are further filtered based on topP with the final token selected using temperature sampling.")
	root_cmd.PersistentFlags().
		Int32Var(&user_opts.MaxOutputToken, "limit", saved_opts.MaxOutputToken, "Sets the maximum number of tokens to include in a candidate.")
	root_cmd.PersistentFlags().
		BoolVar(&user_opts.StreamMode, "stream", saved_opts.StreamMode, "Enable text stream effect (like Gemini, chatGPT, etc) but can not render markdown")
}

func Execute() {
	err := root_cmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}
