# TokenPaginationRequestParameters

## Properties

| Name          | Type                  | Description                                                                                                                                                                                           | Notes                       |
| ------------- | --------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| **PageSize**  | Pointer to **int64**  | Items per Page This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [optional] [default to 250] |
| **PageToken** | Pointer to **string** | Next Page Token The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination).                           | [optional] [default to "1"] |

## Methods

### NewTokenPaginationRequestParameters

`func NewTokenPaginationRequestParameters() *TokenPaginationRequestParameters`

NewTokenPaginationRequestParameters instantiates a new
TokenPaginationRequestParameters object This constructor will assign default
values to properties that have it defined, and makes sure properties required by
API are set, but the set of arguments will change when the set of required
properties is changed

### NewTokenPaginationRequestParametersWithDefaults

`func NewTokenPaginationRequestParametersWithDefaults() *TokenPaginationRequestParameters`

NewTokenPaginationRequestParametersWithDefaults instantiates a new
TokenPaginationRequestParameters object This constructor will only assign
default values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetPageSize

`func (o *TokenPaginationRequestParameters) GetPageSize() int64`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *TokenPaginationRequestParameters) GetPageSizeOk() (*int64, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetPageSize

`func (o *TokenPaginationRequestParameters) SetPageSize(v int64)`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *TokenPaginationRequestParameters) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### GetPageToken

`func (o *TokenPaginationRequestParameters) GetPageToken() string`

GetPageToken returns the PageToken field if non-nil, zero value otherwise.

### GetPageTokenOk

`func (o *TokenPaginationRequestParameters) GetPageTokenOk() (*string, bool)`

GetPageTokenOk returns a tuple with the PageToken field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetPageToken

`func (o *TokenPaginationRequestParameters) SetPageToken(v string)`

SetPageToken sets PageToken field to given value.

### HasPageToken

`func (o *TokenPaginationRequestParameters) HasPageToken() bool`

HasPageToken returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
