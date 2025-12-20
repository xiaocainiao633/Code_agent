package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/services"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware JWT认证中间件
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从请求头获取token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "未提供认证令牌")
				return
			}

			// 解析Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "认证令牌格式错误")
				return
			}

			tokenString := parts[1]

			// 验证token
			user, err := authService.ValidateToken(tokenString)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "无效的认证令牌")
				return
			}

			// 将用户信息存入context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// respondWithError 返回错误响应
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
