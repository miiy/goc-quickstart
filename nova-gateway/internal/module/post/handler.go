package post

import (
	"context"
	"strings"

	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
)

func (m *Module) get(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	resp, err := m.postClient.GetPost(c.Request.Context(), &postv1.GetPostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if err := m.enrichPostAuthors(c.Request.Context(), resp.Post); err != nil {
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

	resp, err := m.postClient.CreatePost(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if err := m.enrichPostAuthors(c.Request.Context(), resp.Post); err != nil {
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

	resp, err := m.postClient.UpdatePost(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if err := m.enrichPostAuthors(c.Request.Context(), resp.Post); err != nil {
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

	resp, err := m.postClient.DeletePost(c.Request.Context(), &postv1.DeletePostRequest{Id: id})
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

	resp, err := m.postClient.ListPosts(c.Request.Context(), &postv1.ListPostsRequest{
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
	if err := m.enrichPostAuthors(c.Request.Context(), resp.Posts...); err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) enrichPostAuthors(ctx context.Context, posts ...*postv1.Post) error {
	ids := uniqueAuthorIDs(posts)
	if len(ids) == 0 || m.userClient == nil {
		return nil
	}

	resp, err := m.userClient.BatchGetUsers(ctx, &userv1.BatchGetUsersRequest{Ids: ids})
	if err != nil {
		return err
	}

	names := make(map[int64]string, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		name := strings.TrimSpace(user.GetNickname())
		if name == "" {
			name = strings.TrimSpace(user.GetUsername())
		}
		if name != "" {
			names[user.GetId()] = name
		}
	}
	for _, post := range posts {
		if post == nil {
			continue
		}
		post.AuthorName = names[post.GetAuthorId()]
	}
	return nil
}

func uniqueAuthorIDs(posts []*postv1.Post) []int64 {
	seen := make(map[int64]struct{}, len(posts))
	ids := make([]int64, 0, len(posts))
	for _, post := range posts {
		if post == nil || post.GetAuthorId() <= 0 {
			continue
		}
		id := post.GetAuthorId()
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

type authorUserClient interface {
	BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error)
}

var _ authorUserClient = userv1.UserServiceClient(nil)
