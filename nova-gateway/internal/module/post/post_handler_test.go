package post

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakePostUserClient struct {
	gotIDs []int64
}

func (f *fakePostUserClient) BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error) {
	f.gotIDs = append([]int64(nil), in.GetIds()...)
	return &userv1.BatchGetUsersResponse{
		Users: []*userv1.User{
			{Id: 1, Username: "alice", Nickname: "Alice Nick", Avatar: "avatars/alice.png"},
			{Id: 2, Username: "bob", Avatar: "avatars/bob.png"},
			{Id: 7, Username: "carol", Nickname: "Carol Nick", Avatar: "avatars/carol.png"},
		},
	}, nil
}

func TestEnrichPostUsersBatchGetsUsers(t *testing.T) {
	userClient := &fakePostUserClient{}
	api := &PostsAPI{userClient: userClient}
	posts := []*postv1.Post{
		{Id: 1, UserId: 1},
		{Id: 2, UserId: 2},
		{Id: 3, UserId: 1},
		{Id: 4},
		nil,
	}

	postUsers, err := api.postUsersByID(context.Background(), posts...)
	if err != nil {
		t.Fatal(err)
	}

	if want := []int64{1, 2}; !reflect.DeepEqual(userClient.gotIDs, want) {
		t.Fatalf("batch ids = %v, want %v", userClient.gotIDs, want)
	}
	if postUsers[1].Nickname != "Alice Nick" || postUsers[1].Avatar != "/uploads/avatars/alice.png" {
		t.Fatalf("post user = %+v, want Alice Nick/avatar", postUsers[1])
	}
	if postUsers[1].Username != "alice" {
		t.Fatalf("post user username = %q, want alice", postUsers[1].Username)
	}
	if postUsers[2].Username != "bob" || postUsers[2].Nickname != "bob" || postUsers[2].Avatar != "/uploads/avatars/bob.png" {
		t.Fatalf("post user = %+v, want bob/avatar", postUsers[2])
	}
}

type fakePostClient struct {
	getPost              *postv1.Post
	listPosts            []*postv1.Post
	createReq            *postv1.CreatePostRequest
	updateReq            *postv1.UpdatePostRequest
	listReq              *postv1.ListPostsRequest
	listCategoriesCalled bool
}

func (f *fakePostClient) GetPost(ctx context.Context, in *postv1.GetPostRequest, opts ...grpc.CallOption) (*postv1.GetPostResponse, error) {
	if f.getPost != nil {
		return &postv1.GetPostResponse{Post: f.getPost}, nil
	}
	return &postv1.GetPostResponse{Post: testProtoPost(in.GetId())}, nil
}

func (f *fakePostClient) CreatePost(ctx context.Context, in *postv1.CreatePostRequest, opts ...grpc.CallOption) (*postv1.CreatePostResponse, error) {
	f.createReq = in
	return &postv1.CreatePostResponse{Post: testProtoPost(42)}, nil
}

func (f *fakePostClient) UpdatePost(ctx context.Context, in *postv1.UpdatePostRequest, opts ...grpc.CallOption) (*postv1.UpdatePostResponse, error) {
	f.updateReq = in
	post := testProtoPost(in.GetId())
	post.Title = in.GetPost().GetTitle()
	post.CoverUrl = in.GetPost().GetCoverUrl()
	return &postv1.UpdatePostResponse{Post: post}, nil
}

func (f *fakePostClient) DeletePost(ctx context.Context, in *postv1.DeletePostRequest, opts ...grpc.CallOption) (*postv1.DeletePostResponse, error) {
	return &postv1.DeletePostResponse{}, nil
}

func (f *fakePostClient) ListPosts(ctx context.Context, in *postv1.ListPostsRequest, opts ...grpc.CallOption) (*postv1.ListPostsResponse, error) {
	f.listReq = in
	posts := f.listPosts
	if len(posts) == 0 {
		posts = []*postv1.Post{testProtoPost(42)}
	}
	return &postv1.ListPostsResponse{
		Total:       1,
		TotalPages:  1,
		PageSize:    in.GetPageSize(),
		CurrentPage: in.GetPage(),
		Posts:       posts,
	}, nil
}

