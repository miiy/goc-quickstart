package auth

import (
	"net/http"
	"strconv"

	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/sessions"
)

type AuthFormView struct {
	template.ViewData
	Email    string
	Flashes  []sessions.Flash
	Username string
}

func RegisterForm(c *gin.Context) {
	flashes, err := sessions.Flashes(c)
	if err != nil {
		_ = c.Error(err)
	}

	c.HTML(http.StatusOK, "auth/register", AuthFormView{
		ViewData: template.NewFormViewData(c),
		Flashes:  flashes,
	})
}

func Register(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")
	passwordConfirmation := c.PostForm("password_confirmation")

	_, err := authModule.authClient.Register(c.Request.Context(), email, username, password, passwordConfirmation)
	if err != nil {
		c.HTML(http.StatusBadRequest, "auth/register", AuthFormView{
			ViewData: template.NewFormViewData(c),
			Flashes: []sessions.Flash{
				{Level: sessions.FlashLevelError, Message: "注册失败：" + err.Error()},
			},
			Email:    email,
			Username: username,
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
		ViewData: template.NewFormViewData(c),
		Flashes:  flashes,
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	resp, err := authModule.authClient.Login(c.Request.Context(), username, password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "auth/login", AuthFormView{
			ViewData: template.NewFormViewData(c),
			Flashes: []sessions.Flash{
				{Level: sessions.FlashLevelError, Message: "用户名或密码错误"},
			},
			Username: username,
		})
		return
	}

	sessionUsername := resp.User.Username
	if sessionUsername == "" {
		sessionUsername = username
	}
	if err := authModule.sessionManager.SaveLoginSession(c, map[string]any{
		"id":       strconv.FormatInt(int64(resp.User.Id), 10),
		"username": sessionUsername,
	}, resp.AccessToken, resp.ExpiresAt, resp.RefreshToken); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存会话失败")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	accessToken, refreshToken := authModule.sessionManager.Tokens(c)
	if accessToken != "" || refreshToken != "" {
		_ = authModule.authClient.Logout(c.Request.Context(), accessToken, refreshToken)
	}
	authModule.sessionManager.Clear(c)

	c.Redirect(http.StatusFound, "/login")
}
