package post

import (
	"time"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoCreatePostInput(input openapi.CreatePostInput) (*postv1.Post, error) {
	postStatus, err := protoPostStatus(input.Status)
	if err != nil {
		return nil, err
	}

	return &postv1.Post{
		Title:      input.Title,
		Content:    input.Content,
		Status:     postStatus,
		Tags:       input.Tags,
		CategoryId: input.CategoryId,
		CoverUrl:   input.CoverUrl,
	}, nil
}

func protoUpdatePostInput(input openapi.UpdatePostInput) (*postv1.Post, error) {
	postStatus, err := protoPostStatus(input.Status)
	if err != nil {
		return nil, err
	}

	return &postv1.Post{
		Title:      input.Title,
		Content:    input.Content,
		Status:     postStatus,
		Tags:       input.Tags,
		CategoryId: input.CategoryId,
		CoverUrl:   input.CoverUrl,
	}, nil
}

func protoPostStatus(value openapi.PostStatus) (postv1.PostStatus, error) {
	switch value {
	case "", openapi.POST_STATUS_UNSPECIFIED:
		return postv1.PostStatus_POST_STATUS_UNSPECIFIED, nil
	case openapi.POST_STATUS_DRAFT:
		return postv1.PostStatus_POST_STATUS_DRAFT, nil
	case openapi.POST_STATUS_PUBLISHED:
		return postv1.PostStatus_POST_STATUS_PUBLISHED, nil
	default:
		return postv1.PostStatus_POST_STATUS_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "unsupported post status: %s", value)
	}
}

func openapiPostStatus(value postv1.PostStatus) openapi.PostStatus {
	switch value {
	case postv1.PostStatus_POST_STATUS_DRAFT:
		return openapi.POST_STATUS_DRAFT
	case postv1.PostStatus_POST_STATUS_PUBLISHED:
		return openapi.POST_STATUS_PUBLISHED
	default:
		return openapi.POST_STATUS_UNSPECIFIED
	}
}

func protoUpdateMask(fields []string) (*fieldmaskpb.FieldMask, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	seen := make(map[string]struct{})
	paths := make([]string, 0, len(fields))
	for _, path := range fields {
		if path == "" {
			continue
		}
		normalized, ok := normalizeUpdateMaskPath(path)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "unsupported update_fields field: %s", path)
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		paths = append(paths, normalized)
	}
	if len(paths) == 0 {
		return nil, nil
	}
	return &fieldmaskpb.FieldMask{Paths: paths}, nil
}

func normalizeUpdateMaskPath(path string) (string, bool) {
	switch path {
	case "title", "content", "cover_url", "status", "tags", "category_id":
		return path, true
	default:
		return "", false
	}
}

func openapiPosts(posts []*postv1.Post, authorNames map[int64]string) []openapi.Post {
	result := make([]openapi.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, openapiPost(post, authorNames))
	}
	return result
}

func openapiPost(post *postv1.Post, authorNames map[int64]string) openapi.Post {
	if post == nil {
		return openapi.Post{
			Status: openapi.POST_STATUS_UNSPECIFIED,
			Tags:   []string{},
		}
	}
	return openapi.Post{
		Id:         post.GetId(),
		AuthorId:   post.GetAuthorId(),
		Title:      post.GetTitle(),
		Content:    post.GetContent(),
		Status:     openapiPostStatus(post.GetStatus()),
		Tags:       openapiTags(post.GetTags()),
		CategoryId: post.GetCategoryId(),
		CreatedAt:  requiredTimestampTime(post.GetCreatedAt()),
		UpdatedAt:  requiredTimestampTime(post.GetUpdatedAt()),
		DeletedAt:  timestampTime(post.GetDeletedAt()),
		AuthorName: authorNames[post.GetAuthorId()],
		CoverUrl:   post.GetCoverUrl(),
	}
}

func timestampTime(value *timestamppb.Timestamp) *time.Time {
	if value == nil {
		return nil
	}
	t := value.AsTime()
	return &t
}

func requiredTimestampTime(value *timestamppb.Timestamp) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.AsTime()
}

func openapiTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}
