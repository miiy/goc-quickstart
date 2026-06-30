# UpdateUserRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**User** | [**UserInput**](UserInput.md) |  | 
**UpdateFields** | Pointer to **[]string** | Mutable user fields to update. | [optional] 

## Methods

### NewUpdateUserRequest

`func NewUpdateUserRequest(user UserInput, ) *UpdateUserRequest`

NewUpdateUserRequest instantiates a new UpdateUserRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateUserRequestWithDefaults

`func NewUpdateUserRequestWithDefaults() *UpdateUserRequest`

NewUpdateUserRequestWithDefaults instantiates a new UpdateUserRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUser

`func (o *UpdateUserRequest) GetUser() UserInput`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *UpdateUserRequest) GetUserOk() (*UserInput, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *UpdateUserRequest) SetUser(v UserInput)`

SetUser sets User field to given value.


### GetUpdateFields

`func (o *UpdateUserRequest) GetUpdateFields() []string`

GetUpdateFields returns the UpdateFields field if non-nil, zero value otherwise.

### GetUpdateFieldsOk

`func (o *UpdateUserRequest) GetUpdateFieldsOk() (*[]string, bool)`

GetUpdateFieldsOk returns a tuple with the UpdateFields field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdateFields

`func (o *UpdateUserRequest) SetUpdateFields(v []string)`

SetUpdateFields sets UpdateFields field to given value.

### HasUpdateFields

`func (o *UpdateUserRequest) HasUpdateFields() bool`

HasUpdateFields returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


