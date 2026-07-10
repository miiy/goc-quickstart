package user

import (
	"net/http"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *UsersAPI) GetProfile(c *gin.Context) {
	id, ok := authctx.CurrentUserInt64ID(c)
	if !ok || id <= 0 {
		transport.WriteError(c, status.Error(codes.Unauthenticated, "invalid authenticated user"))
		return
	}

	resp, err := api.userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.GetProfileResponse{User: protoToUser(resp.GetUser())})
}

func (api *UsersAPI) UpdateProfile(c *gin.Context) {
	id, ok := authctx.CurrentUserInt64ID(c)
	if !ok || id <= 0 {
		transport.WriteError(c, status.Error(codes.Unauthenticated, "invalid authenticated user"))
		return
	}

	var req openapi.UpdateProfileRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	user, err := openapiToProtoUser(req.User)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	updateMask, err := openapiToProtoUpdateMask(req.UpdateFields)
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
	c.JSON(http.StatusOK, openapi.UpdateProfileResponse{User: protoToUser(resp.GetUser())})
}
