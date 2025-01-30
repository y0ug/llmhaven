package chat

// NewChatParams
func NewChatParams(
	opts ...func(*ChatParams),
) *ChatParams {
	p := &ChatParams{}
	p.Update(opts...)
	return p
}

func (p *ChatParams) Update(opts ...func(*ChatParams)) {
	for _, opt := range opts {
		opt(p)
	}
}

// WithModel sets the model for BaseChatMessageNewParams
func WithModel(model string) func(*ChatParams) {
	return func(p *ChatParams) {
		p.Model = model
	}
}

// WithMaxTokens sets the max tokens for BaseChatMessageNewParams
func WithMaxTokens(tokens int) func(*ChatParams) {
	return func(p *ChatParams) {
		p.MaxTokens = tokens
	}
}

// WithTemperature sets the temperature for BaseChatMessageNewParams
func WithTemperature(temp float64) func(*ChatParams) {
	return func(p *ChatParams) {
		p.Temperature = temp
	}
}

// WithMessages sets the messages for BaseChatMessageNewParams
func WithMessages(
	messages ...*ChatMessage,
) func(*ChatParams) {
	return func(p *ChatParams) {
		p.Messages = messages
	}
}

// WithTools sets the tools/functions for BaseChatMessageNewParams
func WithTools(tools ...Tool) func(*ChatParams) {
	return func(p *ChatParams) {
		p.Tools = tools
	}
}

func NewMessage(role string, content ...*MessageContent) *ChatMessage {
	return &ChatMessage{
		Role:    role,
		Content: content,
	}
}

func NewSystemMessage(text string) *ChatMessage {
	return NewMessage("system", NewTextContent(text))
}

func NewUserMessage(text string) *ChatMessage {
	return NewMessage("user", NewTextContent(text))
}
