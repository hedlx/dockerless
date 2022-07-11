# Lambda

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Docker** | [**Docker**](Docker.md) |  | 
**Id** | **string** |  | 
**Name** | **string** |  | 
**CreatedAt** | **int64** |  | 
**UpdatedAt** | **int64** |  | 
**Runtime** | **string** |  | 
**LambdaType** | **string** |  | 

## Methods

### NewLambda

`func NewLambda(docker Docker, id string, name string, createdAt int64, updatedAt int64, runtime string, lambdaType string, ) *Lambda`

NewLambda instantiates a new Lambda object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLambdaWithDefaults

`func NewLambdaWithDefaults() *Lambda`

NewLambdaWithDefaults instantiates a new Lambda object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDocker

`func (o *Lambda) GetDocker() Docker`

GetDocker returns the Docker field if non-nil, zero value otherwise.

### GetDockerOk

`func (o *Lambda) GetDockerOk() (*Docker, bool)`

GetDockerOk returns a tuple with the Docker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDocker

`func (o *Lambda) SetDocker(v Docker)`

SetDocker sets Docker field to given value.


### GetId

`func (o *Lambda) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Lambda) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Lambda) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *Lambda) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Lambda) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Lambda) SetName(v string)`

SetName sets Name field to given value.


### GetCreatedAt

`func (o *Lambda) GetCreatedAt() int64`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Lambda) GetCreatedAtOk() (*int64, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Lambda) SetCreatedAt(v int64)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *Lambda) GetUpdatedAt() int64`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Lambda) GetUpdatedAtOk() (*int64, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Lambda) SetUpdatedAt(v int64)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetRuntime

`func (o *Lambda) GetRuntime() string`

GetRuntime returns the Runtime field if non-nil, zero value otherwise.

### GetRuntimeOk

`func (o *Lambda) GetRuntimeOk() (*string, bool)`

GetRuntimeOk returns a tuple with the Runtime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuntime

`func (o *Lambda) SetRuntime(v string)`

SetRuntime sets Runtime field to given value.


### GetLambdaType

`func (o *Lambda) GetLambdaType() string`

GetLambdaType returns the LambdaType field if non-nil, zero value otherwise.

### GetLambdaTypeOk

`func (o *Lambda) GetLambdaTypeOk() (*string, bool)`

GetLambdaTypeOk returns a tuple with the LambdaType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLambdaType

`func (o *Lambda) SetLambdaType(v string)`

SetLambdaType sets LambdaType field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


