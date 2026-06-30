package post

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/media"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/logger"
	"github.com/unknwon/paginater"
)

type PostsHandler struct {
	logger         logger.Logger
	postClient     *client.PostClient
	fileClient     *client.FileClient
	sessionManager *websession.Manager
}

func NewPostsHandler(logger logger.Logger, postClient *client.PostClient, fileClient *client.FileClient, sessionManager *websession.Manager) *PostsHandler {
	return &PostsHandler{
		logger:         logger,
		postClient:     postClient,
		fileClient:     fileClient,
		sessionManager: sessionManager,
	}
}

type PostListViewData struct {
	template.ViewData
	Posts       []apiclient.Post
	Total       int64
	CurrentPage int32
	TotalPages  int32
	PageSize    int32
	Pager       *paginater.Paginater
	Error       string
}

type PostDetailViewData struct {
	template.ViewData
	Post      *apiclient.Post
	CanManage bool
	Error     string
}

type PostFormData struct {
	template.ViewData
	Post  *apiclient.Post
	Error string
}

func (h *PostsHandler) index(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	resp, err := h.postClient.ListPosts(c.Request.Context(), int32(page), int32(pageSize))
	if err != nil {
		template.InternalError(c)
		return
	}

	c.HTML(http.StatusOK, "post/list", newPostListViewData(c, resp, int32(page), int32(pageSize)))
}

func (h *PostsHandler) pages(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Param("page"), 10, 32)
	if page < 1 {
		page = 1
	}
	pageSize := int32(10)

	resp, err := h.postClient.ListPosts(c.Request.Context(), int32(page), pageSize)
	if err != nil {
		template.InternalError(c)
		return
	}

	c.HTML(http.StatusOK, "post/list", newPostListViewData(c, resp, int32(page), pageSize))
}

func newPostListViewData(c *gin.Context, resp *apiclient.ListPostsResponse, page, pageSize int32) PostListViewData {
	view := PostListViewData{
		ViewData:    template.NewViewData(c),
		CurrentPage: page,
		PageSize:    pageSize,
		Pager:       paginater.New(0, int(pageSize), int(page), 5),
	}
	if resp != nil {
		view.Posts = postsForView(resp.Posts)
		view.Total = resp.Total
		view.CurrentPage = resp.CurrentPage
		view.TotalPages = resp.TotalPages
		view.PageSize = resp.PageSize
		view.Pager = paginater.New(int(resp.Total), int(resp.PageSize), int(resp.CurrentPage), 5)
	}
	return view
}

func (h *PostsHandler) show(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	p, err := h.postClient.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}

	viewData := template.NewViewData(c)
	if _, ok := authctx.CurrentUser(c); ok {
		viewData = template.NewFormViewData(c)
	}
	c.HTML(http.StatusOK, "post/detail", PostDetailViewData{
		ViewData:  viewData,
		Post:      postForView(p),
		CanManage: h.canManagePost(c, p),
	})
}

func (h *PostsHandler) create(c *gin.Context) {
	c.HTML(http.StatusOK, "post/create", PostFormData{
		ViewData: template.NewFormViewData(c),
	})
}

func (h *PostsHandler) store(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	coverURL, err := h.uploadPostCoverIfPresent(c)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		h.renderCreateFormError(c, title, content, err)
		return
	}

	_, err = h.postClient.CreatePost(c.Request.Context(), title, content, coverURL)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		h.renderCreateFormErrorWithCover(c, title, content, coverURL, err)
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

func (h *PostsHandler) edit(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	p, err := h.postClient.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}
	if !h.canManagePost(c, p) {
		c.Status(http.StatusForbidden)
		return
	}

	c.HTML(http.StatusOK, "post/edit", PostFormData{
		ViewData: template.NewFormViewData(c),
		Post:     postForView(p),
	})
}

