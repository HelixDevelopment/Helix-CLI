package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
	"dev.helix.code/internal/database"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database database.Config `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Workers  WorkersConfig  `mapstructure:"workers"`
	Tasks    TasksConfig    `mapstructure:"tasks"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Address         string `mapstructure:"address"`
	Port            int    `mapstructure:"port"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	IdleTimeout     int    `mapstructure:"idle_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	JWTSecret          string `mapstructure:"jwt_secret"`
	TokenExpiry        int    `mapstructure:"token_expiry"`
	SessionExpiry      int    `mapstructure:"session_expiry"`
	BcryptCost         int    `mapstructure:"bcrypt_cost"`
}

// WorkersConfig represents worker configuration
type WorkersConfig struct {
	HealthCheckInterval int `mapstructure:"health_check_interval"`
	HealthTTL           int `mapstructure:"health_ttl"`
	MaxConcurrentTasks  int `mapstructure:"max_concurrent_tasks"`
}

// TasksConfig represents task configuration
type TasksConfig struct {
	MaxRetries         int `mapstructure:"max_retries"`
	CheckpointInterval int `mapstructure:"checkpoint_interval"`
	CleanupInterval    int `mapstructure:"cleanup_interval"`
}

// LLMConfig represents LLM configuration
type LLMConfig struct {
	DefaultProvider string            `mapstructure:"default_provider"`
	Providers       map[string]string `mapstructure:"providers"`
	MaxTokens       int               `mapstructure:"max_tokens"`
	Temperature     float64           `mapstructure:"temperature"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Set default values
	setDefaults()

	// Find config file
	configPath := findConfigFile()
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Use default config locations
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config/")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME/.config/helixcode/")
		viper.AddConfigPath("/etc/helixcode/")
	}

	// Read in environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HELIX")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
		// Config file not found, but we can continue with defaults
		fmt.Println("⚠️  No config file found, using defaults and environment variables")
	} else {
		fmt.Printf("📁 Using config file: %s\n", viper.ConfigFileUsed())
	}

	// Unmarshal config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %v", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 60)
	viper.SetDefault("server.shutdown_timeout", 30)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "helixcode")
	viper.SetDefault("database.dbname", "helixcode")
	viper.SetDefault("database.sslmode", "disable")

	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "default-secret-change-in-production")
	viper.SetDefault("auth.token_expiry", 86400) // 24 hours
	viper.SetDefault("auth.session_expiry", 604800) // 7 days
	viper.SetDefault("auth.bcrypt_cost", 12)

	// Workers defaults
	viper.SetDefault("workers.health_check_interval", 30)
	viper.SetDefault("workers.health_ttl", 120)
	viper.SetDefault("workers.max_concurrent_tasks", 10)

	// Tasks defaults
	viper.SetDefault("tasks.max_retries", 3)
	viper.SetDefault("tasks.checkpoint_interval", 300)
	viper.SetDefault("tasks.cleanup_interval", 3600)

	// LLM defaults
	viper.SetDefault("llm.default_provider", "local")
	viper.SetDefault("llm.max_tokens", 4096)
	viper.SetDefault("llm.temperature", 0.7)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.output", "stdout")
}

// findConfigFile searches for config file in various locations
func findConfigFile() string {
	// Check environment variable first
	if configPath := os.Getenv("HELIX_CONFIG"); configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Check common locations
	locations := []string{
		"./config/config.yaml",
		"./config.yaml",
		"$HOME/.config/helixcode/config.yaml",
		"/etc/helixcode/config.yaml",
	}

	for _, location := range locations {
		if expanded := os.ExpandEnv(location); expanded != location {
			if _, err := os.Stat(expanded); err == nil {
				return expanded
			}
		}
	}

	return ""
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	// Server validation
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	// Database validation
	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Auth validation
	if cfg.Auth.JWTSecret == "" || cfg.Auth.JWTSecret == "default-secret-change-in-production" {
		return fmt.Errorf("JWT secret must be set and not use default value")
	}

	// Workers validation
	if cfg.Workers.HealthCheckInterval < 1 {
		return fmt.Errorf("health check interval must be positive")
	}
	if cfg.Workers.MaxConcurrentTasks < 1 {
		return fmt.Errorf("max concurrent tasks must be positive")
	}

	// Tasks validation
	if cfg.Tasks.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	// LLM validation
	if cfg.LLM.MaxTokens < 1 {
		return fmt.Errorf("max tokens must be positive")
	}
	if cfg.LLM.Temperature < 0 || cfg.LLM.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	return nil
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Create default config content
	configContent := `# HelixCode Server Configuration

server:
  address: "0.0.0.0"
  port: 8080
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 60
  shutdown_timeout: 30

database:
  host: "localhost"
  port: 5432
  user: "helixcode"
  password: "" # Set via HELIX_DATABASE_PASSWORD environment variable
  dbname: "helixcode"
  sslmode: "disable"

auth:
  jwt_secret: "" # Set via HELIX_AUTH_JWT_SECRET environment variable
  token_expiry: 86400
  session_expiry: 604800
  bcrypt_cost: 12

workers:
  health_check_interval: 30
  health_ttl: 120
  max_concurrent_tasks: 10

tasks:
  max_retries: 3
  checkpoint_interval: 300
  cleanup_interval: 3600

llm:
  default_provider: "local"
  providers:
    local: "http://localhost:11434"
    openai: "" # Set API key via environment variable
  max_tokens: 4096
  temperature: 0.7

logging:
  level: "info"
  format: "text"
  output: "stdout"
`

	// Write config file
	if err := os.WriteFile(path, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetEnvOrDefault gets an environment variable with a default value
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvIntOrDefault gets an environment variable as integer with a default value
func GetEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}