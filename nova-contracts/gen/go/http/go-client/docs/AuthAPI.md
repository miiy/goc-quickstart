# \AuthAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ChangePassword**](AuthAPI.md#ChangePassword) | **Put** /auth/password | Change password
[**EmailCheck**](AuthAPI.md#EmailCheck) | **Post** /auth/register/check-email | Check email availability
[**Login**](AuthAPI.md#Login) | **Post** /auth/login | Login with username or email
[**Logout**](AuthAPI.md#Logout) | **Post** /auth/logout | Logout
[**MpLogin**](AuthAPI.md#MpLogin) | **Post** /auth/mp/login | Login with WeChat Mini Program code
[**PhoneAuth**](AuthAPI.md#PhoneAuth) | **Post** /auth/phone/login | Login with phone and SMS code
[**PhoneCheck**](AuthAPI.md#PhoneCheck) | **Post** /auth/register/check-phone | Check phone availability
[**RefreshToken**](AuthAPI.md#RefreshToken) | **Post** /auth/token/refresh | Refresh access token
[**Register**](AuthAPI.md#Register) | **Post** /auth/register | Register a user
[**SendSmsCode**](AuthAPI.md#SendSmsCode) | **Post** /auth/sms/send-code | Send SMS code
[**UsernameCheck**](AuthAPI.md#UsernameCheck) | **Post** /auth/register/check-username | Check username availability



## ChangePassword

> map[string]interface{} ChangePassword(ctx).ChangePasswordRequest(changePasswordRequest).Execute()

Change password

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	changePasswordRequest := *openapiclient.NewChangePasswordRequest("OldPassword_example", "NewPassword_example", "NewPasswordConfirmation_example") // ChangePasswordRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.ChangePassword(context.Background()).ChangePasswordRequest(changePasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.ChangePassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ChangePassword`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.ChangePassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiChangePasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changePasswordRequest** | [**ChangePasswordRequest**](ChangePasswordRequest.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EmailCheck

> EmailCheckResponse EmailCheck(ctx).EmailCheckRequest(emailCheckRequest).Execute()

Check email availability

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	emailCheckRequest := *openapiclient.NewEmailCheckRequest("Value_example") // EmailCheckRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.EmailCheck(context.Background()).EmailCheckRequest(emailCheckRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.EmailCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `EmailCheck`: EmailCheckResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.EmailCheck`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEmailCheckRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **emailCheckRequest** | [**EmailCheckRequest**](EmailCheckRequest.md) |  | 

### Return type

[**EmailCheckResponse**](EmailCheckResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Login

> TokenResponse Login(ctx).LoginRequest(loginRequest).Execute()

Login with username or email

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	loginRequest := *openapiclient.NewLoginRequest("Username_example", "Password_example") // LoginRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.Login(context.Background()).LoginRequest(loginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.Login``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Login`: TokenResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.Login`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLoginRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginRequest** | [**LoginRequest**](LoginRequest.md) |  | 

### Return type

[**TokenResponse**](TokenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Logout

> map[string]interface{} Logout(ctx).LogoutRequest(logoutRequest).Execute()

Logout

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	logoutRequest := *openapiclient.NewLogoutRequest() // LogoutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.Logout(context.Background()).LogoutRequest(logoutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.Logout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Logout`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.Logout`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLogoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutRequest** | [**LogoutRequest**](LogoutRequest.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MpLogin

> TokenResponse MpLogin(ctx).MpLoginRequest(mpLoginRequest).Execute()

Login with WeChat Mini Program code

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	mpLoginRequest := *openapiclient.NewMpLoginRequest("Code_example") // MpLoginRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.MpLogin(context.Background()).MpLoginRequest(mpLoginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.MpLogin``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `MpLogin`: TokenResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.MpLogin`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMpLoginRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mpLoginRequest** | [**MpLoginRequest**](MpLoginRequest.md) |  | 

### Return type

[**TokenResponse**](TokenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PhoneAuth

> TokenResponse PhoneAuth(ctx).PhoneAuthRequest(phoneAuthRequest).Execute()

Login with phone and SMS code

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	phoneAuthRequest := *openapiclient.NewPhoneAuthRequest("Phone_example", "Code_example") // PhoneAuthRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.PhoneAuth(context.Background()).PhoneAuthRequest(phoneAuthRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.PhoneAuth``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PhoneAuth`: TokenResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.PhoneAuth`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPhoneAuthRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **phoneAuthRequest** | [**PhoneAuthRequest**](PhoneAuthRequest.md) |  | 

### Return type

[**TokenResponse**](TokenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PhoneCheck

> PhoneCheckResponse PhoneCheck(ctx).PhoneCheckRequest(phoneCheckRequest).Execute()

Check phone availability

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	phoneCheckRequest := *openapiclient.NewPhoneCheckRequest("Value_example") // PhoneCheckRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.PhoneCheck(context.Background()).PhoneCheckRequest(phoneCheckRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.PhoneCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PhoneCheck`: PhoneCheckResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.PhoneCheck`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPhoneCheckRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **phoneCheckRequest** | [**PhoneCheckRequest**](PhoneCheckRequest.md) |  | 

### Return type

[**PhoneCheckResponse**](PhoneCheckResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshToken

> TokenResponse RefreshToken(ctx).RefreshTokenRequest(refreshTokenRequest).Execute()

Refresh access token

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	refreshTokenRequest := *openapiclient.NewRefreshTokenRequest("RefreshToken_example") // RefreshTokenRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.RefreshToken(context.Background()).RefreshTokenRequest(refreshTokenRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.RefreshToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshToken`: TokenResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.RefreshToken`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRefreshTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshTokenRequest** | [**RefreshTokenRequest**](RefreshTokenRequest.md) |  | 

### Return type

[**TokenResponse**](TokenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Register

> RegisterResponse Register(ctx).RegisterRequest(registerRequest).Execute()

Register a user

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	registerRequest := *openapiclient.NewRegisterRequest("Email_example", "Username_example", "Password_example", "PasswordConfirmation_example") // RegisterRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.Register(context.Background()).RegisterRequest(registerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.Register``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Register`: RegisterResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.Register`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRegisterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registerRequest** | [**RegisterRequest**](RegisterRequest.md) |  | 

### Return type

[**RegisterResponse**](RegisterResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SendSmsCode

> map[string]interface{} SendSmsCode(ctx).SendSmsCodeRequest(sendSmsCodeRequest).Execute()

Send SMS code

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	sendSmsCodeRequest := *openapiclient.NewSendSmsCodeRequest("Phone_example") // SendSmsCodeRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.SendSmsCode(context.Background()).SendSmsCodeRequest(sendSmsCodeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.SendSmsCode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SendSmsCode`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.SendSmsCode`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSendSmsCodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sendSmsCodeRequest** | [**SendSmsCodeRequest**](SendSmsCodeRequest.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UsernameCheck

> UsernameCheckResponse UsernameCheck(ctx).UsernameCheckRequest(usernameCheckRequest).Execute()

Check username availability

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

func main() {
	usernameCheckRequest := *openapiclient.NewUsernameCheckRequest("Value_example") // UsernameCheckRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthAPI.UsernameCheck(context.Background()).UsernameCheckRequest(usernameCheckRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthAPI.UsernameCheck``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UsernameCheck`: UsernameCheckResponse
	fmt.Fprintf(os.Stdout, "Response from `AuthAPI.UsernameCheck`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUsernameCheckRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **usernameCheckRequest** | [**UsernameCheckRequest**](UsernameCheckRequest.md) |  | 

### Return type

[**UsernameCheckResponse**](UsernameCheckResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

