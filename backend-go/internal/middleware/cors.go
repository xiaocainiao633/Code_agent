package middleware

import (
	"net/http"
	"strings"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
	"github.com/xiaocainiao633/Code_agent/backend-go/internal/utils"
)

// CORSMiddleware CORS中间件
func CORSMiddleware(cfg *config.CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// 记录CORS相关日志
			utils.Debug("CORS request from origin: %s, method: %s, path: %s", origin, r.Method, r.URL.Path)
			
			// 设置CORS头
			if origin != "" {
				// 检查来源是否在允许列表中
				allowed := false
				for _, allowedOrigin := range cfg.AllowedOrigins {
					if allowedOrigin == "*" {
						allowed = true
						break
					}
					if allowedOrigin == origin {
						allowed = true
						break
					}
				}
				
				// 如果是预检请求，总是允许
				if r.Method == "OPTIONS" {
					allowed = true
				}
				
				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
					
					if cfg.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}
					
					// 设置预检请求缓存时间
					w.Header().Set("Access-Control-Max-Age", "86400")
				}
			}
			
			// 处理预检请求
			if r.Method == "OPTIONS" {
				utils.Debug("Handling OPTIONS preflight request")
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// 继续处理请求
			next.ServeHTTP(w, r)
		})
	}
}