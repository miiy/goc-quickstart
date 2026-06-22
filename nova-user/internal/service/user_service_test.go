package service

import (
	"context"
	"testing"

	pb "github.com/miiy/goc-quickstart/nova-user/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-user/internal/entity"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockUserRepository implements repository.UserRepository for testing
type MockUserRepository struct {
	users  map[int64]*entity.User
	nextID int64
	err    error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[int64]*entity.User),
		nextID: 1,
	}
}

func authenticatedContext(userID int64) context.Context {
	return gocauth.InjectAuthenticatedUser(context.Background(), &gocauth.AuthenticatedUser{
		ID:       userID,
		Username: "testuser",
	})
}

func testUser(id int64) *entity.User {
	user := &entity.User{
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
		Status:   1,
	}
	user.ID = id
	return user
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.err != nil {
		return m.err
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, id int64, user *entity.User, columns ...string) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	if _, ok := m.users[id]; !ok {
		return 0, nil
	}
	user.ID = id
	user.Username = m.users[id].Username
	m.users[id] = user
	return 1, nil
}

func (m *MockUserRepository) First(ctx context.Context, id int64, columns ...string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) FindByIDs(ctx context.Context, ids []int64, columns ...string) ([]*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	users := make([]*entity.User, 0, len(ids))
	for _, id := range ids {
		if user, ok := m.users[id]; ok {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) List(ctx context.Context, page, pageSize int64, columns ...string) ([]*entity.User, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	var result []*entity.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, int64(len(result)), nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	if _, ok := m.users[id]; !ok {
		return 0, nil
	}
	delete(m.users, id)
	return 1, nil
}

func TestUserService_GetUser(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()

	// Create a test user
	repo.users[1] = testUser(1)

	service := NewUserServiceServer(logger, repo).(*UserService)

	tests := []struct {
		ctx     context.Context
		name    string
		req     *pb.GetUserRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			ctx:     authenticatedContext(1),
			name:    "get existing user",
			req:     &pb.GetUserRequest{Id: 1},
			wantErr: false,
		},
		{
			ctx:     authenticatedContext(999),
			name:    "get non-existing user",
			req:     &pb.GetUserRequest{Id: 999},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			ctx:     authenticatedContext(1),
			name:    "get another user rejected",
			req:     &pb.GetUserRequest{Id: 2},
			wantErr: true,
			errCode: codes.PermissionDenied,
		},
		{
			ctx:     context.Background(),
			name:    "get without authenticated user",
			req:     &pb.GetUserRequest{Id: 1},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			ctx:     context.Background(),
			name:    "invalid id",
			req:     &pb.GetUserRequest{Id: 0},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetUser(tt.ctx, tt.req)
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
				if resp == nil {
					t.Error("expected response")
				}
			}
		})
	}
}

func TestUserService_BatchGetUsers(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()
	repo.users[1] = testUser(1)
	repo.users[2] = &entity.User{Username: "bob", Status: 1}
	repo.users[2].ID = 2

	service := NewUserServiceServer(logger, repo).(*UserService)

	resp, err := service.BatchGetUsers(context.Background(), &pb.BatchGetUsersRequest{
		Ids: []int64{1, 2, 1},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(resp.Users))
	}
	if resp.Users[0].Id != 1 || resp.Users[0].Nickname != "Test User" {
		t.Fatalf("unexpected first user: %+v", resp.Users[0])
	}
	if resp.Users[1].Id != 2 || resp.Users[1].Username != "bob" {
		t.Fatalf("unexpected second user: %+v", resp.Users[1])
	}

	_, err = service.BatchGetUsers(context.Background(), &pb.BatchGetUsersRequest{Ids: []int64{0, -1}})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()

	// Create a test user
	repo.users[1] = testUser(1)

	service := NewUserServiceServer(logger, repo).(*UserService)

	tests := []struct {
		ctx     context.Context
		name    string
		req     *pb.UpdateUserRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			ctx:  authenticatedContext(1),
			name: "update existing user",
			req: &pb.UpdateUserRequest{
				Id: 1,
				User: &pb.User{
					Nickname: "Updated Nickname",
					Email:    "updated@example.com",
				},
			},
			wantErr: false,
		},
		{
			ctx:  authenticatedContext(999),
			name: "update non-existing user",
			req: &pb.UpdateUserRequest{
				Id: 999,
				User: &pb.User{
					Nickname: "Nickname",
				},
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			ctx:  authenticatedContext(2),
			name: "update another user rejected",
			req: &pb.UpdateUserRequest{
				Id: 1,
				User: &pb.User{
					Nickname: "Updated Nickname",
				},
			},
			wantErr: true,
			errCode: codes.PermissionDenied,
		},
		{
			ctx:  context.Background(),
			name: "update without authenticated user",
			req: &pb.UpdateUserRequest{
				Id: 1,
				User: &pb.User{
					Nickname: "Updated Nickname",
				},
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			ctx:     authenticatedContext(1),
			name:    "update with nil user",
			req:     &pb.UpdateUserRequest{Id: 1, User: nil},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.UpdateUser(tt.ctx, tt.req)
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
				if resp == nil {
					t.Error("expected response")
				}
			}
		})
	}
}

func TestUserService_UpdateUserNormalizesAvatarObjectKey(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()
	repo.users[1] = testUser(1)

	service := NewUserServiceServer(logger, repo).(*UserService)

	resp, err := service.UpdateUser(authenticatedContext(1), &pb.UpdateUserRequest{
		Id: 1,
		User: &pb.User{
			Avatar: "http://127.0.0.1:8080/uploads/avatars/2026/06/c783e33713eb8123da7de5bc800f5f9a.png",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "avatars/2026/06/c783e33713eb8123da7de5bc800f5f9a.png"
	if resp.GetUser().GetAvatar() != want {
		t.Fatalf("avatar = %q, want %q", resp.GetUser().GetAvatar(), want)
	}
	if repo.users[1].Avatar != want {
		t.Fatalf("stored avatar = %q, want %q", repo.users[1].Avatar, want)
	}
}

func TestUserService_ListUsers(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()

	// Create test users
	for i := 1; i <= 5; i++ {
		repo.users[int64(i)] = testUser(int64(i))
	}

	service := NewUserServiceServer(logger, repo).(*UserService)

	_, err := service.ListUsers(authenticatedContext(1), &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("expected permission denied, got %v", err)
	}

	_, err = service.ListUsers(context.Background(), &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected unauthenticated, got %v", err)
	}
}
