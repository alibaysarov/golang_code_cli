package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/alibaysarov/coding-agent-cli/domain"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const SYSTEM_PROMPT = `
Ответь на вопрос пользователя и верни ответ СТРОГО в формате JSON,
без markdown-обёртки, без пояснений вне JSON, по следующей схеме:

{
  "files": [
    {
      "filePath": string, // путь до файла, если он упоминается в вопросе, иначе - придумай название файла сам
      "response": string  // содержимое/ответ для этого файла
    }
  ]
}

Код дели на файлы по SOLID не пиши разные уровни абстракции в 1 файл
Ненужно писать длинные ф-ии дроби их на более мелкие и вызывай внутри основной
Если для ответа достаточно одного файла — верни массив с одним элементом.
Если нужно несколько файлов — верни несколько элементов массива, объедини в одну папку если нужно.

`

const maxHistoryMessages = 20

type Ollama struct {
	llm     *ollama.LLM
	system  llms.MessageContent
	history []llms.MessageContent
	mu      sync.Mutex
}

func New() (*Ollama, error) {
	model := os.Getenv("MODEL")
	if model == "" {
		return nil, errors.New("Model not defined check .env")
	}
	llmClient, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		return nil, err
	}

	system := llms.TextParts(llms.ChatMessageTypeSystem, SYSTEM_PROMPT)

	return &Ollama{
		llm:     llmClient,
		system:  system,
		history: []llms.MessageContent{},
	}, nil
}

func (o *Ollama) AskStructured(ctx context.Context, userPrompt string, maxRetries int) (*domain.CodeResponse, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.history = append(o.history, llms.TextParts(llms.ChatMessageTypeHuman, userPrompt))
	o.trimHistory()

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		completion, err := o.llm.GenerateContent(ctx, o.buildMessages(), llms.WithJSONMode())
		if err != nil {
			return nil, fmt.Errorf("ошибка запроса к модели: %w", err)
		}

		raw := completion.Choices[0].Content

		var response domain.CodeResponse
		if err := json.Unmarshal([]byte(raw), &response); err != nil {
			lastErr = fmt.Errorf("невалидный JSON: %w (ответ: %s)", err, raw)
			o.history = append(o.history,
				llms.TextParts(llms.ChatMessageTypeAI, raw),
				llms.TextParts(llms.ChatMessageTypeHuman,
					fmt.Sprintf("Это не валидный JSON (%v). Верни ТОЛЬКО корректный JSON по схеме, без пояснений.", err)),
			)
			o.trimHistory()
			continue
		}

		o.history = append(o.history, llms.TextParts(llms.ChatMessageTypeAI, raw))
		o.trimHistory()
		return &response, nil
	}

	return nil, fmt.Errorf("не удалось получить валидный JSON после %d попыток: %w", maxRetries+1, lastErr)
}

func (o *Ollama) trimHistory() {
	if len(o.history) <= maxHistoryMessages {
		return
	}
	// обрезаем с начала, оставляя "хвост"
	excess := len(o.history) - maxHistoryMessages
	o.history = o.history[excess:]
}

// buildMessages собирает финальный список сообщений для запроса:
// системный промпт + (возможно урезанная) история.
func (o *Ollama) buildMessages() []llms.MessageContent {
	messages := make([]llms.MessageContent, 0, len(o.history)+1)
	messages = append(messages, o.system)
	messages = append(messages, o.history...)
	return messages
}
