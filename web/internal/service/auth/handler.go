package auth

import (
	"net/http"

	"github.com/miiy/goc-quickstart/web/internal/template"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
	"github.com/miiy/goc/gin/sessions"
)

type AuthFormView struct {
	template.ViewData
	Email    string
	Flashes  []sessions.Flash
	Username string
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
		ViewData: template.NewFormViewData(c),
		Flashes:  flashes,
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

	resp, err := authModule.client.Login(c.Request.Context(), username, password)
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
	if err := saveLoginSession(c, map[string]any{
		"id":       int64(resp.User.Id),
		"username": sessionUsername,
	}, resp.AccessToken, resp.ExpiresAt); err != nil {
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

	session := sessions.Default(c)
	accessToken, _ := session.Get(SessionKeyAccessToken).(string)
	if accessToken != "" {
		_ = authModule.client.Logout(c.Request.Context(), accessToken)
	}
	clearSession(c)

	c.Redirect(http.StatusFound, "/login")
}

func Profile(c *gin.Context) {
	session := sessions.Default(c)
	user, _ := gocauthmid.SessionUser(session.Get(gocauthmid.SessionKeyAuthUser))
	c.HTML(http.StatusOK, "auth/profile", ProfileView{
		ViewData: template.NewFormViewData(c),
		User:     user,
	})
}

func saveLoginSession(c *gin.Context, user map[string]any, accessToken string, expiresAt string) error {
	clearSession(c)

	session, err := authModule.sessionStore.New(requestWithoutSessionCookie(c.Request, authModule.sessionName), authModule.sessionName)
	if err != nil {
		return err
	}
	session.Values[gocauthmid.SessionKeyAuthUser] = user
	session.Values[SessionKeyAccessToken] = accessToken
	if expiresAt != "" {
		session.Values[SessionKeyAccessExpiresAt] = expiresAt
	}
	return authModule.sessionStore.Save(c.Request, c.Writer, session)
}

func clearSession(c *gin.Context) {
	session, err := authModule.sessionStore.New(c.Request, authModule.sessionName)
	if err != nil {
		_ = c.Error(err)
		return
	}
	for key := range session.Values {
		delete(session.Values, key)
	}
	session.Options.MaxAge = -1
	_ = authModule.sessionStore.Save(c.Request, c.Writer, session)
}

func requestWithoutSessionCookie(r *http.Request, sessionName string) *http.Request {
	req := r.Clone(r.Context())
	req.Header.Del("Cookie")
	for _, cookie := range r.Cookies() {
		if cookie.Name != sessionName {
			req.AddCookie(cookie)
		}
	}
	return req
}
