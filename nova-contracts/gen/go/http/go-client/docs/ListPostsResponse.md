# ListPostsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Total** | **int64** |  | 
**TotalPages** | **int32** |  | 
**PageSize** | **int32** |  | 
**CurrentPage** | **int32** |  | 
**Posts** | [**[]Post**](Post.md) |  | 

## Methods

### NewListPostsResponse

`func NewListPostsResponse(total int64, totalPages int32, pageSize int32, currentPage int32, posts []Post, ) *ListPostsResponse`

NewListPostsResponse instantiates a new ListPostsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListPostsResponseWithDefaults

`func NewListPostsResponseWithDefaults() *ListPostsResponse`

NewListPostsResponseWithDefaults instantiates a new ListPostsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTotal

`func (o *ListPostsResponse) GetTotal() int64`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *ListPostsResponse) GetTotalOk() (*int64, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *ListPostsResponse) SetTotal(v int64)`

SetTotal sets Total field to given value.


### GetTotalPages

`func (o *ListPostsResponse) GetTotalPages() int32`

GetTotalPages returns the TotalPages field if non-nil, zero value otherwise.

### GetTotalPagesOk

`func (o *ListPostsResponse) GetTotalPagesOk() (*int32, bool)`

GetTotalPagesOk returns a tuple with the TotalPages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalPages

`func (o *ListPostsResponse) SetTotalPages(v int32)`

SetTotalPages sets TotalPages field to given value.


### GetPageSize

`func (o *ListPostsResponse) GetPageSize() int32`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *ListPostsResponse) GetPageSizeOk() (*int32, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *ListPostsResponse) SetPageSize(v int32)`

SetPageSize sets PageSize field to given value.


### GetCurrentPage

`func (o *ListPostsResponse) GetCurrentPage() int32`

GetCurrentPage returns the CurrentPage field if non-nil, zero value otherwise.

### GetCurrentPageOk

`func (o *ListPostsResponse) GetCurrentPageOk() (*int32, bool)`

GetCurrentPageOk returns a tuple with the CurrentPage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentPage

`func (o *ListPostsResponse) SetCurrentPage(v int32)`

SetCurrentPage sets CurrentPage field to given value.


### GetPosts

`func (o *ListPostsResponse) GetPosts() []Post`

GetPosts returns the Posts field if non-nil, zero value otherwise.

### GetPostsOk

`func (o *ListPostsResponse) GetPostsOk() (*[]Post, bool)`

GetPostsOk returns a tuple with the Posts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosts

`func (o *ListPostsResponse) SetPosts(v []Post)`

SetPosts sets Posts field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


