package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
)

// TaskRepository 任务数据库访问层
type TaskRepository struct{}

// NewTaskRepository 创建任务仓库
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Create 创建任务
func (r *TaskRepository) Create(task *models.Task) error {
	paramsJSON, err := json.Marshal(task.Params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	query := `
		INSERT INTO tasks (id, type, status, name, description, params, progress, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = DB.Exec(query,
		task.ID,
		task.Type,
		task.Status,
		task.Name,
		task.Description,
		string(paramsJSON),
		task.Progress,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// Update 更新任务
func (r *TaskRepository) Update(task *models.Task) error {
	var resultJSON, errorJSON sql.NullString

	if task.Result != nil {
		resultBytes, err := json.Marshal(task.Result)
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}
		resultJSON = sql.NullString{String: string(resultBytes), Valid: true}
	}

	if task.Error != "" {
		errorJSON = sql.NullString{String: task.Error, Valid: true}
	}

	query := `
		UPDATE tasks 
		SET status = ?, progress = ?, result = ?, error = ?, updated_at = ?, started_at = ?, completed_at = ?
		WHERE id = ?
	`

	_, err := DB.Exec(query,
		task.Status,
		task.Progress,
		resultJSON,
		errorJSON,
		task.UpdatedAt,
		task.StartedAt,
		task.CompletedAt,
		task.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// GetByID 根据ID获取任务
func (r *TaskRepository) GetByID(taskID string) (*models.Task, error) {
	query := `
		SELECT id, type, status, name, description, params, result, error, progress, 
		       created_at, updated_at, started_at, completed_at
		FROM tasks
		WHERE id = ?
	`

	var task models.Task
	var paramsJSON, resultJSON, errorJSON sql.NullString
	var startedAt, completedAt sql.NullTime

	err := DB.QueryRow(query, taskID).Scan(
		&task.ID,
		&task.Type,
		&task.Status,
		&task.Name,
		&task.Description,
		&paramsJSON,
		&resultJSON,
		&errorJSON,
		&task.Progress,
		&task.CreatedAt,
		&task.UpdatedAt,
		&startedAt,
		&completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// 解析JSON字段
	if paramsJSON.Valid {
		if err := json.Unmarshal([]byte(paramsJSON.String), &task.Params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal params: %w", err)
		}
	}

	if resultJSON.Valid {
		if err := json.Unmarshal([]byte(resultJSON.String), &task.Result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	if errorJSON.Valid {
		task.Error = errorJSON.String
	}

	if startedAt.Valid {
		task.StartedAt = &startedAt.Time
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

// List 获取任务列表
func (r *TaskRepository) List(limit, offset int) ([]models.Task, error) {
	query := `
		SELECT id, type, status, name, description, params, result, error, progress,
		       created_at, updated_at, started_at, completed_at
		FROM tasks
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var paramsJSON, resultJSON, errorJSON sql.NullString
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Type,
			&task.Status,
			&task.Name,
			&task.Description,
			&paramsJSON,
			&resultJSON,
			&errorJSON,
			&task.Progress,
			&task.CreatedAt,
			&task.UpdatedAt,
			&startedAt,
			&completedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		// 解析JSON字段
		if paramsJSON.Valid {
			if err := json.Unmarshal([]byte(paramsJSON.String), &task.Params); err != nil {
				continue // 跳过解析失败的任务
			}
		}

		if resultJSON.Valid {
			if err := json.Unmarshal([]byte(resultJSON.String), &task.Result); err != nil {
				continue
			}
		}

		if errorJSON.Valid {
			task.Error = errorJSON.String
		}

		if startedAt.Valid {
			task.StartedAt = &startedAt.Time
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Delete 删除任务
func (r *TaskRepository) Delete(taskID string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := DB.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found: %s", taskID)
	}

	return nil
}

// DeleteOldTasks 删除旧任务
func (r *TaskRepository) DeleteOldTasks(before time.Time) (int64, error) {
	query := `DELETE FROM tasks WHERE completed_at < ? AND completed_at IS NOT NULL`

	result, err := DB.Exec(query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old tasks: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// Count 统计任务数量
func (r *TaskRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM tasks`

	var count int
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	return count, nil
}

// CountByStatus 按状态统计任务数量
func (r *TaskRepository) CountByStatus(status models.TaskStatus) (int, error) {
	query := `SELECT COUNT(*) FROM tasks WHERE status = ?`

	var count int
	err := DB.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by status: %w", err)
	}

	return count, nil
}
