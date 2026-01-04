package middleware

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/common/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	appService "lotteryBackend/service/app"
)

var h5AuthService = appService.H5AuthService{}

// H5Auth H5用户认证中间件
func H5Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取Token
		token := getH5Token(c)
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或登录已过期", c)
			c.Abort()
			return
		}

		// 2. 验证Token
		claims, err := h5AuthService.ValidateToken(token)
		if err != nil {
			global.GVA_LOG.Error("H5 Token验证失败", zap.Error(err))
			response.FailWithDetailed(gin.H{"reload": true}, "登录已过期，请重新登录", c)
			c.Abort()
			return
		}

		// 3. 检查用户是否存在
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, claims.UserId).Error; err != nil {
			response.FailWithDetailed(gin.H{"reload": true}, "用户不存在", c)
			c.Abort()
			return
		}

		// 4. 设置用户信息到上下文
		c.Set("h5User", &user)
		c.Set("h5UserId", user.ID)
		c.Set("h5OpenId", user.OpenId)
		c.Set("h5Claims", claims)

		c.Next()
	}
}

// H5AuthOptional H5用户认证中间件（可选，不强制登录）
func H5AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getH5Token(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := h5AuthService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, claims.UserId).Error; err != nil {
			c.Next()
			return
		}

		c.Set("h5User", &user)
		c.Set("h5UserId", user.ID)
		c.Set("h5OpenId", user.OpenId)
		c.Set("h5Claims", claims)

		c.Next()
	}
}

// H5RegisteredRequired 要求用户已报名且通过审核
func H5RegisteredRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("h5User")
		if !exists {
			response.FailWithMessage("请先登录", c)
			c.Abort()
			return
		}

		h5User := user.(*annual.AnnualUser)

		// 检查是否已报名
		if h5User.RealName == "" {
			response.FailWithDetailed(gin.H{"needRegister": true}, "请先完成报名", c)
			c.Abort()
			return
		}

		// 检查审核状态
		if *h5User.Status == 0 {
			response.FailWithDetailed(gin.H{"pending": true}, "您的报名正在审核中", c)
			c.Abort()
			return
		}

		if *h5User.Status == 2 {
			response.FailWithDetailed(gin.H{"rejected": true}, "您的报名未通过审核", c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// getH5Token 获取H5 Token
func getH5Token(c *gin.Context) string {
	// 优先从Header获取
	token := c.GetHeader("x-token")
	if token != "" {
		return token
	}

	// 其次从Authorization获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 支持 "Bearer token" 格式
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		return authHeader
	}

	// 最后从Query参数获取（用于WebSocket连接）
	token = c.Query("token")
	if token != "" {
		return token
	}

	return ""
}

// GetH5UserId 从上下文获取用户ID
func GetH5UserId(c *gin.Context) uint {
	userId, exists := c.Get("h5UserId")
	if !exists {
		return 0
	}
	return userId.(uint)
}

// GetH5User 从上下文获取用户信息
func GetH5User(c *gin.Context) *annual.AnnualUser {
	user, exists := c.Get("h5User")
	if !exists {
		return nil
	}
	return user.(*annual.AnnualUser)
}
