# CreateRuntime

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Dockerfile** | **string** |  | 
**Name** | Pointer to **string** |  | [optional] 

## Methods

### NewCreateRuntime

`func NewCreateRuntime(dockerfile string, ) *CreateRuntime`

NewCreateRuntime instantiates a new CreateRuntime object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateRuntimeWithDefaults

`func NewCreateRuntimeWithDefaults() *CreateRuntime`

NewCreateRuntimeWithDefaults instantiates a new CreateRuntime object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDockerfile

`func (o *CreateRuntime) GetDockerfile() string`

GetDockerfile returns the Dockerfile field if non-nil, zero value otherwise.

### GetDockerfileOk

`func (o *CreateRuntime) GetDockerfileOk() (*string, bool)`

GetDockerfileOk returns a tuple with the Dockerfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDockerfile

`func (o *CreateRuntime) SetDockerfile(v string)`

SetDockerfile sets Dockerfile field to given value.


### GetName

`func (o *CreateRuntime) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateRuntime) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateRuntime) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CreateRuntime) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


