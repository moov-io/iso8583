# \Iso8583MessageApi

All URIs are relative to *https://local.moov.io:8208*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Convert**](Iso8583MessageApi.md#Convert) | **Post** /convert | Convert iso8583 message
[**Health**](Iso8583MessageApi.md#Health) | **Get** /health | health iso8583 service
[**Print**](Iso8583MessageApi.md#Print) | **Post** /print | Print iso8583 message with specific format
[**Validator**](Iso8583MessageApi.md#Validator) | **Post** /validator | Validate iso8583 message



## Convert

> *os.File Convert(ctx, optional)

Convert iso8583 message

Convert from original iso8583 message to new iso8583 message

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ConvertOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ConvertOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **format** | **optional.String**| converting message type | [default to json]
 **input** | **optional.Interface of *os.File****optional.*os.File**| iso8583 message | 
 **spec** | **optional.Interface of *os.File****optional.*os.File**| message configuration file | 

### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/octet-stream, application/json, application/xml

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Health

> Success Health(ctx, )

health iso8583 service

Check the iso8583 service to check if running

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**Success**](Success.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Print

> *os.File Print(ctx, optional)

Print iso8583 message with specific format

Print iso8583 message with requested format.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***PrintOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a PrintOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **format** | **optional.String**| print iso8583 type | [default to json]
 **input** | **optional.Interface of *os.File****optional.*os.File**| iso8583 message | 
 **spec** | **optional.Interface of *os.File****optional.*os.File**| message configuration file | 

### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/octet-stream, application/json, application/xml

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Validator

> Success Validator(ctx, optional)

Validate iso8583 message

Validation iso8583 message.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ValidatorOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ValidatorOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **input** | **optional.Interface of *os.File****optional.*os.File**| iso8583 message | 
 **spec** | **optional.Interface of *os.File****optional.*os.File**| message configuration file | 

### Return type

[**Success**](Success.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

