package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// PythonAgentClient Python AI Agent客户端
type PythonAgentClient struct {
	baseURL    string
	httpClient *http.Client
	retryCount int
}

// NewPythonAgentClient 创建Python AI Agent客户端
func NewPythonAgentClient(config *config.PythonAgentConfig) *PythonAgentClient {
	return &PythonAgentClient{
		baseURL: fmt.Sprintf("http://%s:%s", config.Host, config.Port),
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		retryCount: config.RetryCount,
	}
}

// AnalyzeCode 调用代码分析API
func (c *PythonAgentClient) AnalyzeCode(ctx context.Context, code string, language string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"code":     code,
		"language": language,
	}

	return c.callAPI(ctx, "/api/v1/analyze", request)
}

// AnalyzePython2Code 调用Python2专项分析API
func (c *PythonAgentClient) AnalyzePython2Code(ctx context.Context, code string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"code":     code,
		"language": "python",
	}

	return c.callAPI(ctx, "/api/v1/analyze/python2", request)
}

// ConvertCode 调用代码转换API
func (c *PythonAgentClient) ConvertCode(ctx context.Context, code string, fromVersion string, toVersion string, options map[string]interface{}) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"code":        code,
		"language":    "python", // 添加语言参数
		"conversion_type": "python_2_to_3", // 修复转换类型
		"from_version": fromVersion,
		"to_version":  toVersion,
		"options":     options,
	}

	return c.callAPI(ctx, "/api/v1/convert", request)
}

// GenerateTests 调用测试生成API
func (c *PythonAgentClient) GenerateTests(ctx context.Context, code string, testType string, framework string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"code":      code,
		"language":  "python", // 添加语言参数
		"test_type": testType,
		"framework": framework,
	}

	return c.callAPI(ctx, "/api/v1/generate-tests", request)
}

// HealthCheck 检查Python Agent健康状态
func (c *PythonAgentClient) HealthCheck(ctx context.Context) (map[string]interface{}, error) {
	return c.callAPI(ctx, "/api/v1/health", nil)
}

// callAPI 调用API的通用方法
func (c *PythonAgentClient) callAPI(ctx context.Context, endpoint string, requestData interface{}) (map[string]interface{}, error) {
	url := c.baseURL + endpoint

	var reqBody []byte
	var err error

	if requestData != nil {
		reqBody, err = json.Marshal(requestData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	}

	// 重试逻辑
	var lastErr error
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			utils.Warn("Retrying API call to %s (attempt %d/%d)", endpoint, attempt, c.retryCount)
			time.Sleep(time.Duration(attempt) * time.Second) // 指数退避
		}

		result, err := c.doRequest(ctx, url, reqBody)
		if err == nil {
			return result, nil
		}

		lastErr = err
		utils.Error("API call to %s failed (attempt %d): %v", endpoint, attempt+1, err)

		// 检查是否应该重试
		if !c.shouldRetry(err) {
			break
		}
	}

	return nil, fmt.Errorf("API call failed after %d attempts: %w", c.retryCount+1, lastErr)
}

// doRequest 执行HTTP请求
func (c *PythonAgentClient) doRequest(ctx context.Context, url string, body []byte) (map[string]interface{}, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "CodeSage-Go-Backend/1.0")

	utils.Debug("Calling Python Agent API: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	utils.Debug("Python Agent API response: status=%d, body=%s", resp.StatusCode, string(bodyBytes))

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// shouldRetry 判断是否应该重试
func (c *PythonAgentClient) shouldRetry(err error) bool {
	// 网络错误重试
	if err, ok := err.(interface{ Timeout() bool }); ok && err.Timeout() {
		return true
	}

	// HTTP状态码判断
	if err, ok := err.(interface{ StatusCode() int }); ok {
		statusCode := err.StatusCode()
		// 5xx错误重试
		if statusCode >= 500 && statusCode < 600 {
			return true
		}
		// 429 Too Many Requests 重试
		if statusCode == 429 {
			return true
		}
	}

	// 包含特定错误信息的重试
	errStr := err.Error()
	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"temporary failure",
		"no such host",
	}

	for _, retryable := range retryableErrors {
		if contains(errStr, retryable) {
			return true
		}
	}

	return false
}

// contains 检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsHelper(s, substr)
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}