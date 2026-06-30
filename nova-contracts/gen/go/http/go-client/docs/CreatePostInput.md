# CreatePostInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** |  | 
**Content** | **string** |  | 
**Status** | Pointer to [**PostStatus**](PostStatus.md) |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**CategoryId** | Pointer to **int64** |  | [optional] 
**CoverUrl** | Pointer to **string** |  | [optional] 

## Methods

### NewCreatePostInput

`func NewCreatePostInput(title string, content string, ) *CreatePostInput`

NewCreatePostInput instantiates a new CreatePostInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePostInputWithDefaults

`func NewCreatePostInputWithDefaults() *CreatePostInput`

NewCreatePostInputWithDefaults instantiates a new CreatePostInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *CreatePostInput) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *CreatePostInput) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *CreatePostInput) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetContent

`func (o *CreatePostInput) GetContent() string`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *CreatePostInput) GetContentOk() (*string, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *CreatePostInput) SetContent(v string)`

SetContent sets Content field to given value.


### GetStatus

`func (o *CreatePostInput) GetStatus() PostStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CreatePostInput) GetStatusOk() (*PostStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CreatePostInput) SetStatus(v PostStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CreatePostInput) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTags

`func (o *CreatePostInput) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *CreatePostInput) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *CreatePostInput) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *CreatePostInput) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetCategoryId

`func (o *CreatePostInput) GetCategoryId() int64`

GetCategoryId returns the CategoryId field if non-nil, zero value otherwise.

### GetCategoryIdOk

`func (o *CreatePostInput) GetCategoryIdOk() (*int64, bool)`

GetCategoryIdOk returns a tuple with the CategoryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategoryId

`func (o *CreatePostInput) SetCategoryId(v int64)`

SetCategoryId sets CategoryId field to given value.

### HasCategoryId

`func (o *CreatePostInput) HasCategoryId() bool`

HasCategoryId returns a boolean if a field has been set.

### GetCoverUrl

`func (o *CreatePostInput) GetCoverUrl() string`

GetCoverUrl returns the CoverUrl field if non-nil, zero value otherwise.

### GetCoverUrlOk

`func (o *CreatePostInput) GetCoverUrlOk() (*string, bool)`

GetCoverUrlOk returns a tuple with the CoverUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCoverUrl

`func (o *CreatePostInput) SetCoverUrl(v string)`

SetCoverUrl sets CoverUrl field to given value.

### HasCoverUrl

`func (o *CreatePostInput) HasCoverUrl() bool`

HasCoverUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


