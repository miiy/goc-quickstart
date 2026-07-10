package media

import "strings"

const uploadsURLPrefix = "/uploads"

func UploadsURL(objectKey string) string {
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return ""
	}
	if strings.HasPrefix(objectKey, "http://") || strings.HasPrefix(objectKey, "https://") {
		return objectKey
	}
	if objectKey == uploadsURLPrefix || strings.HasPrefix(objectKey, uploadsURLPrefix+"/") {
		return objectKey
	}
	return uploadsURLPrefix + "/" + strings.TrimLeft(objectKey, "/")
}

func FileURL(url, objectKey string) string {
	url = strings.TrimSpace(url)
	if url != "" {
		return url
	}
	return UploadsURL(objectKey)
}
