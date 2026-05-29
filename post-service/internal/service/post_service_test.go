package service

import (
	"context"
	"testing"

	pb "github.com/miiy/goc-quickstart/post-service/gen/go/blog/post/v1"
	"github.com/miiy/goc-quickstart/post-service/internal/entity"
	"github.com/miiy/goc-quickstart/post-service/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// MockPostRepository implements repository.PostRepository for testing
type MockPostRepository struct {
	posts  map[int64]*entity.Post
	nextID int64
	err    error
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int64]*entity.Post),
		nextID: 1,
	}
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
	if _, ok := m.posts[id]; !ok {
		return 0, nil
	}
	post.ID = id
	m.posts[id] = post
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
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func TestPostService_GetPost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create a test post
	repo.posts[1] = &entity.Post{
		AuthorId:   1,
		Title:      "Test Post",
		Content:    "Test Content",
		Status:     1,
		CategoryId: 1,
	}

	service := NewPostServiceServer(logger, repo).(*PostService)

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
	service := NewPostServiceServer(logger, repo).(*PostService)

	tests := []struct {
		name    string
		req     *pb.CreatePostRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "create valid post",
			req: &pb.CreatePostRequest{
				Post: &pb.Post{
					AuthorId:   1,
					Title:      "New Post",
					Content:    "New Content",
					Status:     pb.PostStatus_POST_STATUS_DRAFT,
					CategoryId: 1,
					Tags:       []string{"tag1", "tag2"},
				},
			},
			wantErr: false,
		},
		{
			name: "create post with empty title",
			req: &pb.CreatePostRequest{
				Post: &pb.Post{
					AuthorId: 1,
					Title:    "",
					Content:  "Content",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.CreatePost(context.Background(), tt.req)
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
			}
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create a test post
	repo.posts[1] = &entity.Post{
		AuthorId:   1,
		Title:      "Test Post",
		Content:    "Test Content",
		Status:     1,
		CategoryId: 1,
	}

	service := NewPostServiceServer(logger, repo).(*PostService)

	tests := []struct {
		name    string
		req     *pb.UpdatePostRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "update existing post",
			req: &pb.UpdatePostRequest{
				Id: 1,
				Post: &pb.Post{
					Title:   "Updated Title",
					Content: "Updated Content",
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existing post",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.UpdatePost(context.Background(), tt.req)
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

func TestPostService_DeletePost(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()

	// Create a test post
	repo.posts[1] = &entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "Test Content",
	}

	service := NewPostServiceServer(logger, repo).(*PostService)

	tests := []struct {
		name    string
		req     *pb.DeletePostRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "delete existing post",
			req:     &pb.DeletePostRequest{Id: 1},
			wantErr: false,
		},
		{
			name:    "delete non-existing post",
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
			_, err := service.DeletePost(context.Background(), tt.req)
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
			AuthorId:   1,
			Title:      "Test Post",
			Content:    "Test Content",
			Status:     1,
			CategoryId: 1,
		}
	}

	service := NewPostServiceServer(logger, repo).(*PostService)

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