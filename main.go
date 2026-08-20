package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	_ "embed"
	"strings"
	"syscall"
	"time"

	"github.com/alibaysarov/coding-agent-cli/files"
	"github.com/alibaysarov/coding-agent-cli/llm"
	userinput "github.com/alibaysarov/coding-agent-cli/user_input"
	"github.com/joho/godotenv"
)

func agentLoop(ctx context.Context, model *llm.Ollama) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			userInput := userinput.HandleInput()
			if userInput == "" {
				continue
			}
			if err := prompt(ctx, model, userInput); err != nil {
				fmt.Println("ошибка:", err)
			}
		}
	}
}

func prompt(ctx context.Context, model *llm.Ollama, userPrompt string) error {
	start := time.Now()
	response, err := model.AskStructured(ctx, userPrompt, 2)
	if err != nil {
		fmt.Println("Time:", time.Since(start))
		return err
	}

	fmt.Println("Time:", time.Since(start))
	fmt.Printf("Распарсенная структура: %+v\n", *response)

	command, err := userinput.Confirm()
	if err != nil {
		log.Fatalln("неправильная команда", err)
	}

	if command {
		for _, file := range response.Files {
			resultPath := fmt.Sprintf("testFiles/%s", file.FilePath)
			if err := files.WriteFile(resultPath, file.Response); err != nil {
				log.Fatalln("ошибка при записи в файл")
			}
		}
		return nil
	}

	addInput := userinput.HandleInput()
	return prompt(ctx, model, addInput)
}

//go:embed .env
var embeddedEnv []byte

func loadEnv() error {
	envMap, err := godotenv.Parse(strings.NewReader(string(embeddedEnv)))
	if err != nil {
		return err
	}
	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	err := loadEnv()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	llm, err := llm.New()

	if err != nil {
		log.Fatalf("не удалось создать клиент ollama: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	
	go agentLoop(ctx, llm)

	select {
	case <-ctx.Done():
		fmt.Println("done...")
		return
	}
}
