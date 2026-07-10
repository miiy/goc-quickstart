package post

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
)

func TestListCategoriesWritesOpenAPIResponse(t *testing.T) {
	postClient := &fakePostClient{}
	api := &PostsAPI{postClient: postClient}

	r := gin.New()
	r.GET("/api/v1/categories", api.ListCategories)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if !postClient.listCategoriesCalled {
		t.Fatal("list categories request was not sent")
	}
	var body struct {
		Categories []struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
			Path string `json:"path"`
		} `json:"categories"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Categories) != 1 || body.Categories[0].Name != "Engineering" {
		t.Fatalf("unexpected categories: %+v", body.Categories)
	}
}
