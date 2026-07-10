package service

import (
	"context"
	"testing"

	pb "github.com/miiy/goc-quickstart/nova-post/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
)

func (m *MockPostRepository) ListCategories(ctx context.Context) ([]*entity.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func TestPostService_ListCategories(t *testing.T) {
	logger := zap.NewNop()
	repo := NewMockPostRepository()
	repo.categories = []*entity.Category{
		{Model: gorm.Model{ID: 1}, Name: "Engineering", Path: "/engineering"},
	}
	service := newTestPostService(logger, repo)

	resp, err := service.ListCategories(context.Background(), &pb.ListCategoriesRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.GetCategories()) != 1 || resp.GetCategories()[0].GetName() != "Engineering" {
		t.Fatalf("unexpected categories: %+v", resp.GetCategories())
	}
}
