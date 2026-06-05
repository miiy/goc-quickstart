package auth

import "github.com/miiy/goc/gin"

func (m *Module) RegisterPublicRouter(r *gin.Engine) {
	g := r.Group("/api/v1/auth")
	g.POST("/register", m.register)
	g.POST("/register/username_check", m.usernameCheck)
	g.POST("/register/email_check", m.emailCheck)
	g.POST("/register/phone_check", m.phoneCheck)
	g.POST("/login", m.login)
	g.POST("/send_sms_code", m.sendSMSCode)
	g.POST("/phone_auth", m.phoneAuth)
	g.POST("/mplogin", m.mpLogin)
}

func (m *Module) RegisterProtectedRouter(r gin.IRouter) {
	r.POST("/auth/token/refresh", m.refreshToken)
	r.POST("/auth/logout", m.logout)
}
