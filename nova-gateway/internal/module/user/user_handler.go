package user

import (
	"net/http"
	"strings"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersAPI struct {
	userClient userv1.UserServiceClient
}

func NewUsersAPI(userClient userv1.UserServiceClient) openapi.UsersAPI {
	return &UsersAPI{userClient: userClient}
}

func (api *UsersAPI) GetUser(c *gin.Context) {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		transport.WriteError(c, status.Error(codes.InvalidArgument, "username is required"))
		return
	}

	resp, err := api.userClient.GetUserByUsername(c.Request.Context(), &userv1.GetUserByUsernameRequest{Username: username})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.GetUserResponse{User: protoToPublicUser(resp.GetUser())})
}

func (api *UsersAPI) ListUsers(c *gin.Context) {
	page, ok := transport.Int32Query(c, "page")
	if !ok {
		return
	}
	pageSize, ok := transport.Int32Query(c, "page_size")
	if !ok {
		return
	}

	resp, err := api.userClient.ListUsers(c.Request.Context(), &userv1.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.ListUsersResponse{
		Total:       resp.GetTotal(),
		TotalPages:  resp.GetTotalPages(),
		PageSize:    resp.GetPageSize(),
		CurrentPage: resp.GetCurrentPage(),
		Users:       protoToUsers(resp.GetUsers()),
	})
}
