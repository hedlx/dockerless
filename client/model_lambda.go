/*
core

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// Lambda struct for Lambda
type Lambda struct {
	Docker Docker `json:"docker"`
	Id string `json:"id"`
	Name string `json:"name"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	Runtime string `json:"runtime"`
	LambdaType string `json:"lambda_type"`
}

// NewLambda instantiates a new Lambda object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLambda(docker Docker, id string, name string, createdAt int64, updatedAt int64, runtime string, lambdaType string) *Lambda {
	this := Lambda{}
	this.Id = id
	this.Name = name
	this.CreatedAt = createdAt
	this.UpdatedAt = updatedAt
	this.Runtime = runtime
	this.LambdaType = lambdaType
	return &this
}

// NewLambdaWithDefaults instantiates a new Lambda object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLambdaWithDefaults() *Lambda {
	this := Lambda{}
	return &this
}

// GetDocker returns the Docker field value
func (o *Lambda) GetDocker() Docker {
	if o == nil {
		var ret Docker
		return ret
	}

	return o.Docker
}

// GetDockerOk returns a tuple with the Docker field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetDockerOk() (*Docker, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Docker, true
}

// SetDocker sets field value
func (o *Lambda) SetDocker(v Docker) {
	o.Docker = v
}

// GetId returns the Id field value
func (o *Lambda) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Lambda) SetId(v string) {
	o.Id = v
}

// GetName returns the Name field value
func (o *Lambda) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Lambda) SetName(v string) {
	o.Name = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *Lambda) GetCreatedAt() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetCreatedAtOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *Lambda) SetCreatedAt(v int64) {
	o.CreatedAt = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *Lambda) GetUpdatedAt() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetUpdatedAtOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *Lambda) SetUpdatedAt(v int64) {
	o.UpdatedAt = v
}

// GetRuntime returns the Runtime field value
func (o *Lambda) GetRuntime() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Runtime
}

// GetRuntimeOk returns a tuple with the Runtime field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetRuntimeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Runtime, true
}

// SetRuntime sets field value
func (o *Lambda) SetRuntime(v string) {
	o.Runtime = v
}

// GetLambdaType returns the LambdaType field value
func (o *Lambda) GetLambdaType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LambdaType
}

// GetLambdaTypeOk returns a tuple with the LambdaType field value
// and a boolean to check if the value has been set.
func (o *Lambda) GetLambdaTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LambdaType, true
}

// SetLambdaType sets field value
func (o *Lambda) SetLambdaType(v string) {
	o.LambdaType = v
}

func (o Lambda) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["docker"] = o.Docker
	}
	if true {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["name"] = o.Name
	}
	if true {
		toSerialize["created_at"] = o.CreatedAt
	}
	if true {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if true {
		toSerialize["runtime"] = o.Runtime
	}
	if true {
		toSerialize["lambda_type"] = o.LambdaType
	}
	return json.Marshal(toSerialize)
}

type NullableLambda struct {
	value *Lambda
	isSet bool
}

func (v NullableLambda) Get() *Lambda {
	return v.value
}

func (v *NullableLambda) Set(val *Lambda) {
	v.value = val
	v.isSet = true
}

func (v NullableLambda) IsSet() bool {
	return v.isSet
}

func (v *NullableLambda) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLambda(val *Lambda) *NullableLambda {
	return &NullableLambda{value: val, isSet: true}
}

func (v NullableLambda) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLambda) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