func (f *fakePostClient) ListCategories(ctx context.Context, in *postv1.ListCategoriesRequest, opts ...grpc.CallOption) (*postv1.ListCategoriesResponse, error) {
	f.listCategoriesCalled = true
	return &postv1.ListCategoriesResponse{
		Categories: []*postv1.Category{
			{Id: 1, Name: "Engineering", Path: "/engineering"},
		},
	}, nil
}

func TestCreatePostUsesOpenAPIRequestAndResponse(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.POST("/api/v1/posts", api.CreatePost)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(`{"post":{"title":"new title","summary":"short","content":"body","status":"published","tags":["go"],"category_id":9,"cover_url":"post-covers/cover.png"}}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if postClient.createReq == nil {
		t.Fatal("create request was not sent")
	}
	gotPost := postClient.createReq.GetPost()
	if gotPost.GetUserId() != 0 {
		t.Fatalf("user id = %d, want 0", gotPost.GetUserId())
	}
	if gotPost.GetTitle() != "new title" || gotPost.GetSummary() != "short" || gotPost.GetContent() != "body" || gotPost.GetCategoryId() != 9 || gotPost.GetCoverUrl() != "post-covers/cover.png" {
		t.Fatalf("unexpected proto post: %+v", gotPost)
	}
	if gotPost.GetStatus() != postv1.PostStatus_POST_STATUS_PUBLISHED {
		t.Fatalf("status = %s", gotPost.GetStatus())
	}

	var body struct {
		Post struct {
			Id     string `json:"id"`
			UserId int64  `json:"user_id"`
			User   struct {
				Username string `json:"username"`
				Nickname string `json:"nickname"`
				Avatar   string `json:"avatar"`
			} `json:"user"`
			CategoryId int64    `json:"category_id"`
			Status     string   `json:"status"`
			Tags       []string `json:"tags"`
			CoverUrl   string   `json:"cover_url"`
		} `json:"post"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Post.Id != encodePostID(42) || body.Post.UserId != 7 || body.Post.CategoryId != 9 {
		t.Fatalf("unexpected ids: %+v", body.Post)
	}
	if body.Post.User.Nickname != "" {
		t.Fatalf("user nickname = %q, want empty without user client", body.Post.User.Nickname)
	}
	if body.Post.User.Avatar != "" {
		t.Fatalf("user avatar = %q, want empty without user client", body.Post.User.Avatar)
	}
	if body.Post.Status != "published" {
		t.Fatalf("status = %q", body.Post.Status)
	}
	if body.Post.CoverUrl != "/uploads/post-covers/cover.png" {
		t.Fatalf("cover url = %q, want /uploads/post-covers/cover.png", body.Post.CoverUrl)
	}
	if !reflect.DeepEqual(body.Post.Tags, []string{"go"}) {
		t.Fatalf("tags = %v", body.Post.Tags)
	}
}

func TestGetPostAssemblesUserNicknameInGateway(t *testing.T) {
	postClient := &fakePostClient{}
	userClient := &fakePostUserClient{}
	api := &PostsAPI{postClient: postClient, userClient: userClient}

	r := gin.New()
	r.GET("/api/v1/posts/:id", api.GetPost)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+encodePostID(42), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if want := []int64{7}; !reflect.DeepEqual(userClient.gotIDs, want) {
		t.Fatalf("batch ids = %v, want %v", userClient.gotIDs, want)
	}
	var body struct {
		Post struct {
			User struct {
				Username string `json:"username"`
				Nickname string `json:"nickname"`
				Avatar   string `json:"avatar"`
			} `json:"user"`
			CoverUrl string `json:"cover_url"`
		} `json:"post"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Post.User.Nickname != "Carol Nick" {
		t.Fatalf("user nickname = %q, want Carol Nick", body.Post.User.Nickname)
	}
	if body.Post.User.Username != "carol" {
		t.Fatalf("user username = %q, want carol", body.Post.User.Username)
	}
	if body.Post.User.Avatar != "/uploads/avatars/carol.png" {
		t.Fatalf("user avatar = %q, want /uploads/avatars/carol.png", body.Post.User.Avatar)
	}
	if body.Post.CoverUrl != "/uploads/post-covers/cover.png" {
		t.Fatalf("cover url = %q, want /uploads/post-covers/cover.png", body.Post.CoverUrl)
	}
}

func TestGetPostIncludesCanManageForOwner(t *testing.T) {
	tests := []struct {
		name      string
		authUser  *gocauth.AuthenticatedUser
		canManage bool
	}{
		{name: "anonymous"},
		{name: "other user", authUser: &gocauth.AuthenticatedUser{ID: "8", Username: "alice"}},
		{name: "owner", authUser: &gocauth.AuthenticatedUser{ID: "7", Username: "carol"}, canManage: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postClient := &fakePostClient{}
			api := &PostsAPI{postClient: postClient}

			r := gin.New()
			r.GET("/api/v1/posts/:id", func(c *gin.Context) {
				if tt.authUser != nil {
					ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), tt.authUser)
					c.Request = c.Request.WithContext(ctx)
				}
				api.GetPost(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+encodePostID(42), nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
			}
			var body struct {
				Post struct {
					CanManage bool `json:"can_manage"`
				} `json:"post"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Post.CanManage != tt.canManage {
				t.Fatalf("can_manage = %v, want %v", body.Post.CanManage, tt.canManage)
			}
		})
	}
}

