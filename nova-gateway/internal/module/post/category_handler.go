package post

import (
	"net/http"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
)

func (api *PostsAPI) ListCategories(c *gin.Context) {
	resp, err := api.postClient.ListCategories(c.Request.Context(), &postv1.ListCategoriesRequest{})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.ListCategoriesResponse{Categories: protoToCategories(resp.GetCategories())})
}
