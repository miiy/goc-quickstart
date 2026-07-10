package user

import (
	"net/http"
	"strings"

	userv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/user/v1"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
)

type UserHandler struct{}

func NewUserHandler(_ userv1.UserServiceClient, _ *websession.Manager) *UserHandler {
	return &UserHandler{}
}

type ProfileView struct {
	template.ViewData
}

type UserShowView struct {
	template.ViewData
	Username string
}

func (h *UserHandler) Show(c *gin.Context) {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		template.NotFound(c)
		return
	}
	c.HTML(http.StatusOK, "user/show", UserShowView{
		ViewData: template.NewViewData(c),
		Username: username,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	c.HTML(http.StatusOK, "user/profile", ProfileView{
		ViewData: template.NewFormViewData(c),
	})
}
