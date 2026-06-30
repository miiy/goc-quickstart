package user

import (
	"time"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoUserInput(input openapi.UserInput) (*userv1.User, error) {
	userStatus, err := protoUserStatus(input.Status)
	if err != nil {
		return nil, err
	}

	return &userv1.User{
		Nickname: input.Nickname,
		Avatar:   input.Avatar,
		Email:    input.Email,
		Phone:    input.Phone,
		Status:   userStatus,
	}, nil
}

func protoUserStatus(value openapi.UserStatus) (userv1.UserStatus, error) {
	switch value {
	case "", openapi.USER_STATUS_UNSPECIFIED:
		return userv1.UserStatus_USER_STATUS_UNSPECIFIED, nil
	case openapi.USER_STATUS_ACTIVE:
		return userv1.UserStatus_USER_STATUS_ACTIVE, nil
	case openapi.USER_STATUS_DISABLED:
		return userv1.UserStatus_USER_STATUS_DISABLE, nil
	default:
		return userv1.UserStatus_USER_STATUS_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "unsupported user status: %s", value)
	}
}

func openapiUserStatus(value userv1.UserStatus) openapi.UserStatus {
	switch value {
	case userv1.UserStatus_USER_STATUS_ACTIVE:
		return openapi.USER_STATUS_ACTIVE
	case userv1.UserStatus_USER_STATUS_DISABLE:
		return openapi.USER_STATUS_DISABLED
	default:
		return openapi.USER_STATUS_UNSPECIFIED
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
	case "nickname", "avatar", "email", "phone", "status":
		return path, true
	default:
		return "", false
	}
}

func openapiUsers(users []*userv1.User) []openapi.User {
	result := make([]openapi.User, 0, len(users))
	for _, user := range users {
		result = append(result, OpenAPIUser(user))
	}
	return result
}

func OpenAPIUser(user *userv1.User) openapi.User {
	if user == nil {
		return openapi.User{Status: openapi.USER_STATUS_UNSPECIFIED}
	}
	return openapi.User{
		Id:                user.GetId(),
		Username:          user.GetUsername(),
		Nickname:          user.GetNickname(),
		Avatar:            user.GetAvatar(),
		Email:             user.GetEmail(),
		EmailVerifiedTime: timestampTime(user.GetEmailVerifiedTime()),
		Phone:             user.GetPhone(),
		Status:            openapiUserStatus(user.GetStatus()),
		CreatedAt:         requiredTimestampTime(user.GetCreatedAt()),
		UpdatedAt:         requiredTimestampTime(user.GetUpdatedAt()),
		DeletedAt:         timestampTime(user.GetDeletedAt()),
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
