package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`    // 密码不返回给前端
	Role           string    `json:"role"` // admin, user, guest, local
	Avatar         string    `json:"avatar,omitempty"`
	GithubID       string    `json:"github_id,omitempty"`
	GithubUsername string    `json:"github_username,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	Location       string    `json:"location,omitempty"`
	Occupation     string    `json:"occupation,omitempty"`
	Company        string    `json:"company,omitempty"`
	Website        string    `json:"website,omitempty"`
	Twitter        string    `json:"twitter,omitempty"`
	GithubURL      string    `json:"github_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UserRegisterRequest 注册请求
type UserRegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// UserLoginRequest 登录请求
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginResponse 登录响应
type UserLoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Location   string `json:"location,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	Company    string `json:"company,omitempty"`
	Website    string `json:"website,omitempty"`
	Twitter    string `json:"twitter,omitempty"`
	GithubURL  string `json:"github_url,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// VerifyResetCodeRequest 验证重置码请求
type VerifyResetCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// GithubCallbackRequest GitHub 回调请求
type GithubCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// GithubBindRequest GitHub 绑定请求
type GithubBindRequest struct {
	GithubID       string `json:"github_id"`
	GithubUsername string `json:"github_username"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
}
