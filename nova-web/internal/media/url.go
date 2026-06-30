package media

import "strings"

const uploadsURLPrefix = "/uploads"

func UploadsURL(objectKey string) string {
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return ""
	}
	return uploadsURLPrefix + "/" + strings.TrimLeft(objectKey, "/")
}
