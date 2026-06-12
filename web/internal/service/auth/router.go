package auth

import (
	"github.com/miiy/goc/gin"
)

func Router(r *gin.Engine) {
	r.GET("/register", RegisterForm)
	r.POST("/register", Register)
	r.GET("/login", LoginForm)
	r.POST("/login", Login)

	g := r.Group("/auth")
	{
		g.GET("/logout", Logout)
		g.POST("/logout", Logout)
		g.GET("/profile", ProfileRedirect)
	}
}

func Templates() map[string][]string {
	return map[string][]string{
		"auth/register": {"layout/layout.html", "layout/header.html", "layout/footer.html", "auth/register.html"},
		"auth/login":    {"layout/layout.html", "layout/header.html", "layout/footer.html", "auth/login.html"},
	}
}
