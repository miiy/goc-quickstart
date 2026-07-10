package post

import (
	"context"
	htmltemplate "html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	postv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/user/v1"
	webtemplate "github.com/miiy/goc-quickstart/nova-web/internal/template"
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/gin/sessions"
	pkgTemplate "github.com/miiy/goc/gin/template"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListRendersPostsOnServer(t *testing.T) {
	postClient := &fakePostReadClient{
		listResp: &postv1.ListPostsResponse{
			Total:       1,
			TotalPages:  1,
			PageSize:    25,
			CurrentPage: 2,
			Posts: []*postv1.Post{
				{
					Id:        101,
					UserId:    7,
					Title:     "Server Rendered List Post",
					Content:   "list body",
					Status:    postv1.PostStatus_POST_STATUS_PUBLISHED,
					CreatedAt: timestamppb.New(time.Date(2026, 1, 2, 3, 4, 0, 0, time.UTC)),
				},
			},
		},
	}
	userClient := &fakePostUserClient{
		resp: &userv1.BatchGetUsersResponse{
			Users: []*userv1.User{{Id: 7, Username: "alice", Nickname: "Alice"}},
		},
	}
	r := newPostHandlerTestRouter(t, &PostsHandler{postClient: postClient, userClient: userClient}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/posts?page=2&page_size=25", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("GET /posts status = %d, want %d: %s", w.Code, http.StatusOK, w.Body.String())
	}
	if len(postClient.listReqs) != 1 {
		t.Fatalf("ListPosts calls = %d, want 1", len(postClient.listReqs))
	}
	req := postClient.listReqs[0]
	if req.GetPage() != 2 || req.GetPageSize() != 25 || req.GetStatus() != postv1.PostStatus_POST_STATUS_PUBLISHED {
		t.Fatalf("ListPosts request = page %d size %d status %s", req.GetPage(), req.GetPageSize(), req.GetStatus())
	}
	if len(userClient.reqs) != 1 || len(userClient.reqs[0].GetIds()) != 1 || userClient.reqs[0].GetIds()[0] != 7 {
		t.Fatalf("BatchGetUsers requests = %#v, want user id 7", userClient.reqs)
	}

	body := w.Body.String()
	for _, want := range []string{"Server Rendered List Post", "Alice", "2026-01-02 11:04"} {
		if !strings.Contains(body, want) {
			t.Fatalf("SSR list body missing %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "data-post-list-page") {
		t.Fatal("list response must not mount a client-rendered page")
	}
}

func TestListPagesRouteUsesListHandler(t *testing.T) {
	postClient := &fakePostReadClient{
		listResp: &postv1.ListPostsResponse{
			Total:       0,
			TotalPages:  0,
			PageSize:    10,
			CurrentPage: 3,
		},
	}
	r := newPostHandlerTestRouter(t, &PostsHandler{postClient: postClient}, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/posts/pages/3", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("GET /posts/pages/3 status = %d, want %d: %s", w.Code, http.StatusOK, w.Body.String())
	}
	if len(postClient.listReqs) != 1 {
		t.Fatalf("ListPosts calls = %d, want 1", len(postClient.listReqs))
	}
	req := postClient.listReqs[0]
	if req.GetPage() != 3 || req.GetPageSize() != 10 || req.GetStatus() != postv1.PostStatus_POST_STATUS_PUBLISHED {
		t.Fatalf("ListPosts request = page %d size %d status %s", req.GetPage(), req.GetPageSize(), req.GetStatus())
	}
}

func TestShowRendersPostOnServer(t *testing.T) {
	postClient := &fakePostReadClient{
		getResp: &postv1.GetPostResponse{
			Post: &postv1.Post{
				Id:        202,
				UserId:    42,
				Title:     "Server Rendered Detail Post",
				Content:   "server rendered detail body",
				Status:    postv1.PostStatus_POST_STATUS_PUBLISHED,
				CreatedAt: timestamppb.New(time.Date(2026, 1, 2, 3, 4, 0, 0, time.UTC)),
			},
		},
	}
	userClient := &fakePostUserClient{
		resp: &userv1.BatchGetUsersResponse{
			Users: []*userv1.User{{Id: 42, Username: "bob", Nickname: "Bob"}},
		},
	}
	r := newPostHandlerTestRouter(t, &PostsHandler{postClient: postClient, userClient: userClient}, &gocauth.AuthenticatedUser{ID: "42", Username: "bob"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/posts/"+encodePostID(202), nil))

	if w.Code != http.StatusOK {
		t.Fatalf("GET /posts/:id status = %d, want %d: %s", w.Code, http.StatusOK, w.Body.String())
	}
	if len(postClient.getReqs) != 1 || postClient.getReqs[0].GetId() != 202 {
		t.Fatalf("GetPost requests = %#v, want id 202", postClient.getReqs)
	}
	if len(userClient.reqs) != 1 || len(userClient.reqs[0].GetIds()) != 1 || userClient.reqs[0].GetIds()[0] != 42 {
		t.Fatalf("BatchGetUsers requests = %#v, want user id 42", userClient.reqs)
	}

	body := w.Body.String()
	for _, want := range []string{"Server Rendered Detail Post", "server rendered detail body", "Bob", "2026-01-02 11:04"} {
		if !strings.Contains(body, want) {
			t.Fatalf("SSR detail body missing %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "data-post-detail-page") {
		t.Fatal("detail response must not mount a client-rendered detail page")
	}
	if !strings.Contains(body, `data-post-actions data-post-id="`+encodePostID(202)+`"`) {
		t.Fatal("detail response should reserve React for owner actions only")
	}
}

func newPostHandlerTestRouter(t *testing.T, h *PostsHandler, user *gocauth.AuthenticatedUser) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	webtemplate.SetDefaultSite(webtemplate.SiteData{Name: "nova-web", Locale: "zh-CN"})

	renderer := pkgTemplate.NewTemplate()
	renderer.AddFunc("formatTime", webtemplate.NewFormatTimeFunc("Asia/Shanghai"))
	renderer.AddFunc("viteClient", func() htmltemplate.HTML { return "" })
	renderer.AddFunc("viteEntry", func(string) (htmltemplate.HTML, error) { return "", nil })
	renderer.AddFunc("alertType", webtemplate.FlashLevelClass)
	renderer.AddTemplate(resourceTemplate.FS, Templates())

	r := gin.New()
	r.HTMLRender = renderer.Render
	r.Use(sessions.Middleware("sess", sessions.NewCookieStore("test-secret")))
	if user != nil {
		r.Use(func(c *gin.Context) {
			authctx.SetUser(c, user)
			c.Next()
		})
	}
	r.GET("/posts", h.list)
	r.GET("/posts/pages/:page", h.list)
	r.GET("/posts/:id", h.show)
	return r
}

type fakePostReadClient struct {
	listResp *postv1.ListPostsResponse
	listErr  error
	listReqs []*postv1.ListPostsRequest
	getResp  *postv1.GetPostResponse
	getErr   error
	getReqs  []*postv1.GetPostRequest
}

func (f *fakePostReadClient) ListPosts(_ context.Context, in *postv1.ListPostsRequest, _ ...grpc.CallOption) (*postv1.ListPostsResponse, error) {
	req := *in
	f.listReqs = append(f.listReqs, &req)
	return f.listResp, f.listErr
}

func (f *fakePostReadClient) GetPost(_ context.Context, in *postv1.GetPostRequest, _ ...grpc.CallOption) (*postv1.GetPostResponse, error) {
	req := *in
	f.getReqs = append(f.getReqs, &req)
	return f.getResp, f.getErr
}

func (f *fakePostReadClient) CreatePost(context.Context, *postv1.CreatePostRequest, ...grpc.CallOption) (*postv1.CreatePostResponse, error) {
	panic("unexpected CreatePost call")
}

func (f *fakePostReadClient) UpdatePost(context.Context, *postv1.UpdatePostRequest, ...grpc.CallOption) (*postv1.UpdatePostResponse, error) {
	panic("unexpected UpdatePost call")
}

func (f *fakePostReadClient) DeletePost(context.Context, *postv1.DeletePostRequest, ...grpc.CallOption) (*postv1.DeletePostResponse, error) {
	panic("unexpected DeletePost call")
}

func (f *fakePostReadClient) ListCategories(context.Context, *postv1.ListCategoriesRequest, ...grpc.CallOption) (*postv1.ListCategoriesResponse, error) {
	panic("unexpected ListCategories call")
}

type fakePostUserClient struct {
	resp *userv1.BatchGetUsersResponse
	err  error
	reqs []*userv1.BatchGetUsersRequest
}

func (f *fakePostUserClient) BatchGetUsers(_ context.Context, in *userv1.BatchGetUsersRequest, _ ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error) {
	req := &userv1.BatchGetUsersRequest{Ids: append([]int64(nil), in.GetIds()...)}
	f.reqs = append(f.reqs, req)
	return f.resp, f.err
}

func (f *fakePostUserClient) GetUser(context.Context, *userv1.GetUserRequest, ...grpc.CallOption) (*userv1.GetUserResponse, error) {
	panic("unexpected GetUser call")
}

func (f *fakePostUserClient) GetUserByUsername(context.Context, *userv1.GetUserByUsernameRequest, ...grpc.CallOption) (*userv1.GetUserByUsernameResponse, error) {
	panic("unexpected GetUserByUsername call")
}

func (f *fakePostUserClient) UpdateUser(context.Context, *userv1.UpdateUserRequest, ...grpc.CallOption) (*userv1.UpdateUserResponse, error) {
	panic("unexpected UpdateUser call")
}

func (f *fakePostUserClient) ListUsers(context.Context, *userv1.ListUsersRequest, ...grpc.CallOption) (*userv1.ListUsersResponse, error) {
	panic("unexpected ListUsers call")
}
