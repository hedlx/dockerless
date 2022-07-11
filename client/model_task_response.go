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

// TaskResponse struct for TaskResponse
type TaskResponse struct {
	Task string `json:"task"`
}

// NewTaskResponse instantiates a new TaskResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskResponse(task string) *TaskResponse {
	this := TaskResponse{}
	this.Task = task
	return &this
}

// NewTaskResponseWithDefaults instantiates a new TaskResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskResponseWithDefaults() *TaskResponse {
	this := TaskResponse{}
	return &this
}

// GetTask returns the Task field value
func (o *TaskResponse) GetTask() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Task
}

// GetTaskOk returns a tuple with the Task field value
// and a boolean to check if the value has been set.
func (o *TaskResponse) GetTaskOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Task, true
}

// SetTask sets field value
func (o *TaskResponse) SetTask(v string) {
	o.Task = v
}

func (o TaskResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["task"] = o.Task
	}
	return json.Marshal(toSerialize)
}

type NullableTaskResponse struct {
	value *TaskResponse
	isSet bool
}

func (v NullableTaskResponse) Get() *TaskResponse {
	return v.value
}

func (v *NullableTaskResponse) Set(val *TaskResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskResponse(val *TaskResponse) *NullableTaskResponse {
	return &NullableTaskResponse{value: val, isSet: true}
}

func (v NullableTaskResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


