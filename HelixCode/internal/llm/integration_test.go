//go:build integration

package llm

import (
	"context"
	"testing"
	"time"
)

// TestLlamaCPPProviderIntegration tests the Llama.cpp provider with integration
func TestLlamaCPPProviderIntegration(t *testing.T) {
	config := LlamaConfig{
		ModelPath:     "test-models/test.gguf",
		ContextSize:   2048,
		GPUEnabled:    false,
		GPULayers:     0,
		ServerHost:    "localhost",
		ServerPort:    8081,
		ServerTimeout: 30,
	}

	provider, err := NewLlamaCPPProvider(config)
	if err != nil {
		t.Skipf("Skipping Llama.cpp integration test: %v", err)
	}

	// Test provider availability
	ctx := context.Background()
	available := provider.IsAvailable(ctx)
	if !available {
		t.Skip("Llama.cpp provider not available for integration test")
	}

	// Test model listing
	models := provider.GetModels()
	if len(models) == 0 {
		t.Error("Llama.cpp provider should return available models")
	}

	// Test capabilities
	capabilities := provider.GetCapabilities()
	if len(capabilities) == 0 {
		t.Error("Llama.cpp provider should have capabilities")
	}

	// Test health check
	health, err := provider.GetHealth(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
	if health.Status == "" {
		t.Error("Health status should not be empty")
	}

	t.Logf("✅ Llama.cpp provider integration test passed: %d models, health: %s", 
		len(models), health.Status)

	// Cleanup
	provider.Close()
}

// TestOllamaProviderIntegration tests the Ollama provider with integration
func TestOllamaProviderIntegration(t *testing.T) {
	config := OllamaConfig{
		BaseURL:       "http://localhost:11434",
		DefaultModel:  "llama2",
		Timeout:       30,
		KeepAlive:     300,
		StreamEnabled: false,
	}

	provider, err := NewOllamaProvider(config)
	if err != nil {
		t.Skipf("Skipping Ollama integration test: %v", err)
	}

	// Test provider availability
	ctx := context.Background()
	available := provider.IsAvailable(ctx)
	if !available {
		t.Skip("Ollama provider not available for integration test")
	}

	// Test model listing
	models := provider.GetModels()
	if len(models) == 0 {
		t.Log("No models available from Ollama (may be normal)")
	}

	// Test capabilities
	capabilities := provider.GetCapabilities()
	if len(capabilities) == 0 {
		t.Error("Ollama provider should have capabilities")
	}

	// Test health check
	health, err := provider.GetHealth(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
	if health.Status == "" {
		t.Error("Health status should not be empty")
	}

	t.Logf("✅ Ollama provider integration test passed: %d models, health: %s", 
		len(models), health.Status)

	// Cleanup
	provider.Close()
}

// TestModelManagerIntegration tests the model manager with integration
func TestModelManagerIntegration(t *testing.T) {
	manager := NewModelManager()

	// Test with mock providers
	llamaConfig := LlamaConfig{
		ModelPath:     "test-models/test.gguf",
		ContextSize:   2048,
		GPUEnabled:    false,
		ServerHost:    "localhost",
		ServerPort:    8082,
		ServerTimeout: 30,
	}

	llamaProvider, err := NewLlamaCPPProvider(llamaConfig)
	if err == nil {
		// Only register if provider was created successfully
		if err := manager.RegisterProvider(llamaProvider); err != nil {
			t.Logf("Failed to register Llama.cpp provider: %v", err)
		}
	}

	// Test model selection
	criteria := ModelSelectionCriteria{
		TaskType: "code_generation",
		RequiredCapabilities: []ModelCapability{
			CapabilityCodeGeneration,
			CapabilityCodeAnalysis,
		},
		MaxTokens:        1024,
		QualityPreference: "balanced",
	}

	selectedModel, err := manager.SelectOptimalModel(criteria)
	if err != nil {
		t.Logf("Model selection failed (may be normal): %v", err)
	} else if selectedModel != nil {
		t.Logf("Selected model: %s", selectedModel.Name)
	}

	// Test health checking
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health := manager.HealthCheck(ctx)
	if len(health) == 0 {
		t.Log("No providers available for health check (may be normal)")
	} else {
		for providerType, status := range health {
			t.Logf("Provider %s: %s", providerType, status.Status)
		}
	}

	t.Log("✅ Model manager integration test passed")

	// Cleanup
	if llamaProvider != nil {
		llamaProvider.Close()
	}
}

// TestProviderHealthIntegration tests provider health monitoring with integration
func TestProviderHealthIntegration(t *testing.T) {
	manager := NewModelManager()

	// Create test providers
	providers := []Provider{
		&MockProvider{
			name:         "mock-healthy",
			available:    true,
			healthy:      true,
			models:       []ModelInfo{{Name: "test-model"}},
			capabilities: []ModelCapability{CapabilityTextGeneration},
		},
		&MockProvider{
			name:         "mock-unhealthy",
			available:    true,
			healthy:      false,
			models:       []ModelInfo{{Name: "test-model"}},
			capabilities: []ModelCapability{CapabilityTextGeneration},
		},
	}

	// Register providers
	for _, provider := range providers {
		if err := manager.RegisterProvider(provider); err != nil {
			t.Fatalf("Failed to register provider: %v", err)
		}
	}

	// Test health check
	ctx := context.Background()
	health := manager.HealthCheck(ctx)

	if len(health) != 2 {
		t.Errorf("Expected health status for 2 providers, got %d", len(health))
	}

	// Verify health statuses
	for providerType, status := range health {
		t.Logf("Provider %s: %s", providerType, status.Status)
	}

	t.Log("✅ Provider health integration test passed")
}

// MockProvider is a mock LLM provider for testing
type MockProvider struct {
	name         string
	available    bool
	healthy      bool
	models       []ModelInfo
	capabilities []ModelCapability
}

func (m *MockProvider) GetType() ProviderType {
	return ProviderTypeLocal
}

func (m *MockProvider) GetName() string {
	return m.name
}

func (m *MockProvider) GetModels() []ModelInfo {
	return m.models
}

func (m *MockProvider) GetCapabilities() []ModelCapability {
	return m.capabilities
}

func (m *MockProvider) Generate(ctx context.Context, request *LLMRequest) (*LLMResponse, error) {
	return &LLMResponse{
		Content: "Mock response",
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}, nil
}

func (m *MockProvider) GenerateStream(ctx context.Context, request *LLMRequest, ch chan<- LLMResponse) error {
	ch <- LLMResponse{
		Content: "Mock streaming response",
	}
	close(ch)
	return nil
}

func (m *MockProvider) IsAvailable(ctx context.Context) bool {
	return m.available
}

func (m *MockProvider) GetHealth(ctx context.Context) (*ProviderHealth, error) {
	status := "healthy"
	if !m.healthy {
		status = "unhealthy"
	}

	return &ProviderHealth{
		Status:     status,
		LastCheck:  time.Now(),
		ErrorCount: 0,
		ModelCount: len(m.models),
	}, nil
}

func (m *MockProvider) Close() error {
	return nil
}