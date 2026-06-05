package post

import (
	postv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/post/v1"
	"github.com/miiy/goc-quickstart/api-gateway/internal/transport"

	"github.com/miiy/goc/gin"
)

func (m *Module) get(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	resp, err := m.client.GetPost(c.Request.Context(), &postv1.GetPostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) create(c *gin.Context) {
	var req postv1.CreatePostRequest
	if !transport.BindProto(c, &req) {
		return
	}

	resp, err := m.client.CreatePost(c.Request.Context(), &req)
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

	var req postv1.UpdatePostRequest
	if !transport.BindProto(c, &req) {
		return
	}
	req.Id = id

	resp, err := m.client.UpdatePost(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) delete(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	resp, err := m.client.DeletePost(c.Request.Context(), &postv1.DeletePostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) list(c *gin.Context) {
	authorID, ok := transport.Int64Query(c, "author_id", "authorId")
	if !ok {
		return
	}
	categoryID, ok := transport.Int64Query(c, "category_id", "categoryId")
	if !ok {
		return
	}
	page, ok := transport.Int32Query(c, "page", "")
	if !ok {
		return
	}
	pageSize, ok := transport.Int32Query(c, "page_size", "pageSize")
	if !ok {
		return
	}

	resp, err := m.client.ListPosts(c.Request.Context(), &postv1.ListPostsRequest{
		AuthorId:   authorID,
		CategoryId: categoryID,
		Tag:        transport.QueryValue(c, "tag", ""),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}
