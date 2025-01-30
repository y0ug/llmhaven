# LLMHaven

LLMHaven is a Go library that provides unified access to various Large Language Model (LLM) providers through a single, consistent interface. It simplifies the process of integrating and switching between different LLM services in your Go applications.

## Features

- Unified interface for multiple LLM providers
- Support for streaming responses
- Built-in retry mechanisms and error handling
- Provider-specific configuration options
- Extensible architecture for adding new providers

## Currently Supported Providers

- OpenAI
- Anthropic (Claude)
- Google (Gemini)
- DeepSeek
- OpenRouter

## Installation

```bash
go get github.com/y0ug/llmhaven
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/y0ug/llmhaven"
    "github.com/y0ug/llmhaven/chat"
)

func main() {
    // Create a provider (e.g., OpenAI)
    provider, err := llmhaven.New("openai")
    if err != nil {
        panic(err)
    }

    // Create chat parameters
    params := chat.NewChatParams(
        chat.WithModel("gpt-3.5-turbo"),
        chat.WithMaxTokens(100),
        chat.WithTemperature(0.7),
        chat.WithMessages(
            chat.NewUserMessage("Hello, how are you?"),
        ),
    )

    // Send the request
    response, err := provider.Send(context.Background(), params)
    if err != nil {
        panic(err)
    }

    fmt.Println(response.Choice[0].Content[0].String())
}
```

### Streaming Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/y0ug/llmhaven"
    "github.com/y0ug/llmhaven/chat"
)

func main() {
    provider, _ := llmhaven.New("anthropic")
    
    params := chat.NewChatParams(
        chat.WithModel("claude-3-sonnet"),
        chat.WithMessages(
            chat.NewUserMessage("Write a short story."),
        ),
    )

    stream, err := provider.Stream(context.Background(), params)
    if err != nil {
        panic(err)
    }

    for stream.Next() {
        evt := stream.Current()
        if evt.Type == "text_delta" {
            fmt.Print(evt.Delta)
        }
    }

    if err := stream.Err(); err != nil {
        panic(err)
    }
}
```

### Provider-Specific Configuration

```go
package main

import (
    "github.com/y0ug/llmhaven"
    "github.com/y0ug/llmhaven/http/options"
)

func main() {
    // Configure with custom options
    provider, _ := llmhaven.New("openai",
        options.WithAuthToken("your-api-key"),
        options.WithMaxRetries(3),
        options.WithRequestTimeout(30*time.Second),
    )
    // ... use provider
}
```

## Environment Variables

The library supports the following environment variables for API authentication:

- `OPENAI_API_KEY` - OpenAI API key
- `ANTHROPIC_API_KEY` - Anthropic API key
- `GEMINI_API_KEY` - Google Gemini API key
- `DEEPSEEK_API_KEY` - DeepSeek API key
- `OPENROUTER_API_KEY` - OpenRouter API key

## Advanced Features

- Support for function calling/tools
- Message caching
- Rate limiting
- Custom middleware support
- Detailed usage statistics
- Error handling and retries

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License

Copyright (c) 2024 y0ug

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Acknowledgments

This library is inspired by various LLM provider SDKs and aims to provide a unified interface for easy integration in Go applications.
