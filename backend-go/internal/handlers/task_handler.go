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

// TaskHandler 任务处理HTTP处理器
type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// Handle 处理任务相关操作
func (h *TaskHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 解析任务ID从路径
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks")
	path = strings.TrimPrefix(path, "/")
	
	if path == "" || path == "/" {
		// 没有任务ID，处理集合操作
		h.handleCollection(w, r)
		return
	}
	
	taskID := path
	utils.Debug("Task operation request for ID: %s, Method: %s", taskID, r.Method)
	
	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, taskID)
	case http.MethodDelete:
		h.cancelTask(w, r, taskID)
	default:
		utils.Warn("Invalid method for task operation: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleDetail 处理任务详情操作
func (h *TaskHandler) HandleDetail(w http.ResponseWriter, r *http.Request) {
	// 解析任务ID从路径
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	if path == "" || path == "/" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}
	
	// 检查是否是结果请求
	if strings.HasSuffix(path, "/result") {
		taskID := strings.TrimSuffix(path, "/result")
		h.getTaskResult(w, r, taskID)
		return
	}
	
	// 普通任务详情
	taskID := path
	utils.Debug("Task detail request for ID: %s, Method: %s", taskID, r.Method)
	
	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, taskID)
	default:
		utils.Warn("Invalid method for task detail: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCollection 处理任务集合操作
func (h *TaskHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodGet:
		h.listTasks(w, r)
	default:
		utils.Warn("Invalid method for task collection: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// createTask 创建新任务
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	utils.Info("Task creation request received")
	
	var request models.TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error("Failed to decode task creation request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// 验证请求参数
	if request.Type == "" {
		utils.Warn("Task type not specified")
		http.Error(w, "Task type is required", http.StatusBadRequest)
		return
	}
	
	if request.Name == "" {
		utils.Warn("Task name not specified")
		http.Error(w, "Task name is required", http.StatusBadRequest)
		return
	}
	
	utils.Info("Creating task: %s (type: %s)", request.Name, request.Type)
	
	// 创建任务
	task, err := h.taskService.CreateTask(request.Type, request.Name, request.Description, request.Params)
	if err != nil {
		utils.Error("Failed to create task: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := models.TaskCreateResponse{
		TaskID:  task.ID,
		Message: fmt.Sprintf("Task %s created successfully", task.Name),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode task creation response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Info("Task created successfully: %s (ID: %s)", task.Name, task.ID)
}

// listTasks 获取任务列表
func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	utils.Debug("Task list request received")
	
	tasks, err := h.taskService.ListTasks()
	if err != nil {
		utils.Error("Failed to list tasks: %v", err)
		http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}
	
	response := models.TaskListResponse{
		Tasks: tasks,
		Total: len(tasks),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode task list response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Debug("Task list response sent. Tasks: %d", len(tasks))
}

// getTask 获取单个任务
func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, taskID string) {
	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		utils.Error("Failed to get task %s: %v", taskID, err)
		http.Error(w, fmt.Sprintf("Task not found: %s", taskID), http.StatusNotFound)
		return
	}
	
	response := models.TaskResponse{
		Task: *task,
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode task response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Debug("Task response sent for ID: %s", taskID)
}

// getTaskResult 获取任务结果
func (h *TaskHandler) getTaskResult(w http.ResponseWriter, r *http.Request, taskID string) {
	result, err := h.taskService.GetTaskResult(taskID)
	if err != nil {
		utils.Error("Failed to get task result %s: %v", taskID, err)
		http.Error(w, fmt.Sprintf("Failed to get task result: %v", err), http.StatusBadRequest)
		return
	}
	
	response := models.TaskResultResponse{
		TaskID: taskID,
		Result: result,
		Status: models.TaskStatusCompleted,
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode task result response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Debug("Task result response sent for ID: %s", taskID)
}

// cancelTask 取消任务
func (h *TaskHandler) cancelTask(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.taskService.CancelTask(taskID)
	if err != nil {
		utils.Error("Failed to cancel task %s: %v", taskID, err)
		http.Error(w, fmt.Sprintf("Failed to cancel task: %v", err), http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"task_id": taskID,
		"message": fmt.Sprintf("Task %s cancelled successfully", taskID),
		"status":  "cancelled",
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode cancel response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Info("Task cancelled successfully: %s", taskID)
}