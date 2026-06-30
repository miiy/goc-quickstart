# UserInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Nickname** | Pointer to **string** |  | [optional] 
**Avatar** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**Phone** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**UserStatus**](UserStatus.md) |  | [optional] 

## Methods

### NewUserInput

`func NewUserInput() *UserInput`

NewUserInput instantiates a new UserInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserInputWithDefaults

`func NewUserInputWithDefaults() *UserInput`

NewUserInputWithDefaults instantiates a new UserInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNickname

`func (o *UserInput) GetNickname() string`

GetNickname returns the Nickname field if non-nil, zero value otherwise.

### GetNicknameOk

`func (o *UserInput) GetNicknameOk() (*string, bool)`

GetNicknameOk returns a tuple with the Nickname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNickname

`func (o *UserInput) SetNickname(v string)`

SetNickname sets Nickname field to given value.

### HasNickname

`func (o *UserInput) HasNickname() bool`

HasNickname returns a boolean if a field has been set.

### GetAvatar

`func (o *UserInput) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *UserInput) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *UserInput) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.

### HasAvatar

`func (o *UserInput) HasAvatar() bool`

HasAvatar returns a boolean if a field has been set.

### GetEmail

`func (o *UserInput) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *UserInput) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *UserInput) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *UserInput) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetPhone

`func (o *UserInput) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *UserInput) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *UserInput) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *UserInput) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetStatus

`func (o *UserInput) GetStatus() UserStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UserInput) GetStatusOk() (*UserStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UserInput) SetStatus(v UserStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *UserInput) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


