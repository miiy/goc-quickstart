# File

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** |  | 
**OwnerId** | **int64** |  | 
**OwnerType** | **string** |  | 
**Scene** | [**FileScene**](FileScene.md) |  | 
**ObjectKey** | **string** | Stable relative storage key, for example avatars/2026/06/a.png. | 
**Url** | **string** |  | 
**MimeType** | **string** |  | 
**Size** | **int64** |  | 
**Checksum** | **string** |  | 
**Status** | [**FileStatus**](FileStatus.md) |  | 
**CreatedBy** | **int64** |  | 
**CreatedAt** | **time.Time** |  | 
**UpdatedAt** | **time.Time** |  | 
**DeletedAt** | Pointer to **NullableTime** |  | [optional] 

## Methods

### NewFile

`func NewFile(id int64, ownerId int64, ownerType string, scene FileScene, objectKey string, url string, mimeType string, size int64, checksum string, status FileStatus, createdBy int64, createdAt time.Time, updatedAt time.Time, ) *File`

NewFile instantiates a new File object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFileWithDefaults

`func NewFileWithDefaults() *File`

NewFileWithDefaults instantiates a new File object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *File) GetId() int64`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *File) GetIdOk() (*int64, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *File) SetId(v int64)`

SetId sets Id field to given value.


### GetOwnerId

`func (o *File) GetOwnerId() int64`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *File) GetOwnerIdOk() (*int64, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *File) SetOwnerId(v int64)`

SetOwnerId sets OwnerId field to given value.


### GetOwnerType

`func (o *File) GetOwnerType() string`

GetOwnerType returns the OwnerType field if non-nil, zero value otherwise.

### GetOwnerTypeOk

`func (o *File) GetOwnerTypeOk() (*string, bool)`

GetOwnerTypeOk returns a tuple with the OwnerType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerType

`func (o *File) SetOwnerType(v string)`

SetOwnerType sets OwnerType field to given value.


### GetScene

`func (o *File) GetScene() FileScene`

GetScene returns the Scene field if non-nil, zero value otherwise.

### GetSceneOk

`func (o *File) GetSceneOk() (*FileScene, bool)`

GetSceneOk returns a tuple with the Scene field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScene

`func (o *File) SetScene(v FileScene)`

SetScene sets Scene field to given value.


### GetObjectKey

`func (o *File) GetObjectKey() string`

GetObjectKey returns the ObjectKey field if non-nil, zero value otherwise.

### GetObjectKeyOk

`func (o *File) GetObjectKeyOk() (*string, bool)`

GetObjectKeyOk returns a tuple with the ObjectKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectKey

`func (o *File) SetObjectKey(v string)`

SetObjectKey sets ObjectKey field to given value.


### GetUrl

`func (o *File) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *File) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *File) SetUrl(v string)`

SetUrl sets Url field to given value.


### GetMimeType

`func (o *File) GetMimeType() string`

GetMimeType returns the MimeType field if non-nil, zero value otherwise.

### GetMimeTypeOk

`func (o *File) GetMimeTypeOk() (*string, bool)`

GetMimeTypeOk returns a tuple with the MimeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMimeType

`func (o *File) SetMimeType(v string)`

SetMimeType sets MimeType field to given value.


### GetSize

`func (o *File) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *File) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *File) SetSize(v int64)`

SetSize sets Size field to given value.


### GetChecksum

`func (o *File) GetChecksum() string`

GetChecksum returns the Checksum field if non-nil, zero value otherwise.

### GetChecksumOk

`func (o *File) GetChecksumOk() (*string, bool)`

GetChecksumOk returns a tuple with the Checksum field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChecksum

`func (o *File) SetChecksum(v string)`

SetChecksum sets Checksum field to given value.


### GetStatus

`func (o *File) GetStatus() FileStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *File) GetStatusOk() (*FileStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *File) SetStatus(v FileStatus)`

SetStatus sets Status field to given value.


### GetCreatedBy

`func (o *File) GetCreatedBy() int64`

GetCreatedBy returns the CreatedBy field if non-nil, zero value otherwise.

### GetCreatedByOk

`func (o *File) GetCreatedByOk() (*int64, bool)`

GetCreatedByOk returns a tuple with the CreatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedBy

`func (o *File) SetCreatedBy(v int64)`

SetCreatedBy sets CreatedBy field to given value.


### GetCreatedAt

`func (o *File) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *File) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *File) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *File) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *File) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *File) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetDeletedAt

`func (o *File) GetDeletedAt() time.Time`

GetDeletedAt returns the DeletedAt field if non-nil, zero value otherwise.

### GetDeletedAtOk

`func (o *File) GetDeletedAtOk() (*time.Time, bool)`

GetDeletedAtOk returns a tuple with the DeletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeletedAt

`func (o *File) SetDeletedAt(v time.Time)`

SetDeletedAt sets DeletedAt field to given value.

### HasDeletedAt

`func (o *File) HasDeletedAt() bool`

HasDeletedAt returns a boolean if a field has been set.

### SetDeletedAtNil

`func (o *File) SetDeletedAtNil(b bool)`

 SetDeletedAtNil sets the value for DeletedAt to be an explicit nil

### UnsetDeletedAt
`func (o *File) UnsetDeletedAt()`

UnsetDeletedAt ensures that no value is present for DeletedAt, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


