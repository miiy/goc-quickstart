package service

import (
	"context"
	"database/sql"
	"strconv"
	"testing"
	"time"

	pb "github.com/miiy/goc-quickstart/nova-post/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc-quickstart/nova-post/internal/repository"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// MockPostRepository implements repository.PostRepository for testing
type MockPostRepository struct {
	posts      map[int64]*entity.Post
	categories []*entity.Category
	nextID     int64
	err        error
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int64]*entity.Post),
		nextID: 1,
	}
}

func authenticatedContext(userID int64) context.Context {
	return gocauth.InjectAuthenticatedUser(context.Background(), &gocauth.AuthenticatedUser{
		ID:       strconv.FormatInt(userID, 10),
		Username: "alice",
	})
}

func newTestPostService(logger *zap.Logger, repo *MockPostRepository) *PostService {
	return NewPostServiceServer(logger, repo, repo).(*PostService)
}

func (m *MockPostRepository) Create(ctx context.Context, post *entity.Post) error {
	if m.err != nil {
		return m.err
	}
	post.ID = m.nextID
	m.nextID++
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Update(ctx context.Context, id int64, post *entity.Post, columns ...string) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	existing, ok := m.posts[id]
	if !ok {
		return 0, nil
	}
	if len(columns) == 0 {
		post.ID = id
		m.posts[id] = post
		return 1, nil
	}
	for _, column := range columns {
		switch column {
		case "title":
			existing.Title = post.Title
		case "summary":
			existing.Summary = post.Summary
		case "cover_url":
			existing.CoverUrl = post.CoverUrl
		case "content":
			existing.Content = post.Content
		case "status":
			existing.Status = post.Status
		case "tags":
			existing.Tags = post.Tags
		case "category_id":
			existing.CategoryId = post.CategoryId
		case "published_at":
			existing.PublishedAt = post.PublishedAt
		}
	}
	return 1, nil
}

func (m *MockPostRepository) First(ctx context.Context, id int64, columns ...string) (*entity.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	if post, ok := m.posts[id]; ok {
		return post, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockPostRepository) Delete(ctx context.Context, id int64) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	if _, ok := m.posts[id]; !ok {
		return 0, nil
	}
	delete(m.posts, id)
	return 1, nil
}

func (m *MockPostRepository) List(ctx context.Context, filter *repository.ListFilter, page, pageSize int64, columns ...string) ([]*entity.Post, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	var result []*entity.Post
	for _, p := range m.posts {
		if filter.Status > 0 && p.Status != filter.Status {
			continue
		}
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func TestPostService_GetPost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create a test post
	repo.posts[1] = &entity.Post{
		UserId:     1,
		Title:      "Test Post",
		Content:    "Test Content",
		Status:     1,
		CategoryId: 1,
	}

	service := newTestPostService(logger, repo)

	tests := []struct {
		name    string
		req     *pb.GetPostRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "get existing post",
			req:     &pb.GetPostRequest{Id: 1},
			wantErr: false,
		},
		{
			name:    "get non-existing post",
			req:     &pb.GetPostRequest{Id: 999},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name:    "invalid id",
			req:     &pb.GetPostRequest{Id: 0},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetPost(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil || resp.Post == nil {
					t.Error("expected response with post")
				}
			}
		})
	}
}

func TestPostService_CreatePost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	service := newTestPostService(logger, repo)

	tests := []struct {
		name       string
		ctx        context.Context
		req        *pb.CreatePostRequest
		wantErr    bool
		errCode    codes.Code
		wantUserID int64
	}{
		{
			name: "create valid post",
			ctx:  authenticatedContext(42),
			req: &pb.CreatePostRequest{
				Post: &pb.Post{
					UserId:     999,
					Title:      "New Post",
					Summary:    "Short summary",
					CoverUrl:   "http://cdn.test/cover.png",
					Content:    "New Content",
					Status:     pb.PostStatus_POST_STATUS_DRAFT,
					CategoryId: 1,
					Tags:       []string{"tag1", "tag2"},
				},
			},
			wantErr:    false,
			wantUserID: 42,
		},
		{
			name: "create post with empty title",
			req: &pb.CreatePostRequest{
				Post: &pb.Post{
					UserId:  1,
					Title:   "",
					Content: "Content",
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name:    "create post with nil post",
			req:     &pb.CreatePostRequest{Post: nil},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "create post without authenticated user",
			ctx:  context.Background(),
			req: &pb.CreatePostRequest{
				Post: &pb.Post{
					Title:   "New Post",
					Content: "New Content",
				},
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}
			resp, err := service.CreatePost(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil || resp.Post == nil {
					t.Error("expected response with post")
				}
				if resp.Post.Title != tt.req.Post.Title {
					t.Errorf("expected title %s, got %s", tt.req.Post.Title, resp.Post.Title)
				}
				if resp.Post.UserId != tt.wantUserID {
					t.Errorf("expected user id %d, got %d", tt.wantUserID, resp.Post.UserId)
				}
				if resp.Post.CoverUrl != tt.req.Post.CoverUrl {
					t.Errorf("expected cover url %s, got %s", tt.req.Post.CoverUrl, resp.Post.CoverUrl)
				}
				if resp.Post.Summary != tt.req.Post.Summary {
					t.Errorf("expected summary %s, got %s", tt.req.Post.Summary, resp.Post.Summary)
				}
			}
		})
	}
}

