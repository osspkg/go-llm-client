// Package messages provides typed Anthropic Messages API operations.
package messages

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

//go:generate easyjson -all types.go

// Request creates a message or message stream.
type Request struct {
	// Model is the provider model identifier.
	Model string `json:"model"`
	// Messages contains the user and assistant conversation turns.
	Messages []Message `json:"messages"`
	// MaxTokens is the maximum number of tokens the model may generate.
	MaxTokens int `json:"max_tokens"`
	// System is a system string or an array of system content blocks.
	System Content `json:"system,omitempty"`
	// Stream requests server-sent events instead of one buffered response.
	Stream bool `json:"stream,omitempty"`
	// Temperature controls sampling randomness.
	Temperature float64 `json:"temperature,omitempty"`
	// TopP enables nucleus sampling over the cumulative token probability.
	TopP float64 `json:"top_p,omitempty"`
	// TopK limits sampling to the K most likely tokens.
	TopK int `json:"top_k,omitempty"`
	// StopSequences are caller-defined strings that stop generation.
	StopSequences []string `json:"stop_sequences,omitempty"`
	// Tools declares functions or server tools available to the model.
	Tools []Tool `json:"tools,omitempty"`
	// ToolChoice controls whether and which declared tool the model must use.
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
	// Thinking configures extended thinking and its token budget.
	Thinking *ThinkingConfig `json:"thinking,omitempty"`
	// OutputConfig configures structured output when supported by the model.
	OutputConfig *OutputConfig `json:"output_config,omitempty"`
	// Metadata contains caller-defined request metadata.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Message is one conversation turn. Content accepts a string or content block array.
type Message struct {
	// Role identifies the speaker, normally user or assistant.
	Role string `json:"role"`
	// Content is text or a provider-defined array of content blocks.
	Content Content `json:"content"`
}

// Content is the Anthropic string-or-content-block-array union.
// Use TextContent or BlocksContent to construct it.
//
//easyjson:skip
type Content []byte

// TextContent creates a text message or system prompt.
func TextContent(text string) Content {
	data, err := json.Marshal(text)
	if err != nil {
		return nil
	}
	return Content(data)
}

// BlocksContent creates a content value from typed content blocks.
func BlocksContent(blocks ...ContentBlock) Content {
	data, err := json.Marshal(blocks)
	if err != nil {
		return nil
	}
	return Content(data)
}

// MarshalJSON validates and returns the content union unchanged.
func (content Content) MarshalJSON() ([]byte, error) {
	if len(content) == 0 {
		return []byte("null"), nil
	}
	if !json.Valid(content) {
		return nil, errors.New("invalid anthropic content JSON")
	}
	return append([]byte(nil), content...), nil
}

// UnmarshalJSON decodes a text or typed content-block array.
func (content *Content) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*content = nil
		return nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*content = TextContent(text)
		return nil
	}
	var blocks []ContentBlock
	if err := json.Unmarshal(data, &blocks); err != nil {
		return errors.New("anthropic content must be a string or content block array")
	}
	*content = BlocksContent(blocks...)
	return nil
}

// MarshalEasyJSON implements easyjson.Marshaler for the content union.
func (content Content) MarshalEasyJSON(writer *jwriter.Writer) {
	data, err := content.MarshalJSON()
	writer.Raw(data, err)
}

// UnmarshalEasyJSON implements easyjson.Unmarshaler for the content union.
func (content *Content) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	data := lexer.Raw()
	if lexer.Ok() {
		lexer.AddError(content.UnmarshalJSON(data))
	}
}

// DecodeText decodes a string content value.
func (content Content) DecodeText() (string, error) {
	var value string
	if err := json.Unmarshal(content, &value); err != nil {
		return "", err
	}
	return value, nil
}

// DecodeBlocks decodes an array content value into typed content block envelopes.
func (content Content) DecodeBlocks() ([]ContentBlock, error) {
	var value []ContentBlock
	if err := json.Unmarshal(content, &value); err != nil {
		return nil, err
	}
	return value, nil
}

// ContentBlock is a typed discriminator envelope for Anthropic content blocks.
type ContentBlock struct {
	// Type identifies text, image, document, thinking, tool_use, or tool_result.
	Type string `json:"type"`
	// Text contains text block content.
	Text string `json:"text,omitempty"`
	// Source contains typed image or document source data.
	Source json.RawMessage `json:"source,omitempty"`
	// Thinking contains model thinking text.
	Thinking string `json:"thinking,omitempty"`
	// Signature authenticates thinking content.
	Signature string `json:"signature,omitempty"`
	// Data contains redacted thinking data.
	Data string `json:"data,omitempty"`
	// ID identifies tool-use blocks.
	ID string `json:"id,omitempty"`
	// Name identifies the tool to invoke.
	Name string `json:"name,omitempty"`
	// Input contains tool arguments defined by the tool schema.
	Input json.RawMessage `json:"input,omitempty"`
	// ToolUseID links a tool result to its invocation.
	ToolUseID string `json:"tool_use_id,omitempty"`
	// Content contains nested tool-result content.
	Content json.RawMessage `json:"content,omitempty"`
	// IsError reports a failed tool execution.
	IsError bool `json:"is_error,omitempty"`
}

