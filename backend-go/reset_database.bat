@echo off
echo ========================================
echo 重置数据库
echo ========================================
echo.
echo 警告：此操作将删除所有用户数据！
echo.
pause

if exist data\codesage.db (
    del data\codesage.db
    echo 数据库已删除
) else (
    echo 数据库文件不存在
)

echo.
echo 请重新启动后端服务，数据库将自动重新创建
echo.
pause
