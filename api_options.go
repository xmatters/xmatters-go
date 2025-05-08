package xmatters

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Option is a functional option for configuring the XMattersAPI client
type Option func(*XMattersAPI) error

// WithBaseURL overrides the default base URL used for API calls
func WithBaseURL(baseURL string) Option {
	return func(xmatters *XMattersAPI) error {
		xmatters.BaseURL = StringPtr(fmt.Sprintf("%v%v", baseURL, defaultBasePath))
		return nil
	}
}

// HTTPClient accepts a custom *http.Client for making XMattersAPI calls
func WithHTTPClient(client *http.Client) Option {
	return func(xmatters *XMattersAPI) error {
		xmatters.httpClient = client
		return nil
	}
}

// Headers allows you to set custom HTTP headers when making XMattersAPI calls
func WithHeaders(headers http.Header) Option {
	return func(xmatters *XMattersAPI) error {
		xmatters.headers = headers
		return nil
	}
}

// WithRateLimit applies a non-default rate limit to client API requests
// If not specified the default of 4rps will be applied.
func WithRateLimit(rps float64) Option {
	return func(xmatters *XMattersAPI) error {
		// because ratelimiter doesnt do any windowing
		// setting burst makes it difficult to enforce a fixed rate
		// so setting it equal to 1 this effectively disables bursting
		// this doesn't check for sensible values, ultimately the xmatters will enforce that the value is ok
		xmatters.rateLimiter = rate.NewLimiter(rate.Limit(rps), 1)
		return nil
	}
}

// WithRetryPolicy applies a non-default number of retries and min/max retry delays
// This will be used when the client exponentially backs off after errored requests.
func WithRetryPolicy(maxRetries int, minRetryDelaySecs int, maxRetryDelaySecs int) Option {
	// seconds is very granular for a minimum delay - but this is only in case of failure
	return func(xmatters *XMattersAPI) error {
		xmatters.retryPolicy = RetryPolicy{
			MaxRetries:    maxRetries,
			MinRetryDelay: time.Duration(minRetryDelaySecs) * time.Second,
			MaxRetryDelay: time.Duration(maxRetryDelaySecs) * time.Second,
		}
		return nil
	}
}

// Debug is an option for configuring the XMattersAPI client to enable or disable debugging mode.
// When debugging is enabled, additional information and logs may be output to aid in troubleshooting.
// Use this option by passing a pointer to a boolean indicating whether debugging should be enabled.
// Example usage:
//
//	client := NewXMattersAPI(Debug(true))
//	// Debug mode is now enabled for the client.
func Debug(debug bool) Option {
	return func(xmatters *XMattersAPI) error {
		xmatters.Debug = &debug
		return nil
	}
}

// parseOptions parses the supplied options functions and returns a configured *XMattersAPI instance
func (xmatters *XMattersAPI) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our XMattersAPI type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(xmatters)
		if err != nil {
			return err
		}
	}

	return nil
}
