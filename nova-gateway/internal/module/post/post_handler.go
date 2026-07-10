package post

import (
	"context"
	"net/http"
	"strings"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/media"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	id, err := decodePostID(c.Param("id"))
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	resp, err := api.postClient.GetPost(c.Request.Context(), &postv1.GetPostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if !publicPostVisible(resp.GetPost()) {
		transport.WriteError(c, status.Error(codes.NotFound, "post not found"))
		return
	}
	postUsersByID, err := api.postUsersByID(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.GetPostResponse{Post: protoToPost(resp.GetPost(), postUsersByID, optionalCurrentUserID(c))})
}

func (api *PostsAPI) GetUserPost(c *gin.Context) {
	userID, ok := api.currentUserIDForUsername(c)
	if !ok {
		return
	}

	id, err := decodePostID(c.Param("id"))
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	resp, err := api.postClient.GetPost(c.Request.Context(), &postv1.GetPostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if resp.GetPost() == nil || resp.GetPost().GetUserId() != userID {
		transport.WriteError(c, status.Error(codes.NotFound, "post not found"))
		return
	}
	postUsersByID, err := api.postUsersByID(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.GetPostResponse{Post: protoToPost(resp.GetPost(), postUsersByID, userID)})
}

func (api *PostsAPI) ListUserPosts(c *gin.Context) {
	userID, ok := api.currentUserIDForUsername(c)
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
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	postUsersByID, err := api.postUsersByID(c.Request.Context(), resp.Posts...)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.ListPostsResponse{
		Total:       resp.GetTotal(),
		TotalPages:  resp.GetTotalPages(),
		PageSize:    resp.GetPageSize(),
		CurrentPage: resp.GetCurrentPage(),
		Posts:       protoToPosts(resp.GetPosts(), postUsersByID, userID),
	})
}

func (api *PostsAPI) CreatePost(c *gin.Context) {
	var req openapi.CreatePostRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	post, err := openapiToProtoCreatePostInput(req.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	resp, err := api.postClient.CreatePost(c.Request.Context(), &postv1.CreatePostRequest{Post: post})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	postUsersByID, err := api.postUsersByID(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.CreatePostResponse{Post: protoToPost(resp.GetPost(), postUsersByID, optionalCurrentUserID(c))})
}

func (api *PostsAPI) currentUserIDForUsername(c *gin.Context) (int64, bool) {
	authUser, ok := authctx.CurrentUser(c)
	if !ok {
		transport.WriteError(c, status.Error(codes.Unauthenticated, "invalid authenticated user"))
		return 0, false
	}
	username := strings.TrimSpace(c.Param("username"))
	if username == "" || username != strings.TrimSpace(authUser.Username) {
		transport.WriteError(c, status.Error(codes.NotFound, "post not found"))
		return 0, false
	}
	userID, err := authUser.Int64ID()
	if err != nil || userID <= 0 {
		transport.WriteError(c, status.Error(codes.Unauthenticated, "invalid authenticated user"))
		return 0, false
	}
	return userID, true
}

func (api *PostsAPI) UpdatePost(c *gin.Context) {
	id, err := decodePostID(c.Param("id"))
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	var req openapi.UpdatePostRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	post, err := openapiToProtoUpdatePostInput(req.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	updateMask, err := openapiToProtoUpdateMask(req.UpdateFields)
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
	postUsersByID, err := api.postUsersByID(c.Request.Context(), resp.Post)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.UpdatePostResponse{Post: protoToPost(resp.GetPost(), postUsersByID, optionalCurrentUserID(c))})
}

func (api *PostsAPI) DeletePost(c *gin.Context) {
	id, err := decodePostID(c.Param("id"))
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	_, err = api.postClient.DeletePost(c.Request.Context(), &postv1.DeletePostRequest{Id: id})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{})
}

func (api *PostsAPI) ListPosts(c *gin.Context) {
	userID, ok := transport.Int64Query(c, "user_id")
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
		UserId:     userID,
		CategoryId: categoryID,
		Tag:        c.Query("tag"),
		Page:       page,
		PageSize:   pageSize,
		Status:     postv1.PostStatus_POST_STATUS_PUBLISHED,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	postUsersByID, err := api.postUsersByID(c.Request.Context(), resp.Posts...)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.ListPostsResponse{
		Total:       resp.GetTotal(),
		TotalPages:  resp.GetTotalPages(),
		PageSize:    resp.GetPageSize(),
		CurrentPage: resp.GetCurrentPage(),
		Posts:       protoToPosts(resp.GetPosts(), postUsersByID, 0),
	})
}

func optionalCurrentUserID(c *gin.Context) int64 {
	authUser, ok := authctx.CurrentUser(c)
	if !ok {
		return 0
	}
	userID, err := authUser.Int64ID()
	if err != nil || userID <= 0 {
		return 0
	}
	return userID
}

func (api *PostsAPI) postUsersByID(ctx context.Context, posts ...*postv1.Post) (map[int64]openapi.PostUser, error) {
	ids := uniqueUserIDs(posts)
	if len(ids) == 0 || api.userClient == nil {
		return nil, nil
	}

	resp, err := api.userClient.BatchGetUsers(ctx, &userv1.BatchGetUsersRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	postUsers := make(map[int64]openapi.PostUser, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		nickname := strings.TrimSpace(user.GetNickname())
		if nickname == "" {
			nickname = strings.TrimSpace(user.GetUsername())
		}
		postUsers[user.GetId()] = openapi.PostUser{
			Username: strings.TrimSpace(user.GetUsername()),
			Nickname: nickname,
			Avatar:   media.UploadsURL(user.GetAvatar()),
		}
	}
	return postUsers, nil
}

// Public reads intentionally hide drafts and unpublished posts behind 404.
func publicPostVisible(post *postv1.Post) bool {
	return post != nil && post.GetStatus() == postv1.PostStatus_POST_STATUS_PUBLISHED
}

func uniqueUserIDs(posts []*postv1.Post) []int64 {
	seen := make(map[int64]struct{}, len(posts))
	ids := make([]int64, 0, len(posts))
	for _, post := range posts {
		if post == nil || post.GetUserId() <= 0 {
			continue
		}
		id := post.GetUserId()
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
