package chat_completion

import (
	"context"
	"davisbento/chatgpt-clone/chat-ms/internal/domain/entity"
	"davisbento/chatgpt-clone/chat-ms/internal/domain/gateway"
	"errors"

	"github.com/sashabaranov/go-openai"
)

type ChatCompletionConfigInputDTO struct {
	Model                string
	ModelMaxTokens       int
	Temperature          float32
	TopP                 float32
	N                    int
	Stop                 []string
	MaxTokens            int
	PresencePenalty      float32
	FrequencyPenalty     float32
	InitialSystemMessage string
}

type ChatCompletionInputDTO struct {
	ChatID      string                       `json:"chat_id,omitempty"`
	UserID      string                       `json:"user_id"`
	UserMessage string                       `json:"user_message"`
	Config      ChatCompletionConfigInputDTO `json:"config"`
}

type ChatCompletionOutputDTO struct {
	ChatID  string `json:"chat_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type ChatCompletionUseCase struct {
	ChatGateway  gateway.ChatGateway
	OpenAiClient *openai.Client
}

func NewChatCompletionUseCase(chatGateway gateway.ChatGateway, openAiClient *openai.Client) *ChatCompletionUseCase {
	return &ChatCompletionUseCase{
		ChatGateway:  chatGateway,
		OpenAiClient: openAiClient,
	}
}

func (uc *ChatCompletionUseCase) Execute(ctx context.Context, input ChatCompletionInputDTO) (*ChatCompletionOutputDTO, error) {
	chat, err := uc.ChatGateway.FindChatById(ctx, input.ChatID)

	if err != nil {
		if err.Error() == "chat not found" {
			// create new chat on entity (locally)
			chat, err = createNewChat(input)
			if err != nil {
				return nil, errors.New("error creating new chat: " + err.Error())
			}

			// persist chat
			err = uc.ChatGateway.CreateChat(ctx, chat)
			if err != nil {
				return nil, errors.New("error persisting new chat: " + err.Error())
			}
		} else {
			return nil, errors.New("error finding chat: " + err.Error())
		}
	}

	// create new message on entity (locally)
	userMessage, err := entity.NewMessage("user", input.UserMessage, chat.Config.Model)
	if err != nil {
		return nil, errors.New("error creating new message: " + err.Error())
	}

	err = chat.AddMessage(userMessage)
	if err != nil {
		return nil, errors.New("error adding message to chat: " + err.Error())
	}

	messages := []openai.ChatCompletionMessage{}

	for _, msg := range chat.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// call openai api
	resp, err := uc.OpenAiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:            chat.Config.Model.Name,
		Messages:         messages,
		MaxTokens:        chat.Config.MaxTokens,
		Temperature:      chat.Config.Temperature,
		TopP:             chat.Config.TopP,
		N:                chat.Config.N,
		Stop:             chat.Config.Stop,
		PresencePenalty:  chat.Config.PresencePenalty,
		FrequencyPenalty: chat.Config.FrequencyPenalty,
		Stream:           false,
	})

	if err != nil {
		return nil, errors.New("error calling openai api (chat completion): " + err.Error())
	}

	assistant, err := entity.NewMessage("assistant", resp.Choices[0].Message.Content, chat.Config.Model)
	if err != nil {
		return nil, errors.New("error creating assistant message: " + err.Error())
	}
	err = chat.AddMessage(assistant)
	if err != nil {
		return nil, errors.New("error adding assistant message to chat: " + err.Error())
	}

	// persist chat
	err = uc.ChatGateway.SaveChat(ctx, chat)
	if err != nil {
		return nil, errors.New("error persisting chat: " + err.Error())
	}

	return &ChatCompletionOutputDTO{
		ChatID:  chat.ID,
		UserID:  chat.UserID,
		Content: resp.Choices[0].Message.Content,
	}, nil
}

func createNewChat(input ChatCompletionInputDTO) (*entity.Chat, error) {
	model := entity.NewModel(input.Config.Model, input.Config.ModelMaxTokens)

	chatConfig := &entity.ChatConfig{
		Temperature:      input.Config.Temperature,
		TopP:             input.Config.TopP,
		N:                input.Config.N,
		Stop:             input.Config.Stop,
		MaxTokens:        input.Config.MaxTokens,
		PresencePenalty:  input.Config.PresencePenalty,
		FrequencyPenalty: input.Config.FrequencyPenalty,
		Model:            model,
	}

	initialMessage, err := entity.NewMessage("system", input.Config.InitialSystemMessage, model)

	if err != nil {
		return nil, errors.New("error creating initial message: " + err.Error())
	}

	chat, err := entity.NewChat(input.UserID, initialMessage, chatConfig)
	if err != nil {
		return nil, errors.New("error creating new chat: " + err.Error())
	}

	return chat, nil
}
