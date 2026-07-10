package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorStatus struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorStatus `json:"error"`
}

type HTTPError struct {
	StatusCode int
	Response   ErrorResponse
	Body       string
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return newHTTPError(statusCode, int32(statusCode), message, nil)
}

func newHTTPError(statusCode int, code int32, message string, body []byte) *HTTPError {
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(statusCode)
	}
	if message == "" {
		message = fmt.Sprintf("status: %d", statusCode)
	}
	return &HTTPError{
		StatusCode: statusCode,
		Response: ErrorResponse{
			Error: ErrorStatus{
				Code:    code,
				Message: message,
			},
		},
		Body: string(body),
	}
}

func (e *HTTPError) Error() string {
	if e == nil {
		return ""
	}
	if e.Response.Error.Message != "" {
		return e.Response.Error.Message
	}
	if e.Body == "" {
		return fmt.Sprintf("status: %d", e.StatusCode)
	}
	return fmt.Sprintf("status: %d, body: %s", e.StatusCode, e.Body)
}

func (e *HTTPError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	if body := []byte(e.Body); json.Valid(body) {
		return body, nil
	}
	return json.Marshal(e.Response)
}

func IsStatus(err error, statusCode int) bool {
	var httpErr *HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == statusCode
}

func FromGRPCError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	return newHTTPError(grpcHTTPStatus(st.Code()), int32(st.Code()), st.Message(), nil)
}

// grpcHTTPStatus maps gRPC status codes to HTTP status codes.
// The mapping follows grpc-gateway's runtime.HTTPStatusFromCode behavior.
func grpcHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal, codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
