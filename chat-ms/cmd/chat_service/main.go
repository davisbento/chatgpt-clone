package main

import (
	"database/sql"
	"davisbento/chatgpt-clone/chat-ms/configs"
	"davisbento/chatgpt-clone/chat-ms/internal/infra/grpc/server"
	"davisbento/chatgpt-clone/chat-ms/internal/infra/repository"
	"davisbento/chatgpt-clone/chat-ms/internal/infra/web"
	"davisbento/chatgpt-clone/chat-ms/internal/infra/web/webserver"
	chatcompletion "davisbento/chatgpt-clone/chat-ms/internal/usecase/chat_completion"
	"davisbento/chatgpt-clone/chat-ms/internal/usecase/chat_completion_stream"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sashabaranov/go-openai"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conn, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := repository.NewChatRepositoryMySQL(conn)
	client := openai.NewClient(configs.OpenAIApiKey)

	chatConfig := chatcompletion.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	chatConfigStream := chat_completion_stream.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	useCase := chatcompletion.NewChatCompletionUseCase(repo, client)

	streamChannel := make(chan chat_completion_stream.ChatCompletionOutputDTO)
	useCaseStream := chat_completion_stream.NewChatCompletionStreamUseCase(repo, client, streamChannel)

	fmt.Println("Starting gRPC server on port " + configs.GRPCServerPort)
	grpcServer := server.NewGRPCServer(
		*useCaseStream,
		chatConfigStream,
		configs.GRPCServerPort,
		configs.AuthToken,
		streamChannel,
	)
	go grpcServer.Start()

	webServer := webserver.NewWebServer(":" + configs.WebServerPort)
	webServerChatHandler := web.NewWebChatGPTHandler(*useCase, chatConfig, configs.AuthToken)
	webServer.AddHandler("/chat", webServerChatHandler.Handle)

	fmt.Println("Server running on port " + configs.WebServerPort)
	webServer.Start()
}
