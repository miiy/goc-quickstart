package server

import (
	"context"

	postv1 "github.com/miiy/goc-quickstart/post-service/gen/go/shop/post/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postServer struct {
	postv1.UnimplementedPostServiceServer
}

func NewPostServiceServer() postv1.PostServiceServer {
	return &postServer{}
}

func (s *postServer) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.GetPostResponse, error) {

	return &postv1.GetPostResponse{
		Post: &postv1.Post{
			Id:         1,
			Title:      "title",
			Content:    "content",
			Status:     postv1.PostStatus_POST_STATUS_PUBLISHED,
			CreateTime: nil,
			UpdateTime: nil,
			DeleteTime: nil,
		},
	}, nil
}

func (s *postServer) CreatePost(ctx context.Context, req *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error) {
	return &postv1.CreatePostResponse{}, nil
}

func (s *postServer) GetPostError(ctx context.Context, req *postv1.GetPostErrorRequest) (*postv1.GetPostErrorResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid parameters")
}
