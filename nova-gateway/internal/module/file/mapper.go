package file

import (
	"time"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func openapiFile(file *filev1.File) openapi.File {
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
		Scene:     openapiFileScene(file.GetScene()),
		ObjectKey: file.GetObjectKey(),
		Url:       file.GetUrl(),
		MimeType:  file.GetMimeType(),
		Size:      file.GetSize(),
		Checksum:  file.GetChecksum(),
		Status:    openapiFileStatus(file.GetStatus()),
		CreatedBy: file.GetCreatedBy(),
		CreatedAt: requiredTimestampTime(file.GetCreatedAt()),
		UpdatedAt: requiredTimestampTime(file.GetUpdatedAt()),
		DeletedAt: timestampTime(file.GetDeletedAt()),
	}
}

func openapiFileScene(value filev1.FileScene) openapi.FileScene {
	switch value {
	case filev1.FileScene_FILE_SCENE_AVATAR:
		return openapi.FILE_SCENE_AVATAR
	case filev1.FileScene_FILE_SCENE_POST_COVER:
		return openapi.FILE_SCENE_POST_COVER
	default:
		return openapi.FILE_SCENE_UNSPECIFIED
	}
}

func openapiFileStatus(value filev1.FileStatus) openapi.FileStatus {
	switch value {
	case filev1.FileStatus_FILE_STATUS_ACTIVE:
		return openapi.FILE_STATUS_ACTIVE
	case filev1.FileStatus_FILE_STATUS_DELETED:
		return openapi.FILE_STATUS_DELETED
	default:
		return openapi.FILE_STATUS_UNSPECIFIED
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