func TestGetPostHidesDrafts(t *testing.T) {
	postClient := &fakePostClient{getPost: testProtoPost(42)}
	postClient.getPost.Status = postv1.PostStatus_POST_STATUS_DRAFT
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.GET("/api/v1/posts/:id", api.GetPost)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+encodePostID(42), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetUserPostAllowsOwnerToReadDraft(t *testing.T) {
	postClient := &fakePostClient{getPost: testProtoPost(42)}
	postClient.getPost.Status = postv1.PostStatus_POST_STATUS_DRAFT
	userClient := &fakePostUserClient{}
	api := &PostsAPI{postClient: postClient, userClient: userClient}

	r := gin.New()
	r.GET("/api/v1/users/:username/posts/:id", func(c *gin.Context) {
		ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), &gocauth.AuthenticatedUser{ID: "7", Username: "carol"})
		c.Request = c.Request.WithContext(ctx)
		api.GetUserPost(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/carol/posts/"+encodePostID(42), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Post struct {
			Status   string `json:"status"`
			CoverUrl string `json:"cover_url"`
			User     struct {
				Username string `json:"username"`
			} `json:"user"`
		} `json:"post"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Post.Status != "draft" {
		t.Fatalf("status = %q, want draft", body.Post.Status)
	}
	if body.Post.User.Username != "carol" {
		t.Fatalf("username = %q, want carol", body.Post.User.Username)
	}
	if body.Post.CoverUrl != "/uploads/post-covers/cover.png" {
		t.Fatalf("cover url = %q, want /uploads/post-covers/cover.png", body.Post.CoverUrl)
	}
}

func TestGetUserPostHidesPostsFromOtherUsers(t *testing.T) {
	tests := []struct {
		name     string
		authUser *gocauth.AuthenticatedUser
		path     string
	}{
		{
			name:     "username mismatch",
			authUser: &gocauth.AuthenticatedUser{ID: "7", Username: "carol"},
			path:     "/api/v1/users/alice/posts/" + encodePostID(42),
		},
		{
			name:     "owner mismatch",
			authUser: &gocauth.AuthenticatedUser{ID: "8", Username: "carol"},
			path:     "/api/v1/users/carol/posts/" + encodePostID(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postClient := &fakePostClient{getPost: testProtoPost(42)}
			postClient.getPost.Status = postv1.PostStatus_POST_STATUS_DRAFT
			api := &PostsAPI{postClient: postClient}

			r := gin.New()
			r.GET("/api/v1/users/:username/posts/:id", func(c *gin.Context) {
				ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), tt.authUser)
				c.Request = c.Request.WithContext(ctx)
				api.GetUserPost(c)
			})

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestListUserPostsAllowsOwnerToListDrafts(t *testing.T) {
	draft := testProtoPost(42)
	draft.Status = postv1.PostStatus_POST_STATUS_DRAFT
	postClient := &fakePostClient{listPosts: []*postv1.Post{draft}}
	userClient := &fakePostUserClient{}
	api := &PostsAPI{postClient: postClient, userClient: userClient}

	r := gin.New()
	r.GET("/api/v1/users/:username/posts", func(c *gin.Context) {
		ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), &gocauth.AuthenticatedUser{ID: "7", Username: "carol"})
		c.Request = c.Request.WithContext(ctx)
		api.ListUserPosts(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/carol/posts?page=2&page_size=10", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if postClient.listReq == nil || postClient.listReq.GetUserId() != 7 || postClient.listReq.GetPage() != 2 || postClient.listReq.GetPageSize() != 10 {
		t.Fatalf("unexpected list request: %+v", postClient.listReq)
	}
	if postClient.listReq.GetStatus() != postv1.PostStatus_POST_STATUS_UNSPECIFIED {
		t.Fatalf("status filter = %s, want unspecified", postClient.listReq.GetStatus())
	}
	var body struct {
		Posts []struct {
			Status    string `json:"status"`
			CanManage bool   `json:"can_manage"`
		} `json:"posts"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Posts) != 1 || body.Posts[0].Status != "draft" {
		t.Fatalf("unexpected posts: %+v", body.Posts)
	}
	if !body.Posts[0].CanManage {
		t.Fatalf("can_manage = false, want true")
	}
}

