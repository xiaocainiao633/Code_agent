package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/middleware"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/services"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode register request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("User registration attempt: %s", req.Username)

	user, err := h.authService.Register(&req)
	if err != nil {
		utils.Error("Registration failed for user %s: %v", req.Username, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("User registered successfully: %s (ID: %d)", user.Username, user.ID)
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "注册成功",
		"user":    user,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode login request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("User login attempt: %s", req.Username)

	response, err := h.authService.Login(&req)
	if err != nil {
		utils.Error("Login failed for user %s: %v", req.Username, err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Info("User logged in successfully: %s (ID: %d)", response.User.Username, response.User.ID)
	respondWithJSON(w, http.StatusOK, response)
}

// GetProfile 获取用户资料
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从context获取用户信息
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未授权")
		return
	}

	utils.Debug("Get profile for user: %s (ID: %d)", user.Username, user.ID)
	respondWithJSON(w, http.StatusOK, user)
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从context获取用户信息
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode update profile request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("Update profile for user: %s (ID: %d)", user.Username, user.ID)

	updatedUser, err := h.authService.UpdateProfile(user.ID, &req)
	if err != nil {
		utils.Error("Failed to update profile for user %d: %v", user.ID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("Profile updated successfully for user: %s (ID: %d)", updatedUser.Username, updatedUser.ID)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "资料更新成功",
		"user":    updatedUser,
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从context获取用户信息
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode change password request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("Change password for user: %s (ID: %d)", user.Username, user.ID)

	if err := h.authService.ChangePassword(user.ID, &req); err != nil {
		utils.Error("Failed to change password for user %d: %v", user.ID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("Password changed successfully for user: %s (ID: %d)", user.Username, user.ID)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "密码修改成功",
	})
}

// GetUserByID 根据ID获取用户（管理员功能）
func (h *AuthHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从URL路径获取用户ID
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	userID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.Error("Failed to get user %d: %v", userID, err)
		respondWithError(w, http.StatusNotFound, "用户不存在")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// respondWithJSON 返回JSON响应
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		utils.Error("Failed to marshal JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError 返回错误响应
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// ForgotPassword 忘记密码
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode forgot password request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("Forgot password request for email: %s", req.Email)

	resetCode, err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		utils.Error("Forgot password failed for email %s: %v", req.Email, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("Reset code generated for email: %s", req.Email)

	// 在实际应用中，不应该返回重置码，而是发送邮件
	// 这里为了演示方便，直接返回重置码
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "重置码已发送到您的邮箱（演示模式：直接返回）",
		"code":    resetCode, // 生产环境应该删除这行
	})
}

// VerifyResetCode 验证重置码
func (h *AuthHandler) VerifyResetCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.VerifyResetCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode verify reset code request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("Verify reset code for email: %s", req.Email)

	_, err := h.authService.VerifyResetCode(req.Email, req.Code)
	if err != nil {
		utils.Error("Verify reset code failed for email %s: %v", req.Email, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "验证码正确",
	})
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode reset password request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("Reset password for email: %s", req.Email)

	if err := h.authService.ResetPassword(req.Email, req.Code, req.NewPassword); err != nil {
		utils.Error("Reset password failed for email %s: %v", req.Email, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("Password reset successfully for email: %s", req.Email)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "密码重置成功",
	})
}

// BindGithub 绑定 GitHub 账号
func (h *AuthHandler) BindGithub(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从context获取用户信息
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "未授权")
		return
	}

	var req models.GithubBindRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode bind github request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("Bind GitHub for user: %s (ID: %d)", user.Username, user.ID)

	if err := h.authService.BindGithub(user.ID, req.GithubID, req.GithubUsername); err != nil {
		utils.Error("Bind GitHub failed for user %d: %v", user.ID, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.Info("GitHub bound successfully for user: %s (ID: %d)", user.Username, user.ID)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "GitHub 账号绑定成功",
	})
}

// GithubLogin GitHub 登录
func (h *AuthHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		GithubID string `json:"github_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Failed to decode github login request: %v", err)
		respondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	utils.Info("GitHub login attempt for GitHub ID: %s", req.GithubID)

	response, err := h.authService.LoginWithGithub(req.GithubID)
	if err != nil {
		utils.Error("GitHub login failed for GitHub ID %s: %v", req.GithubID, err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Info("GitHub login successful for user: %s (ID: %d)", response.User.Username, response.User.ID)
	respondWithJSON(w, http.StatusOK, response)
}
