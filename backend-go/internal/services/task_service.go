package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/database"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/websocket"
)

// TaskService 任务调度服务
type TaskService struct {
	config            *config.TaskSchedulerConfig
	tasks             map[string]*models.Task
	taskQueue         chan *models.Task
	workers           int
	runningTasks      map[string]bool
	mu                sync.RWMutex
	wg                sync.WaitGroup
	stopCh            chan struct{}
	pythonAgentClient *PythonAgentClient
	gitService        *GitService
	wsHub             *websocket.Hub
	taskRepo          *database.TaskRepository
}

// NewTaskService 创建任务服务
func NewTaskService(cfg *config.Config) *TaskService {
	return &TaskService{
		config:            &cfg.TaskScheduler,
		tasks:             make(map[string]*models.Task),
		taskQueue:         make(chan *models.Task, cfg.TaskScheduler.MaxConcurrentTasks*2),
		workers:           cfg.TaskScheduler.MaxConcurrentTasks,
		runningTasks:      make(map[string]bool),
		stopCh:            make(chan struct{}),
		pythonAgentClient: NewPythonAgentClient(&cfg.PythonAgent),
		gitService:        NewGitService(&cfg.Git),
		wsHub:             nil, // 将在外部设置
		taskRepo:          database.NewTaskRepository(),
	}
}

// SetWebSocketHub 设置WebSocket Hub
func (s *TaskService) SetWebSocketHub(hub *websocket.Hub) {
	s.wsHub = hub
}

// Start 启动任务调度器
func (s *TaskService) Start() {
	utils.Info("Starting task scheduler with %d workers", s.workers)

	// 启动工作协程
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	// 启动任务清理协程
	s.wg.Add(1)
	go s.cleanupWorker()

	utils.Info("Task scheduler started successfully")
}

// Stop 停止任务调度器
func (s *TaskService) Stop() {
	utils.Info("Stopping task scheduler")
	close(s.stopCh)
	s.wg.Wait()
	utils.Info("Task scheduler stopped")
}

