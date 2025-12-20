package models

import (
	"time"
)

// File 文件模型
type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Extension   string    `json:"extension"`
	Content     string    `json:"content,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	Files []FileUpload `json:"files"`
}

// FileUpload 单个文件上传
type FileUpload struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Files []File `json:"files"`
	Total int    `json:"total"`
}

// FileDeleteResponse 文件删除响应
type FileDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// BatchFileProcessRequest 批量文件处理请求
type BatchFileProcessRequest struct {
	FileIDs  []string          `json:"file_ids"`
	Options  map[string]interface{} `json:"options"`
}

// BatchFileProcessResponse 批量文件处理响应
type BatchFileProcessResponse struct {
	TaskID  string `json:"task_id"`
	Message string `json:"message"`
}