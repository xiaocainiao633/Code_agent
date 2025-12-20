package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// FileService 文件处理服务
type FileService struct {
	config      *config.FileProcessorConfig
	allowedExts map[string]bool
	tempDir     string
}

// NewFileService 创建文件服务
func NewFileService(cfg *config.Config) *FileService {
	allowedExts := make(map[string]bool)
	for _, ext := range cfg.Git.AllowedExtensions {
		allowedExts[ext] = true
	}

	return &FileService{
		config:      &cfg.FileProcessor,
		allowedExts: allowedExts,
		tempDir:     cfg.FileProcessor.TempDir,
	}
}

// UploadFiles 处理文件上传
func (s *FileService) UploadFiles(files []*multipart.FileHeader) ([]models.File, error) {
	utils.Info("Starting file upload process for %d files", len(files))
	
	var uploadedFiles []models.File
	
	for _, fileHeader := range files {
		file, err := s.processUploadedFile(fileHeader)
		if err != nil {
			utils.Error("Failed to process file %s: %v", fileHeader.Filename, err)
			return nil, fmt.Errorf("failed to process file %s: %w", fileHeader.Filename, err)
		}
		
		uploadedFiles = append(uploadedFiles, file)
		utils.Info("Successfully uploaded file: %s (ID: %s)", file.Name, file.ID)
	}
	
	utils.Info("File upload completed. Total files: %d", len(uploadedFiles))
	return uploadedFiles, nil
}

// processUploadedFile 处理单个上传文件
func (s *FileService) processUploadedFile(fileHeader *multipart.FileHeader) (models.File, error) {
	// 验证文件大小
	if fileHeader.Size > s.config.MaxUploadSize {
		return models.File{}, fmt.Errorf("file size %d exceeds maximum allowed size %d", fileHeader.Size, s.config.MaxUploadSize)
	}
	
	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !s.isAllowedExtension(ext) {
		return models.File{}, fmt.Errorf("file extension %s is not allowed", ext)
	}
	
	// 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return models.File{}, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()
	
	// 创建临时文件
	tempFileName := fmt.Sprintf("upload_%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	tempFilePath := filepath.Join(s.tempDir, tempFileName)
	
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return models.File{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer dst.Close()
	
	// 复制文件内容
	content, err := io.ReadAll(src)
	if err != nil {
		os.Remove(tempFilePath) // 清理临时文件
		return models.File{}, fmt.Errorf("failed to read file content: %w", err)
	}
	
	// 写入临时文件
	if _, err := dst.Write(content); err != nil {
		os.Remove(tempFilePath) // 清理临时文件
		return models.File{}, fmt.Errorf("failed to write temp file: %w", err)
	}
	
	// 创建文件模型
	file := models.File{
		ID:          generateFileID(),
		Name:        fileHeader.Filename,
		Path:        tempFilePath,
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Extension:   ext,
		Content:     string(content),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	return file, nil
}

// GetFile 获取文件信息
func (s *FileService) GetFile(fileID string) (*models.File, error) {
	utils.Debug("Getting file info for ID: %s", fileID)
	
	// TODO: 从存储中获取文件信息
	// 这里应该实现文件存储和检索逻辑
	return nil, fmt.Errorf("file not found: %s", fileID)
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(fileID string) error {
	utils.Info("Deleting file: %s", fileID)
	
	// TODO: 实现文件删除逻辑
	// 1. 从存储中获取文件路径
	// 2. 删除物理文件
	// 3. 从索引中移除
	
	return fmt.Errorf("file deletion not implemented yet")
}

// ListFiles 列出所有文件
func (s *FileService) ListFiles() ([]models.File, error) {
	utils.Debug("Listing all files")
	
	// TODO: 实现文件列表逻辑
	// 这里应该返回存储中的所有文件
	
	return []models.File{}, nil
}

// BatchProcessFiles 批量处理文件
func (s *FileService) BatchProcessFiles(fileIDs []string, options map[string]interface{}) (string, error) {
	utils.Info("Starting batch processing for %d files", len(fileIDs))
	
	// TODO: 实现批量处理逻辑
	// 1. 验证所有文件ID
	// 2. 创建处理任务
	// 3. 返回任务ID
	
	return "", fmt.Errorf("batch processing not implemented yet")
}

// CleanupTempFiles 清理临时文件
func (s *FileService) CleanupTempFiles() error {
	utils.Info("Starting temporary file cleanup")
	
	// 清理超过保留时间的临时文件
	cutoffTime := time.Now().Add(-1 * s.config.CleanupInterval)
	
	err := filepath.Walk(s.tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		
		// 检查文件是否过期
		if info.ModTime().Before(cutoffTime) {
			utils.Debug("Cleaning up old temp file: %s", path)
			if err := os.Remove(path); err != nil {
				utils.Error("Failed to remove temp file %s: %v", path, err)
				return err
			}
		}
		
		return nil
	})
	
	if err != nil {
		utils.Error("Cleanup failed: %v", err)
		return err
	}
	
	utils.Info("Temporary file cleanup completed")
	return nil
}

// isAllowedExtension 检查文件扩展名是否允许
func (s *FileService) isAllowedExtension(ext string) bool {
	return s.allowedExts[ext]
}

// generateFileID 生成文件ID
func generateFileID() string {
	return fmt.Sprintf("file_%d", time.Now().UnixNano())
}

// ValidateFileContent 验证文件内容
func (s *FileService) ValidateFileContent(content string, fileType string) error {
	utils.Debug("Validating file content for type: %s", fileType)
	
	// 基本的文件内容验证
	if len(content) == 0 {
		return fmt.Errorf("file content is empty")
	}
	
	// 检查文件大小（内容长度）
	if int64(len(content)) > s.config.MaxUploadSize {
		return fmt.Errorf("file content size exceeds maximum allowed size")
	}
	
	// TODO: 根据文件类型进行更详细的验证
	// - Python文件语法检查
	// - JavaScript语法检查
	// - 编码检测等
	
	return nil
}