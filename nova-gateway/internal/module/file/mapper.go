package file

import (
	"time"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/media"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoToFile(file *filev1.File) openapi.File {
	if file == nil {
		return openapi.File{
			Scene:  openapi.FILE_SCENE_UNSPECIFIED,
			Status: openapi.FILE_STATUS_UNSPECIFIED,
		}
	}
	return openapi.File{
		Id:        file.GetId(),
		OwnerId:   file.GetOwnerId(),
		OwnerType: file.GetOwnerType(),
		Scene:     protoToFileScene(file.GetScene()),
		ObjectKey: file.GetObjectKey(),
		Url:       media.FileURL(file.GetUrl(), file.GetObjectKey()),
		MimeType:  file.GetMimeType(),
		Size:      file.GetSize(),
		Checksum:  file.GetChecksum(),
		Status:    protoToFileStatus(file.GetStatus()),
		CreatedBy: file.GetCreatedBy(),
		CreatedAt: requiredTimestampTime(file.GetCreatedAt()),
		UpdatedAt: requiredTimestampTime(file.GetUpdatedAt()),
		DeletedAt: timestampTime(file.GetDeletedAt()),
	}
}

func protoToFileScene(value filev1.FileScene) openapi.FileScene {
	switch value {
	case filev1.FileScene_FILE_SCENE_AVATAR:
		return openapi.FILE_SCENE_AVATAR
	case filev1.FileScene_FILE_SCENE_POST_COVER:
		return openapi.FILE_SCENE_POST_COVER
	case filev1.FileScene_FILE_SCENE_POST_CONTENT:
		return openapi.FILE_SCENE_POST_CONTENT
	default:
		return openapi.FILE_SCENE_UNSPECIFIED
	}
}

func protoToFileStatus(value filev1.FileStatus) openapi.FileStatus {
	switch value {
	case filev1.FileStatus_FILE_STATUS_ACTIVE:
		return openapi.FILE_STATUS_ACTIVE
	case filev1.FileStatus_FILE_STATUS_DELETED:
		return openapi.FILE_STATUS_DELETED
	default:
		return openapi.FILE_STATUS_UNSPECIFIED
	}
}

func protoToUser(user *userv1.User) openapi.User {
	if user == nil {
		return openapi.User{Status: openapi.USER_STATUS_UNSPECIFIED}
	}
	return openapi.User{
		Id:                user.GetId(),
		Username:          user.GetUsername(),
		Nickname:          user.GetNickname(),
		Avatar:            media.UploadsURL(user.GetAvatar()),
		Email:             user.GetEmail(),
		EmailVerifiedTime: timestampTime(user.GetEmailVerifiedTime()),
		Phone:             user.GetPhone(),
		Status:            protoToUserStatus(user.GetStatus()),
		CreatedAt:         requiredTimestampTime(user.GetCreatedAt()),
		UpdatedAt:         requiredTimestampTime(user.GetUpdatedAt()),
		DeletedAt:         timestampTime(user.GetDeletedAt()),
	}
}

func protoToUserStatus(value userv1.UserStatus) openapi.UserStatus {
	switch value {
	case userv1.UserStatus_USER_STATUS_ACTIVE:
		return openapi.USER_STATUS_ACTIVE
	case userv1.UserStatus_USER_STATUS_DISABLE:
		return openapi.USER_STATUS_DISABLED
	default:
		return openapi.USER_STATUS_UNSPECIFIED
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
