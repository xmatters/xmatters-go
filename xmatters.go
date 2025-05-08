// Package xmatters provides a Go client for interacting with the XMatters REST API.
//
// xMatters is a communication platform that enables enterprises to manage and automate communication
// with their employees, customers, and other stakeholders during incidents and other critical events.
// This package allows you to programmatically access and utilize the xMatters REST API to integrate xMatters
// functionality into your Go applications.
//
// API Documentation:
//
//	https://help.xmatters.com/xmapi/
//
// Usage:
//
//	// Create a new XMattersAPI client with your API Token
//	apiToken := "your-api-token"
//	xmattersClient, err := xmatters.NewWithAPIToken(&apiToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use the client to interact with the xMatters REST API
//	// For example, retrieve information about users:
//	users, err := xmattersClient.GetPersonList(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(users)
package xmatters

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/motemen/go-loghttp"
	"golang.org/x/time/rate"
)

const (
	// xMatters_go constants
	defaultBasePath    = "/api/xm/1"
	ContentJSON        = "application/json"
	StatusOK           = 200
	StatusCreated      = 201
	StatusNoContent    = 204
	StatusUnauthorized = 401
)

var (
	Version       string = "1"
	AuthTypeBasic string = "Basic"
	AuthTypeOAuth string = "OAuth"
)

// XMattersAPI represents the configuration options for interacting with the xMatters API.
type XMattersAPI struct {
	Username    *string
	Password    *string
	Token       *string
	AuthType    *string
	BaseURL     *string
	UserAgent   *string
	headers     http.Header
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	retryPolicy RetryPolicy
	Debug       *bool
}

// RetryPolicy specifies number of retries and min/max retry delays
// This config is used when the client exponentially backs off after errored requests.
type RetryPolicy struct {
	MaxRetries    int
	MinRetryDelay time.Duration
	MaxRetryDelay time.Duration
}

// newClient builds and configures a new instance of the XMattersAPI client with customizable options.
func newClient(hostname string, opts ...Option) (*XMattersAPI, error) {
	// Initialize the default HTTP client.
	// The retryablehttp package provides a client that automatically retries failed requests.
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = &loghttp.Transport{}

	// Initialize the XMattersAPI client with the base URL, user agent, and HTTP client.
	// The headers field is initialized as an empty http.Header map.
	xmatters := &XMattersAPI{
		BaseURL:    StringPtr(fmt.Sprintf("https://%v%v", hostname, defaultBasePath)),
		UserAgent:  StringPtr(fmt.Sprintf("xmatters-go/%v", Version)),
		httpClient: retryablehttp.NewClient().StandardClient(),
		headers:    make(http.Header),
	}

	// Process any additional options provided to the client.
	err := xmatters.parseOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("options parsing failed: %w", err)
	}
	return xmatters, nil
}

// NewWithBasicAuth creates a new instance of XMattersAPI with the provided URL, credentials, and options.
func NewWithBasicAuth(hostname, username, password *string, opts ...Option) (*XMattersAPI, error) {
	// Ensure that the hostname, username, and password are provided
	if hostname == nil {
		return nil, ErrNoHostname
	}
	// Initialize empty XMattersAPI client struct and error variable
	var xmatters *XMattersAPI
	var err error

	// Encode combined user/pass auth string
	base64AuthString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", *username, *password)))

	// Create a new XMattersAPI client with the provided hostname and options
	if xmatters, err = newClient(*hostname, opts...); err != nil {
		return nil, err
	}
	// Set required client object properties
	xmatters.AuthType = &AuthTypeBasic
	xmatters.Username = username
	xmatters.Password = password
	xmatters.headers.Add("Authorization", "Basic "+base64AuthString)

	return xmatters, nil
}

// NewWithToken creates a new instance of XMattersAPI with the provided URL, API token, and options.
func NewWithToken(hostname, token *string, opts ...Option) (*XMattersAPI, error) {
	// Ensure that the hostname and token are provided
	if hostname == nil {
		return nil, ErrNoHostname
	}
	// Initialize empty XMattersAPI client struct and error variable
	var xmatters *XMattersAPI
	var err error

	// Create a new XMattersAPI client with the provided hostname and options
	if xmatters, err = newClient(*hostname, opts...); err != nil {
		return nil, err
	}
	// Set required client object properties
	xmatters.AuthType = &AuthTypeOAuth
	xmatters.Token = token
	xmatters.headers.Add("Authorization", "Bearer "+*xmatters.Token)

	return xmatters, nil
}

// Request performs an HTTP request with the specified method, URI, content type, and request body.
// It returns the response body as a byte slice or an error, if any.
func (xmatters *XMattersAPI) Request(httpMethod, uri, contentType string, body interface{}) ([]byte, error) {
	// Initialize the request body and error variable
	var reqBody io.Reader
	var err error

	// Check for any provided body content and create the io.Reader
	// The body content must be type assertable to io.Reader or []byte, or able to be marshalled to JSON
	if body != nil {
		if r, ok := body.(io.Reader); ok {
			reqBody = r
		} else if paramBytes, ok := body.([]byte); ok {
			reqBody = bytes.NewReader(paramBytes)
		} else {
			var jsonBody []byte
			jsonBody, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("error marshalling body to JSON: %w", err)
			}
			reqBody = bytes.NewReader(jsonBody)
		}
	}

	// Create the HTTP request with the specified method, URI, and request body
	request, err := http.NewRequest(httpMethod, *xmatters.BaseURL+uri, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation failed: %w", err)
	}

	// Set necessary headers
	requestHeaders := make(http.Header)
	requestHeaders.Set("Content-Type", contentType)
	requestHeaders.Set("User-Agent", *xmatters.UserAgent)
	copyHeader(requestHeaders, xmatters.headers)
	request.Header = requestHeaders

	// Perform the request.
	response, err := xmatters.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer response.Body.Close()

	// Return error if no body content is returned
	if response.StatusCode == StatusNoContent {
		return nil, ErrNoContent // Return a generic 204 xMattersError struct
	}

	// If the response status code is 401, return an unauthorized error.
	if response.StatusCode == StatusUnauthorized {
		return nil, ErrInavlidCredentials
	}

	// Read the response body.
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body: %w", err)
	}

	// If the response status code is not 200 or 201, return an error.
	if response.StatusCode != StatusOK && response.StatusCode != StatusCreated {
		return nil, newXMattersError(respBody)
	}

	return respBody, nil
}

// buildURI assembles the base path and queries for API requests.
func buildURI(path string, options interface{}) string {
	v, _ := query.Values(options)
	groupsAttr := v.Get("groups")
	v.Del("groups")

	rawQuery := v.Encode()
	if groupsAttr != "" {
		rawQuery += "&groups=" + groupsAttr
	}

	return (&url.URL{Path: path, RawQuery: rawQuery}).String()
}

// copyHeader copies the headers from the source http.Header to the target http.Header.
// Note: The function overwrites any existing headers in the target with the corresponding headers from the source.
func copyHeader(target, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}

// Helper function to get a string pointer
func StringPtr(value string) *string {
	return &value
}
