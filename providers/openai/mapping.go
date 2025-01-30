package openai

import (
	"encoding/json"

	"github.com/y0ug/llmhaven/chat"
)

func MessageToOpenAI(
	m ...*chat.ChatMessage,
) []ChatCompletionMessageParam {
	userMessages := make([]ChatCompletionMessageParam, 0)
	ignored := map[int]bool{}

	for i, msg := range m {
		// for _, content := range msg.Content {
		content := msg.Content[0]
		if ignored[i] {
			continue
		}

		switch content.Type {
		case chat.ContentTypeToolUse:
			// For toolCalls we need to process all of them in one time
			userMessages = append(userMessages, ChatCompletionMessageParam{
				Role:      "assistant",
				ToolCalls: MessageContentToToolCall(msg.Content...),
			})

		case chat.ContentTypeToolResult:
			for _, c := range msg.Content {
				userMessages = append(userMessages, ChatCompletionMessageParam{
					Role:       "tool",
					Content:    c.Content,
					ToolCallID: c.ToolUseID,
				})
			}
		default:
			userMessages = append(userMessages, ChatCompletionMessageParam{
				Role:    msg.Role,
				Content: content.String(),
			})
		}
	}
	return userMessages
}

func ToolCallToMessageContent(t ToolCall) *chat.MessageContent {
	// var args map[string]interface{}
	// _ = json.Unmarshal([]byte(t.Function.Arguments), &args)
	return chat.NewToolUseContent(t.ID, t.Function.Name, json.RawMessage(t.Function.Arguments))
}

func MessageContentToToolCall(t ...*chat.MessageContent) []ToolCall {
	d := make([]ToolCall, 0)
	for _, content := range t {
		if content.Type == chat.ContentTypeToolUse {
			d = append(d, ToolCall{
				ID:   content.ID,
				Type: "function",
				Function: FunctionCall{
					Name:      content.Name,
					Arguments: string(content.Input),
				},
			})
		}
		if content.Type == chat.ContentTypeToolResult {
		}
	}
	return d
}

func ToolsToOpenAI(tools ...chat.Tool) []Tool {
	result := make([]Tool, 0)
	for _, tool := range tools {
		var desc *string
		if tool.Description != nil {
			descCopy := *tool.Description
			desc = &descCopy
			if len(*desc) > 512 {
				foo := descCopy[:512]
				desc = &foo
			}
		}
		aiTool := Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        tool.Name,
				Description: desc,
				Parameters:  tool.InputSchema,
			},
		}
		result = append(result, aiTool)

	}
	return result
}

func ToStopReason(reason string) string {
	match := map[string]string{
		"stop":          "end_turn",
		"length":        "max_tokens",
		"stop_sequence": "stop_sequence", // Stop is same as stop_sequence we dont handle it
		"tool_calls":    "tool_use",
	}
	if r, ok := match[reason]; ok {
		return r
	} else {
		return reason
	}
}

func ToChatResponse(cc *ChatCompletion) *chat.ChatResponse {
	cm := &chat.ChatResponse{}
	cm.ID = cc.ID
	cm.Model = cc.Model
	cm.Usage = &chat.ChatUsage{}
	cm.Usage.InputTokens = cc.Usage.PromptTokens
	cm.Usage.OutputTokens = cc.Usage.CompletionTokens
	cm.Usage.OutputReasoningTokens = cc.Usage.CompletionTokensDetails.ReasoningTokens
	cm.Usage.InputCachedTokens = cc.Usage.PromptTokensDetails.CachedTokens
	cm.Usage.InputAudioTokens = cc.Usage.PromptTokensDetails.AudioTokens
	cm.Usage.OutputAudioTokens = cc.Usage.CompletionTokensDetails.AudioTokens

	for _, choice := range cc.Choices {
		c := chat.ChatChoice{}
		for _, call := range choice.Message.ToolCalls {
			c.Content = append(
				c.Content,
				ToolCallToMessageContent(call),
			)
		}

		if choice.Message.Content != "" {
			c.Content = append(c.Content, chat.NewTextContent(choice.Message.Content))
		}

		// Role is not choice is our model
		c.Role = choice.Message.Role

		// The reason the model stopped generating tokens. This will be `stop` if the model
		// hit a natural stop point or a provided stop sequence, `length` if the maximum
		// number of tokens specified in the request was reached, `content_filter` if
		// content was omitted due to a flag from our content filters, `tool_calls` if the
		// model called a tool, or `function_call` (deprecated) if the model called a
		// function.
		c.StopReason = ToStopReason(choice.FinishReason)
		cm.Choice = append(cm.Choice, c)
	}
	return cm
}

func ToChatCompletionNewParams(
	params chat.ChatParams,
) ChatCompletionNewParams {
	return ChatCompletionNewParams{
		Model:               params.Model,
		MaxCompletionTokens: &params.MaxTokens,
		Temperature:         params.Temperature,
		N:                   params.N,
		Messages:            MessageToOpenAI(params.Messages...),
		Tools:               ToolsToOpenAI(params.Tools...),
	}
}
