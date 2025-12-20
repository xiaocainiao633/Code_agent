package models

import (
	"time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeAnalysis   TaskType = "analysis"
	TaskTypeConvert    TaskType = "convert"
	TaskTypeTest       TaskType = "test"
	TaskTypeBatch      TaskType = "batch"
	TaskTypeGitClone   TaskType = "git_clone"
	TaskTypeGitAnalyze TaskType = "git_analyze"
	TaskTypeGitHistory TaskType = "git_history"
	TaskTypeGitDiff    TaskType = "git_diff"
)

// Task 任务模型
type Task struct {
	ID          string                 `json:"id"`
	Type        TaskType               `json:"type"`
	Status      TaskStatus             `json:"status"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Params      map[string]interface{} `json:"params"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Progress    int                    `json:"progress"` // 0-100
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// TaskCreateRequest 任务创建请求
type TaskCreateRequest struct {
	Type        TaskType               `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Params      map[string]interface{} `json:"params"`
}

// TaskCreateResponse 任务创建响应
type TaskCreateResponse struct {
	TaskID  string `json:"task_id"`
	Message string `json:"message"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	Task Task `json:"task"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks []Task `json:"tasks"`
	Total int    `json:"total"`
}

// TaskResultResponse 任务结果响应
type TaskResultResponse struct {
	TaskID string                 `json:"task_id"`
	Result map[string]interface{} `json:"result"`
	Status TaskStatus             `json:"status"`
}