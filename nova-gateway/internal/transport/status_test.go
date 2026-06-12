package transport

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestGRPCHTTPStatus(t *testing.T) {
	tests := []struct {
		code codes.Code
		want int
	}{
		{code: codes.OK, want: http.StatusOK},
		{code: codes.Canceled, want: 499},
		{code: codes.InvalidArgument, want: http.StatusBadRequest},
		{code: codes.NotFound, want: http.StatusNotFound},
		{code: codes.AlreadyExists, want: http.StatusConflict},
		{code: codes.PermissionDenied, want: http.StatusForbidden},
		{code: codes.ResourceExhausted, want: http.StatusTooManyRequests},
		{code: codes.Unimplemented, want: http.StatusNotImplemented},
		{code: codes.Unauthenticated, want: http.StatusUnauthorized},
		{code: codes.Unavailable, want: http.StatusServiceUnavailable},
		{code: codes.DataLoss, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			if got := grpcHTTPStatus(tt.code); got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
