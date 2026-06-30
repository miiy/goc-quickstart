# UpdatePostInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | Pointer to **string** |  | [optional] 
**Content** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**PostStatus**](PostStatus.md) |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**CategoryId** | Pointer to **int64** |  | [optional] 
**CoverUrl** | Pointer to **string** |  | [optional] 

## Methods

### NewUpdatePostInput

`func NewUpdatePostInput() *UpdatePostInput`

NewUpdatePostInput instantiates a new UpdatePostInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdatePostInputWithDefaults

`func NewUpdatePostInputWithDefaults() *UpdatePostInput`

NewUpdatePostInputWithDefaults instantiates a new UpdatePostInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *UpdatePostInput) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *UpdatePostInput) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *UpdatePostInput) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *UpdatePostInput) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetContent

`func (o *UpdatePostInput) GetContent() string`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *UpdatePostInput) GetContentOk() (*string, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *UpdatePostInput) SetContent(v string)`

SetContent sets Content field to given value.

### HasContent

`func (o *UpdatePostInput) HasContent() bool`

HasContent returns a boolean if a field has been set.

### GetStatus

`func (o *UpdatePostInput) GetStatus() PostStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdatePostInput) GetStatusOk() (*PostStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdatePostInput) SetStatus(v PostStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *UpdatePostInput) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTags

`func (o *UpdatePostInput) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *UpdatePostInput) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *UpdatePostInput) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *UpdatePostInput) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetCategoryId

`func (o *UpdatePostInput) GetCategoryId() int64`

GetCategoryId returns the CategoryId field if non-nil, zero value otherwise.

### GetCategoryIdOk

`func (o *UpdatePostInput) GetCategoryIdOk() (*int64, bool)`

GetCategoryIdOk returns a tuple with the CategoryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategoryId

`func (o *UpdatePostInput) SetCategoryId(v int64)`

SetCategoryId sets CategoryId field to given value.

### HasCategoryId

`func (o *UpdatePostInput) HasCategoryId() bool`

HasCategoryId returns a boolean if a field has been set.

### GetCoverUrl

`func (o *UpdatePostInput) GetCoverUrl() string`

GetCoverUrl returns the CoverUrl field if non-nil, zero value otherwise.

### GetCoverUrlOk

`func (o *UpdatePostInput) GetCoverUrlOk() (*string, bool)`

GetCoverUrlOk returns a tuple with the CoverUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCoverUrl

`func (o *UpdatePostInput) SetCoverUrl(v string)`

SetCoverUrl sets CoverUrl field to given value.

### HasCoverUrl

`func (o *UpdatePostInput) HasCoverUrl() bool`

HasCoverUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


