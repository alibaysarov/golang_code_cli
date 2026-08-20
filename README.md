# Coding Agent CLI

## Русский

### Название

**Coding Agent CLI** — консольный coding agent на Go, который использует локальную модель Ollama для генерации и изменения файлов проекта.

### Описание

Приложение принимает запрос пользователя в интерактивном режиме, отправляет его в Ollama и ожидает структурированный JSON-ответ со списком файлов. Перед записью изменений приложение показывает полученную структуру и запрашивает подтверждение `y/n`. После подтверждения файлы создаются или перезаписываются в каталоге `testFiles/`.

В проекте используются:

- Go 1.26.3 или новее;
- Ollama с установленной локальной моделью;
- `langchaingo` для обращения к Ollama;
- `godotenv` для чтения `.env`.

### Установка и сборка

#### Общая подготовка

1. Установите Go с [официального сайта](https://go.dev/dl/).
2. Установите [Ollama](https://ollama.com/download).
3. Запустите Ollama и загрузите модель, например:

   ```bash
   ollama serve
   ollama pull llama3.2
   ```

4. В корне проекта создайте или измените `.env`:

   ```env
   MODEL=llama3.2
   ```

   Значение `MODEL` должно совпадать с именем модели, установленной в Ollama.
5. Загрузите зависимости:

   ```bash
   go mod download
   ```

#### Windows

В PowerShell из корня проекта выполните:

```powershell
go mod download
go build -o coding-agent-cli.exe .
.\coding-agent-cli.exe
```

Или используйте стандартную команду Go:

```powershell
go run .
```

#### Ubuntu

```bash
go mod download
go build -o coding-agent-cli .
./coding-agent-cli
```

Для установки команды в `/usr/local/bin`:

```bash
go build -o coding-agent-cli .
sudo install -m 755 coding-agent-cli /usr/local/bin/coding-agent-cli
coding-agent-cli
```

Также доступен target из `makefile`:

```bash
make build
```

Он создаёт бинарник `test_cli` и копирует его в `/usr/local/bin/`.

#### macOS

```bash
go mod download
go build -o coding-agent-cli .
./coding-agent-cli
```

Для запуска без отдельной сборки:

```bash
go run .
```

### Использование

1. Запустите приложение из корня проекта, где находится `.env`.
2. Введите запрос, например `Создай файл hello.js с HTTP-сервером`.
3. Проверьте распарсенную структуру ответа.
4. Введите `y`, чтобы применить изменения, или `n`, чтобы отказаться.
5. Сгенерированные файлы появятся в каталоге `testFiles/`.

Для завершения нажмите `Ctrl+C`.

### Конфигурация и безопасность

В текущей реализации `.env` подключается директивой `go:embed`, поэтому значения из `.env` встраиваются в собранный бинарник. Не храните в этом файле пароли, токены или другие секреты и не публикуйте бинарники, собранные с конфиденциальной конфигурацией.

Ollama по умолчанию должен быть доступен локально по адресу `http://localhost:11434`.

### Структура проекта

- `main.go` — запуск приложения и основной цикл агента;
- `llm/` — клиент Ollama и обработка структурированного ответа;
- `domain/` — модели ответа и генерируемых файлов;
- `user_input/` — чтение запросов и подтверждений пользователя;
- `files/` — создание каталогов и запись файлов;
- `testFiles/` — каталог для результатов генерации.

## English

### Name

**Coding Agent CLI** is a Go command-line coding agent that uses a local Ollama model to generate and modify project files.

### Description

The application accepts prompts interactively, sends them to Ollama, and expects a structured JSON response containing generated files. Before applying changes, it displays the parsed response and asks for `y/n` confirmation. After confirmation, files are created or overwritten in the `testFiles/` directory.

The project uses:

- Go 1.26.3 or later;
- Ollama with a locally installed model;
- `langchaingo` for Ollama integration;
- `godotenv` for loading `.env`.

### Installation and build

#### Common setup

1. Install Go from the [official website](https://go.dev/dl/).
2. Install [Ollama](https://ollama.com/download).
3. Start Ollama and pull a model, for example:

   ```bash
   ollama serve
   ollama pull llama3.2
   ```

4. Create or update `.env` in the project root:

   ```env
   MODEL=llama3.2
   ```

   `MODEL` must match the name of a model installed in Ollama.
5. Download Go dependencies:

   ```bash
   go mod download
   ```

#### Windows

Run these commands in PowerShell from the project root:

```powershell
go mod download
go build -o coding-agent-cli.exe .
.\coding-agent-cli.exe
```

You can also run the project without creating a binary:

```powershell
go run .
```

#### Ubuntu

```bash
go mod download
go build -o coding-agent-cli .
./coding-agent-cli
```

To install the command system-wide:

```bash
go build -o coding-agent-cli .
sudo install -m 755 coding-agent-cli /usr/local/bin/coding-agent-cli
coding-agent-cli
```

The repository also includes a Makefile target:

```bash
make build
```

It creates a `test_cli` binary and copies it to `/usr/local/bin/`.

#### macOS

```bash
go mod download
go build -o coding-agent-cli .
./coding-agent-cli
```

To run without a separate build step:

```bash
go run .
```

### Usage

1. Run the application from the project root containing `.env`.
2. Enter a prompt, for example: `Create hello.js with an HTTP server`.
3. Review the parsed response structure.
4. Enter `y` to apply the changes or `n` to reject them.
5. Generated files will be written to `testFiles/`.

Press `Ctrl+C` to stop the application.

### Configuration and security

The current implementation loads `.env` with `go:embed`, which embeds its values into the compiled binary. Do not store passwords, tokens, or other secrets in this file, and do not distribute binaries built with confidential configuration.

By default, Ollama must be available locally at `http://localhost:11434`.

### Project structure

- `main.go` — application entry point and agent loop;
- `llm/` — Ollama client and structured response handling;
- `domain/` — response and generated-file models;
- `user_input/` — prompt and confirmation input;
- `files/` — directory creation and file writing;
- `testFiles/` — output directory for generated files.