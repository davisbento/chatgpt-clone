package web

import (
	chatCompletion "davisbento/whats-gpt/chat-ms/internal/usecase/chat_completion"
)

type WebChatGPTHandler struct {
	CompletionUseCase chatCompletion.ChatCompletionUseCase
	Config            chatCompletion.ChatCompletionConfigInputDTO
	AuthToken         string
}

func NewWebChatGPTHandler(completionUseCase chatCompletion.ChatCompletionUseCase, config chatCompletion.ChatCompletionConfigInputDTO, authToken string) *WebChatGPTHandler {
	return &WebChatGPTHandler{
		CompletionUseCase: completionUseCase,
		Config:            config,
		AuthToken:         authToken,
	}
}
