package auth

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(public, protected gin.IRouter) {
	api := m.authAPI

	publicGroup := public.Group("/auth")
	publicGroup.POST("/register", api.Register)
	publicGroup.POST("/register/check-username", api.UsernameCheck)
	publicGroup.POST("/register/check-email", api.EmailCheck)
	publicGroup.POST("/register/check-phone", api.PhoneCheck)
	publicGroup.POST("/login", api.Login)
	publicGroup.POST("/sms/send-code", api.SendSmsCode)
	publicGroup.POST("/phone/login", api.PhoneAuth)
	publicGroup.POST("/mp/login", api.MpLogin)
	publicGroup.POST("/token/refresh", api.RefreshToken)

	protectedGroup := protected.Group("/auth")
	protectedGroup.PUT("/password", api.ChangePassword)
	protectedGroup.POST("/logout", api.Logout)
}
