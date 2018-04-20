package snapchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

// RequestResponse is used to encapsulate important info from an http request
type RequestResponse struct {
	// Header contains the http headers that were returned by the api
	Header http.Header
	// StatusCode contains the status code returned by the api
	StatusCode int
	// RequestURL contains the url of the request that was executed
	RequestURL *url.URL
}

// do is used to executed http responses and unmarshal the results into the provided interface
func (cli *Client) do(ctx context.Context, request *http.Request, target interface{}) error {
	responseObj := RequestResponse{RequestURL: request.URL, StatusCode: -1}
	request.Header.Set("User-Agent", `Snapchat Ads API Go SDK `+cli.version)

	response, err := ctxhttp.Do(ctx, cli.client, request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response != nil {
		responseObj.StatusCode = response.StatusCode
		responseObj.Header = response.Header
	}

	if statusErr := getErrorFromStatusCode(response.StatusCode); statusErr != nil {
		return statusErr
	}

	err = json.NewDecoder(response.Body).Decode(target)
	return err
}

// createRequest is used to get an http request object
func (cli *Client) createRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	if (method == "POST" || method == "PUT") && body == nil {
		body = bytes.NewReader([]byte{})
	}

	path = fmt.Sprintf(`%s/%s/%s`, cli.host, cli.version, path)
	request, err := http.NewRequest(method, path, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	return request, nil
}

func getErrorFromStatusCode(statusCode int) error {
	switch statusCode {
	case 401:
		return new(ErrBadRequest)
	case 402:
		return new(ErrUnauthorized)
	case 403:
		return new(ErrForbidden)
	case 404:
		return new(ErrNotFound)
	case 405:
		return new(ErrMethodNotAllowed)
	case 406:
		return new(ErrNotAcceptable)
	case 410:
		return new(ErrGone)
	case 429:
		return new(ErrTooManyRequests)
	case 500:
		return new(ErrInternalServerError)
	case 503:
		return new(ErrServiceUnavailable)
	default:
		if statusCode < 200 || statusCode >= 400 {
			return fmt.Errorf(`%d status code returned from snapchat api`, statusCode)
		}
	}
	return nil
}

type ErrUnauthorized struct{}

func (err *ErrUnauthorized) Error() string {
	return "401: unauthorized"
}

type ErrBadRequest struct{}

func (err *ErrBadRequest) Error() string {
	return "400: bad request"
}

type ErrForbidden struct{}

func (err *ErrForbidden) Error() string {
	return "403: forbidden"
}

type ErrNotFound struct{}

func (err *ErrNotFound) Error() string {
	return "404: not found"
}

type ErrMethodNotAllowed struct{}

func (err *ErrMethodNotAllowed) Error() string {
	return "405: method not allowed"
}

type ErrNotAcceptable struct{}

func (err *ErrNotAcceptable) Error() string {
	return "406: not acceptable"
}

type ErrGone struct{}

func (err *ErrGone) Error() string {
	return "410: gone"
}

type ErrTooManyRequests struct{}

func (err *ErrTooManyRequests) Error() string {
	return "429: too many requests"
}

type ErrInternalServerError struct{}

func (err *ErrInternalServerError) Error() string {
	return "500: internal server error"
}

type ErrServiceUnavailable struct{}

func (err *ErrServiceUnavailable) Error() string {
	return "503: service unavailable"
}
