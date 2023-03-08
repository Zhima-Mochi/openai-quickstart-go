package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GenerateHandler struct {
	client *openai.Client
}

func NewGenerateHandler(client *openai.Client) *GenerateHandler {
	return &GenerateHandler{
		client: client,
	}
}

// TODO
func (handler *GenerateHandler) ChatCompletionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// swagger:parameters TextCompletionRequest
type TextCompletionRequest struct {
	// The name of the OpenAI model to use.
	//
	// Example: gpt3
	//
	// required: false
	Model string `json:"model"`

	// The prompt to generate completions for, encoded as a string.
	//
	// Example: "Once upon a time"
	//
	// required: true
	// in: body
	Prompt string `json:"prompt,omitempty"`

	// Sampling temperature to use. Higher values means the model will take more risks.
	//
	// Example: 0.5
	//
	// required: false
	Temperature *float32 `json:"temperature,omitempty"`

	// The maximum number of tokens(common sequences of characters found in text) to generate in the completion.
	//
	// Example: 10
	//
	// required: false
	MaxTokens int `json:"max_tokens,omitempty"`
}

// swagger:route POST /generate/text text-completion
//
// Generate text completion based on prompt.
//
// Generates text completion based on the given prompt using the OpenAI GPT-3 API.
//
// Responses:
//
//	200: string
//	400: errorResponse
//	500: errorResponse
func (handler *GenerateHandler) TextCompletionHandler(c *gin.Context) {
	var textCompletionRequest TextCompletionRequest

	if err := c.ShouldBindJSON(&textCompletionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.Trim(textCompletionRequest.Prompt, " ")) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please enter a valid animal"})
		return
	}

	var temperature float32 = 0.6
	if textCompletionRequest.Temperature != nil {
		temperature = *textCompletionRequest.Temperature
	}

	var model string = openai.GPT3TextDavinci003
	if textCompletionRequest.Model != "" {
		model = textCompletionRequest.Model
	}

	req := openai.CompletionRequest{
		Model:       model,
		Prompt:      generatePrompt(textCompletionRequest.Prompt),
		Temperature: temperature,
	}

	resp, err := handler.client.CreateCompletion(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred during your request."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": resp.Choices[0].Text})
}

func generatePrompt(animal string) string {
	capitalizedAnimal := cases.Title(language.English).String(strings.ToLower(animal))
	return fmt.Sprintf("Suggest three names for an animal that is a superhero.\nAnimal: Cat\nNames: Captain Sharpclaw, Agent Fluffball, The Incredible Feline\nAnimal: Dog\nNames: Ruff the Protector, Wonder Canine, Sir Barks-a-Lot\nAnimal: %s\nNames:", capitalizedAnimal)
}
