package services

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/database"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/models"
)

var (
	ErrUserExists         = errors.New("用户名或邮箱已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrInvalidUsername    = errors.New("用户名格式不正确")
	ErrInvalidEmail       = errors.New("邮箱格式不正确")
	ErrInvalidPassword    = errors.New("密码格式不正确")
	ErrPasswordMismatch   = errors.New("两次输入的密码不一致")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrOldPasswordWrong   = errors.New("原密码错误")
	ErrInvalidResetToken  = errors.New("重置令牌无效或已过期")
	ErrGithubNotBound     = errors.New("GitHub 账号未绑定")
)

// JWT密钥 - 生产环境应该从环境变量读取
var jwtSecret = []byte("codesage-secret-key-change-in-production")

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register 用户注册
func (s *AuthService) Register(req *models.UserRegisterRequest) (*models.User, error) {
	// 验证输入
	if err := s.validateRegisterInput(req); err != nil {
		return nil, err
	}

	// 检查用户是否已存在
	exists, err := s.userExists(req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 插入用户
	result, err := database.DB.Exec(
		`INSERT INTO users (username, email, password, role, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, ?)`,
		req.Username, req.Email, string(hashedPassword), "user", time.Now(), time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// 返回用户信息
	user := &models.User{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *models.UserLoginRequest) (*models.UserLoginResponse, error) {
	// 验证输入
	if req.Username == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// 查询用户
	var user models.User
	var hashedPassword string
	var avatar, phone, bio, location, occupation, company, website, twitter, githubURL sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, username, email, password, role, avatar, phone, bio, location, occupation, company, website, twitter, github_url, created_at, updated_at 
		 FROM users WHERE username = ?`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.Role, &avatar, &phone, &bio, &location, &occupation, &company, &website, &twitter, &githubURL, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 处理可能为 NULL 的字段
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}
	if location.Valid {
		user.Location = location.String
	}
	if occupation.Valid {
		user.Occupation = occupation.String
	}
	if company.Valid {
		user.Company = company.String
	}
	if website.Valid {
		user.Website = website.String
	}
	if twitter.Valid {
		user.Twitter = twitter.String
	}
	if githubURL.Valid {
		user.GithubURL = githubURL.String
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成JWT token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.UserLoginResponse{
		User:  &user,
		Token: token,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	var avatar, phone, bio, location, occupation, company, website, twitter, githubURL sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, username, email, role, avatar, phone, bio, location, occupation, company, website, twitter, github_url, created_at, updated_at 
		 FROM users WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &avatar, &phone, &bio, &location, &occupation, &company, &website, &twitter, &githubURL, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 处理可能为 NULL 的字段
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}
	if location.Valid {
		user.Location = location.String
	}
	if occupation.Valid {
		user.Occupation = occupation.String
	}
	if company.Valid {
		user.Company = company.String
	}
	if website.Valid {
		user.Website = website.String
	}
	if twitter.Valid {
		user.Twitter = twitter.String
	}
	if githubURL.Valid {
		user.GithubURL = githubURL.String
	}

	return &user, nil
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID int64, req *models.UpdateProfileRequest) (*models.User, error) {
	// 获取当前用户信息
	currentUser, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 验证输入
	if req.Username != "" && req.Username != currentUser.Username {
		if err := s.validateUsername(req.Username); err != nil {
			return nil, err
		}
		// 检查用户名是否已被其他用户使用
		var count int
		err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND id != ?", req.Username, userID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if count > 0 {
			return nil, errors.New("用户名已被使用")
		}
	}
	if req.Email != "" && req.Email != currentUser.Email {
		if err := s.validateEmail(req.Email); err != nil {
			return nil, err
		}
		// 检查邮箱是否已被其他用户使用
		var count int
		err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", req.Email, userID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if count > 0 {
			return nil, errors.New("邮箱已被使用")
		}
	}

	// 构建更新SQL
	query := "UPDATE users SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Username != "" {
		query += ", username = ?"
		args = append(args, req.Username)
	}
	if req.Email != "" {
		query += ", email = ?"
		args = append(args, req.Email)
	}
	if req.Avatar != "" {
		query += ", avatar = ?"
		args = append(args, req.Avatar)
	}
	if req.Phone != "" {
		query += ", phone = ?"
		args = append(args, req.Phone)
	}
	if req.Bio != "" {
		query += ", bio = ?"
		args = append(args, req.Bio)
	}
	if req.Location != "" {
		query += ", location = ?"
		args = append(args, req.Location)
	}
	if req.Occupation != "" {
		query += ", occupation = ?"
		args = append(args, req.Occupation)
	}
	if req.Company != "" {
		query += ", company = ?"
		args = append(args, req.Company)
	}
	if req.Website != "" {
		query += ", website = ?"
		args = append(args, req.Website)
	}
	if req.Twitter != "" {
		query += ", twitter = ?"
		args = append(args, req.Twitter)
	}
	if req.GithubURL != "" {
		query += ", github_url = ?"
		args = append(args, req.GithubURL)
	}

	query += " WHERE id = ?"
	args = append(args, userID)

	// 执行更新
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// 返回更新后的用户信息
	return s.GetUserByID(userID)
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, req *models.ChangePasswordRequest) error {
	// 验证新密码
	if err := s.validatePassword(req.NewPassword); err != nil {
		return err
	}

	// 获取当前密码
	var currentPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&currentPassword)
	if err != nil {
		return fmt.Errorf("failed to get user password: %w", err)
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(req.OldPassword)); err != nil {
		return ErrOldPasswordWrong
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	_, err = database.DB.Exec(
		"UPDATE users SET password = ?, updated_at = ? WHERE id = ?",
		string(hashedPassword), time.Now(), userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// ValidateToken 验证JWT token
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["user_id"].(float64))
		return s.GetUserByID(userID)
	}

	return nil, errors.New("invalid token")
}

// generateToken 生成JWT token
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// validateRegisterInput 验证注册输入
func (s *AuthService) validateRegisterInput(req *models.UserRegisterRequest) error {
	if err := s.validateUsername(req.Username); err != nil {
		return err
	}
	if err := s.validateEmail(req.Email); err != nil {
		return err
	}
	if err := s.validatePassword(req.Password); err != nil {
		return err
	}
	if req.Password != req.ConfirmPassword {
		return ErrPasswordMismatch
	}
	return nil
}

// validateUsername 验证用户名
func (s *AuthService) validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	// 只允许字母、数字、下划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return ErrInvalidUsername
	}
	return nil
}

// validateEmail 验证邮箱
func (s *AuthService) validateEmail(email string) error {
	// 简单的邮箱格式验证
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return ErrInvalidEmail
	}
	return nil
}

// validatePassword 验证密码
func (s *AuthService) validatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return ErrInvalidPassword
	}
	return nil
}

// userExists 检查用户是否存在
func (s *AuthService) userExists(username, email string) (bool, error) {
	var count int
	err := database.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE username = ? OR email = ?",
		username, email,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ForgotPassword 忘记密码 - 发送重置码到邮箱
func (s *AuthService) ForgotPassword(email string) (string, error) {
	// 验证邮箱格式
	if err := s.validateEmail(email); err != nil {
		return "", err
	}

	// 检查用户是否存在
	var userID int64
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to query user: %w", err)
	}

	// 生成6位数字重置码
	resetCode := fmt.Sprintf("%06d", time.Now().Unix()%1000000)

	// 生成重置令牌（用于验证）
	resetToken := fmt.Sprintf("%s-%d", resetCode, time.Now().Unix())

	// 设置过期时间（15分钟）
	expiresAt := time.Now().Add(15 * time.Minute)

	// 保存重置令牌到数据库
	_, err = database.DB.Exec(
		"UPDATE users SET reset_token = ?, reset_token_expires = ?, updated_at = ? WHERE id = ?",
		resetToken, expiresAt, time.Now(), userID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	// 返回重置码（实际应用中应该发送邮件）
	// TODO: 集成邮件服务发送重置码
	return resetCode, nil
}

// VerifyResetCode 验证重置码
func (s *AuthService) VerifyResetCode(email, code string) (string, error) {
	var resetToken string
	var expiresAt time.Time

	err := database.DB.QueryRow(
		"SELECT reset_token, reset_token_expires FROM users WHERE email = ?",
		email,
	).Scan(&resetToken, &expiresAt)

	if err == sql.ErrNoRows {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to query reset token: %w", err)
	}

	// 检查令牌是否过期
	if time.Now().After(expiresAt) {
		return "", ErrInvalidResetToken
	}

	// 验证重置码
	if !regexp.MustCompile(`^` + code + `-`).MatchString(resetToken) {
		return "", ErrInvalidResetToken
	}

	return resetToken, nil
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(email, code, newPassword string) error {
	// 验证重置码
	_, err := s.VerifyResetCode(email, code)
	if err != nil {
		return err
	}

	// 验证新密码
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码并清除重置令牌
	_, err = database.DB.Exec(
		"UPDATE users SET password = ?, reset_token = NULL, reset_token_expires = NULL, updated_at = ? WHERE email = ?",
		string(hashedPassword), time.Now(), email,
	)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}

// BindGithub 绑定 GitHub 账号
func (s *AuthService) BindGithub(userID int64, githubID, githubUsername string) error {
	// 检查 GitHub ID 是否已被其他用户绑定
	var existingUserID int64
	err := database.DB.QueryRow("SELECT id FROM users WHERE github_id = ?", githubID).Scan(&existingUserID)
	if err == nil && existingUserID != userID {
		return errors.New("该 GitHub 账号已被其他用户绑定")
	}

	// 绑定 GitHub 账号
	_, err = database.DB.Exec(
		"UPDATE users SET github_id = ?, github_username = ?, updated_at = ? WHERE id = ?",
		githubID, githubUsername, time.Now(), userID,
	)
	if err != nil {
		return fmt.Errorf("failed to bind github: %w", err)
	}

	return nil
}

// LoginWithGithub GitHub 登录
func (s *AuthService) LoginWithGithub(githubID string) (*models.UserLoginResponse, error) {
	// 查询绑定了该 GitHub ID 的用户
	var user models.User
	var avatar sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, username, email, role, avatar, github_id, github_username, created_at, updated_at 
		 FROM users WHERE github_id = ?`,
		githubID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &avatar, &user.GithubID, &user.GithubUsername, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrGithubNotBound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 处理可能为 NULL 的 avatar 字段
	if avatar.Valid {
		user.Avatar = avatar.String
	}

	// 生成JWT token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.UserLoginResponse{
		User:  &user,
		Token: token,
	}, nil
}

// GetUserByGithubID 根据 GitHub ID 获取用户
func (s *AuthService) GetUserByGithubID(githubID string) (*models.User, error) {
	var user models.User
	var avatar sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, username, email, role, avatar, github_id, github_username, created_at, updated_at 
		 FROM users WHERE github_id = ?`,
		githubID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &avatar, &user.GithubID, &user.GithubUsername, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 处理可能为 NULL 的 avatar 字段
	if avatar.Valid {
		user.Avatar = avatar.String
	}

	return &user, nil
}