func (h *PostsHandler) update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	title := c.PostForm("title")
	content := c.PostForm("content")

	p, err := h.postClient.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}
	if !h.canManagePost(c, p) {
		c.Status(http.StatusForbidden)
		return
	}

	coverURL, err := h.uploadPostCoverIfPresent(c)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		p.Title = title
		p.Content = content
		h.renderEditFormError(c, p, err)
		return
	}

	// nil keeps cover_url out of the update mask when no new cover is uploaded.
	var coverURLPtr *string
	if coverURL != "" {
		coverURLPtr = &coverURL
	}

	_, err = h.postClient.UpdatePost(c.Request.Context(), id, title, content, coverURLPtr)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		p.Title = title
		p.Content = content
		if coverURLPtr != nil {
			p.CoverUrl = *coverURLPtr
		}
		h.renderEditFormError(c, p, err)
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+c.Param("id"))
}

func (h *PostsHandler) post(c *gin.Context) {
	switch c.PostForm("_method") {
	case "PUT":
		h.update(c)
	case "DELETE":
		h.destroy(c)
	default:
		c.Status(http.StatusMethodNotAllowed)
	}
}

func (h *PostsHandler) destroy(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	p, err := h.postClient.GetPost(c.Request.Context(), id)
	if err != nil {
		template.NotFound(c)
		return
	}
	if !h.canManagePost(c, p) {
		c.Status(http.StatusForbidden)
		return
	}

	err = h.postClient.DeletePost(c.Request.Context(), id)
	if err != nil {
		if h.handleAuthError(c, err) {
			return
		}
		template.InternalError(c)
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

func (h *PostsHandler) canManagePost(c *gin.Context, post *apiclient.Post) bool {
	if post == nil || post.AuthorId <= 0 {
		return false
	}
	userID, ok := authctx.CurrentUserInt64ID(c)
	return ok && userID == post.AuthorId
}

func (h *PostsHandler) handleAuthError(c *gin.Context, err error) bool {
	if !client.IsStatus(err, http.StatusUnauthorized) {
		return false
	}
	h.sessionManager.Clear(c)
	c.Redirect(http.StatusFound, "/login")
	return true
}

func (h *PostsHandler) renderCreateFormError(c *gin.Context, title, content string, err error) {
	h.renderCreateFormErrorWithCover(c, title, content, "", err)
}

func (h *PostsHandler) renderCreateFormErrorWithCover(c *gin.Context, title, content, coverURL string, err error) {
	c.HTML(postFormErrorStatus(err), "post/create", PostFormData{
		ViewData: template.NewFormViewData(c),
		Post: &apiclient.Post{
			Title:    title,
			CoverUrl: media.UploadsURL(coverURL),
			Content:  content,
		},
		Error: err.Error(),
	})
}

func (h *PostsHandler) renderEditFormError(c *gin.Context, post *apiclient.Post, err error) {
	c.HTML(postFormErrorStatus(err), "post/edit", PostFormData{
		ViewData: template.NewFormViewData(c),
		Post:     postForView(post),
		Error:    err.Error(),
	})
}

func postsForView(posts []apiclient.Post) []apiclient.Post {
	if len(posts) == 0 {
		return posts
	}
	viewPosts := append([]apiclient.Post(nil), posts...)
	for i := range viewPosts {
		viewPosts[i].CoverUrl = media.UploadsURL(viewPosts[i].CoverUrl)
	}
	return viewPosts
}

func postForView(post *apiclient.Post) *apiclient.Post {
	if post == nil {
		return nil
	}
	viewPost := *post
	viewPost.CoverUrl = media.UploadsURL(viewPost.CoverUrl)
	return &viewPost
}

func (h *PostsHandler) uploadPostCoverIfPresent(c *gin.Context) (string, error) {
	file, header, err := c.Request.FormFile("cover")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		if errors.Is(err, multipart.ErrMessageTooLarge) {
			return "", client.NewHTTPError(http.StatusRequestEntityTooLarge, "cover image is too large")
		}
		return "", err
	}
	defer file.Close()
	if header == nil || header.Filename == "" {
		return "", nil
	}
	if h.fileClient == nil {
		return "", client.NewHTTPError(http.StatusBadGateway, "file service not configured")
	}

	resp, err := h.fileClient.UploadPostCover(c.Request.Context(), header.Filename, file)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", client.NewHTTPError(http.StatusBadGateway, "empty file object key")
	}
	if resp.ObjectKey == "" {
		return "", client.NewHTTPError(http.StatusBadGateway, "empty file object key")
	}
	return resp.ObjectKey, nil
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
