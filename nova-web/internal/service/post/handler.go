package post

import (
	"net/http"
	"strconv"

	"github.com/miiy/goc-quickstart/nova-web/client"
	webauth "github.com/miiy/goc-quickstart/nova-web/internal/service/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
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
	Post      *client.PostResponse
	CanManage bool
	Error     string
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
		ViewData:    template.NewViewData(c),
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

	viewData := template.NewViewData(c)
	if c.GetBool("isLoggedIn") {
		viewData = template.NewFormViewData(c)
	}
	c.HTML(http.StatusOK, "post/detail", PostDetailViewData{
		ViewData:  viewData,
		Post:      p,
		CanManage: canManagePost(c, p),
	})
}

func createHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "post/create", PostFormData{
		ViewData: template.NewFormViewData(c),
	})
}

func storeHandler(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	authorID, ok := gocauthmid.GetAuthUserID(c)
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}

	_, err := postModule.client.CreatePost(c.Request.Context(), title, content, authorID)
	if err != nil {
		if handleAuthError(c, err) {
			return
		}
		renderCreateFormError(c, title, content, err)
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
	if !canManagePost(c, p) {
		c.Status(http.StatusForbidden)
		return
	}

	c.HTML(http.StatusOK, "post/edit", PostFormData{
		ViewData: template.NewFormViewData(c),
		Post:     p,
	})
}

func updateHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	title := c.PostForm("title")
	content := c.PostForm("content")

	p, err := postModule.client.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}
	if !canManagePost(c, p) {
		c.Status(http.StatusForbidden)
		return
	}

	_, err = postModule.client.UpdatePost(c.Request.Context(), id, title, content)
	if err != nil {
		if handleAuthError(c, err) {
			return
		}
		p.Title = title
		p.Content = content
		renderEditFormError(c, p, err)
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

	p, err := postModule.client.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}
	if !canManagePost(c, p) {
		c.Status(http.StatusForbidden)
		return
	}

	err = postModule.client.DeletePost(c.Request.Context(), id)
	if err != nil {
		if handleAuthError(c, err) {
			return
		}
		template.InternalError(c)
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

func canManagePost(c *gin.Context, post *client.PostResponse) bool {
	if post == nil || post.AuthorId <= 0 {
		return false
	}
	userID, ok := gocauthmid.GetAuthUserID(c)
	return ok && userID == post.AuthorId
}

func handleAuthError(c *gin.Context, err error) bool {
	if !client.IsStatus(err, http.StatusUnauthorized) {
		return false
	}
	webauth.ClearSession(c)
	c.Redirect(http.StatusFound, "/login")
	return true
}

func renderCreateFormError(c *gin.Context, title, content string, err error) {
	c.HTML(postFormErrorStatus(err), "post/create", PostFormData{
		ViewData: template.NewFormViewData(c),
		Post: &client.PostResponse{
			Title:   title,
			Content: content,
		},
		Error: err.Error(),
	})
}

func renderEditFormError(c *gin.Context, post *client.PostResponse, err error) {
	c.HTML(postFormErrorStatus(err), "post/edit", PostFormData{
		ViewData: template.NewFormViewData(c),
		Post:     post,
		Error:    err.Error(),
	})
}

func postFormErrorStatus(err error) int {
	switch {
	case client.IsStatus(err, http.StatusBadRequest):
		return http.StatusBadRequest
	case client.IsStatus(err, http.StatusForbidden):
		return http.StatusForbidden
	default:
		return http.StatusBadGateway
	}
}
