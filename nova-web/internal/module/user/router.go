package user

import "github.com/miiy/goc/gin"

func Router(r *gin.Engine) {
	g := r.Group("/user")
	{
		g.GET("/profile", Profile)
		g.POST("/profile", UpdateProfile)
		g.POST("/avatar", UploadAvatar)
		g.PUT("/password", ChangePassword)
	}
}

func Templates() map[string][]string {
	return map[string][]string{
		"user/profile": {"layout/layout.html", "layout/header.html", "layout/footer.html", "user/profile.html"},
	}
}
