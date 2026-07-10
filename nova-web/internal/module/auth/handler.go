package auth

import (
	"net/http"
	"strconv"
	"time"

	authv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc-quickstart/nova-web/internal/transport"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

type AuthHandler struct {
	logger          logger.Logger
	authClient      authv1.AuthServiceClient
	sessionManager  *session.Manager
	registerEnabled bool
}

func NewAuthHandler(logger logger.Logger, authClient authv1.AuthServiceClient, sessionManager *session.Manager, registerEnabled bool) *AuthHandler {
	return &AuthHandler{
		logger:          logger,
		authClient:      authClient,
		sessionManager:  sessionManager,
		registerEnabled: registerEnabled,
	}
}

type AuthFormView struct {
	template.ViewData
	Email    string
	Flashes  []sessions.Flash
	Username string
}

func (h *AuthHandler) RegisterForm(c *gin.Context) {
	if h.rejectDisabledRegister(c) {
		return
	}

	flashes, err := sessions.Flashes(c)
	if err != nil {
		_ = c.Error(err)
	}

	c.HTML(http.StatusOK, "auth/register", AuthFormView{
		ViewData: template.NewFormViewData(c),
		Flashes:  flashes,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	if h.rejectDisabledRegister(c) {
		return
	}

	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")
	passwordConfirmation := c.PostForm("password_confirmation")

	_, err := h.authClient.Register(c.Request.Context(), &authv1.RegisterRequest{
		Email:                email,
		Username:             username,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
	})
	if err != nil {
		err = transport.FromGRPCError(err)
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

func (h *AuthHandler) rejectDisabledRegister(c *gin.Context) bool {
	if h.registerEnabled {
		return false
	}
	if err := sessions.AddFlash(c, sessions.FlashLevelError, "注册暂未开放"); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存提示信息失败")
		return true
	}
	c.Redirect(http.StatusFound, "/login")
	return true
}

func (h *AuthHandler) LoginForm(c *gin.Context) {
	flashes, err := sessions.Flashes(c)
	if err != nil {
		_ = c.Error(err)
	}

	c.HTML(http.StatusOK, "auth/login", AuthFormView{
		ViewData: template.NewFormViewData(c),
		Flashes:  flashes,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	resp, err := h.authClient.Login(c.Request.Context(), &authv1.LoginRequest{
		Username: username,
		Password: password,
	})
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

	sessionUsername := resp.GetUser().GetUsername()
	if sessionUsername == "" {
		sessionUsername = username
	}
	if err := h.sessionManager.SaveLoginSession(c, map[string]any{
		"id":       strconv.FormatInt(resp.GetUser().GetId(), 10),
		"username": sessionUsername,
	}, resp.GetAccessToken(), formatAPITime(resp.GetExpiresAt().AsTime()), resp.GetRefreshToken()); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存会话失败")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	accessToken, refreshToken := h.sessionManager.Tokens(c)
	if accessToken != "" || refreshToken != "" {
		_, _ = h.authClient.Logout(c.Request.Context(), &authv1.LogoutRequest{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
	h.sessionManager.Clear(c)

	c.Redirect(http.StatusFound, "/login")
}

func formatAPITime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
