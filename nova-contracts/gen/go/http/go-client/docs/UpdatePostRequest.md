# UpdatePostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Post** | [**UpdatePostInput**](UpdatePostInput.md) |  | 
**UpdateFields** | Pointer to **[]string** | Mutable post fields to update. | [optional] 

## Methods

### NewUpdatePostRequest

`func NewUpdatePostRequest(post UpdatePostInput, ) *UpdatePostRequest`

NewUpdatePostRequest instantiates a new UpdatePostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdatePostRequestWithDefaults

`func NewUpdatePostRequestWithDefaults() *UpdatePostRequest`

NewUpdatePostRequestWithDefaults instantiates a new UpdatePostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPost

`func (o *UpdatePostRequest) GetPost() UpdatePostInput`

GetPost returns the Post field if non-nil, zero value otherwise.

### GetPostOk

`func (o *UpdatePostRequest) GetPostOk() (*UpdatePostInput, bool)`

GetPostOk returns a tuple with the Post field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPost

`func (o *UpdatePostRequest) SetPost(v UpdatePostInput)`

SetPost sets Post field to given value.


### GetUpdateFields

`func (o *UpdatePostRequest) GetUpdateFields() []string`

GetUpdateFields returns the UpdateFields field if non-nil, zero value otherwise.

### GetUpdateFieldsOk

`func (o *UpdatePostRequest) GetUpdateFieldsOk() (*[]string, bool)`

GetUpdateFieldsOk returns a tuple with the UpdateFields field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdateFields

`func (o *UpdatePostRequest) SetUpdateFields(v []string)`

SetUpdateFields sets UpdateFields field to given value.

### HasUpdateFields

`func (o *UpdatePostRequest) HasUpdateFields() bool`

HasUpdateFields returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


