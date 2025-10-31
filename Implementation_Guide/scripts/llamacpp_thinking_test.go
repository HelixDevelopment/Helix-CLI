package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// LlamaCPPThinkingTester tests advanced capabilities of local coding models
type LlamaCPPThinkingTester struct {
	baseURL string
}

func NewLlamaCPPThinkingTester() *LlamaCPPThinkingTester {
	return &LlamaCPPThinkingTester{
		baseURL: "http://localhost:8080",
	}
}

// TestReasoningCapability tests if model can perform chain-of-thought reasoning
func (t *LlamaCPPThinkingTester) TestReasoningCapability(model string, prompt string, t *testing.T) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// Enhanced prompt for reasoning
	reasoningPrompt := fmt.Sprintf(`Think step by step about this problem:

%s

Show your reasoning process clearly with steps.`, prompt)
	
	request := map[string]interface{}{
		"model":       model,
		"prompt":      reasoningPrompt,
		"stream":      false,
		"temperature": 0.3,
		"max_tokens":  1000,
	}
	
	resp, err := t.makeRequest(ctx, "/completion", request)
	if err != nil {
		t.Logf("Reasoning test failed for %s: %v", model, err)
		return false
	}
	
	content := resp["content"].(string)
	
	// Check for reasoning indicators
	reasoningIndicators := []string{"step", "first", "then", "next", "therefore", "conclusion", "reason", "because"}
	reasoningScore := 0
	
	for _, indicator := range reasoningIndicators {
		if strings.Contains(strings.ToLower(content), indicator) {
			reasoningScore++
		}
	}
	
	// Check for structured reasoning
	if strings.Contains(content, "1.") && strings.Contains(content, "2.") {
		reasoningScore += 2
	}
	
	if strings.Contains(content, "Step 1:") || strings.Contains(content, "First,") {
		reasoningScore += 2
	}
	
	t.Logf("Reasoning test for %s: score %d/10", model, reasoningScore)
	
	return reasoningScore >= 5
}

// TestToolCallingCapability tests if model can understand and generate tool calls
func (t *LlamaCPPThinkingTester) TestToolCallingCapability(model string, t *testing.T) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	toolPrompt := `You have access to these tools:

- create_file: Create a new file with content
  Parameters: filename (string), content (string)
  
- run_tests: Execute tests in a directory
  Parameters: directory (string), verbose (boolean)
  
- git_commit: Commit changes to git
  Parameters: message (string), files (array of strings)

When you need to use a tool, respond in this exact format:
TOOL: tool_name
ARGS: {"param1": "value1", "param2": "value2"}

User request: Create a new Go file called "utils.go" with helper functions for string manipulation.

Respond with tool calls if needed:`
	
	request := map[string]interface{}{
		"model":       model,
		"prompt":      toolPrompt,
		"stream":      false,
		"temperature": 0.7,
		"max_tokens":  500,
	}
	
	resp, err := t.makeRequest(ctx, "/completion", request)
	if err != nil {
		t.Logf("Tool calling test failed for %s: %v", model, err)
		return false
	}
	
	content := resp["content"].(string)
	
	// Check for tool call patterns
	toolCallScore := 0
	
	if strings.Contains(content, "TOOL:") {
		toolCallScore += 3
	}
	
	if strings.Contains(content, "create_file") {
		toolCallScore += 2
	}
	
	if strings.Contains(content, "ARGS:") || strings.Contains(content, "{\"filename\"") {
		toolCallScore += 2
	}
	
	// Check for JSON-like arguments
	if strings.Contains(content, "utils.go") || strings.Contains(content, "string manipulation") {
		toolCallScore += 1
	}
	
	t.Logf("Tool calling test for %s: score %d/8", model, toolCallScore)
	
	return toolCallScore >= 4
}

// TestComplexCodeGeneration tests model's ability to generate working code
func (t *LlamaCPPThinkingTester) TestComplexCodeGeneration(model string, complexity string, t *testing.T) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	
	var prompt string
	
	switch complexity {
	case "simple":
		prompt = `Create a Go function that reverses a string. Return only the code without explanations.`
	case "medium":
		prompt = `Create a Go HTTP middleware that logs requests and responses with timing information. Return only the code without explanations.`
	case "complex":
		prompt = `Create a complete Go microservice with:
- REST API endpoints
- Database connection
- Error handling
- Logging
Return only the code without explanations.`
	}
	
	request := map[string]interface{}{
		"model":       model,
		"prompt":      prompt,
		"stream":      false,
		"temperature": 0.3,
		"max_tokens":  1500,
	}
	
	resp, err := t.makeRequest(ctx, "/completion", request)
	if err != nil {
		t.Logf("Code generation test failed for %s: %v", model, err)
		return false
	}
	
	content := resp["content"].(string)
	
	// Extract Go code
	code := t.extractGoCode(content)
	
	if code == "" {
		t.Logf("No valid Go code generated for %s complexity %s", model, complexity)
		return false
	}
	
	// Check code quality indicators
	qualityScore := 0
	
	// Basic Go syntax
	if strings.Contains(code, "package ") {
		qualityScore++
	}
	if strings.Contains(code, "import ") {
		qualityScore++
	}
	if strings.Contains(code, "func ") {
		qualityScore++
	}
	
	// Code structure based on complexity
	switch complexity {
	case "simple":
		if strings.Contains(code, "string") && strings.Contains(code, "range") {
			qualityScore += 2
		}
	case "medium":
		if strings.Contains(code, "http.Handler") && strings.Contains(code, "middleware") {
			qualityScore += 2
		}
		if strings.Contains(code, "time.Now()") || strings.Contains(code, "duration") {
			qualityScore++
		}
	case "complex":
		if strings.Contains(code, "router") || strings.Contains(code, "mux") {
			qualityScore++
		}
		if strings.Contains(code, "database") || strings.Contains(code, "sql") {
			qualityScore++
		}
		if strings.Contains(code, "error") && strings.Contains(code, "if err != nil") {
			qualityScore++
		}
	}
	
	minScore := map[string]int{"simple": 4, "medium": 5, "complex": 6}
	
	passed := qualityScore >= minScore[complexity]
	
	t.Logf("Code generation test for %s (%s): score %d/%d", model, complexity, qualityScore, minScore[complexity])
	
	if passed {
		t.Logf("✅ %s can handle %s code generation", model, complexity)
	} else {
		t.Logf("⚠️  %s struggles with %s code generation", model, complexity)
	}
	
	return passed
}

