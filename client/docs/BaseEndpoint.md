# BaseEndpoint

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Path** | **string** |  | 
**Lambda** | **string** |  | 

## Methods

### NewBaseEndpoint

`func NewBaseEndpoint(name string, path string, lambda string, ) *BaseEndpoint`

NewBaseEndpoint instantiates a new BaseEndpoint object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBaseEndpointWithDefaults

`func NewBaseEndpointWithDefaults() *BaseEndpoint`

NewBaseEndpointWithDefaults instantiates a new BaseEndpoint object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BaseEndpoint) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BaseEndpoint) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BaseEndpoint) SetName(v string)`

SetName sets Name field to given value.


### GetPath

`func (o *BaseEndpoint) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *BaseEndpoint) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *BaseEndpoint) SetPath(v string)`

SetPath sets Path field to given value.


### GetLambda

`func (o *BaseEndpoint) GetLambda() string`

GetLambda returns the Lambda field if non-nil, zero value otherwise.

### GetLambdaOk

`func (o *BaseEndpoint) GetLambdaOk() (*string, bool)`

GetLambdaOk returns a tuple with the Lambda field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLambda

`func (o *BaseEndpoint) SetLambda(v string)`

SetLambda sets Lambda field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


