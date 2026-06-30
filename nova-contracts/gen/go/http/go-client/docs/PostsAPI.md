# \PostsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreatePost**](PostsAPI.md#CreatePost) | **Post** /posts | Create a post
[**DeletePost**](PostsAPI.md#DeletePost) | **Delete** /posts/{id} | Delete a post
[**GetPost**](PostsAPI.md#GetPost) | **Get** /posts/{id} | Get a post
[**ListPosts**](PostsAPI.md#ListPosts) | **Get** /posts | List posts
[**UpdatePost**](PostsAPI.md#UpdatePost) | **Put** /posts/{id} | Update a post



## CreatePost

> CreatePostResponse CreatePost(ctx).CreatePostRequest(createPostRequest).Execute()

Create a post

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	createPostRequest := *openapiclient.NewCreatePostRequest(*openapiclient.NewCreatePostInput("Title_example", "Content_example")) // CreatePostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PostsAPI.CreatePost(context.Background()).CreatePostRequest(createPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PostsAPI.CreatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreatePost`: CreatePostResponse
	fmt.Fprintf(os.Stdout, "Response from `PostsAPI.CreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPostRequest** | [**CreatePostRequest**](CreatePostRequest.md) |  | 

### Return type

[**CreatePostResponse**](CreatePostResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePost

> map[string]interface{} DeletePost(ctx, id).Execute()

Delete a post

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	id := int64(789) // int64 | Post id.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PostsAPI.DeletePost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PostsAPI.DeletePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeletePost`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `PostsAPI.DeletePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Post id. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**map[string]interface{}**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPost

> GetPostResponse GetPost(ctx, id).Execute()

Get a post

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	id := int64(789) // int64 | Post id.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PostsAPI.GetPost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PostsAPI.GetPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPost`: GetPostResponse
	fmt.Fprintf(os.Stdout, "Response from `PostsAPI.GetPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Post id. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetPostResponse**](GetPostResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPosts

> ListPostsResponse ListPosts(ctx).AuthorId(authorId).CategoryId(categoryId).Tag(tag).Page(page).PageSize(pageSize).Execute()

List posts

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	authorId := int64(789) // int64 | Filter by author id. (optional)
	categoryId := int64(789) // int64 | Filter by category id. (optional)
	tag := "tag_example" // string | Filter by tag. (optional)
	page := int32(56) // int32 | Page number. (optional)
	pageSize := int32(56) // int32 | Number of posts per page. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PostsAPI.ListPosts(context.Background()).AuthorId(authorId).CategoryId(categoryId).Tag(tag).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PostsAPI.ListPosts``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPosts`: ListPostsResponse
	fmt.Fprintf(os.Stdout, "Response from `PostsAPI.ListPosts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPostsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **authorId** | **int64** | Filter by author id. | 
 **categoryId** | **int64** | Filter by category id. | 
 **tag** | **string** | Filter by tag. | 
 **page** | **int32** | Page number. | 
 **pageSize** | **int32** | Number of posts per page. | 

### Return type

[**ListPostsResponse**](ListPostsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdatePost

> UpdatePostResponse UpdatePost(ctx, id).UpdatePostRequest(updatePostRequest).Execute()

Update a post

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	id := int64(789) // int64 | Post id.
	updatePostRequest := *openapiclient.NewUpdatePostRequest(*openapiclient.NewUpdatePostInput()) // UpdatePostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PostsAPI.UpdatePost(context.Background(), id).UpdatePostRequest(updatePostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PostsAPI.UpdatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdatePost`: UpdatePostResponse
	fmt.Fprintf(os.Stdout, "Response from `PostsAPI.UpdatePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Post id. | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updatePostRequest** | [**UpdatePostRequest**](UpdatePostRequest.md) |  | 

### Return type

[**UpdatePostResponse**](UpdatePostResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

