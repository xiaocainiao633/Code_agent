package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/database"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/handlers"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/middleware"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/services"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/websocket"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志系统
	if err := utils.InitLogger(&cfg.Logging); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if logger, ok := utils.GlobalLogger.(*utils.SimpleLogger); ok {
			logger.Close()
		}
	}()

	// 确保必要的目录存在
	if err := cfg.EnsureDirectories(); err != nil {
		utils.Fatal("Failed to create directories: %v", err)
	}

	// 初始化数据库
	if err := database.InitDatabase(cfg.Database.Path); err != nil {
		utils.Fatal("Failed to initialize database: %v", err)
	}
	defer database.Close()
	utils.Info("Database initialized successfully at %s", cfg.Database.Path)

	// 初始化服务
	authService := services.NewAuthService()
	fileService := services.NewFileService(cfg)
	taskService := services.NewTaskService(cfg)
	gitService := services.NewGitService(&cfg.Git)

	// 初始化WebSocket Hub
	wsHub := websocket.NewHub(&cfg.WebSocket)
	go wsHub.Run()

	// 设置任务服务的WebSocket Hub
	taskService.SetWebSocketHub(wsHub)

	// 启动任务调度器
	taskService.Start()
	defer taskService.Stop()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      setupRoutes(cfg, authService, fileService, taskService, gitService, wsHub),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 记录启动信息
	utils.Info("Starting CodeSage Go Backend Server")
	utils.Info("Configuration loaded successfully")
	utils.Info("Python Agent URL: %s", cfg.GetPythonAgentURL())

	// 启动服务器
	go func() {
		utils.Info("Starting HTTP server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Server failed to start: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced to shutdown: %v", err)
	}

	utils.Info("Server exited successfully")
}

func setupRoutes(cfg *config.Config, authService *services.AuthService, fileService *services.FileService, taskService *services.TaskService, gitService *services.GitService, wsHub *websocket.Hub) http.Handler {
	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("/api/v1/health", healthHandler)
	mux.HandleFunc("/api/v1/health/detailed", detailedHealthHandler(cfg))

	// 基础路由测试
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "CodeSage Go Backend", "version": "0.1.0", "docs": "/api/v1/health"}`)
	})

	// 认证路由（无需token）
	authHandler := handlers.NewAuthHandler(authService)
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/api/v1/auth/verify-reset-code", authHandler.VerifyResetCode)
	mux.HandleFunc("/api/v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("/api/v1/auth/github/login", authHandler.GithubLogin)

	// 需要认证的路由
	authMiddleware := middleware.AuthMiddleware(authService)

	// 用户资料路由（需要认证）
	mux.Handle("/api/v1/auth/profile", authMiddleware(http.HandlerFunc(authHandler.GetProfile)))
	mux.Handle("/api/v1/auth/profile/update", authMiddleware(http.HandlerFunc(authHandler.UpdateProfile)))
	mux.Handle("/api/v1/auth/password/change", authMiddleware(http.HandlerFunc(authHandler.ChangePassword)))
	mux.Handle("/api/v1/auth/github/bind", authMiddleware(http.HandlerFunc(authHandler.BindGithub)))
	mux.Handle("/api/v1/users/", authMiddleware(http.HandlerFunc(authHandler.GetUserByID)))

	// 文件处理路由（需要认证）
	fileHandler := handlers.NewFileHandler(fileService)
	mux.Handle("/api/v1/files/upload", authMiddleware(http.HandlerFunc(fileHandler.Upload)))
	mux.Handle("/api/v1/files/list", authMiddleware(http.HandlerFunc(fileHandler.List)))
	mux.Handle("/api/v1/files/", authMiddleware(http.HandlerFunc(fileHandler.Handle)))
	mux.Handle("/api/v1/files/batch", authMiddleware(http.HandlerFunc(fileHandler.BatchProcess)))

	// 任务管理路由（需要认证）
	taskHandler := handlers.NewTaskHandler(taskService)
	mux.Handle("/api/v1/tasks", authMiddleware(http.HandlerFunc(taskHandler.Handle)))
	mux.Handle("/api/v1/tasks/", authMiddleware(http.HandlerFunc(taskHandler.HandleDetail)))

	// Git相关路由（需要认证）
	gitHandler := handlers.NewGitHandler(taskService, gitService)
	mux.Handle("/api/v1/git/", authMiddleware(http.HandlerFunc(gitHandler.Handle)))

	// WebSocket路由（需要认证）
	wsHandler := websocket.NewWebSocketHandler(wsHub, &cfg.CORS)
	mux.Handle("/ws/agent/{taskId}", authMiddleware(http.HandlerFunc(wsHandler.HandleTaskWebSocket)))
	mux.Handle("/ws/progress/{taskId}", authMiddleware(http.HandlerFunc(wsHandler.HandleTaskWebSocket)))

	// 应用CORS中间件
	return middleware.CORSMiddleware(&cfg.CORS)(mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.Debug("Health check request from %s", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status": "healthy", "service": "go-backend", "timestamp": "`+time.Now().Format(time.RFC3339)+`"}`)
	utils.Debug("Health check response sent")
}

func detailedHealthHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Debug("Detailed health check request from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "healthy",
			"service": "go-backend",
			"version": "0.1.0",
			"timestamp": "%s",
			"config": {
				"server": {
					"host": "%s",
					"port": "%s"
				},
				"python_agent": {
					"url": "%s"
				}
			}
		}`, time.Now().Format(time.RFC3339), cfg.Server.Host, cfg.Server.Port, cfg.GetPythonAgentURL())
		utils.Debug("Detailed health check response sent")
	}
}