// CreateTask 创建新任务
func (s *TaskService) CreateTask(taskType models.TaskType, name, description string, params map[string]interface{}) (*models.Task, error) {
	utils.Info("Creating new task: %s (type: %s)", name, taskType)

	task := &models.Task{
		ID:          generateTaskID(),
		Type:        taskType,
		Status:      models.TaskStatusPending,
		Name:        name,
		Description: description,
		Params:      params,
		Progress:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存任务
	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()

	// 将任务加入队列
	select {
	case s.taskQueue <- task:
		utils.Info("Task %s added to queue", task.ID)
		return task, nil
	default:
		utils.Warn("Task queue is full, task %s will be processed later", task.ID)
		// 启动后台协程等待队列空间
		go func() {
			s.taskQueue <- task
		}()
		return task, nil
	}
}

// GetTask 获取任务信息
func (s *TaskService) GetTask(taskID string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return task, nil
}

// GetTaskResult 获取任务结果
func (s *TaskService) GetTaskResult(taskID string) (map[string]interface{}, error) {
	task, err := s.GetTask(taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != models.TaskStatusCompleted {
		return nil, fmt.Errorf("task is not completed yet: %s", taskID)
	}

	return task.Result, nil
}

// ListTasks 获取任务列表
func (s *TaskService) ListTasks() ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

// CancelTask 取消任务
func (s *TaskService) CancelTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if task.Status == models.TaskStatusRunning {
		return fmt.Errorf("cannot cancel running task: %s", taskID)
	}

	if task.Status == models.TaskStatusCompleted || task.Status == models.TaskStatusFailed {
		return fmt.Errorf("task already finished: %s", taskID)
	}

	task.Status = models.TaskStatusCancelled
	task.UpdatedAt = time.Now()

	utils.Info("Task cancelled: %s", taskID)
	return nil
}

// worker 工作协程
func (s *TaskService) worker(id int) {
	defer s.wg.Done()
	utils.Info("Task worker %d started", id)

	for {
		select {
		case task := <-s.taskQueue:
			s.processTask(task)
		case <-s.stopCh:
			utils.Info("Task worker %d stopped", id)
			return
		}
	}
}

// processTask 处理任务
func (s *TaskService) processTask(task *models.Task) {
	utils.Info("Worker processing task: %s (type: %s)", task.ID, task.Type)

	// 更新任务状态
	s.updateTaskStatus(task.ID, models.TaskStatusRunning, 0)

	now := time.Now()
	task.StartedAt = &now

	// 标记为运行中
	s.mu.Lock()
	s.runningTasks[task.ID] = true
	s.mu.Unlock()

	// 执行任务
	result, err := s.executeTask(task)

	// 移除运行标记
	s.mu.Lock()
	delete(s.runningTasks, task.ID)
	s.mu.Unlock()

	// 更新任务结果
	if err != nil {
		s.updateTaskError(task.ID, err.Error())
		utils.Error("Task %s failed: %v", task.ID, err)
	} else {
		s.updateTaskResult(task.ID, result)
		utils.Info("Task %s completed successfully", task.ID)
	}
}

// executeTask 执行具体任务
func (s *TaskService) executeTask(task *models.Task) (map[string]interface{}, error) {
	utils.Debug("Executing task: %s (type: %s)", task.ID, task.Type)

	switch task.Type {
	case models.TaskTypeAnalysis:
		return s.executeAnalysisTask(task)
	case models.TaskTypeConvert:
		return s.executeConvertTask(task)
	case models.TaskTypeTest:
		return s.executeTestTask(task)
	case models.TaskTypeBatch:
		return s.executeBatchTask(task)
	case models.TaskTypeGitClone:
		return s.executeGitCloneTask(task)
	case models.TaskTypeGitAnalyze:
		return s.executeGitAnalyzeTask(task)
	case models.TaskTypeGitHistory:
		return s.executeGitHistoryTask(task)
	case models.TaskTypeGitDiff:
		return s.executeGitDiffTask(task)
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

// executeAnalysisTask 执行分析任务
func (s *TaskService) executeAnalysisTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing analysis task: %s", task.ID)

	// 获取代码和语言类型
	code, ok := task.Params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid code parameter")
	}

	language, _ := task.Params["language"].(string)
	if language == "" {
		language = "python"
	}

	// 模拟进度更新
	for i := 0; i <= 80; i += 20 {
		s.updateTaskProgress(task.ID, i)
		time.Sleep(200 * time.Millisecond)
	}

	// 调用Python AI Agent进行分析
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result map[string]interface{}
	var err error

	if language == "python" && task.Params["python2_analysis"] == true {
		result, err = s.pythonAgentClient.AnalyzePython2Code(ctx, code)
	} else {
		result, err = s.pythonAgentClient.AnalyzeCode(ctx, code, language)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to call Python agent: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	return result, nil
}

// executeConvertTask 执行转换任务
func (s *TaskService) executeConvertTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing convert task: %s", task.ID)

	// 获取必要参数
	code, ok := task.Params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid code parameter")
	}

	fromVersion, _ := task.Params["from_version"].(string)
	toVersion, _ := task.Params["to_version"].(string)
	if fromVersion == "" {
		fromVersion = "python2"
	}
	if toVersion == "" {
		toVersion = "python3"
	}

	options, _ := task.Params["options"].(map[string]interface{})

	// 模拟进度更新
	for i := 0; i <= 80; i += 20 {
		s.updateTaskProgress(task.ID, i)
		time.Sleep(300 * time.Millisecond)
	}

	// 调用Python AI Agent进行转换
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := s.pythonAgentClient.ConvertCode(ctx, code, fromVersion, toVersion, options)
	if err != nil {
		return nil, fmt.Errorf("failed to call Python agent: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	return result, nil
}

// executeTestTask 执行测试任务
func (s *TaskService) executeTestTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing test task: %s", task.ID)

	// 获取必要参数
	code, ok := task.Params["code"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid code parameter")
	}

	testType, _ := task.Params["test_type"].(string)
	if testType == "" {
		testType = "unit"
	}

	framework, _ := task.Params["framework"].(string)
	if framework == "" {
		framework = "pytest"
	}

	// 模拟进度更新
	for i := 0; i <= 80; i += 20 {
		s.updateTaskProgress(task.ID, i)
		time.Sleep(200 * time.Millisecond)
	}

	// 调用Python AI Agent生成测试
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	result, err := s.pythonAgentClient.GenerateTests(ctx, code, testType, framework)
	if err != nil {
		return nil, fmt.Errorf("failed to call Python agent: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	return result, nil
}

// executeGitCloneTask 执行Git克隆任务
func (s *TaskService) executeGitCloneTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing Git clone task: %s", task.ID)

	// 获取参数
	remoteURL, ok := task.Params["remote_url"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid remote_url parameter")
	}

	targetPath, ok := task.Params["target_path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid target_path parameter")
	}

	// 模拟进度更新
	s.updateTaskProgress(task.ID, 20)
	time.Sleep(500 * time.Millisecond)

	// 执行克隆
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err := s.gitService.CloneRepository(ctx, remoteURL, targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	result := map[string]interface{}{
		"status":      "success",
		"remote_url":  remoteURL,
		"target_path": targetPath,
		"message":     "Repository cloned successfully",
	}

	return result, nil
}

// executeGitAnalyzeTask 执行Git分析任务
func (s *TaskService) executeGitAnalyzeTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing Git analyze task: %s", task.ID)

	// 获取参数
	repoPath, ok := task.Params["repo_path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid repo_path parameter")
	}

	// 检查是否需要克隆
	if cloneIfNotExists, _ := task.Params["clone_if_not_exists"].(bool); cloneIfNotExists {
		if remoteURL, ok := task.Params["remote_url"].(string); ok {
			// 检查仓库是否存在
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				s.updateTaskProgress(task.ID, 10)
				// 克隆仓库
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
				err := s.gitService.CloneRepository(ctx, remoteURL, repoPath)
				cancel()
				if err != nil {
					return nil, fmt.Errorf("failed to clone repository: %w", err)
				}
			}
		}
	}

	// 模拟进度更新
	for i := 20; i <= 80; i += 20 {
		s.updateTaskProgress(task.ID, i)
		time.Sleep(300 * time.Millisecond)
	}

	// 执行分析
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	result, err := s.gitService.AnalyzeRepository(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze repository: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	return map[string]interface{}{
		"status":   "success",
		"analysis": result,
		"message":  "Repository analyzed successfully",
	}, nil
}

