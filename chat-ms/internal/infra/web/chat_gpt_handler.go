package web

import (
	chatCompletion "davisbento/whats-gpt/chat-ms/internal/usecase/chat_completion"
	"encoding/json"
	"io"
	"net/http"
)

type WebChatGPTHandler struct {
	CompletionUseCase chatCompletion.ChatCompletionUseCase
	Config            chatCompletion.ChatCompletionConfigInputDTO
	AuthToken         string
}

func NewWebChatGPTHandler(
	completionUseCase chatCompletion.ChatCompletionUseCase,
	config chatCompletion.ChatCompletionConfigInputDTO,
	authToken string,
) *WebChatGPTHandler {
	return &WebChatGPTHandler{
		CompletionUseCase: completionUseCase,
		Config:            config,
		AuthToken:         authToken,
	}
}

func (h *WebChatGPTHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Authorization") != h.AuthToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !json.Valid(body) {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var dto chatCompletion.ChatCompletionInputDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.Config = h.Config

	result, err := h.CompletionUseCase.Execute(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}