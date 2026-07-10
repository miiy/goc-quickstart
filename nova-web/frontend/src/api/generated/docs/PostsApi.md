# PostsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createPost**](PostsApi.md#createpostoperation) | **POST** /posts | Create a post |
| [**deletePost**](PostsApi.md#deletepost) | **DELETE** /posts/{id} | Delete a post |
| [**getPost**](PostsApi.md#getpost) | **GET** /posts/{id} | Get a post |
| [**getUserPost**](PostsApi.md#getuserpost) | **GET** /users/{username}/posts/{id} | Get a post owned by the authenticated user |
| [**listCategories**](PostsApi.md#listcategories) | **GET** /categories | List post categories |
| [**listPosts**](PostsApi.md#listposts) | **GET** /posts | List posts |
| [**listUserPosts**](PostsApi.md#listuserposts) | **GET** /users/{username}/posts | List posts owned by the authenticated user |
| [**updatePost**](PostsApi.md#updatepostoperation) | **PUT** /posts/{id} | Update a post |



## createPost

> CreatePostResponse createPost(createPostRequest)

Create a post

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { CreatePostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PostsApi(config);

  const body = {
    // CreatePostRequest
    createPostRequest: ...,
  } satisfies CreatePostOperationRequest;

  try {
    const data = await api.createPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **createPostRequest** | [CreatePostRequest](CreatePostRequest.md) |  | |

### Return type

[**CreatePostResponse**](CreatePostResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Created post. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deletePost

> object deletePost(id)

Delete a post

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { DeletePostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PostsApi(config);

  const body = {
    // string | Post id.
    id: id_example,
  } satisfies DeletePostRequest;

  try {
    const data = await api.deletePost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | Post id. | [Defaults to `undefined`] |

### Return type

**object**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Post deleted. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getPost

> GetPostResponse getPost(id)

Get a post

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { GetPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PostsApi();

  const body = {
    // string | Post id.
    id: id_example,
  } satisfies GetPostRequest;

  try {
    const data = await api.getPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | Post id. | [Defaults to `undefined`] |

### Return type

[**GetPostResponse**](GetPostResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Post detail. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getUserPost

> GetPostResponse getUserPost(username, id)

Get a post owned by the authenticated user

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { GetUserPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PostsApi(config);

  const body = {
    // string | Username.
    username: username_example,
    // string | Post id.
    id: id_example,
  } satisfies GetUserPostRequest;

  try {
    const data = await api.getUserPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | `string` | Username. | [Defaults to `undefined`] |
| **id** | `string` | Post id. | [Defaults to `undefined`] |

### Return type

[**GetPostResponse**](GetPostResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Current user\&#39;s post detail, including drafts. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listCategories

> ListCategoriesResponse listCategories()

List post categories

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { ListCategoriesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PostsApi();

  try {
    const data = await api.listCategories();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ListCategoriesResponse**](ListCategoriesResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Read-only post categories. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listPosts

> ListPostsResponse listPosts(userId, categoryId, tag, page, pageSize)

List posts

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { ListPostsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PostsApi();

  const body = {
    // number | Filter by user id. (optional)
    userId: 789,
    // number | Filter by category id. (optional)
    categoryId: 789,
    // string | Filter by tag. (optional)
    tag: tag_example,
    // number | Page number. (optional)
    page: 56,
    // number | Number of posts per page. (optional)
    pageSize: 56,
  } satisfies ListPostsRequest;

  try {
    const data = await api.listPosts(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **userId** | `number` | Filter by user id. | [Optional] [Defaults to `undefined`] |
| **categoryId** | `number` | Filter by category id. | [Optional] [Defaults to `undefined`] |
| **tag** | `string` | Filter by tag. | [Optional] [Defaults to `undefined`] |
| **page** | `number` | Page number. | [Optional] [Defaults to `undefined`] |
| **pageSize** | `number` | Number of posts per page. | [Optional] [Defaults to `undefined`] |

### Return type

[**ListPostsResponse**](ListPostsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Posts matched by filters. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listUserPosts

> ListPostsResponse listUserPosts(username, page, pageSize)

List posts owned by the authenticated user

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { ListUserPostsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PostsApi(config);

  const body = {
    // string | Username.
    username: username_example,
    // number | Page number. (optional)
    page: 56,
    // number | Number of posts per page. (optional)
    pageSize: 56,
  } satisfies ListUserPostsRequest;

  try {
    const data = await api.listUserPosts(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **username** | `string` | Username. | [Defaults to `undefined`] |
| **page** | `number` | Page number. | [Optional] [Defaults to `undefined`] |
| **pageSize** | `number` | Number of posts per page. | [Optional] [Defaults to `undefined`] |

### Return type

[**ListPostsResponse**](ListPostsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Current user\&#39;s posts, including drafts. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updatePost

> UpdatePostResponse updatePost(id, updatePostRequest)

Update a post

### Example

```ts
import {
  Configuration,
  PostsApi,
} from '';
import type { UpdatePostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new PostsApi(config);

  const body = {
    // string | Post id.
    id: id_example,
    // UpdatePostRequest
    updatePostRequest: ...,
  } satisfies UpdatePostOperationRequest;

  try {
    const data = await api.updatePost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | Post id. | [Defaults to `undefined`] |
| **updatePostRequest** | [UpdatePostRequest](UpdatePostRequest.md) |  | |

### Return type

[**UpdatePostResponse**](UpdatePostResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Updated post. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