// TextBlockValue creates a typed text content block.
func TextBlockValue(text string) ContentBlock { return ContentBlock{Type: "text", Text: text} }

// ImageBlockValue creates an image content block from a typed source.
func ImageBlockValue(source ImageSource) ContentBlock {
	data, err := json.Marshal(source)
	if err != nil {
		return ContentBlock{Type: "image"}
	}
	return ContentBlock{Type: "image", Source: data}
}

// DocumentBlockValue creates a document content block from a typed source.
func DocumentBlockValue(source DocumentSource) ContentBlock {
	data, err := json.Marshal(source)
	if err != nil {
		return ContentBlock{Type: "document"}
	}
	return ContentBlock{Type: "document", Source: data}
}

// ThinkingBlockValue creates a thinking content block.
func ThinkingBlockValue(thinking, signature string) ContentBlock {
	return ContentBlock{Type: "thinking", Thinking: thinking, Signature: signature}
}

// RedactedThinkingBlockValue creates a redacted thinking content block.
func RedactedThinkingBlockValue(data string) ContentBlock {
	return ContentBlock{Type: "redacted_thinking", Data: data}
}

// DecodeImageSource decodes the source of an image content block.
func (block ContentBlock) DecodeImageSource() (ImageSource, error) {
	var value ImageSource
	if err := json.Unmarshal(block.Source, &value); err != nil {
		return ImageSource{}, err
	}
	return value, nil
}

// DecodeDocumentSource decodes the source of a document content block.
func (block ContentBlock) DecodeDocumentSource() (DocumentSource, error) {
	var value DocumentSource
	if err := json.Unmarshal(block.Source, &value); err != nil {
		return DocumentSource{}, err
	}
	return value, nil
}

// DecodeToolResultContent decodes nested tool-result text or blocks.
func (block ContentBlock) DecodeToolResultContent() (Content, error) {
	var value Content
	if err := value.UnmarshalJSON(block.Content); err != nil {
		return nil, err
	}
	return value, nil
}

// ToolUseBlockValue creates a tool-use content block.
func ToolUseBlockValue(id, name string, input json.RawMessage) ContentBlock {
	return ContentBlock{Type: "tool_use", ID: id, Name: name, Input: input}
}

// ToolResultBlockValue creates a tool-result content block.
func ToolResultBlockValue(toolUseID string, content json.RawMessage, isError bool) ContentBlock {
	return ContentBlock{Type: "tool_result", ToolUseID: toolUseID, Content: content, IsError: isError}
}

// TextBlock contains generated or supplied text content.
type TextBlock struct {
	// Type is the discriminator and is normally "text".
	Type string `json:"type"`
	// Text is the visible UTF-8 content.
	Text string `json:"text"`
	// Citations contains source references attached to generated text.
	Citations []Citation `json:"citations,omitempty"`
}

// ImageBlock contains an image supplied to the model.
type ImageBlock struct {
	// Type is the discriminator and is normally "image".
	Type string `json:"type"`
	// Source describes the image bytes or remote media reference.
	Source ImageSource `json:"source"`
}

// DocumentBlock contains a document supplied to the model.
type DocumentBlock struct {
	// Type is the discriminator and is normally "document".
	Type string `json:"type"`
	// Source describes document bytes, text, or a URL.
	Source DocumentSource `json:"source"`
	// Title is the optional display name shown to the model.
	Title string `json:"title,omitempty"`
	// Context gives the model additional instructions about the document.
	Context string `json:"context,omitempty"`
	// Citations enables or configures document citation behavior.
	Citations json.RawMessage `json:"citations,omitempty"`
}

// ThinkingBlock contains model reasoning emitted by a thinking-capable model.
type ThinkingBlock struct {
	// Type is the discriminator and is normally "thinking".
	Type string `json:"type"`
	// Thinking contains the provider's thinking text.
	Thinking string `json:"thinking"`
	// Signature authenticates the thinking block for continuation requests.
	Signature string `json:"signature"`
}

// RedactedThinkingBlock contains encrypted or redacted model reasoning.
type RedactedThinkingBlock struct {
	// Type is the discriminator and is normally "redacted_thinking".
	Type string `json:"type"`
	// Data contains the provider's redacted thinking payload.
	Data string `json:"data"`
}

