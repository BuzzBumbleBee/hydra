/*
 * Ory Hydra API
 *
 * Documentation for all of Ory Hydra's APIs.
 *
 * API version: 1.0.0
 * Contact: hi@ory.sh
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Linger please
var (
	_ context.Context
)

type V0alpha2Api interface {

	/*
			 * PerformOAuth2DeviceFlow The OAuth 2.0 Device Authorize Endpoint
			 * This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows.
		OAuth2 is a very popular protocol and a library for your programming language will exists.

		To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc8628
			 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			 * @return V0alpha2ApiApiPerformOAuth2DeviceFlowRequest
	*/
	PerformOAuth2DeviceFlow(ctx context.Context) V0alpha2ApiApiPerformOAuth2DeviceFlowRequest

	/*
	 * PerformOAuth2DeviceFlowExecute executes the request
	 * @return OAuth2ApiDeviceAuthorizationResponse
	 */
	PerformOAuth2DeviceFlowExecute(r V0alpha2ApiApiPerformOAuth2DeviceFlowRequest) (*OAuth2ApiDeviceAuthorizationResponse, *http.Response, error)
}

// V0alpha2ApiService V0alpha2Api service
type V0alpha2ApiService service

type V0alpha2ApiApiPerformOAuth2DeviceFlowRequest struct {
	ctx        context.Context
	ApiService V0alpha2Api
}

func (r V0alpha2ApiApiPerformOAuth2DeviceFlowRequest) Execute() (*OAuth2ApiDeviceAuthorizationResponse, *http.Response, error) {
	return r.ApiService.PerformOAuth2DeviceFlowExecute(r)
}

/*
 * PerformOAuth2DeviceFlow The OAuth 2.0 Device Authorize Endpoint
 * This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows.
OAuth2 is a very popular protocol and a library for your programming language will exists.

To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc8628
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @return V0alpha2ApiApiPerformOAuth2DeviceFlowRequest
*/
func (a *V0alpha2ApiService) PerformOAuth2DeviceFlow(ctx context.Context) V0alpha2ApiApiPerformOAuth2DeviceFlowRequest {
	return V0alpha2ApiApiPerformOAuth2DeviceFlowRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

/*
 * Execute executes the request
 * @return OAuth2ApiDeviceAuthorizationResponse
 */
func (a *V0alpha2ApiService) PerformOAuth2DeviceFlowExecute(r V0alpha2ApiApiPerformOAuth2DeviceFlowRequest) (*OAuth2ApiDeviceAuthorizationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  *OAuth2ApiDeviceAuthorizationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.PerformOAuth2DeviceFlow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth2/device/auth"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v JsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v JsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