func TestListUserPostsHidesOtherUsers(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.GET("/api/v1/users/:username/posts", func(c *gin.Context) {
		ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), &gocauth.AuthenticatedUser{ID: "7", Username: "carol"})
		c.Request = c.Request.WithContext(ctx)
		api.ListUserPosts(c)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/users/alice/posts", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if postClient.listReq != nil {
		t.Fatalf("ListPosts should not be called: %+v", postClient.listReq)
	}
}

func TestUpdatePostAcceptsUpdateFields(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.PUT("/api/v1/posts/:id", api.UpdatePost)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/posts/"+encodePostID(42), strings.NewReader(`{"post":{"title":"updated","summary":"short","content":"body","category_id":9,"cover_url":"post-covers/cover.png"},"update_fields":["title","summary","cover_url","category_id"]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if postClient.updateReq == nil {
		t.Fatal("update request was not sent")
	}
	if postClient.updateReq.GetId() != 42 {
		t.Fatalf("id = %d, want 42", postClient.updateReq.GetId())
	}
	if got, want := postClient.updateReq.GetUpdateMask().GetPaths(), []string{"title", "summary", "cover_url", "category_id"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("update mask = %v, want %v", got, want)
	}
}

func TestListPostsWritesOpenAPIResponse(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.GET("/api/v1/posts", api.ListPosts)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?user_id=7&category_id=9&tag=go&page=2&page_size=10", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if postClient.listReq == nil || postClient.listReq.GetUserId() != 7 || postClient.listReq.GetCategoryId() != 9 || postClient.listReq.GetTag() != "go" {
		t.Fatalf("unexpected list request: %+v", postClient.listReq)
	}
	if postClient.listReq.GetStatus() != postv1.PostStatus_POST_STATUS_PUBLISHED {
		t.Fatalf("status filter = %s, want published", postClient.listReq.GetStatus())
	}
	var body struct {
		Total int64 `json:"total"`
		Posts []struct {
			Id        string `json:"id"`
			CanManage *bool  `json:"can_manage"`
		} `json:"posts"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total != 1 || len(body.Posts) != 1 || body.Posts[0].Id != encodePostID(42) {
		t.Fatalf("unexpected response: %+v", body)
	}
	if body.Posts[0].CanManage != nil {
		t.Fatalf("public list should not include can_manage, got %v", *body.Posts[0].CanManage)
	}
}

func TestProtoUpdateMaskRejectsUnsupportedField(t *testing.T) {
	if _, err := openapiToProtoUpdateMask([]string{"user_id"}); err == nil {
		t.Fatal("expected unsupported update mask error")
	}
}

func TestProtoUpdateMaskRejectsCamelCase(t *testing.T) {
	if _, err := openapiToProtoUpdateMask([]string{"coverUrl"}); err == nil {
		t.Fatal("expected camelCase update mask error")
	}
}

func TestProtoUpdateMaskAcceptsSummaryAndPublishedAt(t *testing.T) {
	mask, err := openapiToProtoUpdateMask([]string{"summary", "published_at"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := mask.GetPaths(), []string{"summary", "published_at"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("update mask = %v, want %v", got, want)
	}
}

func testProtoPost(id int64) *postv1.Post {
	now := timestamppb.New(time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC))
	return &postv1.Post{
		Id:          id,
		UserId:      7,
		Title:       "new title",
		Summary:     "short",
		Content:     "body",
		Status:      postv1.PostStatus_POST_STATUS_PUBLISHED,
		Tags:        []string{"go"},
		CategoryId:  9,
		PublishedAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
		CoverUrl:    "post-covers/cover.png",
	}
}
