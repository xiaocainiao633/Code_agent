package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/services"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// GitHandler Git相关HTTP处理器
type GitHandler struct {
	taskService *services.TaskService
	gitService  *services.GitService
}

// NewGitHandler 创建Git处理器
func NewGitHandler(taskService *services.TaskService, gitService *services.GitService) *GitHandler {
	return &GitHandler{
		taskService: taskService,
		gitService:  gitService,
	}
}

// Handle 处理Git相关操作
func (h *GitHandler) Handle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/git")
	path = strings.TrimPrefix(path, "/")

	switch {
	case path == "clone":
		h.handleClone(w, r)
	case path == "analyze":
		h.handleAnalyze(w, r)
	case strings.HasPrefix(path, "history/"):
		filePath := strings.TrimPrefix(path, "history/")
		h.handleHistory(w, r, filePath)
	case strings.HasPrefix(path, "diff/"):
		h.handleDiff(w, r)
	default:
		http.Error(w, "Invalid Git operation", http.StatusBadRequest)
	}
}

// handleClone 处理Git克隆
func (h *GitHandler) handleClone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.Info("Git clone request received")

	var request models.GitCloneParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error("Failed to decode git clone request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证参数
	if request.RemoteURL == "" {
		utils.Warn("Remote URL not specified")
		http.Error(w, "Remote URL is required", http.StatusBadRequest)
		return
	}

	if request.TargetPath == "" {
		utils.Warn("Target path not specified")
		http.Error(w, "Target path is required", http.StatusBadRequest)
		return
	}

	utils.Info("Creating git clone task for: %s", request.RemoteURL)

	// 创建任务
	task, err := h.taskService.CreateTask(
		models.TaskTypeGitClone,
		fmt.Sprintf("Clone %s", request.RemoteURL),
		fmt.Sprintf("Clone Git repository from %s to %s", request.RemoteURL, request.TargetPath),
		map[string]interface{}{
			"remote_url":  request.RemoteURL,
			"target_path": request.TargetPath,
		},
	)
	if err != nil {
		utils.Error("Failed to create git clone task: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.TaskCreateResponse{
		TaskID:  task.ID,
		Message: fmt.Sprintf("Git clone task created successfully for %s", request.RemoteURL),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode git clone response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	utils.Info("Git clone task created successfully: %s (ID: %s)", request.RemoteURL, task.ID)
}

// handleAnalyze 处理Git分析
func (h *GitHandler) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.Info("Git analyze request received")

	var request models.GitTaskParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error("Failed to decode git analyze request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证参数
	if request.RepoPath == "" {
		utils.Warn("Repository path not specified")
		http.Error(w, "Repository path is required", http.StatusBadRequest)
		return
	}

	utils.Info("Creating git analyze task for: %s", request.RepoPath)

	// 创建任务
	task, err := h.taskService.CreateTask(
		models.TaskTypeGitAnalyze,
		fmt.Sprintf("Analyze %s", request.RepoPath),
		fmt.Sprintf("Analyze Git repository at %s", request.RepoPath),
		map[string]interface{}{
			"repo_path":        request.RepoPath,
			"remote_url":       request.RemoteURL,
			"clone_if_not_exists": request.CloneIfNotExists,
		},
	)
	if err != nil {
		utils.Error("Failed to create git analyze task: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.TaskCreateResponse{
		TaskID:  task.ID,
		Message: fmt.Sprintf("Git analyze task created successfully for %s", request.RepoPath),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode git analyze response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	utils.Info("Git analyze task created successfully: %s (ID: %s)", request.RepoPath, task.ID)
}

// handleHistory 处理Git文件历史
func (h *GitHandler) handleHistory(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.Info("Git history request received for file: %s", filePath)

	var request models.GitFileHistoryParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error("Failed to decode git history request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证参数
	if request.RepoPath == "" {
		utils.Warn("Repository path not specified")
		http.Error(w, "Repository path is required", http.StatusBadRequest)
		return
	}

	if filePath == "" {
		filePath = request.FilePath
	}

	if filePath == "" {
		utils.Warn("File path not specified")
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	utils.Info("Creating git history task for: %s in %s", filePath, request.RepoPath)

	// 创建任务
	task, err := h.taskService.CreateTask(
		models.TaskTypeGitHistory,
		fmt.Sprintf("History %s", filePath),
		fmt.Sprintf("Get history for file %s in repository %s", filePath, request.RepoPath),
		map[string]interface{}{
			"repo_path": request.RepoPath,
			"file_path": filePath,
		},
	)
	if err != nil {
		utils.Error("Failed to create git history task: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.TaskCreateResponse{
		TaskID:  task.ID,
		Message: fmt.Sprintf("Git history task created successfully for %s", filePath),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode git history response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	utils.Info("Git history task created successfully: %s (ID: %s)", filePath, task.ID)
}

// handleDiff 处理Git差异
func (h *GitHandler) handleDiff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.Info("Git diff request received")

	var request models.GitDiffParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error("Failed to decode git diff request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证参数
	if request.RepoPath == "" {
		utils.Warn("Repository path not specified")
		http.Error(w, "Repository path is required", http.StatusBadRequest)
		return
	}

	if request.FilePath == "" {
		utils.Warn("File path not specified")
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	if request.FromCommit == "" {
		utils.Warn("From commit not specified")
		http.Error(w, "From commit is required", http.StatusBadRequest)
		return
	}

	if request.ToCommit == "" {
		utils.Warn("To commit not specified")
		http.Error(w, "To commit is required", http.StatusBadRequest)
		return
	}

	utils.Info("Creating git diff task for: %s between %s and %s", request.FilePath, request.FromCommit, request.ToCommit)

	// 创建任务
	task, err := h.taskService.CreateTask(
		models.TaskTypeGitDiff,
		fmt.Sprintf("Diff %s", request.FilePath),
		fmt.Sprintf("Get diff for file %s between commits %s and %s", request.FilePath, request.FromCommit, request.ToCommit),
		map[string]interface{}{
			"repo_path":   request.RepoPath,
			"file_path":   request.FilePath,
			"from_commit": request.FromCommit,
			"to_commit":   request.ToCommit,
		},
	)
	if err != nil {
		utils.Error("Failed to create git diff task: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.TaskCreateResponse{
		TaskID:  task.ID,
		Message: fmt.Sprintf("Git diff task created successfully for %s", request.FilePath),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode git diff response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	utils.Info("Git diff task created successfully: %s (ID: %s)", request.FilePath, task.ID)
}