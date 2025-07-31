package middlewares

import (
	"go01_4/config"
	"go01_4/logger"
	"go01_4/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var mlog = logger.GetLogger()

// CtxUserIDKey 用于在上下文中存储用户ID的键
// 这个键可以在中间件和处理器之间传递用户ID信息
const CtxUserIDKey = "userID"

func RecoveryMiddleware() gin.HandlerFunc {
	// 返回一个 Gin 中间件函数（符合 gin.HandlerFunc 类型）
	return func(ctx *gin.Context) {
		// defer 语句：延迟注册一个延迟执行的匿名函数，在当前函数（中间件）执行结束时触发
		defer func() {
			// 捕获 panic 错误：如果有 panic 发生，recover() 会返回错误信息，否则为 nil
			if err := recover(); err != nil {
				// 记录错误日志，包含请求上下文信息
				mlog.WithFields(logrus.Fields{
					"path":   ctx.Request.URL.Path, // 请求路径（如 /api/user）
					"method": ctx.Request.Method,   // 请求方法（如 GET/POST）
					"err":    err,                  // panic 错误详情
					"ip":     ctx.ClientIP(),       // 客户端 IP 地址
				}).Error("服务器内部异常")

				// 向客户端返回 JSON 格式的错误响应
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    30001,     // 自定义错误码
					"message": "服务器内部异常", // 错误提示信息
				})

				// 中止请求处理：阻止后续的中间件或处理器执行
				ctx.Abort()
			}
		}()

		// 执行后续的中间件或路由处理器
		ctx.Next()
	}
}

func JWTCheckUserAuth(cfg *config.BlogConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取token
		tokenString := ctx.GetHeader("token")
		if tokenString == "" {
			utils.Error(ctx, http.StatusUnauthorized, "未找到认证信息")
		}
		// 解析令牌验证token是否，提取存储的UserID
		jwt_UserID, err := utils.GetJwtClaimsUsername(tokenString)
		if err != nil {
			utils.Error(ctx, http.StatusUnauthorized, "token验证异常")
		} else {
			ctx.Set(CtxUserIDKey, jwt_UserID)
			ctx.Next()
		}
	}
}
