package snapchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"golang.org/x/net/context/ctxhttp"
)

// requestResponse is used to encapsulate important info from an http request
type requestResponse struct {
	// header contains the http headers that were returned by the api
	header http.Header
	// statusCode contains the status code returned by the api
	statusCode int
	// requestURL contains the url of the request that was executed
	requestURL *url.URL
}

// do is used to executed http responses and unmarshal the results into the provided interface
func (cli *Client) do(ctx context.Context, request *http.Request, target interface{}) (requestResponse, error) {
	responseObj := requestResponse{requestURL: request.URL, statusCode: -1}
	request.Header.Set("User-Agent", `Snapchat Ads API Go SDK `+cli.version)

	response, err := ctxhttp.Do(ctx, cli.client, request)
	if err != nil {
		return responseObj, err
	}
	defer response.Body.Close()

	if response != nil {
		responseObj.statusCode = response.StatusCode
		responseObj.header = response.Header
	}

	if responseObj.statusCode < 200 || responseObj.statusCode >= 400 {
		return responseObj, errors.New(fmt.Sprintf(`%d status code returned from snapchat api`, responseObj.statusCode))
	}

	err = json.NewDecoder(response.Body).Decode(target)
	return responseObj, err
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
