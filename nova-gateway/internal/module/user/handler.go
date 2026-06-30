package user

import (
	"net/http"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
)

type UsersAPI struct {
	userClient userv1.UserServiceClient
}

func NewUsersAPI(userClient userv1.UserServiceClient) openapi.UsersAPI {
	return &UsersAPI{userClient: userClient}
}

func (api *UsersAPI) GetUser(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	resp, err := api.userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.GetUserResponse{User: OpenAPIUser(resp.GetUser())})
}

func (api *UsersAPI) BatchGetUsers(c *gin.Context) {
	ids, ok := transport.Int64SliceQuery(c, "ids")
	if !ok {
		return
	}

	resp, err := api.userClient.BatchGetUsers(c.Request.Context(), &userv1.BatchGetUsersRequest{Ids: ids})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.BatchGetUsersResponse{Users: openapiUsers(resp.GetUsers())})
}

func (api *UsersAPI) UpdateUser(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	var req openapi.UpdateUserRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	user, err := protoUserInput(req.User)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	updateMask, err := protoUpdateMask(req.UpdateFields)
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	resp, err := api.userClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		Id:         id,
		User:       user,
		UpdateMask: updateMask,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.UpdateUserResponse{User: OpenAPIUser(resp.GetUser())})
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
		Users:       openapiUsers(resp.GetUsers()),
	})
}
