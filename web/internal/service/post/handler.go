package post

import (
	"net/http"
	"strconv"

	"github.com/miiy/goc-quickstart/web/client"
	"github.com/miiy/goc-quickstart/web/internal/template"
	"github.com/miiy/goc/gin"
	"github.com/unknwon/paginater"
)

type PostListViewData struct {
	template.ViewData
	Posts       []*client.PostResponse
	Total       int32
	CurrentPage int32
	TotalPages  int32
	PageSize    int32
	Pager       *paginater.Paginater
	Error       string
}

type PostDetailViewData struct {
	template.ViewData
	Post  *client.PostResponse
	Error string
}

type PostFormData struct {
	template.ViewData
	Post  *client.PostResponse
	Error string
}

func indexHandler(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	resp, err := postModule.client.ListPosts(c.Request.Context(), int32(page), int32(pageSize))
	if err != nil {
		template.InternalError(c)
		return
	}

	c.HTML(http.StatusOK, "post/list", newPostListViewData(c, resp, int32(page), int32(pageSize)))
}

func pagesHandler(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Param("page"), 10, 32)
	if page < 1 {
		page = 1
	}
	pageSize := int32(10)

	resp, err := postModule.client.ListPosts(c.Request.Context(), int32(page), pageSize)
	if err != nil {
		template.InternalError(c)
		return
	}

	c.HTML(http.StatusOK, "post/list", newPostListViewData(c, resp, int32(page), pageSize))
}

func newPostListViewData(c *gin.Context, resp *client.PostListResponse, page, pageSize int32) PostListViewData {
	view := PostListViewData{
		ViewData: template.ViewData{
			IsLoggedIn: c.GetBool("isLoggedIn"),
		},
		CurrentPage: page,
		PageSize:    pageSize,
		Pager:       paginater.New(0, int(pageSize), int(page), 5),
	}
	if resp != nil {
		view.Posts = resp.Posts
		view.Total = resp.Total
		view.CurrentPage = resp.CurrentPage
		view.TotalPages = resp.TotalPages
		view.PageSize = resp.PageSize
		view.Pager = paginater.New(int(resp.Total), int(resp.PageSize), int(resp.CurrentPage), 5)
	}
	return view
}

func showHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	p, err := postModule.client.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}

	c.HTML(http.StatusOK, "post/detail", PostDetailViewData{
		ViewData: template.ViewData{IsLoggedIn: c.GetBool("isLoggedIn")},
		Post:     p,
	})
}

func createHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "post/create", PostFormData{
		ViewData: template.ViewData{IsLoggedIn: c.GetBool("isLoggedIn")},
	})
}

func storeHandler(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	_, err := postModule.client.CreatePost(c.Request.Context(), title, content, 1)
	if err != nil {
		template.InternalError(c)
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

func editHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	p, err := postModule.client.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}

	c.HTML(http.StatusOK, "post/edit", PostFormData{
		ViewData: template.ViewData{IsLoggedIn: c.GetBool("isLoggedIn")},
		Post:     p,
	})
}

func updateHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	title := c.PostForm("title")
	content := c.PostForm("content")

	_, err := postModule.client.UpdatePost(c.Request.Context(), id, title, content)
	if err != nil {
		template.InternalError(c)
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+c.Param("id"))
}

func postHandler(c *gin.Context) {
	switch c.PostForm("_method") {
	case "PUT":
		updateHandler(c)
	case "DELETE":
		destroyHandler(c)
	default:
		c.Status(http.StatusMethodNotAllowed)
	}
}

func destroyHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := postModule.client.DeletePost(c.Request.Context(), id)
	if err != nil {
		template.InternalError(c)
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}