// ToolUseBlock is a model request to invoke a declared tool.
type ToolUseBlock struct {
	// Type is the discriminator and is normally "tool_use".
	Type string `json:"type"`
	// ID uniquely identifies this tool invocation.
	ID string `json:"id"`
	// Name is the declared tool name to invoke.
	Name string `json:"name"`
	// Input contains arguments validated against the tool input schema.
	Input json.RawMessage `json:"input"`
}

// ToolResultBlock contains an application's result for a tool invocation.
type ToolResultBlock struct {
	// Type is the discriminator and is normally "tool_result".
	Type string `json:"type"`
	// ToolUseID links the result to the model's tool invocation.
	ToolUseID string `json:"tool_use_id"`
	// Content is text or nested result blocks returned to the model.
	Content Content `json:"content,omitempty"`
	// IsError tells the model that tool execution failed.
	IsError bool `json:"is_error,omitempty"`
}

// ImageSource identifies image bytes or a remote image location.
type ImageSource struct {
	// Type identifies the source encoding, such as base64 or URL.
	Type string `json:"type"`
	// MediaType is the MIME type when the source contains base64 bytes.
	MediaType string `json:"media_type,omitempty"`
	// Data is the base64-encoded image payload.
	Data string `json:"data,omitempty"`
	// URL is a remote image URL when URL sources are supported.
	URL string `json:"url,omitempty"`
}

// DocumentSource identifies document bytes, text, or a URL.
type DocumentSource struct {
	// Type identifies the source encoding, such as base64, text, or URL.
	Type string `json:"type"`
	// MediaType is the MIME type of a base64 or text document.
	MediaType string `json:"media_type,omitempty"`
	// Data is the base64-encoded document payload.
	Data string `json:"data,omitempty"`
	// Text is the plain-text document payload.
	Text string `json:"text,omitempty"`
	// URL is a remote document URL when URL sources are supported.
	URL string `json:"url,omitempty"`
}

// Citation identifies a source span referenced by generated text.
type Citation struct {
	// Type identifies the citation variant.
	Type string `json:"type"`
	// CitedText is the excerpt associated with the citation.
	CitedText string `json:"cited_text,omitempty"`
	// DocumentIndex is the zero-based document index for document citations.
	DocumentIndex int `json:"document_index,omitempty"`
	// DocumentTitle is the cited document title.
	DocumentTitle string `json:"document_title,omitempty"`
	// StartCharIndex is the inclusive source character offset.
	StartCharIndex int `json:"start_char_index,omitempty"`
	// EndCharIndex is the exclusive source character offset.
	EndCharIndex int `json:"end_char_index,omitempty"`
}

// Tool declares a function the model may invoke.
type Tool struct {
	// Name is the stable tool identifier exposed to the model.
	Name string `json:"name"`
	// Description explains when and how the model should use the tool.
	Description string `json:"description,omitempty"`
	// InputSchema is the JSON Schema for tool arguments.
	InputSchema json.RawMessage `json:"input_schema"`
	// CacheControl controls prompt-cache behavior for this tool definition.
	CacheControl json.RawMessage `json:"cache_control,omitempty"`
}

// ToolChoice selects automatic, any-tool, or named-tool selection.
type ToolChoice struct {
	// Type is auto, any, or tool.
	Type string `json:"type"`
	// Name is required when Type is tool.
	Name string `json:"name,omitempty"`
	// DisableParallelToolUse prevents multiple tool calls in one response.
	DisableParallelToolUse bool `json:"disable_parallel_tool_use,omitempty"`
}

// ThinkingConfig enables extended thinking with a bounded token budget.
type ThinkingConfig struct {
	// Type is normally "enabled".
	Type string `json:"type"`
	// BudgetTokens is the maximum thinking token budget.
	BudgetTokens int `json:"budget_tokens"`
}

// OutputConfig configures structured output behavior.
type OutputConfig struct {
	// Format is the provider-defined output format or JSON schema configuration.
	Format json.RawMessage `json:"format,omitempty"`
}

// Response is a complete assistant message.
type Response struct {
	// ID uniquely identifies the generated message.
	ID string `json:"id"`
	// Type is the resource discriminator and is normally "message".
	Type string `json:"type"`
	// Role identifies the generated speaker, normally assistant.
	Role string `json:"role"`
	// Content contains generated text, thinking, tool use, and other blocks.
	Content Content `json:"content"`
	// Model is the model that produced the message.
	Model string `json:"model"`
	// StopReason explains why generation ended.
	StopReason string `json:"stop_reason,omitempty"`
	// StopSequence is the stop sequence that ended generation, if any.
	StopSequence string `json:"stop_sequence,omitempty"`
	// Usage contains input, output, and cache token accounting.
	Usage Usage `json:"usage"`
	// Container contains provider-defined container metadata when enabled.
	Container json.RawMessage `json:"container,omitempty"`
}

