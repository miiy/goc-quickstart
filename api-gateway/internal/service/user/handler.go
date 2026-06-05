package user

import (
	userv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/user/v1"
	"github.com/miiy/goc-quickstart/api-gateway/internal/transport"

	"github.com/miiy/goc/gin"
)

func (m *Module) get(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	resp, err := m.client.GetUser(c.Request.Context(), &userv1.GetUserRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) update(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	var req userv1.UpdateUserRequest
	if !transport.BindProto(c, &req) {
		return
	}
	req.Id = id

	resp, err := m.client.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) list(c *gin.Context) {
	page, ok := transport.Int32Query(c, "page", "")
	if !ok {
		return
	}
	pageSize, ok := transport.Int32Query(c, "page_size", "pageSize")
	if !ok {
		return
	}

	resp, err := m.client.ListUsers(c.Request.Context(), &userv1.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}
