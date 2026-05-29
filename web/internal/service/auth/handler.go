package auth

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc-quickstart/web/internal/template"
)

func init() {
	gob.Register(&gocauth.AuthenticatedUser{})
}

type AuthFormView struct {
	template.ViewData
	Flashes []sessions.Flash
}

type ProfileView struct {
	template.ViewData
	User *gocauth.AuthenticatedUser
}

func RegisterForm(c *gin.Context) {
	flashes, err := sessions.Flashes(c)
	if err != nil {
		_ = c.Error(err)
	}

	c.HTML(http.StatusOK, "auth/register", AuthFormView{
		ViewData: template.ViewData{
			IsLoggedIn: c.GetBool("isLoggedIn"),
		},
		Flashes: flashes,
	})
}

func Register(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")
	passwordConfirmation := c.PostForm("password_confirmation")

	_, err := authModule.client.Register(c.Request.Context(), email, username, password, passwordConfirmation)
	if err != nil {
		c.HTML(http.StatusBadRequest, "auth/register", AuthFormView{
			ViewData: template.ViewData{
				IsLoggedIn: c.GetBool("isLoggedIn"),
			},
			Flashes: []sessions.Flash{
				{Level: sessions.FlashLevelError, Message: "注册失败：" + err.Error()},
			},
		})
		return
	}

	if err := sessions.AddFlash(c, sessions.FlashLevelSuccess, "注册成功，请登录"); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存提示信息失败")
		return
	}
	c.Redirect(http.StatusFound, "/login")
}

func LoginForm(c *gin.Context) {
	flashes, err := sessions.Flashes(c)
	if err != nil {
		_ = c.Error(err)
	}

	c.HTML(http.StatusOK, "auth/login", AuthFormView{
		ViewData: template.ViewData{
			IsLoggedIn: c.GetBool("isLoggedIn"),
		},
		Flashes: flashes,
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	resp, err := authModule.client.Login(c.Request.Context(), username, password)
	if err != nil {
		if err := sessions.AddFlash(c, sessions.FlashLevelError, "用户名或密码错误"); err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "保存提示信息失败")
			return
		}
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Set authenticated user in session
	user := &gocauth.AuthenticatedUser{
		Username: resp.User.Username,
	}
	session := sessions.Default(c)
	session.Set(gocauthmid.SessionKeyAuthUser, user)
	if err := session.Save(); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存会话失败")
		return
	}

	c.SetCookie("token", resp.AccessToken, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(gocauthmid.SessionKeyAuthUser)
	_ = session.Save()

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func Profile(c *gin.Context) {
	session := sessions.Default(c)
	user, _ := session.Get(gocauthmid.SessionKeyAuthUser).(*gocauth.AuthenticatedUser)
	c.HTML(http.StatusOK, "auth/profile", ProfileView{
		ViewData: template.ViewData{
			IsLoggedIn: c.GetBool("isLoggedIn"),
		},
		User: user,
	})
}

func AuthRequired() gin.HandlerFunc {
	return gocauthmid.SessionAuthenticationMiddleware("/login")
}
