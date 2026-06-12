package service

import (
	"context"
	"testing"

	pb "github.com/miiy/goc-quickstart/nova-post/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc-quickstart/nova-post/internal/repository"
	gocauth "github.com/miiy/goc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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

func authenticatedContext(userID int64) context.Context {
	return gocauth.InjectAuthenticatedUser(context.Background(), &gocauth.AuthenticatedUser{
		ID:       userID,
		Username: "alice",
	})
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
		name         string
		ctx          context.Context
		req          *pb.CreatePostRequest
		wantErr      bool
		errCode      codes.Code
		wantAuthorID int64
	}{
		{
			name: "create valid post",
			ctx:  authenticatedContext(42),
			req: &pb.CreatePostRequest{
				Post: &pb.Post{
					AuthorId:   999,
					Title:      "New Post",
					Content:    "New Content",
					Status:     pb.PostStatus_POST_STATUS_DRAFT,
					CategoryId: 1,
					Tags:       []string{"tag1", "tag2"},
				},
			},
			wantErr:      false,
			wantAuthorID: 42,
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
				if resp.Post.AuthorId != tt.wantAuthorID {
					t.Errorf("expected author id %d, got %d", tt.wantAuthorID, resp.Post.AuthorId)
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
					AuthorId: 999,
					Title:    "Updated Title",
					Content:  "Updated Content",
				},
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
			name: "update by non-author",
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
			name: "update author id is rejected",
			ctx:  authenticatedContext(1),
			req: &pb.UpdatePostRequest{
				Id: 1,
				Post: &pb.Post{
					AuthorId: 2,
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"author_id"}},
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
				if resp.Post.AuthorId != 1 {
					t.Errorf("expected author id to remain 1, got %d", resp.Post.AuthorId)
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
			name:    "delete by non-author",
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
