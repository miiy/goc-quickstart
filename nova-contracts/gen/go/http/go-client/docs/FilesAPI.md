# \FilesAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UploadAvatar**](FilesAPI.md#UploadAvatar) | **Post** /files/upload/avatar | Upload current user&#39;s avatar
[**UploadFile**](FilesAPI.md#UploadFile) | **Post** /files/upload | Upload a file



## UploadAvatar

> UpdateUserResponse UploadAvatar(ctx).Avatar(avatar).Execute()

Upload current user's avatar

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
	avatar := os.NewFile(1234, "some_file") // *os.File | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FilesAPI.UploadAvatar(context.Background()).Avatar(avatar).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FilesAPI.UploadAvatar``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UploadAvatar`: UpdateUserResponse
	fmt.Fprintf(os.Stdout, "Response from `FilesAPI.UploadAvatar`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUploadAvatarRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **avatar** | ***os.File** |  | 

### Return type

[**UpdateUserResponse**](UpdateUserResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UploadFile

> UploadFileResponse UploadFile(ctx).Scene(scene).File(file).Execute()

Upload a file

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
	scene := openapiclient.FileScene("unspecified") // FileScene | 
	file := os.NewFile(1234, "some_file") // *os.File | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.FilesAPI.UploadFile(context.Background()).Scene(scene).File(file).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FilesAPI.UploadFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UploadFile`: UploadFileResponse
	fmt.Fprintf(os.Stdout, "Response from `FilesAPI.UploadFile`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUploadFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scene** | [**FileScene**](FileScene.md) |  | 
 **file** | ***os.File** |  | 

### Return type

[**UploadFileResponse**](UploadFileResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

