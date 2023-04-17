package gateway

import (
	"context"
	"davisbento/whats-gpt/chat-ms/internal/domain/entity"
)

type ChatGateway interface {
	CreateChat(ctx context.Context, chat *entity.Chat) error
	FindChatById(ctx context.Context, id string) (*entity.Chat, error)
	SaveChat(ctx context.Context, chat *entity.Chat) error
}
