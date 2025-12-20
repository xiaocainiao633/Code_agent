@echo off
echo ========================================
echo CodeSage 认证系统测试脚本
echo ========================================
echo.

set BASE_URL=http://localhost:8082/api/v1

echo [1] 测试健康检查...
curl -X GET %BASE_URL%/health
echo.
echo.

echo [2] 测试用户注册...
curl -X POST %BASE_URL%/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\",\"confirmPassword\":\"password123\"}"
echo.
echo.

echo [3] 测试用户登录...
curl -X POST %BASE_URL%/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"password\":\"password123\"}" > temp_token.json
echo.
echo.

echo [4] 提取 Token (请手动复制 token 值)
type temp_token.json
echo.
echo.

echo ========================================
echo 测试完成！
echo 请复制上面的 token 值，然后运行以下命令测试认证：
echo.
echo curl -X GET %BASE_URL%/auth/profile -H "Authorization: Bearer YOUR_TOKEN_HERE"
echo ========================================

pause
