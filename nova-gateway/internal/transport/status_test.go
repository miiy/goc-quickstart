package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestWriteErrorUsesOpenAPIErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	WriteError(c, status.Error(codes.NotFound, "post not found"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	var body openapi.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != int32(codes.NotFound) || body.Error.Message != "post not found" {
		t.Fatalf("unexpected error body: %+v", body)
	}
}

func TestBindJSONWritesBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{`))
	c.Request.Header.Set("Content-Type", "application/json")

	var dst struct {
		Name string `json:"name"`
	}
	if BindJSON(c, &dst) {
		t.Fatal("BindJSON succeeded for invalid json")
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	var body openapi.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != int32(codes.InvalidArgument) || body.Error.Message == "" {
		t.Fatalf("unexpected error body: %+v", body)
	}
}
