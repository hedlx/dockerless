# \RuntimeApi

All URIs are relative to *https://virtserver.swaggerhub.com/hedlx/doless/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateRuntime**](RuntimeApi.md#CreateRuntime) | **Post** /runtime | Create runtime
[**GetRuntime**](RuntimeApi.md#GetRuntime) | **Get** /runtime/{id} | Get runtime
[**ListRuntimes**](RuntimeApi.md#ListRuntimes) | **Get** /runtime | List runtimes



## CreateRuntime

> Runtime CreateRuntime(ctx).CreateRuntime(createRuntime).Execute()

Create runtime

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    createRuntime := *openapiclient.NewCreateRuntime("Dockerfile_example") // CreateRuntime | Create runtime body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RuntimeApi.CreateRuntime(context.Background()).CreateRuntime(createRuntime).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RuntimeApi.CreateRuntime``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateRuntime`: Runtime
    fmt.Fprintf(os.Stdout, "Response from `RuntimeApi.CreateRuntime`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateRuntimeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createRuntime** | [**CreateRuntime**](CreateRuntime.md) | Create runtime body | 

### Return type

[**Runtime**](Runtime.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRuntime

> Runtime GetRuntime(ctx, id).Execute()

Get runtime

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "id_example" // string | Runtime id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RuntimeApi.GetRuntime(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RuntimeApi.GetRuntime``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRuntime`: Runtime
    fmt.Fprintf(os.Stdout, "Response from `RuntimeApi.GetRuntime`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runtime id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRuntimeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Runtime**](Runtime.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRuntimes

> []Runtime ListRuntimes(ctx).Execute()

List runtimes

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RuntimeApi.ListRuntimes(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RuntimeApi.ListRuntimes``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRuntimes`: []Runtime
    fmt.Fprintf(os.Stdout, "Response from `RuntimeApi.ListRuntimes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListRuntimesRequest struct via the builder pattern


### Return type

[**[]Runtime**](Runtime.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

