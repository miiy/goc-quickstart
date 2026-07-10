package service

import (
	"context"

	pb "github.com/miiy/goc-quickstart/nova-post/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc/logger/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListCategories exposes read-only category records without adding category writes.
func (s *PostService) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	categories, err := s.categoryRepo.ListCategories(ctx)
	if err != nil {
		s.logger.Error("repo.ListCategories", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbCategories := make([]*pb.Category, 0, len(categories))
	for _, category := range categories {
		pbCategories = append(pbCategories, categoryToProto(category))
	}

	return &pb.ListCategoriesResponse{Categories: pbCategories}, nil
}

// categoryToProto maps a database category to the public RPC shape.
func categoryToProto(category *entity.Category) *pb.Category {
	if category == nil {
		return &pb.Category{}
	}
	return &pb.Category{
		Id:       category.ID,
		Name:     category.Name,
		ParentId: category.ParentId,
		Path:     category.Path,
	}
}
