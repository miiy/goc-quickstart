package post

import (
	"context"
	"net/http"
	"strings"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
)

type PostsAPI struct {
	postClient postv1.PostServiceClient
	userClient userClient
}

func NewPostsAPI(postClient postv1.PostServiceClient, userClient userv1.UserServiceClient) openapi.PostsAPI {
	return &PostsAPI{
		postClient: postClient,
		userClient: userClient,
	}
}

func (api *PostsAPI) GetPost(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	resp, err := api.postClient.GetPost(c.Request.Context(), &postv1.GetPostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	authorNames, err := api.authorNames(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.GetPostResponse{Post: openapiPost(resp.GetPost(), authorNames)})
}

func (api *PostsAPI) CreatePost(c *gin.Context) {
	var req openapi.CreatePostRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	post, err := protoCreatePostInput(req.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	resp, err := api.postClient.CreatePost(c.Request.Context(), &postv1.CreatePostRequest{Post: post})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	authorNames, err := api.authorNames(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.CreatePostResponse{Post: openapiPost(resp.GetPost(), authorNames)})
}

func (api *PostsAPI) UpdatePost(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	var req openapi.UpdatePostRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	post, err := protoUpdatePostInput(req.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	updateMask, err := protoUpdateMask(req.UpdateFields)
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	resp, err := api.postClient.UpdatePost(c.Request.Context(), &postv1.UpdatePostRequest{
		Id:         id,
		Post:       post,
		UpdateMask: updateMask,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	authorNames, err := api.authorNames(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.UpdatePostResponse{Post: openapiPost(resp.GetPost(), authorNames)})
}

func (api *PostsAPI) DeletePost(c *gin.Context) {
	id, ok := transport.Int64Param(c, "id")
	if !ok {
		return
	}

	_, err := api.postClient.DeletePost(c.Request.Context(), &postv1.DeletePostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{})
}

func (api *PostsAPI) ListPosts(c *gin.Context) {
	authorID, ok := transport.Int64Query(c, "author_id")
	if !ok {
		return
	}
	categoryID, ok := transport.Int64Query(c, "category_id")
	if !ok {
		return
	}
	page, ok := transport.Int32Query(c, "page")
	if !ok {
		return
	}
	pageSize, ok := transport.Int32Query(c, "page_size")
	if !ok {
		return
	}

	resp, err := api.postClient.ListPosts(c.Request.Context(), &postv1.ListPostsRequest{
		AuthorId:   authorID,
		CategoryId: categoryID,
		Tag:        c.Query("tag"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	authorNames, err := api.authorNames(c.Request.Context(), resp.Posts...)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.ListPostsResponse{
		Total:       resp.GetTotal(),
		TotalPages:  resp.GetTotalPages(),
		PageSize:    resp.GetPageSize(),
		CurrentPage: resp.GetCurrentPage(),
		Posts:       openapiPosts(resp.GetPosts(), authorNames),
	})
}

func (api *PostsAPI) authorNames(ctx context.Context, posts ...*postv1.Post) (map[int64]string, error) {
	ids := uniqueAuthorIDs(posts)
	if len(ids) == 0 || api.userClient == nil {
		return nil, nil
	}

	resp, err := api.userClient.BatchGetUsers(ctx, &userv1.BatchGetUsersRequest{Ids: ids})
	if err != nil {
		return nil, err
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
	return names, nil
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

type userClient interface {
	BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error)
}

var _ userClient = userv1.UserServiceClient(nil)
