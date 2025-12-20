package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase 初始化数据库
func InitDatabase(dbPath string) error {
	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables 创建数据库表
func createTables() error {
	// 创建users表
	if err := createUsersTable(); err != nil {
		return err
	}

	// 创建tasks表
	if err := createTasksTable(); err != nil {
		return err
	}

	return nil
}

// createUsersTable 创建用户表
func createUsersTable() error {
	var tableName string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)

	if err == sql.ErrNoRows {
		// 表不存在，创建新表
		userTableSQL := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			avatar TEXT,
			github_id TEXT UNIQUE,
			github_username TEXT,
			reset_token TEXT,
			reset_token_expires DATETIME,
			phone TEXT,
			bio TEXT,
			location TEXT,
			occupation TEXT,
			company TEXT,
			website TEXT,
			twitter TEXT,
			github_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX idx_users_username ON users(username);
		CREATE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_github_id ON users(github_id);
		`

		if _, err := DB.Exec(userTableSQL); err != nil {
			return fmt.Errorf("failed to create users table: %w", err)
		}
	} else if err == nil {
		// 表已存在，检查并添加缺失的列
		if err := migrateUsersTable(); err != nil {
			return fmt.Errorf("failed to migrate users table: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check users table existence: %w", err)
	}

	return nil
}

// createTasksTable 创建任务表
func createTasksTable() error {
	var tableName string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tasks'").Scan(&tableName)

	if err == sql.ErrNoRows {
		// 表不存在，创建新表
		taskTableSQL := `
		CREATE TABLE tasks (
			id TEXT PRIMARY KEY,
			user_id INTEGER,
			type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			name TEXT NOT NULL,
			description TEXT,
			params TEXT,
			result TEXT,
			error TEXT,
			progress INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			started_at DATETIME,
			completed_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		
		CREATE INDEX idx_tasks_user_id ON tasks(user_id);
		CREATE INDEX idx_tasks_status ON tasks(status);
		CREATE INDEX idx_tasks_type ON tasks(type);
		CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
		`

		if _, err := DB.Exec(taskTableSQL); err != nil {
			return fmt.Errorf("failed to create tasks table: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check tasks table existence: %w", err)
	}

	return nil
}

// migrateUsersTable 迁移用户表，添加新字段
func migrateUsersTable() error {
	// 检查并添加 github_id 列
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='github_id'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check github_id column: %w", err)
	}
	if count == 0 {
		// SQLite 不支持在 ALTER TABLE 时添加 UNIQUE 约束，需要分两步
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN github_id TEXT"); err != nil {
			return fmt.Errorf("failed to add github_id column: %w", err)
		}
		// 创建唯一索引
		if _, err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id) WHERE github_id IS NOT NULL"); err != nil {
			return fmt.Errorf("failed to create github_id index: %w", err)
		}
	}

	// 检查并添加 github_username 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='github_username'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check github_username column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN github_username TEXT"); err != nil {
			return fmt.Errorf("failed to add github_username column: %w", err)
		}
	}

	// 检查并添加 reset_token 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='reset_token'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check reset_token column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN reset_token TEXT"); err != nil {
			return fmt.Errorf("failed to add reset_token column: %w", err)
		}
	}

	// 检查并添加 reset_token_expires 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='reset_token_expires'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check reset_token_expires column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN reset_token_expires DATETIME"); err != nil {
			return fmt.Errorf("failed to add reset_token_expires column: %w", err)
		}
	}

	// 检查并添加 phone 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='phone'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check phone column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN phone TEXT"); err != nil {
			return fmt.Errorf("failed to add phone column: %w", err)
		}
	}

	// 检查并添加 bio 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='bio'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check bio column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN bio TEXT"); err != nil {
			return fmt.Errorf("failed to add bio column: %w", err)
		}
	}

	// 检查并添加 location 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='location'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check location column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN location TEXT"); err != nil {
			return fmt.Errorf("failed to add location column: %w", err)
		}
	}

	// 检查并添加 occupation 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='occupation'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check occupation column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN occupation TEXT"); err != nil {
			return fmt.Errorf("failed to add occupation column: %w", err)
		}
	}

	// 检查并添加 company 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='company'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check company column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN company TEXT"); err != nil {
			return fmt.Errorf("failed to add company column: %w", err)
		}
	}

	// 检查并添加 website 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='website'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check website column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN website TEXT"); err != nil {
			return fmt.Errorf("failed to add website column: %w", err)
		}
	}

	// 检查并添加 twitter 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='twitter'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check twitter column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN twitter TEXT"); err != nil {
			return fmt.Errorf("failed to add twitter column: %w", err)
		}
	}

	// 检查并添加 github_url 列
	err = DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='github_url'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check github_url column: %w", err)
	}
	if count == 0 {
		if _, err := DB.Exec("ALTER TABLE users ADD COLUMN github_url TEXT"); err != nil {
			return fmt.Errorf("failed to add github_url column: %w", err)
		}
	}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