// makeRequest helper function for API calls
func (t *LlamaCPPThinkingTester) makeRequest(ctx context.Context, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", t.baseURL+endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	
	return result, nil
}

// extractGoCode extracts Go code from model response
func (t *LlamaCPPThinkingTester) extractGoCode(response string) string {
	// Look for code blocks
	start := strings.Index(response, "```go")
	if start != -1 {
		end := strings.Index(response[start:], "```")
		if end != -1 {
			code := response[start+5 : start+end]
			return strings.TrimSpace(code)
		}
	}
	
	// If no code blocks, try to find package declaration
	lines := strings.Split(response, "\n")
	var codeLines []string
	inCode := false
	
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			inCode = true
		}
		if inCode {
			codeLines = append(codeLines, line)
		}
		// Stop if we hit explanations
		if inCode && (strings.HasPrefix(strings.TrimSpace(line), "//") || 
			strings.Contains(strings.ToLower(line), "explanation") ||
			strings.Contains(strings.ToLower(line), "note:")) {
			break
		}
	}
	
	if len(codeLines) > 0 {
		return strings.Join(codeLines, "\n")
	}
	
	return ""
}

// TestLlamaCPPThinkingAndTooling is the main test function
func TestLlamaCPPThinkingAndTooling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Llama.cpp thinking and tooling tests in short mode")
	}
	
	tester := NewLlamaCPPThinkingTester()
	
	// Test models
	testModels := []string{
		"codellama:7b",
		"codellama:13b",
		"llama3.1:8b", 
		"deepseek-coder:6.7b",
	}
	
	successfulModels := []string{}
	
	for _, model := range testModels {
		t.Run(model, func(t *testing.T) {
			t.Logf("\n=== Testing %s ===", model)
			
			reasoningPass := 0
			toolPass := 0
			codePass := 0
			
			// Test reasoning
			if tester.TestReasoningCapability(model, "Design an algorithm to detect cycles in a linked list", t) {
				reasoningPass++
			}
			
			if tester.TestReasoningCapability(model, "Plan a microservices architecture for an e-commerce platform", t) {
				reasoningPass++
			}
			
			// Test tool calling
			if tester.TestToolCallingCapability(model, t) {
				toolPass++
			}
			
			// Test code generation
			if tester.TestComplexCodeGeneration(model, "simple", t) {
				codePass++
			}
			
			if tester.TestComplexCodeGeneration(model, "medium", t) {
				codePass++
			}
			
			// Evaluate model
			totalScore := reasoningPass + toolPass + codePass
			
			if totalScore >= 4 {
				t.Logf("\n✅ %s PASSED comprehensive testing", model)
				t.Logf("   Reasoning: %d/2, Tooling: %d/1, Code: %d/2", reasoningPass, toolPass, codePass)
				successfulModels = append(successfulModels, model)
			} else {
				t.Logf("\n⚠️  %s has limited capabilities", model)
				t.Logf("   Reasoning: %d/2, Tooling: %d/1, Code: %d/2", reasoningPass, toolPass, codePass)
			}
		})
	}
	
	// Final summary
	t.Logf("\n=== TEST SUMMARY ===")
	if len(successfulModels) > 0 {
		t.Logf("✅ SUCCESSFUL MODELS (Thinking + Tooling):")
		for _, model := range successfulModels {
			t.Logf("   - %s", model)
		}
		t.Logf("\n🎉 %d model(s) ready for HelixCode development!", len(successfulModels))
	} else {
		t.Fatal("❌ No models fully support thinking and tooling capabilities")
	}
}

func main() {
	// Run tests when executed directly
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "TestLlamaCPPThinkingAndTooling",
				F:    TestLlamaCPPThinkingAndTooling,
			},
		},
		[]testing.InternalBenchmark{},
		[]testing.InternalExample{},
	)
}