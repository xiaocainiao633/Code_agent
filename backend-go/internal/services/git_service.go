package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// GitService Git分析服务
type GitService struct {
	config *config.GitConfig
}

// NewGitService 创建Git服务
func NewGitService(cfg *config.GitConfig) *GitService {
	return &GitService{
		config: cfg,
	}
}

// AnalyzeRepository 分析Git仓库
func (s *GitService) AnalyzeRepository(ctx context.Context, repoPath string) (*models.GitAnalysisResult, error) {
	utils.Info("Analyzing Git repository: %s", repoPath)

	// 打开Git仓库
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// 获取仓库信息
	repoInfo, err := s.getRepositoryInfo(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	// 获取提交历史
	commits, err := s.getCommitHistory(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %w", err)
	}

	// 获取贡献者统计
	contributors, err := s.getContributors(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributors: %w", err)
	}

	// 获取分支信息
	branches, err := s.getBranches(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// 分析代码文件
	fileAnalysis, err := s.analyzeCodeFiles(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze code files: %w", err)
	}

	result := &models.GitAnalysisResult{
		Repository:   repoInfo,
		Commits:      commits,
		Contributors: contributors,
		Branches:     branches,
		FileAnalysis: fileAnalysis,
		AnalyzedAt:   time.Now(),
	}

	utils.Info("Git repository analysis completed: %s", repoPath)
	return result, nil
}

// CloneRepository 克隆Git仓库
func (s *GitService) CloneRepository(ctx context.Context, url, targetPath string) error {
	utils.Info("Cloning Git repository: %s to %s", url, targetPath)

	// 检查目标路径是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path already exists: %s", targetPath)
	}

	// 创建父目录
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// 设置克隆选项
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	// 克隆仓库
	_, err := git.PlainCloneContext(ctx, targetPath, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	utils.Info("Git repository cloned successfully: %s", url)
	return nil
}

// getRepositoryInfo 获取仓库基本信息
func (s *GitService) getRepositoryInfo(repo *git.Repository) (*models.RepositoryInfo, error) {
	// 获取HEAD引用
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// 获取远程仓库信息
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	var remoteURL string
	if len(remotes) > 0 {
		remoteURL = remotes[0].Config().URLs[0]
	}

	return &models.RepositoryInfo{
		CurrentBranch: head.Name().Short(),
		CurrentCommit: head.Hash().String(),
		RemoteURL:     remoteURL,
	}, nil
}

// getCommitHistory 获取提交历史
func (s *GitService) getCommitHistory(repo *git.Repository) ([]models.CommitInfo, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var commits []models.CommitInfo
	err = commitIter.ForEach(func(c *object.Commit) error {
		commitInfo := models.CommitInfo{
			Hash:      c.Hash.String(),
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Message:   strings.TrimSpace(c.Message),
			Timestamp: c.Author.When,
		}
		commits = append(commits, commitInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}

// getContributors 获取贡献者统计
func (s *GitService) getContributors(repo *git.Repository) (map[string]models.ContributorStats, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	contributors := make(map[string]models.ContributorStats)

	err = commitIter.ForEach(func(c *object.Commit) error {
		author := c.Author.Name
		email := c.Author.Email
		
		key := fmt.Sprintf("%s <%s>", author, email)
		
		if stats, exists := contributors[key]; exists {
			stats.CommitCount++
			if c.Author.When.Before(stats.FirstCommit) {
				stats.FirstCommit = c.Author.When
			}
			if c.Author.When.After(stats.LastCommit) {
				stats.LastCommit = c.Author.When
			}
			contributors[key] = stats
		} else {
			contributors[key] = models.ContributorStats{
				Name:        author,
				Email:       email,
				CommitCount: 1,
				FirstCommit: c.Author.When,
				LastCommit:  c.Author.When,
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return contributors, nil
}

// getBranches 获取分支信息
func (s *GitService) getBranches(repo *git.Repository) ([]models.BranchInfo, error) {
	var branches []models.BranchInfo

	// 获取所有分支
	refIter, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	err = refIter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branchInfo := models.BranchInfo{
				Name: ref.Name().Short(),
				Hash: ref.Hash().String(),
			}
			
			// 检查是否是当前分支
			head, _ := repo.Head()
			if head != nil && head.Name() == ref.Name() {
				branchInfo.IsCurrent = true
			}
			
			branches = append(branches, branchInfo)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate references: %w", err)
	}

	return branches, nil
}

// analyzeCodeFiles 分析代码文件
func (s *GitService) analyzeCodeFiles(repoPath string) (*models.CodeFileAnalysis, error) {
	analysis := &models.CodeFileAnalysis{
		TotalFiles:   0,
		LanguageStats: make(map[string]int),
		FileList:     []models.FileInfo{},
	}

	// 遍历仓库目录
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和隐藏文件
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// 检查文件扩展名
		ext := filepath.Ext(path)
		allowed := false
		for _, allowedExt := range s.config.AllowedExtensions {
			if ext == allowedExt {
				allowed = true
				break
			}
		}

		if !allowed {
			return nil
		}

		// 检查文件大小
		if info.Size() > s.config.MaxFileSize {
			utils.Warn("File too large, skipping: %s (%d bytes)", path, info.Size())
			return nil
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			utils.Warn("Failed to read file: %s, error: %v", path, err)
			return nil
		}

		// 分析文件
		fileInfo := models.FileInfo{
			Path:     path,
			Name:     info.Name(),
			Size:     info.Size(),
			Language: s.getLanguageByExtension(ext),
			Lines:    len(strings.Split(string(content), "\n")),
		}

		analysis.TotalFiles++
		analysis.LanguageStats[fileInfo.Language]++
		analysis.FileList = append(analysis.FileList, fileInfo)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk repository path: %w", err)
	}

	return analysis, nil
}

// getLanguageByExtension 根据扩展名获取语言
func (s *GitService) getLanguageByExtension(ext string) string {
	languageMap := map[string]string{
		".py":   "Python",
		".js":   "JavaScript",
		".java": "Java",
		".cpp":  "C++",
		".c":    "C",
		".go":   "Go",
		".rs":   "Rust",
		".ts":   "TypeScript",
	}

	if lang, exists := languageMap[ext]; exists {
		return lang
	}
	return "Unknown"
}

// GetFileHistory 获取文件历史
func (s *GitService) GetFileHistory(ctx context.Context, repoPath, filePath string) ([]models.FileHistory, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var history []models.FileHistory

	err = commitIter.ForEach(func(c *object.Commit) error {
		// 检查文件是否存在于该提交中
		file, err := c.File(filePath)
		if err != nil {
			return nil // 文件不存在于此提交中，继续
		}

		content, err := file.Contents()
		if err != nil {
			return err
		}

		historyItem := models.FileHistory{
			CommitHash: c.Hash.String(),
			Author:     c.Author.Name,
			Message:    strings.TrimSpace(c.Message),
			Timestamp:  c.Author.When,
			Lines:      len(strings.Split(content, "\n")),
		}

		history = append(history, historyItem)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return history, nil
}

// GetDiff 获取文件差异
func (s *GitService) GetDiff(ctx context.Context, repoPath, filePath, fromCommit, toCommit string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	fromHash := plumbing.NewHash(fromCommit)
	toHash := plumbing.NewHash(toCommit)

	fromCommitObj, err := repo.CommitObject(fromHash)
	if err != nil {
		return "", fmt.Errorf("failed to get from commit: %w", err)
	}

	toCommitObj, err := repo.CommitObject(toHash)
	if err != nil {
		return "", fmt.Errorf("failed to get to commit: %w", err)
	}

	fromFile, err := fromCommitObj.File(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file from commit %s: %w", fromCommit, err)
	}

	toFile, err := toCommitObj.File(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file from commit %s: %w", toCommit, err)
	}

	fromContent, err := fromFile.Contents()
	if err != nil {
		return "", fmt.Errorf("failed to read from file content: %w", err)
	}

	toContent, err := toFile.Contents()
	if err != nil {
		return "", fmt.Errorf("failed to read to file content: %w", err)
	}

	// 简单的差异比较（可以改进为更复杂的diff算法）
	diff := fmt.Sprintf("--- %s\n+++ %s\n@@ file diff @@\n-%s\n+%s",
		fromCommit[:7], toCommit[:7], fromContent, toContent)

	return diff, nil
}