// executeGitHistoryTask 执行Git历史任务
func (s *TaskService) executeGitHistoryTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing Git history task: %s", task.ID)

	// 获取参数
	repoPath, ok := task.Params["repo_path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid repo_path parameter")
	}

	filePath, ok := task.Params["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid file_path parameter")
	}

	// 模拟进度更新
	s.updateTaskProgress(task.ID, 30)
	time.Sleep(500 * time.Millisecond)

	// 获取文件历史
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	history, err := s.gitService.GetFileHistory(ctx, repoPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file history: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	result := map[string]interface{}{
		"status":    "success",
		"repo_path": repoPath,
		"file_path": filePath,
		"history":   history,
		"count":     len(history),
	}

	return result, nil
}

// executeGitDiffTask 执行Git差异任务
func (s *TaskService) executeGitDiffTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing Git diff task: %s", task.ID)

	// 获取参数
	repoPath, ok := task.Params["repo_path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid repo_path parameter")
	}

	filePath, ok := task.Params["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid file_path parameter")
	}

	fromCommit, ok := task.Params["from_commit"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid from_commit parameter")
	}

	toCommit, ok := task.Params["to_commit"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid to_commit parameter")
	}

	// 模拟进度更新
	s.updateTaskProgress(task.ID, 40)
	time.Sleep(300 * time.Millisecond)

	// 获取差异
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	diff, err := s.gitService.GetDiff(ctx, repoPath, filePath, fromCommit, toCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}

	// 完成进度
	s.updateTaskProgress(task.ID, 100)

	result := map[string]interface{}{
		"status":      "success",
		"repo_path":   repoPath,
		"file_path":   filePath,
		"from_commit": fromCommit,
		"to_commit":   toCommit,
		"diff":        diff,
	}

	return result, nil
}

