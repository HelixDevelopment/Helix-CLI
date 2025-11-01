package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dev.helix.code/internal/hardware"
	"dev.helix.code/internal/llm"
)

// SimpleCLI represents a simplified command-line interface
type SimpleCLI struct {
	modelManager   *llm.ModelManager
	hardwareDetector *hardware.Detector
}

// NewSimpleCLI creates a new simple CLI instance
func NewSimpleCLI() *SimpleCLI {
	return &SimpleCLI{
		modelManager:   llm.NewModelManager(),
		hardwareDetector: hardware.NewDetector(),
	}
}

// Initialize sets up the CLI environment
func (c *SimpleCLI) Initialize() error {
	log.Println("🚀 Initializing Helix CLI...")

	// Detect hardware
	hardwareInfo, err := c.hardwareDetector.Detect()
	if err != nil {
		log.Printf("Warning: Hardware detection failed: %v", err)
	} else {
		log.Printf("✅ Hardware detected: %s CPU, %s GPU, %s RAM", 
			hardwareInfo.CPU.Model, hardwareInfo.GPU.Model, hardwareInfo.Memory.TotalRAM)
		log.Printf("📊 Optimal model size: %s", c.hardwareDetector.GetOptimalModelSize())
	}

	// Initialize basic providers
	if err := c.initializeBasicProviders(); err != nil {
		return fmt.Errorf("failed to initialize providers: %w", err)
	}

	log.Println("✅ Helix CLI initialized successfully")
	return nil
}

// Run executes the CLI with the given arguments
func (c *SimpleCLI) Run(args []string) error {
	if len(args) < 2 {
		return c.showHelp()
	}

	command := args[1]
	switch command {
	case "help", "--help", "-h":
		return c.showHelp()
	case "version", "--version", "-v":
		return c.showVersion()
	case "models":
		return c.listModels()
	case "hardware":
		return c.showHardwareInfo()
	case "health":
		return c.checkHealth()
	default:
		return fmt.Errorf("unknown command: %s. Use 'help' for available commands", command)
	}
}

// initializeBasicProviders sets up basic LLM providers
func (c *SimpleCLI) initializeBasicProviders() error {
	// Initialize a simple local provider for demonstration
	llamaConfig := llm.LlamaConfig{
		ModelPath:     "~/models/llama-2-7b-chat.gguf",
		ContextSize:   4096,
		GPUEnabled:    true,
		GPULayers:     35,
		ServerHost:    "localhost",
		ServerPort:    8080,
		ServerTimeout: 30,
	}

	provider, err := llm.NewLlamaCPPProvider(llamaConfig)
	if err != nil {
		log.Printf("Warning: Failed to initialize Llama.cpp provider: %v", err)
	} else {
		if err := c.modelManager.RegisterProvider(provider); err != nil {
			log.Printf("Warning: Failed to register Llama.cpp provider: %v", err)
		}
	}

	return nil
}

// Command implementations

func (c *SimpleCLI) showHelp() error {
	helpText := `
Helix CLI - AI-Powered Development Assistant (Phase 4 Implementation)

Usage:
  helix <command> [arguments]

Commands:
  help                    Show this help message
  version                 Show version information
  models                  List available AI models
  hardware                Show hardware information
  health                  Check system health

Examples:
  helix models
  helix hardware
  helix health

This is a Phase 4 implementation demonstrating the core architecture.
Advanced features like chat, code generation, and project planning are
coming in future phases.
`
	fmt.Println(helpText)
	return nil
}

func (c *SimpleCLI) showVersion() error {
	fmt.Println("Helix CLI v1.0.0 (Phase 4 Implementation)")
	fmt.Println("Build: development")
	fmt.Println("Go version: go1.24.0")
	fmt.Println("Architecture: Core LLM Integration & Hardware Detection")
	return nil
}

func (c *SimpleCLI) listModels() error {
	models := c.modelManager.GetAvailableModels()
	if len(models) == 0 {
		fmt.Println("No models available. Check your provider configurations.")
		return nil
	}

	fmt.Println("\nAvailable AI Models:")
	fmt.Println("===================")
	for _, model := range models {
		fmt.Printf("\n📦 %s\n", model.Name)
		fmt.Printf("   Provider: %s\n", model.Provider)
		fmt.Printf("   Context: %d tokens\n", model.ContextSize)
		fmt.Printf("   Capabilities: %v\n", model.Capabilities)
		if model.Description != "" {
			fmt.Printf("   Description: %s\n", model.Description)
		}
	}

	fmt.Printf("\nTotal: %d models available\n", len(models))
	
	// Show model selection example
	fmt.Println("\n💡 Model Selection Example:")
	criteria := llm.ModelSelectionCriteria{
		TaskType: "code_generation",
		RequiredCapabilities: []llm.ModelCapability{
			llm.CapabilityCodeGeneration,
			llm.CapabilityCodeAnalysis,
		},
		MaxTokens: 2048,
		QualityPreference: "balanced",
	}
	
	selectedModel, err := c.modelManager.SelectOptimalModel(criteria)
	if err != nil {
		fmt.Printf("   Model selection failed: %v\n", err)
	} else {
		fmt.Printf("   For code generation, recommended model: %s\n", selectedModel.Name)
	}

	return nil
}