// Usage contains Anthropic token accounting.
type Usage struct {
	// InputTokens counts tokens in the request.
	InputTokens int `json:"input_tokens"`
	// OutputTokens counts tokens generated in the response.
	OutputTokens int `json:"output_tokens"`
	// CacheCreationInputTokens counts tokens written to the prompt cache.
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	// CacheReadInputTokens counts tokens read from the prompt cache.
	CacheReadInputTokens int `json:"cache_read_input_tokens,omitempty"`
}

// StreamEvent is one typed SSE envelope from the Messages API.
type StreamEvent struct {
	// Type identifies message_start, content_block_delta, message_stop, or another event.
	Type string `json:"type"`
	// Index identifies the content block affected by a block event.
	Index int `json:"index,omitempty"`
	// Message contains the initial message object for message_start.
	Message json.RawMessage `json:"message,omitempty"`
	// ContentBlock contains the block created by content_block_start.
	ContentBlock json.RawMessage `json:"content_block,omitempty"`
	// Delta contains text, JSON, thinking, signature, or message deltas.
	Delta json.RawMessage `json:"delta,omitempty"`
	// Usage contains usage snapshots attached to message events.
	Usage *Usage `json:"usage,omitempty"`
	// Error contains a provider error event.
	Error json.RawMessage `json:"error,omitempty"`
}

// ErrorPayload describes the error object carried by an error event.
type ErrorPayload struct {
	// Type identifies the provider error category.
	Type string `json:"type"`
	// Message explains why the streamed request failed.
	Message string `json:"message"`
}

// DecodeMessage decodes the message_start payload.
func (event StreamEvent) DecodeMessage() (Response, error) {
	var value Response
	if err := json.Unmarshal(event.Message, &value); err != nil {
		return Response{}, err
	}
	return value, nil
}

// DecodeContentBlock decodes the content_block_start payload.
func (event StreamEvent) DecodeContentBlock() (ContentBlock, error) {
	var value ContentBlock
	if err := json.Unmarshal(event.ContentBlock, &value); err != nil {
		return ContentBlock{}, err
	}
	return value, nil
}

// DecodeContentBlockDelta decodes a text, tool-input, thinking, or signature delta.
func (event StreamEvent) DecodeContentBlockDelta() (ContentBlockDelta, error) {
	var value ContentBlockDelta
	if err := json.Unmarshal(event.Delta, &value); err != nil {
		return ContentBlockDelta{}, err
	}
	return value, nil
}

// DecodeMessageDelta decodes the final stop reason and usage delta.
func (event StreamEvent) DecodeMessageDelta() (MessageDelta, error) {
	var value MessageDelta
	if err := json.Unmarshal(event.Delta, &value); err != nil {
		return MessageDelta{}, err
	}
	return value, nil
}

// DecodeErrorPayload decodes the error object carried by an error event.
func (event StreamEvent) DecodeErrorPayload() (ErrorPayload, error) {
	var value ErrorPayload
	if err := json.Unmarshal(event.Error, &value); err != nil {
		return ErrorPayload{}, err
	}
	return value, nil
}

// ContentBlockDelta describes a typed content block delta payload.
type ContentBlockDelta struct {
	// Type identifies text_delta, input_json_delta, thinking_delta, or signature_delta.
	Type string `json:"type"`
	// Text is an incremental text fragment.
	Text string `json:"text,omitempty"`
	// PartialJSON is an incremental tool-input JSON fragment.
	PartialJSON string `json:"partial_json,omitempty"`
	// Thinking is an incremental thinking fragment.
	Thinking string `json:"thinking,omitempty"`
	// Signature is an incremental thinking signature fragment.
	Signature string `json:"signature,omitempty"`
}

// MessageDelta describes the final reason and usage for a streamed message.
type MessageDelta struct {
	// StopReason explains why the streamed response ended.
	StopReason string `json:"stop_reason,omitempty"`
	// StopSequence is the stop sequence that ended generation.
	StopSequence string `json:"stop_sequence,omitempty"`
	// Usage contains final output token accounting.
	Usage *Usage `json:"usage,omitempty"`
}

// CountTokensRequest asks Anthropic to count input tokens without generating output.
type CountTokensRequest struct {
	// Model is the model whose tokenizer should be used.
	Model string `json:"model"`
	// Messages contains the conversation to count.
	Messages []Message `json:"messages"`
	// System is a system string or array of system blocks.
	System Content `json:"system,omitempty"`
	// Tools contains declarations that contribute to the input token count.
	Tools []Tool `json:"tools,omitempty"`
	// Thinking contains thinking configuration included in the count.
	Thinking *ThinkingConfig `json:"thinking,omitempty"`
}

// CountTokensResponse contains the counted input token total.
type CountTokensResponse struct {
	// InputTokens is the number of tokens Anthropic counted for the request.
	InputTokens int `json:"input_tokens"`
}
