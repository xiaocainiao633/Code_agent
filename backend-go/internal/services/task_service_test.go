package services

import (
	"context"
	"testing"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
)

func TestTaskService_CreateTask(t *testing.T) {
	cfg := &config.Config{
		TaskScheduler: config.TaskSchedulerConfig{
			MaxConcurrentTasks: 5,
			TaskTimeout:        10 * time.Minute,
			ResultRetention:    24 * time.Hour,
		},
		PythonAgent: config.PythonAgentConfig{
			Host:       "localhost",
			Port:       "8000",
			Timeout:    30 * time.Second,
			RetryCount: 3,
		},
		Git: config.GitConfig{
			CloneTimeout:      5 * time.Minute,
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".py", ".js", ".go"},
		},
	}

	taskService := NewTaskService(cfg)
	
	// 测试创建任务
	task, err := taskService.CreateTask(
		models.TaskTypeAnalysis,
		"Test Task",
		"Test description",
		map[string]interface{}{
			"code":     "print('test')",
			"language": "python",
		},
	)
	
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	
	if task.Name != "Test Task" {
		t.Errorf("Expected task name 'Test Task', got '%s'", task.Name)
	}
	
	if task.Type != models.TaskTypeAnalysis {
		t.Errorf("Expected task type 'analysis', got '%s'", task.Type)
	}
	
	if task.Status != models.TaskStatusPending {
		t.Errorf("Expected task status 'pending', got '%s'", task.Status)
	}
}

func TestTaskService_GetTask(t *testing.T) {
	cfg := &config.Config{
		TaskScheduler: config.TaskSchedulerConfig{
			MaxConcurrentTasks: 5,
			TaskTimeout:        10 * time.Minute,
			ResultRetention:    24 * time.Hour,
		},
		PythonAgent: config.PythonAgentConfig{
			Host:       "localhost",
			Port:       "8000",
			Timeout:    30 * time.Second,
			RetryCount: 3,
		},
		Git: config.GitConfig{
			CloneTimeout:      5 * time.Minute,
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".py", ".js", ".go"},
		},
	}

	taskService := NewTaskService(cfg)
	
	// 创建测试任务
	task, _ := taskService.CreateTask(
		models.TaskTypeAnalysis,
		"Test Task",
		"Test description",
		map[string]interface{}{"test": "data"},
	)
	
	// 测试获取任务
	retrievedTask, err := taskService.GetTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	
	if retrievedTask.ID != task.ID {
		t.Errorf("Expected task ID '%s', got '%s'", task.ID, retrievedTask.ID)
	}
	
	// 测试获取不存在的任务
	_, err = taskService.GetTask("non-existent-task")
	if err == nil {
		t.Error("Expected error when getting non-existent task")
	}
}

func TestTaskService_ListTasks(t *testing.T) {
	cfg := &config.Config{
		TaskScheduler: config.TaskSchedulerConfig{
			MaxConcurrentTasks: 5,
			TaskTimeout:        10 * time.Minute,
			ResultRetention:    24 * time.Hour,
		},
		PythonAgent: config.PythonAgentConfig{
			Host:       "localhost",
			Port:       "8000",
			Timeout:    30 * time.Second,
			RetryCount: 3,
		},
		Git: config.GitConfig{
			CloneTimeout:      5 * time.Minute,
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".py", ".js", ".go"},
		},
	}

	taskService := NewTaskService(cfg)
	
	// 创建多个任务
	task1, _ := taskService.CreateTask(models.TaskTypeAnalysis, "Task 1", "Desc 1", map[string]interface{}{"test": "1"})
	task2, _ := taskService.CreateTask(models.TaskTypeConvert, "Task 2", "Desc 2", map[string]interface{}{"test": "2"})
	
	// 等待任务创建完成（避免并发问题）
	time.Sleep(100 * time.Millisecond)
	
	// 测试列出任务
	tasks, err := taskService.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	
	// 验证任务存在（可能包含其他测试创建的任务）
	foundTask1 := false
	foundTask2 := false
	for _, task := range tasks {
		if task.ID == task1.ID {
			foundTask1 = true
		}
		if task.ID == task2.ID {
			foundTask2 = true
		}
	}
	
	if !foundTask1 {
		t.Error("Task 1 not found in task list")
	}
	if !foundTask2 {
		t.Error("Task 2 not found in task list")
	}
}

func TestTaskService_CancelTask(t *testing.T) {
	cfg := &config.Config{
		TaskScheduler: config.TaskSchedulerConfig{
			MaxConcurrentTasks: 5,
			TaskTimeout:        10 * time.Minute,
			ResultRetention:    24 * time.Hour,
		},
		PythonAgent: config.PythonAgentConfig{
			Host:       "localhost",
			Port:       "8000",
			Timeout:    30 * time.Second,
			RetryCount: 3,
		},
		Git: config.GitConfig{
			CloneTimeout:      5 * time.Minute,
			MaxFileSize:       10 * 1024 * 1024,
			AllowedExtensions: []string{".py", ".js", ".go"},
		},
	}

	taskService := NewTaskService(cfg)
	
	// 创建任务
	task, _ := taskService.CreateTask(
		models.TaskTypeAnalysis,
		"Test Task",
		"Test description",
		map[string]interface{}{"test": "data"},
	)
	
	// 测试取消任务
	err := taskService.CancelTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to cancel task: %v", err)
	}
	
	// 验证任务状态
	retrievedTask, _ := taskService.GetTask(task.ID)
	if retrievedTask.Status != models.TaskStatusCancelled {
		t.Errorf("Expected task status 'cancelled', got '%s'", retrievedTask.Status)
	}
}

func TestPythonAgentClient_CallAPI(t *testing.T) {
	cfg := &config.PythonAgentConfig{
		Host:       "localhost",
		Port:       "8000",
		Timeout:    5 * time.Second,
		RetryCount: 1,
	}
	
	client := NewPythonAgentClient(cfg)
	
	// 测试健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 注意：这需要Python AI Agent服务正在运行
	// 如果服务未运行，这个测试会失败
	result, err := client.HealthCheck(ctx)
	if err != nil {
		t.Logf("Python AI Agent health check failed (expected if service not running): %v", err)
		// 不将测试标记为失败，因为服务可能未运行
		return
	}
	
	if result == nil {
		t.Error("Health check result should not be nil")
	}
}