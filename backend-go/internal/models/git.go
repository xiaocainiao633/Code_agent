package models

import (
	"time"
)

// GitAnalysisResult Git分析结果
type GitAnalysisResult struct {
	Repository   *RepositoryInfo         `json:"repository"`
	Commits      []CommitInfo            `json:"commits"`
	Contributors map[string]ContributorStats `json:"contributors"`
	Branches     []BranchInfo            `json:"branches"`
	FileAnalysis *CodeFileAnalysis       `json:"file_analysis"`
	AnalyzedAt   time.Time               `json:"analyzed_at"`
}

// RepositoryInfo 仓库基本信息
type RepositoryInfo struct {
	CurrentBranch string `json:"current_branch"`
	CurrentCommit string `json:"current_commit"`
	RemoteURL     string `json:"remote_url"`
}

// CommitInfo 提交信息
type CommitInfo struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ContributorStats 贡献者统计
type ContributorStats struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	CommitCount int       `json:"commit_count"`
	FirstCommit time.Time `json:"first_commit"`
	LastCommit  time.Time `json:"last_commit"`
}

// BranchInfo 分支信息
type BranchInfo struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	IsCurrent bool   `json:"is_current"`
}

// CodeFileAnalysis 代码文件分析
type CodeFileAnalysis struct {
	TotalFiles    int                  `json:"total_files"`
	LanguageStats map[string]int       `json:"language_stats"`
	FileList      []FileInfo           `json:"file_list"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Language string `json:"language"`
	Lines    int    `json:"lines"`
}

// FileHistory 文件历史
type FileHistory struct {
	CommitHash string    `json:"commit_hash"`
	Author     string    `json:"author"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	Lines      int       `json:"lines"`
}

// GitTaskParams Git任务参数
type GitTaskParams struct {
	RepoPath   string `json:"repo_path"`
	RemoteURL  string `json:"remote_url,omitempty"`
	CloneIfNotExists bool `json:"clone_if_not_exists,omitempty"`
}

// GitCloneParams Git克隆参数
type GitCloneParams struct {
	RemoteURL string `json:"remote_url"`
	TargetPath string `json:"target_path"`
}

// GitFileHistoryParams Git文件历史参数
type GitFileHistoryParams struct {
	RepoPath string `json:"repo_path"`
	FilePath string `json:"file_path"`
}

// GitDiffParams Git差异参数
type GitDiffParams struct {
	RepoPath  string `json:"repo_path"`
	FilePath  string `json:"file_path"`
	FromCommit string `json:"from_commit"`
	ToCommit  string `json:"to_commit"`
}