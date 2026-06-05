package transport

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/miiy/goc/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	protoMarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}
	protoUnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type ErrorResponse struct {
	Error ErrorStatus `json:"error"`
}

type ErrorStatus struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func BindProto(c *gin.Context, msg proto.Message) bool {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		WriteError(c, status.Error(codes.InvalidArgument, err.Error()))
		return false
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return true
	}
	if err := protoUnmarshalOptions.Unmarshal(body, msg); err != nil {
		WriteError(c, status.Errorf(codes.InvalidArgument, "%v", err))
		return false
	}
	return true
}

func WriteProto(c *gin.Context, msg proto.Message) {
	body, err := protoMarshalOptions.Marshal(msg)
	if err != nil {
		WriteError(c, status.Error(codes.Internal, err.Error()))
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}

// WriteUnauthorized writes a 401 response with a standardized error body.
func WriteUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
		Error: ErrorStatus{
			Code:    int32(codes.Unauthenticated),
			Message: message,
		},
	})
}

func WriteError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Internal, err.Error())
	}

	c.AbortWithStatusJSON(grpcHTTPStatus(st.Code()), ErrorResponse{
		Error: ErrorStatus{
			Code:    int32(st.Code()),
			Message: st.Message(),
		},
	})
}

func Int64Param(c *gin.Context, name string) (int64, bool) {
	raw := c.Param(name)
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		WriteError(c, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", name, err))
		return 0, false
	}
	return value, true
}

func Int64Query(c *gin.Context, snakeName string, camelName string) (int64, bool) {
	raw := QueryValue(c, snakeName, camelName)
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		WriteError(c, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", snakeName, err))
		return 0, false
	}
	return value, true
}

func Int32Query(c *gin.Context, snakeName string, camelName string) (int32, bool) {
	raw := QueryValue(c, snakeName, camelName)
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		WriteError(c, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", snakeName, err))
		return 0, false
	}
	return int32(value), true
}

func QueryValue(c *gin.Context, snakeName string, camelName string) string {
	if value, ok := c.GetQuery(snakeName); ok {
		return value
	}
	if camelName != "" {
		if value, ok := c.GetQuery(camelName); ok {
			return value
		}
	}
	return ""
}
