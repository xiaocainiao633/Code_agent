package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/services"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// FileHandler 文件处理HTTP处理器
type FileHandler struct {
	fileService *services.FileService
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// Upload 处理文件上传
func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	utils.Info("File upload request received")
	
	// 验证请求方法
	if r.Method != http.MethodPost {
		utils.Warn("Invalid method for upload: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 解析multipart表单
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		utils.Error("Failed to parse multipart form: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	
	// 获取上传的文件
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		utils.Warn("No files found in upload request")
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}
	
	utils.Info("Processing %d uploaded files", len(files))
	
	// 处理文件上传
	uploadedFiles, err := h.fileService.UploadFiles(files)
	if err != nil {
		utils.Error("File upload failed: %v", err)
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	// 返回成功响应
	response := models.FileListResponse{
		Files: uploadedFiles,
		Total: len(uploadedFiles),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Info("File upload completed successfully. Files: %d", len(uploadedFiles))
}

// List 获取文件列表
func (h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	utils.Debug("File list request received")
	
	if r.Method != http.MethodGet {
		utils.Warn("Invalid method for list: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	files, err := h.fileService.ListFiles()
	if err != nil {
		utils.Error("Failed to list files: %v", err)
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}
	
	response := models.FileListResponse{
		Files: files,
		Total: len(files),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Debug("File list response sent. Files: %d", len(files))
}

// Handle 处理文件相关操作
func (h *FileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 解析文件ID从路径
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/files/")
	if path == "" || path == "/" {
		http.Error(w, "File ID required", http.StatusBadRequest)
		return
	}
	
	fileID := filepath.Base(path)
	utils.Debug("File operation request for ID: %s, Method: %s", fileID, r.Method)
	
	switch r.Method {
	case http.MethodGet:
		h.getFile(w, r, fileID)
	case http.MethodDelete:
		h.deleteFile(w, r, fileID)
	default:
		utils.Warn("Invalid method for file operation: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getFile 获取单个文件
func (h *FileHandler) getFile(w http.ResponseWriter, r *http.Request, fileID string) {
	file, err := h.fileService.GetFile(fileID)
	if err != nil {
		utils.Error("Failed to get file %s: %v", fileID, err)
		http.Error(w, fmt.Sprintf("File not found: %s", fileID), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(file); err != nil {
		utils.Error("Failed to encode file response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Debug("File response sent for ID: %s", fileID)
}

// deleteFile 删除文件
func (h *FileHandler) deleteFile(w http.ResponseWriter, r *http.Request, fileID string) {
	err := h.fileService.DeleteFile(fileID)
	if err != nil {
		utils.Error("Failed to delete file %s: %v", fileID, err)
		http.Error(w, fmt.Sprintf("Failed to delete file: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := models.FileDeleteResponse{
		Success: true,
		Message: fmt.Sprintf("File %s deleted successfully", fileID),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode delete response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Info("File deleted successfully: %s", fileID)
}

// BatchProcess 批量处理文件
func (h *FileHandler) BatchProcess(w http.ResponseWriter, r *http.Request) {
	utils.Info("Batch file process request received")
	
	if r.Method != http.MethodPost {
		utils.Warn("Invalid method for batch process: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request models.BatchFileProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error("Failed to decode batch process request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if len(request.FileIDs) == 0 {
		utils.Warn("No file IDs provided for batch processing")
		http.Error(w, "No file IDs provided", http.StatusBadRequest)
		return
	}
	
	utils.Info("Processing %d files in batch", len(request.FileIDs))
	
	taskID, err := h.fileService.BatchProcessFiles(request.FileIDs, request.Options)
	if err != nil {
		utils.Error("Batch processing failed: %v", err)
		http.Error(w, fmt.Sprintf("Batch processing failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := models.BatchFileProcessResponse{
		TaskID:  taskID,
		Message: fmt.Sprintf("Batch processing started for %d files", len(request.FileIDs)),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.Error("Failed to encode batch process response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	
	utils.Info("Batch processing started with task ID: %s", taskID)
}