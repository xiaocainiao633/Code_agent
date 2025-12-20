package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	PythonAgent   PythonAgentConfig   `mapstructure:"python_agent"`
	Git           GitConfig           `mapstructure:"git"`
	FileProcessor FileProcessorConfig `mapstructure:"file_processor"`
	TaskScheduler TaskSchedulerConfig `mapstructure:"task_scheduler"`
	WebSocket     WebSocketConfig     `mapstructure:"websocket"`
	Logging       LoggingConfig       `mapstructure:"logging"`
	Database      DatabaseConfig      `mapstructure:"database"`
	GithubOAuth   GithubOAuthConfig   `mapstructure:"github_oauth"`
	Email         EmailConfig         `mapstructure:"email"`
	CORS          CORSConfig          `mapstructure:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// PythonAgentConfig Python AI Agent配置
type PythonAgentConfig struct {
	Host       string        `mapstructure:"host"`
	Port       string        `mapstructure:"port"`
	Timeout    time.Duration `mapstructure:"timeout"`
	RetryCount int           `mapstructure:"retry_count"`
}

// GitConfig Git配置
type GitConfig struct {
	CloneTimeout      time.Duration `mapstructure:"clone_timeout"`
	MaxFileSize       int64         `mapstructure:"max_file_size"`
	AllowedExtensions []string      `mapstructure:"allowed_extensions"`
}

// FileProcessorConfig 文件处理器配置
type FileProcessorConfig struct {
	MaxUploadSize   int64         `mapstructure:"max_upload_size"`
	TempDir         string        `mapstructure:"temp_dir"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// TaskSchedulerConfig 任务调度器配置
type TaskSchedulerConfig struct {
	MaxConcurrentTasks int           `mapstructure:"max_concurrent_tasks"`
	TaskTimeout        time.Duration `mapstructure:"task_timeout"`
	ResultRetention    time.Duration `mapstructure:"result_retention"`
}

// WebSocketConfig WebSocket配置
type WebSocketConfig struct {
	PingInterval   time.Duration `mapstructure:"ping_interval"`
	PongTimeout    time.Duration `mapstructure:"pong_timeout"`
	MaxMessageSize int64         `mapstructure:"max_message_size"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// GithubOAuthConfig GitHub OAuth配置
type GithubOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CODEGO")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// 配置文件不存在时使用默认值
		fmt.Println("Config file not found, using defaults")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// 验证配置
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认配置
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// Python Agent默认配置
	viper.SetDefault("python_agent.host", "localhost")
	viper.SetDefault("python_agent.port", "8000")
	viper.SetDefault("python_agent.timeout", "30s")
	viper.SetDefault("python_agent.retry_count", 3)

	// Git默认配置
	viper.SetDefault("git.clone_timeout", "5m")
	viper.SetDefault("git.max_file_size", 10485760) // 10MB
	viper.SetDefault("git.allowed_extensions", []string{".py", ".js", ".java", ".cpp", ".c", ".go", ".rs", ".ts"})

	// 文件处理器默认配置
	viper.SetDefault("file_processor.max_upload_size", 52428800) // 50MB
	viper.SetDefault("file_processor.temp_dir", "./temp")
	viper.SetDefault("file_processor.cleanup_interval", "1h")

	// 任务调度器默认配置
	viper.SetDefault("task_scheduler.max_concurrent_tasks", 10)
	viper.SetDefault("task_scheduler.task_timeout", "10m")
	viper.SetDefault("task_scheduler.result_retention", "24h")

	// WebSocket默认配置
	viper.SetDefault("websocket.ping_interval", "30s")
	viper.SetDefault("websocket.pong_timeout", "60s")
	viper.SetDefault("websocket.max_message_size", 512000) // 500KB

	// 日志默认配置
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "./logs/go-backend.log")

	// 数据库默认配置
	viper.SetDefault("database.path", "./data/codesage.db")

	// GitHub OAuth默认配置
	viper.SetDefault("github_oauth.client_id", "")
	viper.SetDefault("github_oauth.client_secret", "")
	viper.SetDefault("github_oauth.redirect_url", "http://localhost:3000/auth/github/callback")

	// 邮件默认配置
	viper.SetDefault("email.smtp_host", "smtp.gmail.com")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.smtp_user", "")
	viper.SetDefault("email.smtp_password", "")
	viper.SetDefault("email.from_email", "noreply@codesage.dev")
	viper.SetDefault("email.from_name", "CodeSage")

	// CORS默认配置
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"*"})
	viper.SetDefault("cors.allow_credentials", true)
}

// validate 验证配置
func validate(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}
	if cfg.PythonAgent.Host == "" || cfg.PythonAgent.Port == "" {
		return fmt.Errorf("python agent host and port cannot be empty")
	}
	if cfg.FileProcessor.TempDir == "" {
		return fmt.Errorf("file processor temp directory cannot be empty")
	}
	if cfg.Logging.Output == "" {
		return fmt.Errorf("logging output cannot be empty")
	}
	return nil
}

// GetPythonAgentURL 获取Python Agent服务URL
func (c *Config) GetPythonAgentURL() string {
	return fmt.Sprintf("http://%s:%s", c.PythonAgent.Host, c.PythonAgent.Port)
}

// EnsureDirectories 确保必要的目录存在
func (c *Config) EnsureDirectories() error {
	directories := []string{
		c.FileProcessor.TempDir,
		filepath.Dir(c.Logging.Output),
		filepath.Dir(c.Database.Path),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
