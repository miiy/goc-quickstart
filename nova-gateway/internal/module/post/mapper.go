package post

import (
	"time"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/media"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func openapiToProtoCreatePostInput(input openapi.CreatePostInput) (*postv1.Post, error) {
	postStatus, err := openapiToProtoPostStatus(input.Status)
	if err != nil {
		return nil, err
	}

	return &postv1.Post{
		Title:       input.Title,
		Summary:     input.Summary,
		Content:     input.Content,
		Status:      postStatus,
		Tags:        input.Tags,
		CategoryId:  input.CategoryId,
		CoverUrl:    input.CoverUrl,
		PublishedAt: timeToProtoTimestamp(input.PublishedAt),
	}, nil
}

func openapiToProtoUpdatePostInput(input openapi.UpdatePostInput) (*postv1.Post, error) {
	postStatus, err := openapiToProtoPostStatus(input.Status)
	if err != nil {
		return nil, err
	}

	return &postv1.Post{
		Title:       input.Title,
		Summary:     input.Summary,
		Content:     input.Content,
		Status:      postStatus,
		Tags:        input.Tags,
		CategoryId:  input.CategoryId,
		CoverUrl:    input.CoverUrl,
		PublishedAt: timeToProtoTimestamp(input.PublishedAt),
	}, nil
}

func openapiToProtoPostStatus(value openapi.PostStatus) (postv1.PostStatus, error) {
	switch value {
	case "", openapi.POST_STATUS_UNSPECIFIED:
		return postv1.PostStatus_POST_STATUS_UNSPECIFIED, nil
	case openapi.POST_STATUS_DRAFT:
		return postv1.PostStatus_POST_STATUS_DRAFT, nil
	case openapi.POST_STATUS_PUBLISHED:
		return postv1.PostStatus_POST_STATUS_PUBLISHED, nil
	case openapi.POST_STATUS_PENDING_REVIEW:
		return postv1.PostStatus_POST_STATUS_PENDING_REVIEW, nil
	default:
		return postv1.PostStatus_POST_STATUS_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "unsupported post status: %s", value)
	}
}

func protoToPostStatus(value postv1.PostStatus) openapi.PostStatus {
	switch value {
	case postv1.PostStatus_POST_STATUS_DRAFT:
		return openapi.POST_STATUS_DRAFT
	case postv1.PostStatus_POST_STATUS_PUBLISHED:
		return openapi.POST_STATUS_PUBLISHED
	case postv1.PostStatus_POST_STATUS_PENDING_REVIEW:
		return openapi.POST_STATUS_PENDING_REVIEW
	default:
		return openapi.POST_STATUS_UNSPECIFIED
	}
}

func openapiToProtoUpdateMask(fields []string) (*fieldmaskpb.FieldMask, error) {
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
	case "title", "summary", "content", "cover_url", "published_at", "status", "tags", "category_id":
		return path, true
	default:
		return "", false
	}
}

func protoToPosts(posts []*postv1.Post, postUsersByID map[int64]openapi.PostUser, currentUserID int64) []openapi.Post {
	result := make([]openapi.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, protoToPost(post, postUsersByID, currentUserID))
	}
	return result
}

func protoToPost(post *postv1.Post, postUsersByID map[int64]openapi.PostUser, currentUserID int64) openapi.Post {
	if post == nil {
		return openapi.Post{
			Status: openapi.POST_STATUS_UNSPECIFIED,
			Tags:   []string{},
			User:   openapi.PostUser{},
		}
	}
	user := postUsersByID[post.GetUserId()]
	return openapi.Post{
		Id:          encodePostID(post.GetId()),
		UserId:      post.GetUserId(),
		User:        user,
		Title:       post.GetTitle(),
		Summary:     post.GetSummary(),
		Content:     post.GetContent(),
		Status:      protoToPostStatus(post.GetStatus()),
		Tags:        protoToTags(post.GetTags()),
		CategoryId:  post.GetCategoryId(),
		PublishedAt: timestampTime(post.GetPublishedAt()),
		CreatedAt:   requiredTimestampTime(post.GetCreatedAt()),
		UpdatedAt:   requiredTimestampTime(post.GetUpdatedAt()),
		DeletedAt:   timestampTime(post.GetDeletedAt()),
		CoverUrl:    media.UploadsURL(post.GetCoverUrl()),
		CanManage:   canManagePost(currentUserID, post),
	}
}

func canManagePost(currentUserID int64, post *postv1.Post) bool {
	return currentUserID > 0 && post != nil && post.GetUserId() > 0 && currentUserID == post.GetUserId()
}

func protoToCategories(categories []*postv1.Category) []openapi.Category {
	result := make([]openapi.Category, 0, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		result = append(result, openapi.Category{
			Id:       category.GetId(),
			Name:     category.GetName(),
			ParentId: category.GetParentId(),
			Path:     category.GetPath(),
		})
	}
	return result
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

func timeToProtoTimestamp(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}
	return timestamppb.New(*value)
}

func protoToTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}