// executeBatchTask 执行批量任务
func (s *TaskService) executeBatchTask(task *models.Task) (map[string]interface{}, error) {
	utils.Info("Executing batch task: %s", task.ID)

	// 从参数中获取文件ID列表
	fileIDs, ok := task.Params["file_ids"].([]string)
	if !ok {
		return nil, fmt.Errorf("missing file_ids in batch task parameters")
	}

	utils.Info("Processing %d files in batch task: %s", len(fileIDs), task.ID)

	// 模拟批量处理进度
	for i, fileID := range fileIDs {
		progress := int(float64(i+1) / float64(len(fileIDs)) * 100)
		s.updateTaskProgress(task.ID, progress)
		utils.Debug("Processing file %d/%d: %s", i+1, len(fileIDs), fileID)
		time.Sleep(1 * time.Second) // 模拟处理时间
	}

	// 返回批量处理结果
	result := map[string]interface{}{
		"processed_files": len(fileIDs),
		"file_ids":        fileIDs,
		"status":          "completed",
	}

	return result, nil
}

// updateTaskStatus 更新任务状态
func (s *TaskService) updateTaskStatus(taskID string, status models.TaskStatus, progress int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Status = status
		task.Progress = progress
		task.UpdatedAt = time.Now()

		if status == models.TaskStatusCompleted || status == models.TaskStatusFailed {
			now := time.Now()
			task.CompletedAt = &now
		}

		utils.Debug("Task %s status updated: %s (progress: %d%%)", taskID, status, progress)
	}
}

// updateTaskProgress 更新任务进度
func (s *TaskService) updateTaskProgress(taskID string, progress int) {
	s.updateTaskStatus(taskID, s.tasks[taskID].Status, progress)
}

// updateTaskResult 更新任务结果
func (s *TaskService) updateTaskResult(taskID string, result map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Result = result
		task.Status = models.TaskStatusCompleted
		task.Progress = 100
		task.UpdatedAt = time.Now()

		now := time.Now()
		task.CompletedAt = &now

		utils.Info("Task %s completed with result", taskID)
	}
}

// updateTaskError 更新任务错误
func (s *TaskService) updateTaskError(taskID string, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Error = errorMsg
		task.Status = models.TaskStatusFailed
		task.Progress = 100
		task.UpdatedAt = time.Now()

		now := time.Now()
		task.CompletedAt = &now

		utils.Error("Task %s failed with error: %s", taskID, errorMsg)
	}
}

// cleanupWorker 清理工作协程
func (s *TaskService) cleanupWorker() {
	defer s.wg.Done()
	utils.Info("Task cleanup worker started")

	ticker := time.NewTicker(1 * time.Hour) // 每小时清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupOldTasks()
		case <-s.stopCh:
			utils.Info("Task cleanup worker stopped")
			return
		}
	}
}

// cleanupOldTasks 清理旧任务
func (s *TaskService) cleanupOldTasks() {
	utils.Debug("Starting cleanup of old tasks")

	s.mu.Lock()
	defer s.mu.Unlock()

	cutoffTime := time.Now().Add(-s.config.ResultRetention)
	deletedCount := 0

	for taskID, task := range s.tasks {
		if task.CompletedAt != nil && task.CompletedAt.Before(cutoffTime) {
			delete(s.tasks, taskID)
			deletedCount++
		}
	}

	if deletedCount > 0 {
		utils.Info("Cleaned up %d old tasks", deletedCount)
	}
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

// GetRunningTasks 获取正在运行的任务数量
func (s *TaskService) GetRunningTasks() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.runningTasks)
}

// GetQueueSize 获取队列大小
func (s *TaskService) GetQueueSize() int {
	return len(s.taskQueue)
}
