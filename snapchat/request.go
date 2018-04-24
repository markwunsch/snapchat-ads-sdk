package snapchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context/ctxhttp"
)

// do is used to executed http responses and unmarshal the results into the provided interface
func (cli *Client) do(ctx context.Context, request *http.Request, target interface{}) error {
	request.Header.Set("User-Agent", `Snapchat Ads API Go SDK `+cli.version)

	response, err := ctxhttp.Do(ctx, cli.client, request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response != nil {
		if statusErr := getErrorFromStatusCode(response.StatusCode); statusErr != nil {
			return statusErr
		}
		return json.NewDecoder(response.Body).Decode(target)
	}
	return fmt.Errorf(`nil response`)
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

// ErrUnauthorized is the error returned when the api returns a 401 status code
type ErrUnauthorized struct{}

func (err *ErrUnauthorized) Error() string {
	return "401: unauthorized"
}

// ErrBadRequest is the error returned when the api returns a 400 status code
type ErrBadRequest struct{}

func (err *ErrBadRequest) Error() string {
	return "400: bad request"
}

// ErrForbidden is the error returned when the api returns a 403 status code
type ErrForbidden struct{}

func (err *ErrForbidden) Error() string {
	return "403: forbidden"
}

// ErrNotFound is the error returned when the api returns a 404 status code
type ErrNotFound struct{}

func (err *ErrNotFound) Error() string {
	return "404: not found"
}

// ErrMethodNotAllowed is the error returned when the api returns a 405 status code
type ErrMethodNotAllowed struct{}

func (err *ErrMethodNotAllowed) Error() string {
	return "405: method not allowed"
}

// ErrNotAcceptable is the error returned when the api returns a 406 status code
type ErrNotAcceptable struct{}

func (err *ErrNotAcceptable) Error() string {
	return "406: not acceptable"
}

// ErrGone is the error returned when the api returns a 410 status code
type ErrGone struct{}

func (err *ErrGone) Error() string {
	return "410: gone"
}

// ErrTooManyRequests is the error returned when the api returns a 429 status code
type ErrTooManyRequests struct{}

func (err *ErrTooManyRequests) Error() string {
	return "429: too many requests"
}

// ErrInternalServerError is the error returned when the api returns a 500 status code
type ErrInternalServerError struct{}

func (err *ErrInternalServerError) Error() string {
	return "500: internal server error"
}

// ErrServiceUnavailable is the error returned when the api returns a 503 status code
type ErrServiceUnavailable struct{}

func (err *ErrServiceUnavailable) Error() string {
	return "503: service unavailable"
}
