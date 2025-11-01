package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Enhanced LLM Provider Interface with Tool Calling

// ToolGenerationRequest represents a request for generation with tools
type ToolGenerationRequest struct {
	ID          uuid.UUID              `json:"id"`
	Prompt      string                 `json:"prompt"`
	Tools       []Tool                 `json:"tools"`
	MaxTokens   int                    `json:"max_tokens"`
	Temperature float64                `json:"temperature"`
	Stream      bool                   `json:"stream"`
	Context     map[string]interface{} `json:"context"`
}

// ToolGenerationResponse represents the response from tool-based generation
type ToolGenerationResponse struct {
	ID        uuid.UUID              `json:"id"`
	Text      string                 `json:"text"`
	ToolCalls []ToolCall             `json:"tool_calls"`
	Reasoning string                 `json:"reasoning"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ToolStreamChunk represents a streaming chunk for tool-based generation
type ToolStreamChunk struct {
	ID        uuid.UUID              `json:"id"`
	Content   string                 `json:"content"`
	ToolCalls []ToolCall             `json:"tool_calls"`
	Reasoning string                 `json:"reasoning"`
	Done      bool                   `json:"done"`
	Error     string                 `json:"error,omitempty"`
}

// EnhancedLLMProvider extends the base LLMProvider with tool calling capabilities
type EnhancedLLMProvider interface {
	LLMProvider
	GenerateWithTools(ctx context.Context, req ToolGenerationRequest) (*ToolGenerationResponse, error)
	StreamWithTools(ctx context.Context, req ToolGenerationRequest) (<-chan ToolStreamChunk, error)
	ListAvailableTools() []Tool
	RegisterTool(tool Tool) error
}

// ToolCallingProvider implements EnhancedLLMProvider with tool calling support
type ToolCallingProvider struct {
	baseProvider LLMProvider
	tools        map[string]Tool
	reasoningEngine *ReasoningEngine
}

// NewToolCallingProvider creates a new tool calling provider
func NewToolCallingProvider(baseProvider LLMProvider) *ToolCallingProvider {
	return &ToolCallingProvider{
		baseProvider:   baseProvider,
		tools:          make(map[string]Tool),
		reasoningEngine: NewReasoningEngine(baseProvider),
	}
}

// GenerateWithTools performs generation with tool calling support
func (p *ToolCallingProvider) GenerateWithTools(ctx context.Context, req ToolGenerationRequest) (*ToolGenerationResponse, error) {
	startTime := time.Now()

	// Build tool-enhanced prompt
	enhancedPrompt := p.buildToolEnhancedPrompt(req.Prompt, req.Tools)

	// Generate initial response
	genReq := GenerationRequest{
		Prompt:      enhancedPrompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	resp, err := p.baseProvider.Generate(ctx, genReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate with tools: %v", err)
	}

	// Parse tool calls from response
	toolCalls, reasoning := p.extractToolCallsAndReasoning(resp.Text)

	// Execute tool calls if any
	if len(toolCalls) > 0 {
		results, err := p.executeToolCalls(ctx, toolCalls)
		if err != nil {
			log.Printf("Warning: Some tool calls failed: %v", err)
		}

		// Generate final response with tool results
		finalPrompt := p.buildFinalPrompt(req.Prompt, resp.Text, results)
		genReq.Prompt = finalPrompt
		
		finalResp, err := p.baseProvider.Generate(ctx, genReq)
		if err != nil {
			return nil, fmt.Errorf("failed to generate final response: %v", err)
		}
		resp = finalResp
	}

	return &ToolGenerationResponse{
		ID:        uuid.New(),
		Text:      resp.Text,
		ToolCalls: toolCalls,
		Reasoning: reasoning,
		Metadata: map[string]interface{}{
			"duration_ms": time.Since(startTime).Milliseconds(),
			"tools_used":   len(toolCalls),
		},
	}, nil
}

// StreamWithTools performs streaming generation with tool calling support
func (p *ToolCallingProvider) StreamWithTools(ctx context.Context, req ToolGenerationRequest) (<-chan ToolStreamChunk, error) {
	ch := make(chan ToolStreamChunk, 100)

	go func() {
		defer close(ch)

		// Build tool-enhanced prompt
		enhancedPrompt := p.buildToolEnhancedPrompt(req.Prompt, req.Tools)

		// Stream initial response
		streamReq := GenerationRequest{
			Prompt:      enhancedPrompt,
			MaxTokens:   req.MaxTokens,
			Temperature: req.Temperature,
			Stream:      true,
		}

		stream, err := p.baseProvider.Stream(ctx, streamReq)
		if err != nil {
			ch <- ToolStreamChunk{
				ID:    uuid.New(),
				Error: fmt.Sprintf("Failed to start streaming: %v", err),
				Done:  true,
			}
			return
		}

		var fullResponse string
		var toolCalls []ToolCall
		var reasoning string

		for chunk := range stream {
			if chunk.Error != "" {
				ch <- ToolStreamChunk{
					ID:    uuid.New(),
					Error: chunk.Error,
					Done:  true,
				}
				return
			}

			fullResponse += chunk.Content

			// Send streaming chunk
			ch <- ToolStreamChunk{
				ID:        uuid.New(),
				Content:   chunk.Content,
				ToolCalls: []ToolCall{},
				Reasoning: "",
				Done:      false,
			}
		}

		// Parse tool calls after streaming completes
		toolCalls, reasoning = p.extractToolCallsAndReasoning(fullResponse)

		// Execute tool calls if any
		if len(toolCalls) > 0 {
			results, err := p.executeToolCalls(ctx, toolCalls)
			if err != nil {
				log.Printf("Warning: Some tool calls failed: %v", err)
			}

			// Generate final response with tool results
			finalPrompt := p.buildFinalPrompt(req.Prompt, fullResponse, results)
			
			// Stream final response
			finalStreamReq := GenerationRequest{
				Prompt:      finalPrompt,
				MaxTokens:   req.MaxTokens,
				Temperature: req.Temperature,
				Stream:      true,
			}

			finalStream, err := p.baseProvider.Stream(ctx, finalStreamReq)
			if err != nil {
				ch <- ToolStreamChunk{
					ID:    uuid.New(),
					Error: fmt.Sprintf("Failed to stream final response: %v", err),
					Done:  true,
				}
				return
			}

			for chunk := range finalStream {
				if chunk.Error != "" {
					ch <- ToolStreamChunk{
						ID:    uuid.New(),
						Error: chunk.Error,
						Done:  true,
					}
					return
				}

				ch <- ToolStreamChunk{
					ID:        uuid.New(),
					Content:   chunk.Content,
					ToolCalls: toolCalls,
					Reasoning: reasoning,
					Done:      chunk.Done,
				}
			}
		} else {
			// No tool calls, send final chunk
			ch <- ToolStreamChunk{
				ID:        uuid.New(),
				Content:   "",
				ToolCalls: toolCalls,
				Reasoning: reasoning,
				Done:      true,
			}
		}
	}()

	return ch, nil
}

// ListAvailableTools returns all registered tools
func (p *ToolCallingProvider) ListAvailableTools() []Tool {
	tools := make([]Tool, 0, len(p.tools))
	for _, tool := range p.tools {
		tools = append(tools, tool)
	}
	return tools
}

// RegisterTool registers a new tool with the provider
func (p *ToolCallingProvider) RegisterTool(tool Tool) error {
	if _, exists := p.tools[tool.Name]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name)
	}
	p.tools[tool.Name] = tool
	
	// Also register with reasoning engine
	if p.reasoningEngine != nil {
		p.reasoningEngine.RegisterTool(tool)
	}
	
	log.Printf("Tool registered: %s", tool.Name)
	return nil
}

// Implement base LLMProvider interface

func (p *ToolCallingProvider) Generate(ctx context.Context, req GenerationRequest) (*GenerationResponse, error) {
	return p.baseProvider.Generate(ctx, req)
}

func (p *ToolCallingProvider) Stream(ctx context.Context, req GenerationRequest) (<-chan StreamChunk, error) {
	return p.baseProvider.Stream(ctx, req)
}

func (p *ToolCallingProvider) GetModelInfo() ModelInfo {
	return p.baseProvider.GetModelInfo()
}

func (p *ToolCallingProvider) IsHealthy() bool {
	return p.baseProvider.IsHealthy()
}

// Helper methods

func (p *ToolCallingProvider) buildToolEnhancedPrompt(prompt string, tools []Tool) string {
	toolDescriptions := ""
	for _, tool := range tools {
		paramsJSON, _ := json.Marshal(tool.Parameters)
		toolDescriptions += fmt.Sprintf("- %s: %s (parameters: %s)\n", 
			tool.Name, tool.Description, string(paramsJSON))
	}

	return fmt.Sprintf(`You have access to the following tools:
%s

When you need to use a tool, specify it in this format:
TOOL_CALL: {"tool_name": "tool_name", "arguments": {...}}

After using tools, provide your final answer.

User request: %s

Your response:`, toolDescriptions, prompt)
}

func (p *ToolCallingProvider) extractToolCallsAndReasoning(text string) ([]ToolCall, string) {
	var toolCalls []ToolCall
	reasoning := ""

	// Simple parsing for tool calls
	// In a real implementation, you would use more sophisticated parsing
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "TOOL_CALL:") {
			// Extract JSON from tool call
			jsonStart := strings.Index(line, "{")
			jsonEnd := strings.LastIndex(line, "}")
			if jsonStart != -1 && jsonEnd != -1 {
				jsonStr := line[jsonStart:jsonEnd+1]
				var toolCall ToolCall
				if err := json.Unmarshal([]byte(jsonStr), &toolCall); err == nil {
					toolCalls = append(toolCalls, toolCall)
				}
			}
		} else if !strings.Contains(line, "TOOL_CALL:") {
			// Collect reasoning (non-tool-call lines)
			reasoning += line + "\n"
		}
	}

	return toolCalls, strings.TrimSpace(reasoning)
}

func (p *ToolCallingProvider) executeToolCalls(ctx context.Context, toolCalls []ToolCall) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	
	for _, toolCall := range toolCalls {
		tool, exists := p.tools[toolCall.ToolName]
		if !exists {
			results[toolCall.ToolName] = fmt.Sprintf("Tool not found: %s", toolCall.ToolName)
			continue
		}

		result, err := tool.Handler(ctx, toolCall.Arguments)
		if err != nil {
			results[toolCall.ToolName] = fmt.Sprintf("Tool error: %v", err)
		} else {
			results[toolCall.ToolName] = result
		}
	}

	return results, nil
}

func (p *ToolCallingProvider) buildFinalPrompt(originalPrompt, initialResponse string, toolResults map[string]interface{}) string {
	resultsStr := ""
	for toolName, result := range toolResults {
		resultsStr += fmt.Sprintf("- %s: %v\n", toolName, result)
	}

	return fmt.Sprintf(`Original request: %s

Initial response: %s

Tool execution results:
%s

Based on the tool results, provide your final answer:`, 
		originalPrompt, initialResponse, resultsStr)
}