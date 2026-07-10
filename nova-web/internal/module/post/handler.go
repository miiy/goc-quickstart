package post

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	postv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-web/internal/media"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc-quickstart/nova-web/internal/transport"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/sqids"
	"github.com/unknwon/paginater"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsHandler struct {
	postClient postv1.PostServiceClient
	userClient userv1.UserServiceClient
}

func NewPostsHandler(postClient postv1.PostServiceClient, userClient userv1.UserServiceClient) *PostsHandler {
	return &PostsHandler{
		postClient: postClient,
		userClient: userClient,
	}
}

type PostListViewData struct {
	template.ViewData
	Posts       []PostView
	Total       int64
	CurrentPage int32
	TotalPages  int32
	PageSize    int32
	Pager       *paginater.Paginater
	Error       string
}

type PostDetailViewData struct {
	template.ViewData
	Post      *PostView
	CanManage bool
}

type PostView struct {
	Id        string
	UserId    int64
	User      PostUserView
	Title     string
	Summary   string
	Content   string
	CoverUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostUserView struct {
	Username string
	Nickname string
	Avatar   string
}

func (h *PostsHandler) list(c *gin.Context) {
	page := int32Param(c.DefaultQuery("page", "1"), 1)
	if pathPage := strings.TrimSpace(c.Param("page")); pathPage != "" {
		page = int32Param(pathPage, 1)
	}
	pageSize := int32Param(c.DefaultQuery("page_size", "10"), 10)
	if pageSize > 100 {
		pageSize = 100
	}

	if h.postClient == nil {
		template.InternalError(c)
		return
	}

	resp, err := h.postClient.ListPosts(c.Request.Context(), &postv1.ListPostsRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   postv1.PostStatus_POST_STATUS_PUBLISHED,
	})
	if err != nil {
		template.InternalError(c)
		return
	}

	view := PostListViewData{
		ViewData:    template.NewViewData(c),
		CurrentPage: page,
		PageSize:    pageSize,
		Pager:       paginater.New(0, int(pageSize), int(page), 5),
	}
	if resp == nil {
		c.HTML(http.StatusOK, "post/list", view)
		return
	}

	usersByID, err := h.postUsersByID(c.Request.Context(), resp.GetPosts()...)
	if err != nil {
		view.Error = err.Error()
		c.HTML(http.StatusOK, "post/list", view)
		return
	}

	currentPage := resp.GetCurrentPage()
	if currentPage < 1 {
		currentPage = page
	}
	respPageSize := resp.GetPageSize()
	if respPageSize < 1 {
		respPageSize = pageSize
	}

	view.Posts = postsForView(resp.GetPosts(), usersByID)
	view.Total = resp.GetTotal()
	view.CurrentPage = currentPage
	view.TotalPages = resp.GetTotalPages()
	view.PageSize = respPageSize
	view.Pager = paginater.New(int(resp.GetTotal()), int(respPageSize), int(currentPage), 5)
	c.HTML(http.StatusOK, "post/list", view)
}

func (h *PostsHandler) show(c *gin.Context) {
	id, err := decodePostID(c.Param("id"))
	if err != nil {
		template.NotFound(c)
		return
	}
	if h.postClient == nil {
		template.InternalError(c)
		return
	}

	resp, err := h.postClient.GetPost(c.Request.Context(), &postv1.GetPostRequest{Id: id})
	if err != nil {
		if transport.IsStatus(transport.FromGRPCError(err), http.StatusNotFound) {
			template.NotFound(c)
			return
		}
		template.InternalError(c)
		return
	}
	if !publicPostVisible(resp.GetPost()) {
		template.NotFound(c)
		return
	}

	usersByID, err := h.postUsersByID(c.Request.Context(), resp.GetPost())
	if err != nil {
		template.InternalError(c)
		return
	}

	post := postForView(resp.GetPost(), usersByID)
	canManage := false
	if post != nil && post.UserId > 0 {
		userID, ok := authctx.CurrentUserInt64ID(c)
		canManage = ok && userID == post.UserId
	}
	c.HTML(http.StatusOK, "post/detail", PostDetailViewData{
		ViewData:  template.NewViewData(c),
		Post:      post,
		CanManage: canManage,
	})
}

func (h *PostsHandler) create(c *gin.Context) {
	c.HTML(http.StatusOK, "post/create", template.NewFormViewData(c))
}

func (h *PostsHandler) edit(c *gin.Context) {
	if strings.TrimSpace(c.Param("id")) == "" {
		template.NotFound(c)
		return
	}
	c.HTML(http.StatusOK, "post/edit", template.NewFormViewData(c))
}

func (h *PostsHandler) postUsersByID(ctx context.Context, posts ...*postv1.Post) (map[int64]PostUserView, error) {
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
	if len(ids) == 0 || h.userClient == nil {
		return nil, nil
	}

	resp, err := h.userClient.BatchGetUsers(ctx, &userv1.BatchGetUsersRequest{Ids: ids})
	if err != nil {
		return nil, transport.FromGRPCError(err)
	}

	users := make(map[int64]PostUserView, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		if user == nil {
			continue
		}
		nickname := strings.TrimSpace(user.GetNickname())
		if nickname == "" {
			nickname = strings.TrimSpace(user.GetUsername())
		}
		users[user.GetId()] = PostUserView{
			Username: strings.TrimSpace(user.GetUsername()),
			Nickname: nickname,
			Avatar:   media.UploadsURL(strings.TrimSpace(user.GetAvatar())),
		}
	}
	return users, nil
}

func postsForView(posts []*postv1.Post, usersByID map[int64]PostUserView) []PostView {
	if len(posts) == 0 {
		return nil
	}
	viewPosts := make([]PostView, 0, len(posts))
	for _, post := range posts {
		if !publicPostVisible(post) {
			continue
		}
		viewPosts = append(viewPosts, *postForView(post, usersByID))
	}
	return viewPosts
}

func postForView(post *postv1.Post, usersByID map[int64]PostUserView) *PostView {
	if post == nil {
		return nil
	}
	createdAt := time.Time{}
	if post.GetCreatedAt() != nil {
		createdAt = post.GetCreatedAt().AsTime()
	}
	updatedAt := time.Time{}
	if post.GetUpdatedAt() != nil {
		updatedAt = post.GetUpdatedAt().AsTime()
	}
	return &PostView{
		Id:        encodePostID(post.GetId()),
		UserId:    post.GetUserId(),
		User:      usersByID[post.GetUserId()],
		Title:     post.GetTitle(),
		Summary:   post.GetSummary(),
		Content:   post.GetContent(),
		CoverUrl:  media.UploadsURL(post.GetCoverUrl()),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func publicPostVisible(post *postv1.Post) bool {
	return post != nil && post.GetStatus() == postv1.PostStatus_POST_STATUS_PUBLISHED
}

func int32Param(raw string, fallback int32) int32 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 32)
	if err != nil || parsed < 1 {
		return fallback
	}
	return int32(parsed)
}

var postIDEncoder = sqids.MustNew()

func encodePostID(id int64) string {
	if id <= 0 {
		return ""
	}
	return postIDEncoder.MustEncode(id)
}

func decodePostID(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, status.Error(codes.InvalidArgument, "invalid post id")
	}

	id, err := postIDEncoder.Decode(raw)
	if err != nil || id <= 0 {
		return 0, status.Error(codes.InvalidArgument, "invalid post id")
	}
	return id, nil
}