func (c *SimpleCLI) showHardwareInfo() error {
	hardwareInfo, err := c.hardwareDetector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect hardware: %w", err)
	}

	fmt.Println("\nHardware Information:")
	fmt.Println("====================")
	fmt.Printf("CPU: %s (%s)\n", hardwareInfo.CPU.Model, hardwareInfo.CPU.Architecture)
	fmt.Printf("Cores: %d\n", hardwareInfo.CPU.Cores)
	fmt.Printf("GPU: %s (%s)\n", hardwareInfo.GPU.Model, hardwareInfo.GPU.Vendor)
	if hardwareInfo.GPU.VRAM != "" {
		fmt.Printf("VRAM: %s\n", hardwareInfo.GPU.VRAM)
	}
	fmt.Printf("RAM: %s\n", hardwareInfo.Memory.TotalRAM)
	fmt.Printf("Platform: %s/%s\n", hardwareInfo.Platform.OS, hardwareInfo.Platform.Architecture)
	
	optimalSize := c.hardwareDetector.GetOptimalModelSize()
	fmt.Printf("\n💡 Recommended model size: %s\n", optimalSize)
	
	flags := c.hardwareDetector.GetCompilationFlags()
	if len(flags) > 0 {
		fmt.Printf("🔧 Compilation flags: %v\n", flags)
	}

	// Show hardware compatibility for different model sizes
	fmt.Println("\n📊 Hardware Compatibility:")
	modelSizes := []string{"3B", "7B", "13B", "34B", "70B"}
	for _, size := range modelSizes {
		canRun := c.hardwareDetector.CanRunModel(size)
		status := "❌"
		if canRun {
			status = "✅"
		}
		fmt.Printf("   %s %s models: %s\n", status, size, 
			map[bool]string{true: "Supported", false: "Not supported"}[canRun])
	}

	return nil
}

func (c *SimpleCLI) checkHealth() error {
	fmt.Println("\nSystem Health Check:")
	fmt.Println("===================")

	// Check hardware
	hardwareInfo, err := c.hardwareDetector.Detect()
	if err != nil {
		fmt.Printf("❌ Hardware detection: FAILED (%v)\n", err)
	} else {
		fmt.Printf("✅ Hardware detection: OK (%s CPU, %s GPU)\n", 
			hardwareInfo.CPU.Model, hardwareInfo.GPU.Model)
	}

	// Check providers
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health := c.modelManager.HealthCheck(ctx)
	if len(health) == 0 {
		fmt.Println("❌ No LLM providers available")
	} else {
		for providerType, status := range health {
			if status.Status == "healthy" {
				fmt.Printf("✅ %s: HEALTHY (latency: %v)\n", providerType, status.Latency)
			} else {
				fmt.Printf("❌ %s: %s\n", providerType, status.Status)
			}
		}
	}

	// Check models
	models := c.modelManager.GetAvailableModels()
	if len(models) == 0 {
		fmt.Println("❌ No AI models available")
	} else {
		fmt.Printf("✅ Models: %d available\n", len(models))
	}

	// System summary
	fmt.Println("\n📈 System Summary:")
	optimalSize := c.hardwareDetector.GetOptimalModelSize()
	fmt.Printf("   Optimal model size: %s\n", optimalSize)
	fmt.Printf("   Available models: %d\n", len(models))
	fmt.Printf("   Available providers: %d\n", len(health))
	
	if len(models) > 0 && len(health) > 0 {
		fmt.Println("\n🎉 System is ready for AI-powered development!")
	} else {
		fmt.Println("\n⚠️  System needs configuration for full functionality")
	}

	return nil
}

// Main function
func main() {
	// Create CLI instance
	cli := NewSimpleCLI()

	// Initialize CLI
	if err := cli.Initialize(); err != nil {
		log.Fatalf("Failed to initialize CLI: %v", err)
	}

	// Handle interrupt signals
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigCh
		log.Println("\n🛑 Received interrupt signal, shutting down...")
		cancel()
	}()

	// Run CLI with command-line arguments
	if err := cli.Run(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}