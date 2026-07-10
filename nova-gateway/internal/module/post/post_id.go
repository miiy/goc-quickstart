package post

import (
	"strings"

	"github.com/miiy/goc/sqids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var postIDEncoder = sqids.MustNew()

func encodePostID(id int64) string {
	if id <= 0 {
		return ""
	}
	return postIDEncoder.MustEncode(id)
}

func decodePostID(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, status.Error(codes.InvalidArgument, "invalid post id")
	}

	id, err := postIDEncoder.Decode(raw)
	if err != nil || id <= 0 {
		return 0, status.Error(codes.InvalidArgument, "invalid post id")
	}
	return id, nil
}
