package service

import (
	"context"
	"testing"

	pb "github.com/miiy/goc-quickstart/user-service/gen/go/blog/user/v1"
	"github.com/miiy/goc-quickstart/user-service/internal/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
	repo.users[1] = &entity.User{
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
		Status:   1,
	}

	service := NewUserServiceServer(logger, repo).(*UserService)

	tests := []struct {
		name    string
		req     *pb.GetUserRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "get existing user",
			req:     &pb.GetUserRequest{Id: 1},
			wantErr: false,
		},
		{
			name:    "get non-existing user",
			req:     &pb.GetUserRequest{Id: 999},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name:    "invalid id",
			req:     &pb.GetUserRequest{Id: 0},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetUser(context.Background(), tt.req)
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

func TestUserService_UpdateUser(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()

	// Create a test user
	repo.users[1] = &entity.User{
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
		Status:   1,
	}

	service := NewUserServiceServer(logger, repo).(*UserService)

	tests := []struct {
		name    string
		req     *pb.UpdateUserRequest
		wantErr bool
		errCode codes.Code
	}{
		{
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
			name:    "update with nil user",
			req:     &pb.UpdateUserRequest{Id: 1, User: nil},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.UpdateUser(context.Background(), tt.req)
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

func TestUserService_ListUsers(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockUserRepository()

	// Create test users
	for i := 1; i <= 5; i++ {
		repo.users[int64(i)] = &entity.User{
			Username: "testuser",
			Nickname: "Test User",
			Email:    "test@example.com",
			Status:   1,
		}
	}

	service := NewUserServiceServer(logger, repo).(*UserService)

	resp, err := service.ListUsers(context.Background(), &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected response")
	}

	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}

	if len(resp.Users) != 5 {
		t.Errorf("expected 5 users, got %d", len(resp.Users))
	}
}