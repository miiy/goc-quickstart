package transport

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
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

func BindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		WriteBadRequest(c, err.Error())
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
	WriteOpenAPIError(c, http.StatusUnauthorized, int32(codes.Unauthenticated), message)
}

func WriteBadRequest(c *gin.Context, message string) {
	WriteOpenAPIError(c, http.StatusBadRequest, int32(codes.InvalidArgument), message)
}

func WriteOpenAPIError(c *gin.Context, httpStatus int, code int32, message string) {
	c.AbortWithStatusJSON(httpStatus, openapi.ErrorResponse{
		Error: openapi.ErrorStatus{
			Code:    code,
			Message: message,
		},
	})
}

func WriteError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Internal, err.Error())
	}

	WriteOpenAPIError(c, grpcHTTPStatus(st.Code()), int32(st.Code()), st.Message())
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

func Int64Query(c *gin.Context, name string) (int64, bool) {
	raw := c.Query(name)
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		WriteError(c, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", name, err))
		return 0, false
	}
	return value, true
}

func Int64SliceQuery(c *gin.Context, name string) ([]int64, bool) {
	rawValues := c.QueryArray(name)
	if len(rawValues) == 0 {
		return nil, true
	}

	values := make([]int64, 0, len(rawValues))
	for _, raw := range rawValues {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			value, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				WriteError(c, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", name, err))
				return nil, false
			}
			values = append(values, value)
		}
	}
	return values, true
}

func Int32Query(c *gin.Context, name string) (int32, bool) {
	raw := c.Query(name)
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		WriteError(c, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", name, err))
		return 0, false
	}
	return int32(value), true
}