func TestPostService_CreatePublishedPostDefaultsPublishedAt(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	service := newTestPostService(logger, repo)

	resp, err := service.CreatePost(authenticatedContext(42), &pb.CreatePostRequest{
		Post: &pb.Post{
			Title:   "Published Post",
			Content: "Body",
			Status:  pb.PostStatus_POST_STATUS_PUBLISHED,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetPost().GetPublishedAt() == nil {
		t.Fatal("expected published_at to be defaulted")
	}
}

func TestPostService_CreatePostNormalizesUploadsCoverURL(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	service := newTestPostService(logger, repo)

	resp, err := service.CreatePost(authenticatedContext(42), &pb.CreatePostRequest{
		Post: &pb.Post{
			Title:    "Post With Cover",
			Content:  "Body",
			Status:   pb.PostStatus_POST_STATUS_DRAFT,
			CoverUrl: "/uploads/post-covers/2026/07/cover.png",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := resp.GetPost().GetCoverUrl(), "post-covers/2026/07/cover.png"; got != want {
		t.Fatalf("cover url = %q, want %q", got, want)
	}
	if got, want := repo.posts[resp.GetPost().GetId()].CoverUrl, "post-covers/2026/07/cover.png"; got != want {
		t.Fatalf("stored cover url = %q, want %q", got, want)
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create a test post
	repo.posts[1] = &entity.Post{
		UserId:     1,
		Title:      "Test Post",
		CoverUrl:   "http://cdn.test/old.png",
		Content:    "Test Content",
		Status:     1,
		CategoryId: 1,
	}

	service := newTestPostService(logger, repo)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *pb.UpdatePostRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "update existing post",
			ctx:  authenticatedContext(1),
			req: &pb.UpdatePostRequest{
				Id: 1,
				Post: &pb.Post{
					UserId:   999,
					Title:    "Updated Title",
					Summary:  "Updated Summary",
					CoverUrl: "http://cdn.test/new.png",
					Content:  "Updated Content",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"title", "summary", "content", "cover_url"}},
			},
			wantErr: false,
		},
		{
			name: "update non-existing post",
			ctx:  authenticatedContext(1),
			req: &pb.UpdatePostRequest{
				Id: 999,
				Post: &pb.Post{
					Title: "Title",
				},
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name:    "update with nil post",
			req:     &pb.UpdatePostRequest{Id: 1, Post: nil},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "update without authenticated user",
			ctx:  context.Background(),
			req: &pb.UpdatePostRequest{
				Id: 1,
				Post: &pb.Post{
					Title: "Updated Title",
				},
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name: "update by non-user",
			ctx:  authenticatedContext(2),
			req: &pb.UpdatePostRequest{
				Id: 1,
				Post: &pb.Post{
					Title: "Updated Title",
				},
			},
			wantErr: true,
			errCode: codes.PermissionDenied,
		},
		{
			name: "update user id is rejected",
			ctx:  authenticatedContext(1),
			req: &pb.UpdatePostRequest{
				Id: 1,
				Post: &pb.Post{
					UserId: 2,
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"user_id"}},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}
			resp, err := service.UpdatePost(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil || resp.Post == nil {
					t.Error("expected response with post")
				}
				if resp.Post.UserId != 1 {
					t.Errorf("expected user id to remain 1, got %d", resp.Post.UserId)
				}
				if resp.Post.CoverUrl != "http://cdn.test/new.png" {
					t.Errorf("expected cover url to update, got %s", resp.Post.CoverUrl)
				}
				if resp.Post.Summary != "Updated Summary" {
					t.Errorf("expected summary to update, got %s", resp.Post.Summary)
				}
				if resp.Post.Status != pb.PostStatus_POST_STATUS_PENDING_REVIEW {
					t.Errorf("expected status to become pending review, got %s", resp.Post.Status)
				}
			}
		})
	}
}

func TestPostService_UpdatePostNormalizesUploadsCoverURL(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	repo.posts[1] = &entity.Post{
		UserId:   1,
		Title:    "Draft",
		CoverUrl: "post-covers/old.png",
		Content:  "Body",
		Status:   entity.PostStatusDraft,
	}
	service := newTestPostService(logger, repo)

	resp, err := service.UpdatePost(authenticatedContext(1), &pb.UpdatePostRequest{
		Id: 1,
		Post: &pb.Post{
			CoverUrl: "http://127.0.0.1:8081/uploads/post-covers/2026/07/cover.png",
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"cover_url"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := resp.GetPost().GetCoverUrl(), "post-covers/2026/07/cover.png"; got != want {
		t.Fatalf("cover url = %q, want %q", got, want)
	}
	if got, want := repo.posts[1].CoverUrl, "post-covers/2026/07/cover.png"; got != want {
		t.Fatalf("stored cover url = %q, want %q", got, want)
	}
}

func TestPostService_UpdatePostForcesPendingReview(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	repo.posts[1] = &entity.Post{
		UserId:  1,
		Title:   "Draft",
		Content: "Body",
		Status:  entity.PostStatusDraft,
	}
	service := newTestPostService(logger, repo)

	resp, err := service.UpdatePost(authenticatedContext(1), &pb.UpdatePostRequest{
		Id: 1,
		Post: &pb.Post{
			Status: pb.PostStatus_POST_STATUS_PUBLISHED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetPost().GetStatus() != pb.PostStatus_POST_STATUS_PENDING_REVIEW {
		t.Fatalf("status = %s, want pending review", resp.GetPost().GetStatus())
	}
	if resp.GetPost().GetPublishedAt() != nil {
		t.Fatal("expected pending review update not to default published_at")
	}
}

func TestPostService_DeletePost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create a test post
	repo.posts[1] = &entity.Post{
		UserId:  1,
		Title:   "Test Post",
		Content: "Test Content",
	}

	service := newTestPostService(logger, repo)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *pb.DeletePostRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "delete without authenticated user",
			ctx:     context.Background(),
			req:     &pb.DeletePostRequest{Id: 1},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:    "delete by non-user",
			ctx:     authenticatedContext(2),
			req:     &pb.DeletePostRequest{Id: 1},
			wantErr: true,
			errCode: codes.PermissionDenied,
		},
		{
			name:    "delete existing post",
			ctx:     authenticatedContext(1),
			req:     &pb.DeletePostRequest{Id: 1},
			wantErr: false,
		},
		{
			name:    "delete non-existing post",
			ctx:     authenticatedContext(1),
			req:     &pb.DeletePostRequest{Id: 999},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name:    "delete with invalid id",
			req:     &pb.DeletePostRequest{Id: 0},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}
			_, err := service.DeletePost(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPostService_ListPosts(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create test posts
	for i := 1; i <= 3; i++ {
		repo.posts[int64(i)] = &entity.Post{
			UserId:     1,
			Title:      "Test Post",
			Content:    "Test Content",
			Status:     1,
			CategoryId: 1,
		}
	}

	service := newTestPostService(logger, repo)

	resp, err := service.ListPosts(context.Background(), &pb.ListPostsRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected response")
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}

	if len(resp.Posts) != 3 {
		t.Errorf("expected 3 posts, got %d", len(resp.Posts))
	}
}

func TestPostService_ListPostsFiltersStatus(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	repo.posts[1] = &entity.Post{UserId: 1, Title: "Draft", Status: entity.PostStatusDraft}
	repo.posts[2] = &entity.Post{UserId: 1, Title: "Published", Status: entity.PostStatusPublished, PublishedAt: sql.NullTime{Time: time.Now(), Valid: true}}
	service := newTestPostService(logger, repo)

	resp, err := service.ListPosts(context.Background(), &pb.ListPostsRequest{
		Status: pb.PostStatus_POST_STATUS_PUBLISHED,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetTotal() != 1 || len(resp.GetPosts()) != 1 || resp.GetPosts()[0].GetTitle() != "Published" {
		t.Fatalf("unexpected published list: %+v", resp)
	}
}
