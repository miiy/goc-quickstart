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
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeAuthorUserClient struct {
	gotIDs []int64
}

func (f *fakeAuthorUserClient) BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error) {
	f.gotIDs = append([]int64(nil), in.GetIds()...)
	return &userv1.BatchGetUsersResponse{
		Users: []*userv1.User{
			{Id: 1, Username: "alice", Nickname: "Alice Nick"},
			{Id: 2, Username: "bob"},
			{Id: 7, Username: "carol", Nickname: "Carol Nick"},
		},
	}, nil
}

func TestEnrichPostAuthorsBatchGetsUsers(t *testing.T) {
	userClient := &fakeAuthorUserClient{}
	api := &PostsAPI{userClient: userClient}
	posts := []*postv1.Post{
		{Id: 1, AuthorId: 1},
		{Id: 2, AuthorId: 2},
		{Id: 3, AuthorId: 1},
		{Id: 4},
		nil,
	}

	names, err := api.authorNames(context.Background(), posts...)
	if err != nil {
		t.Fatal(err)
	}

	if want := []int64{1, 2}; !reflect.DeepEqual(userClient.gotIDs, want) {
		t.Fatalf("batch ids = %v, want %v", userClient.gotIDs, want)
	}
	if names[1] != "Alice Nick" {
		t.Fatalf("author name = %q, want Alice Nick", names[1])
	}
	if names[2] != "bob" {
		t.Fatalf("author name = %q, want bob", names[2])
	}
}

type fakePostClient struct {
	createReq *postv1.CreatePostRequest
	updateReq *postv1.UpdatePostRequest
	listReq   *postv1.ListPostsRequest
}

func (f *fakePostClient) GetPost(ctx context.Context, in *postv1.GetPostRequest, opts ...grpc.CallOption) (*postv1.GetPostResponse, error) {
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
	return &postv1.ListPostsResponse{
		Total:       1,
		TotalPages:  1,
		PageSize:    in.GetPageSize(),
		CurrentPage: in.GetPage(),
		Posts:       []*postv1.Post{testProtoPost(42)},
	}, nil
}

func TestCreatePostUsesOpenAPIRequestAndResponse(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.POST("/api/v1/posts", api.CreatePost)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(`{"post":{"title":"new title","content":"body","status":"published","tags":["go"],"category_id":9,"cover_url":"post-covers/cover.png"}}`))
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
	if gotPost.GetAuthorId() != 0 {
		t.Fatalf("author id = %d, want 0", gotPost.GetAuthorId())
	}
	if gotPost.GetTitle() != "new title" || gotPost.GetContent() != "body" || gotPost.GetCategoryId() != 9 || gotPost.GetCoverUrl() != "post-covers/cover.png" {
		t.Fatalf("unexpected proto post: %+v", gotPost)
	}
	if gotPost.GetStatus() != postv1.PostStatus_POST_STATUS_PUBLISHED {
		t.Fatalf("status = %s", gotPost.GetStatus())
	}

	var body struct {
		Post struct {
			Id         int64    `json:"id"`
			AuthorId   int64    `json:"author_id"`
			AuthorName string   `json:"author_name"`
			CategoryId int64    `json:"category_id"`
			Status     string   `json:"status"`
			Tags       []string `json:"tags"`
		} `json:"post"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Post.Id != 42 || body.Post.AuthorId != 7 || body.Post.CategoryId != 9 {
		t.Fatalf("unexpected ids: %+v", body.Post)
	}
	if body.Post.AuthorName != "" {
		t.Fatalf("author name = %q, want empty without user client", body.Post.AuthorName)
	}
	if body.Post.Status != "published" {
		t.Fatalf("status = %q", body.Post.Status)
	}
	if !reflect.DeepEqual(body.Post.Tags, []string{"go"}) {
		t.Fatalf("tags = %v", body.Post.Tags)
	}
}

func TestGetPostAssemblesAuthorNameInGateway(t *testing.T) {
	postClient := &fakePostClient{}
	userClient := &fakeAuthorUserClient{}
	api := &PostsAPI{postClient: postClient, userClient: userClient}

	r := gin.New()
	r.GET("/api/v1/posts/:id", api.GetPost)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/42", nil)
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
			AuthorName string `json:"author_name"`
		} `json:"post"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Post.AuthorName != "Carol Nick" {
		t.Fatalf("author name = %q, want Carol Nick", body.Post.AuthorName)
	}
}

func TestUpdatePostAcceptsUpdateFields(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.PUT("/api/v1/posts/:id", api.UpdatePost)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/posts/42", strings.NewReader(`{"post":{"title":"updated","content":"body","category_id":9,"cover_url":"post-covers/cover.png"},"update_fields":["title","cover_url","category_id"]}`))
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
	if got, want := postClient.updateReq.GetUpdateMask().GetPaths(), []string{"title", "cover_url", "category_id"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("update mask = %v, want %v", got, want)
	}
}

func TestListPostsWritesOpenAPIResponse(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.GET("/api/v1/posts", api.ListPosts)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?author_id=7&category_id=9&tag=go&page=2&page_size=10", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if postClient.listReq == nil || postClient.listReq.GetAuthorId() != 7 || postClient.listReq.GetCategoryId() != 9 || postClient.listReq.GetTag() != "go" {
		t.Fatalf("unexpected list request: %+v", postClient.listReq)
	}
	var body struct {
		Total int64 `json:"total"`
		Posts []struct {
			Id int64 `json:"id"`
		} `json:"posts"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total != 1 || len(body.Posts) != 1 || body.Posts[0].Id != 42 {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestProtoUpdateMaskRejectsUnsupportedField(t *testing.T) {
	if _, err := protoUpdateMask([]string{"author_id"}); err == nil {
		t.Fatal("expected unsupported update mask error")
	}
}

func TestProtoUpdateMaskRejectsCamelCase(t *testing.T) {
	if _, err := protoUpdateMask([]string{"coverUrl"}); err == nil {
		t.Fatal("expected camelCase update mask error")
	}
}

func testProtoPost(id int64) *postv1.Post {
	now := timestamppb.New(time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC))
	return &postv1.Post{
		Id:         id,
		AuthorId:   7,
		Title:      "new title",
		Content:    "body",
		Status:     postv1.PostStatus_POST_STATUS_PUBLISHED,
		Tags:       []string{"go"},
		CategoryId: 9,
		CreatedAt:  now,
		UpdatedAt:  now,
		CoverUrl:   "post-covers/cover.png",
	}
}
