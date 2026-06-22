package user

import (
	"net/http"
	"strings"

	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
	"github.com/miiy/goc/gin/sessions"
)

type ProfileView struct {
	template.ViewData
	User          *client.UserResponse
	Flashes       []sessions.Flash
	ProfileError  string
	PasswordError string
}

func Profile(c *gin.Context) {
	user, err := loadCurrentProfile(c)
	if err != nil {
		if handleAuthError(c, err) {
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
	renderProfile(c, http.StatusOK, user, flashes, "", "")
}

func UpdateProfile(c *gin.Context) {
	authUser, ok := currentAuthUser(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	nickname := strings.TrimSpace(c.PostForm("nickname"))
	email := strings.TrimSpace(c.PostForm("email"))
	_, err := userModule.userClient.UpdateProfile(c.Request.Context(), authUser.ID, nickname, email)
	if err != nil {
		if handleAuthError(c, err) {
			return
		}
		renderProfile(c, profileErrorStatus(err), &client.UserResponse{
			Id:       client.Int64String(authUser.ID),
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

func UploadAvatar(c *gin.Context) {
	if _, ok := currentAuthUser(c); !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		if handleAuthError(c, err) {
			return
		}
		user, loadErr := loadCurrentProfile(c)
		if loadErr != nil {
			user = nil
		}
		renderProfile(c, http.StatusBadRequest, user, nil, err.Error(), "")
		return
	}
	defer file.Close()

	if _, err := userModule.fileClient.UploadAvatar(c.Request.Context(), header.Filename, file); err != nil {
		if handleAuthError(c, err) {
			return
		}
		user, loadErr := loadCurrentProfile(c)
		if loadErr != nil {
			user = nil
		}
		renderProfile(c, profileErrorStatus(err), user, nil, err.Error(), "")
		return
	}

	if err := sessions.AddFlash(c, sessions.FlashLevelSuccess, "头像已更新"); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存提示信息失败")
		return
	}
	c.Redirect(http.StatusFound, "/user/profile")
}

func ChangePassword(c *gin.Context) {
	if _, ok := currentAuthUser(c); !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	err := userModule.authClient.ChangePassword(
		c.Request.Context(),
		c.PostForm("old_password"),
		c.PostForm("new_password"),
		c.PostForm("new_password_confirmation"),
	)
	if err != nil {
		if handleAuthError(c, err) {
			return
		}
		user, loadErr := loadCurrentProfile(c)
		if loadErr != nil {
			user = nil
		}
		renderProfile(c, profileErrorStatus(err), user, nil, "", err.Error())
		return
	}

	if err := sessions.AddFlash(c, sessions.FlashLevelSuccess, "密码已更新"); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "保存提示信息失败")
		return
	}
	c.Redirect(http.StatusFound, "/user/profile")
}

func loadCurrentProfile(c *gin.Context) (*client.UserResponse, error) {
	authUser, ok := currentAuthUser(c)
	if !ok {
		return nil, &client.HTTPError{StatusCode: http.StatusUnauthorized, Message: "unauthenticated"}
	}
	return userModule.userClient.GetUser(c.Request.Context(), authUser.ID)
}

func currentAuthUser(c *gin.Context) (*gocauth.AuthenticatedUser, bool) {
	if user, ok := c.Get("currentUser"); ok {
		if authUser, ok := user.(*gocauth.AuthenticatedUser); ok && authUser != nil {
			return authUser, true
		}
		if authUser, ok := gocauthmid.SessionUser(user); ok {
			return authUser, true
		}
	}

	session := sessions.Default(c)
	return gocauthmid.SessionUser(session.Get(gocauthmid.SessionKeyAuthUser))
}

func handleAuthError(c *gin.Context, err error) bool {
	if !client.IsStatus(err, http.StatusUnauthorized) {
		return false
	}
	userModule.sessionManager.Clear(c)
	c.Redirect(http.StatusFound, "/login")
	return true
}

func renderProfile(c *gin.Context, status int, user *client.UserResponse, flashes []sessions.Flash, profileError, passwordError string) {
	c.HTML(status, "user/profile", ProfileView{
		ViewData:      template.NewFormViewData(c),
		User:          user,
		Flashes:       flashes,
		ProfileError:  profileError,
		PasswordError: passwordError,
	})
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
