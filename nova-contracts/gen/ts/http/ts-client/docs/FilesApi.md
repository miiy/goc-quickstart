# FilesApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**uploadAvatar**](FilesApi.md#uploadavatar) | **POST** /files/upload/avatar | Upload current user\&#39;s avatar |
| [**uploadFile**](FilesApi.md#uploadfile) | **POST** /files/upload | Upload a file |



## uploadAvatar

> UpdateProfileResponse uploadAvatar(avatar)

Upload current user\&#39;s avatar

### Example

```ts
import {
  Configuration,
  FilesApi,
} from '';
import type { UploadAvatarRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FilesApi(config);

  const body = {
    // Blob
    avatar: BINARY_DATA_HERE,
  } satisfies UploadAvatarRequest;

  try {
    const data = await api.uploadAvatar(body);
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
| **avatar** | `Blob` |  | [Defaults to `undefined`] |

### Return type

[**UpdateProfileResponse**](UpdateProfileResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `multipart/form-data`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Updated user with avatar. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## uploadFile

> UploadFileResponse uploadFile(scene, file)

Upload a file

### Example

```ts
import {
  Configuration,
  FilesApi,
} from '';
import type { UploadFileRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: BearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new FilesApi(config);

  const body = {
    // FileScene
    scene: ...,
    // Blob
    file: BINARY_DATA_HERE,
  } satisfies UploadFileRequest;

  try {
    const data = await api.uploadFile(body);
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
| **scene** | `FileScene` |  | [Defaults to `undefined`] [Enum: unspecified, avatar, post_cover, post_content] |
| **file** | `Blob` |  | [Defaults to `undefined`] |

### Return type

[**UploadFileResponse**](UploadFileResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `multipart/form-data`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Stored file record. |  -  |
| **0** | Error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

