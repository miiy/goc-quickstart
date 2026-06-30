package user

import (
	"net/http"
	"strings"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/media"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

type UserHandler struct {
	logger         logger.Logger
	authClient     *client.AuthClient
	userClient     *client.UserClient
	fileClient     *client.FileClient
	sessionManager *websession.Manager
}

func NewUserHandler(logger logger.Logger, authClient *client.AuthClient, userClient *client.UserClient, fileClient *client.FileClient, sessionManager *websession.Manager) *UserHandler {
	return &UserHandler{
		logger:         logger,
		authClient:     authClient,
		userClient:     userClient,
		fileClient:     fileClient,
		sessionManager: sessionManager,
	}
}

type ProfileView struct {
	template.ViewData
	User          *apiclient.User
	Flashes       []sessions.Flash
	ProfileError  string
	PasswordError string
}

func (h *UserHandler) Profile(c *gin.Context) {
	user, err := h.loadCurrentUser(c)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		c.HTML(http.StatusBadGateway, "user/profile", ProfileView{
			ViewData:     template.NewFormViewData(c),
			ProfileError: err.Error(),
		})
		return
	}

	flashes, err := sessions.Flashes(c)
	if err != nil {
		_ = c.Error(err)
	}
	h.renderProfile(c, http.StatusOK, user, flashes, "", "")
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	authUser, ok := authctx.CurrentUser(c)
	if !ok {
		h.handleAuthError(c, client.NewHTTPError(http.StatusUnauthorized, "unauthenticated"))
		return
	}

	authUserID, err := authUser.Int64ID()
	if err != nil {
		h.handleAuthError(c, client.NewHTTPError(http.StatusUnauthorized, err.Error()))
		return
	}

	nickname := strings.TrimSpace(c.PostForm("nickname"))
	email := strings.TrimSpace(c.PostForm("email"))
	_, err = h.userClient.UpdateUser(c.Request.Context(), authUserID, nickname, email)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		h.renderProfile(c, profileErrorStatus(err), &apiclient.User{
			Id:       authUserID,
			Username: authUser.Username,
			Nickname: nickname,
			Email:    email,
		}, nil, err.Error(), "")
		return
	}

	if err := sessions.AddFlash(c, sessions.FlashLevelSuccess, "资料已更新"); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存提示信息失败")
		return
	}
	c.Redirect(http.StatusFound, "/user/profile")
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		h.redirectProfileWithFlash(c, sessions.FlashLevelError, err.Error())
		return
	}
	defer file.Close()

	if _, err := h.fileClient.UploadAvatar(c.Request.Context(), header.Filename, file); err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		h.redirectProfileWithFlash(c, sessions.FlashLevelError, err.Error())
		return
	}

	h.redirectProfileWithFlash(c, sessions.FlashLevelSuccess, "头像已更新")
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	err := h.authClient.ChangePassword(
		c.Request.Context(),
		c.PostForm("old_password"),
		c.PostForm("new_password"),
		c.PostForm("new_password_confirmation"),
	)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		h.redirectProfileWithFlash(c, sessions.FlashLevelError, err.Error())
		return
	}

	h.redirectProfileWithFlash(c, sessions.FlashLevelSuccess, "密码已更新")
}

func (h *UserHandler) loadCurrentUser(c *gin.Context) (*apiclient.User, error) {
	authUser, ok := authctx.CurrentUser(c)
	if !ok {
		return nil, client.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
	}

	authUserID, err := authUser.Int64ID()
	if err != nil {
		return nil, client.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	return h.userClient.GetUser(c.Request.Context(), authUserID)
}

func (h *UserHandler) handleAuthError(c *gin.Context, err error) bool {
	if !client.IsStatus(err, http.StatusUnauthorized) {
		return false
	}
	h.sessionManager.Clear(c)
	c.Redirect(http.StatusFound, "/login")
	return true
}

func (h *UserHandler) renderProfile(c *gin.Context, status int, user *apiclient.User, flashes []sessions.Flash, profileError, passwordError string) {
	if user != nil {
		viewUser := *user
		viewUser.Avatar = media.UploadsURL(viewUser.Avatar)
		user = &viewUser
	}
	c.HTML(status, "user/profile", ProfileView{
		ViewData:      template.NewFormViewData(c),
		User:          user,
		Flashes:       flashes,
		ProfileError:  profileError,
		PasswordError: passwordError,
	})
}

func (h *UserHandler) redirectProfileWithFlash(c *gin.Context, level, message string) {
	if err := sessions.AddFlash(c, level, message); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存提示信息失败")
		return
	}
	c.Redirect(http.StatusFound, "/user/profile")
}

func profileErrorStatus(err error) int {
	switch {
	case client.IsStatus(err, http.StatusBadRequest):
		return http.StatusBadRequest
	case client.IsStatus(err, http.StatusForbidden):
		return http.StatusForbidden
	default:
		return http.StatusBadGateway
	}
}
