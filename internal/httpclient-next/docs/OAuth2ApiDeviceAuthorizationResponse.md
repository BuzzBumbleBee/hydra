# OAuth2ApiDeviceAuthorizationResponse

## Properties

| Name                        | Type                  | Description                                                                                                                                                                                | Notes      |
| --------------------------- | --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------- |
| **DeviceCode**              | Pointer to **string** | The device verification code.                                                                                                                                                              | [optional] |
| **ExpiresIn**               | Pointer to **int64**  | The lifetime in seconds of the \&quot;device_code\&quot; and \&quot;user_code\&quot;.                                                                                                      | [optional] |
| **Interval**                | Pointer to **int64**  | The minimum amount of time in seconds that the client SHOULD wait between polling requests to the token endpoint. If no value is provided, clients MUST use 5 as the default.              | [optional] |
| **UserCode**                | Pointer to **string** | The end-user verification code.                                                                                                                                                            | [optional] |
| **VerificationUri**         | Pointer to **string** | The end-user verification URI on the authorization server. The URI should be short and easy to remember as end users will be asked to manually type it into their user agent.              | [optional] |
| **VerificationUriComplete** | Pointer to **string** | A verification URI that includes the \&quot;user_code\&quot; (or other information with the same function as the \&quot;user_code\&quot;), which is designed for non-textual transmission. | [optional] |

## Methods

### NewOAuth2ApiDeviceAuthorizationResponse

`func NewOAuth2ApiDeviceAuthorizationResponse() *OAuth2ApiDeviceAuthorizationResponse`

NewOAuth2ApiDeviceAuthorizationResponse instantiates a new
OAuth2ApiDeviceAuthorizationResponse object This constructor will assign default
values to properties that have it defined, and makes sure properties required by
API are set, but the set of arguments will change when the set of required
properties is changed

### NewOAuth2ApiDeviceAuthorizationResponseWithDefaults

`func NewOAuth2ApiDeviceAuthorizationResponseWithDefaults() *OAuth2ApiDeviceAuthorizationResponse`

NewOAuth2ApiDeviceAuthorizationResponseWithDefaults instantiates a new
OAuth2ApiDeviceAuthorizationResponse object This constructor will only assign
default values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetDeviceCode

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetDeviceCode() string`

GetDeviceCode returns the DeviceCode field if non-nil, zero value otherwise.

### GetDeviceCodeOk

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetDeviceCodeOk() (*string, bool)`

GetDeviceCodeOk returns a tuple with the DeviceCode field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetDeviceCode

`func (o *OAuth2ApiDeviceAuthorizationResponse) SetDeviceCode(v string)`

SetDeviceCode sets DeviceCode field to given value.

### HasDeviceCode

`func (o *OAuth2ApiDeviceAuthorizationResponse) HasDeviceCode() bool`

HasDeviceCode returns a boolean if a field has been set.

### GetExpiresIn

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetExpiresIn() int64`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetExpiresInOk() (*int64, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *OAuth2ApiDeviceAuthorizationResponse) SetExpiresIn(v int64)`

SetExpiresIn sets ExpiresIn field to given value.

### HasExpiresIn

`func (o *OAuth2ApiDeviceAuthorizationResponse) HasExpiresIn() bool`

HasExpiresIn returns a boolean if a field has been set.

### GetInterval

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetInterval() int64`

GetInterval returns the Interval field if non-nil, zero value otherwise.

### GetIntervalOk

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetIntervalOk() (*int64, bool)`

GetIntervalOk returns a tuple with the Interval field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetInterval

`func (o *OAuth2ApiDeviceAuthorizationResponse) SetInterval(v int64)`

SetInterval sets Interval field to given value.

### HasInterval

`func (o *OAuth2ApiDeviceAuthorizationResponse) HasInterval() bool`

HasInterval returns a boolean if a field has been set.

### GetUserCode

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetUserCode() string`

GetUserCode returns the UserCode field if non-nil, zero value otherwise.

### GetUserCodeOk

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetUserCodeOk() (*string, bool)`

GetUserCodeOk returns a tuple with the UserCode field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetUserCode

`func (o *OAuth2ApiDeviceAuthorizationResponse) SetUserCode(v string)`

SetUserCode sets UserCode field to given value.

### HasUserCode

`func (o *OAuth2ApiDeviceAuthorizationResponse) HasUserCode() bool`

HasUserCode returns a boolean if a field has been set.

### GetVerificationUri

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetVerificationUri() string`

GetVerificationUri returns the VerificationUri field if non-nil, zero value
otherwise.

### GetVerificationUriOk

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetVerificationUriOk() (*string, bool)`

GetVerificationUriOk returns a tuple with the VerificationUri field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetVerificationUri

`func (o *OAuth2ApiDeviceAuthorizationResponse) SetVerificationUri(v string)`

SetVerificationUri sets VerificationUri field to given value.

### HasVerificationUri

`func (o *OAuth2ApiDeviceAuthorizationResponse) HasVerificationUri() bool`

HasVerificationUri returns a boolean if a field has been set.

### GetVerificationUriComplete

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetVerificationUriComplete() string`

GetVerificationUriComplete returns the VerificationUriComplete field if non-nil,
zero value otherwise.

### GetVerificationUriCompleteOk

`func (o *OAuth2ApiDeviceAuthorizationResponse) GetVerificationUriCompleteOk() (*string, bool)`

GetVerificationUriCompleteOk returns a tuple with the VerificationUriComplete
field if it's non-nil, zero value otherwise and a boolean to check if the value
has been set.

### SetVerificationUriComplete

`func (o *OAuth2ApiDeviceAuthorizationResponse) SetVerificationUriComplete(v string)`

SetVerificationUriComplete sets VerificationUriComplete field to given value.

### HasVerificationUriComplete

`func (o *OAuth2ApiDeviceAuthorizationResponse) HasVerificationUriComplete() bool`

HasVerificationUriComplete returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
