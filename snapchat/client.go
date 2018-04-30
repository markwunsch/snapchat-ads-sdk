package snapchat

import (
	"context"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

const (
	// DefaultSnapchatHost is the default host address to use for calls to the snapchat ads api
	DefaultSnapchatHost = `https://adsapi.snapchat.com`
	// DefaultSnapchatVersion is the default version to use for calls to the snapchat ads api
	DefaultSnapchatVersion = `v1`
)

type service struct {
	client *Client
}

// Client is used to perform all operations with the snapchat ads api
type Client struct {
	// host holds the address of the snapchat server
	host string
	// client handles http requests
	client *http.Client
	// version of the snapchat api to use
	version string
	// customHTTPHeaders can be set by user to include additional headers
	customHTTPHeaders map[string]string
	// customVersion is set to true if the user specified a custom version
	customVersion bool
	// accessToken is the user's access token to use for authorization
	accessToken string
	// Users is the service used to get the authenticated user
	Users *UserService
	// Organizations is the service used to interact with organizations
	Organizations *OrganizationService
	// AdAccounts is the service used to interact with ad accounts
	AdAccounts *AdAccountService
	// Campaigns is the service used to interact with campaigns
	Campaigns *CampaignService
	// FundingSources is the service used to interact with campaigns
	FundingSources *FundingSourceService
	// AdSquads is the service used to interact with ad squads
	AdSquads *AdSquadService
}

// NewClient creates a new instance of Client with any optional functions applied
func NewClient(optFns ...func(*Client) error) (*Client, error) {
	client, err := defaultHTTPClient()
	if err != nil {
		return nil, err
	}
	c := &Client{
		host:    DefaultSnapchatHost,
		version: DefaultSnapchatVersion,
		client:  client,
	}
	c.Users = &UserService{client: c}
	c.Organizations = &OrganizationService{client: c}
	c.AdAccounts = &AdAccountService{client: c}
	c.Campaigns = &CampaignService{client: c}
	c.FundingSources = &FundingSourceService{client: c}
	c.AdSquads = &AdSquadService{client: c}

	for _, fn := range optFns {
		if err := fn(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithAccessToken allows the user to provide an access token to use for authorization
func WithAccessToken(ctx context.Context, accessToken string) func(*Client) error {
	return func(c *Client) error {
		c.accessToken = accessToken
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		tc.Timeout = time.Minute
		c.client = tc
		return nil
	}
}

// UpdateAccessToken can be used to update an outdated access token
func (cli *Client) UpdateAccessToken(ctx context.Context, accessToken string) {
	cli.accessToken = accessToken
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli.client = tc
}

// WithHTTPClient allows user to provide a custom client
func WithHTTPClient(client *http.Client) func(*Client) error {
	return func(c *Client) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
}

// WithHost allows the user to use a custom snapchat ads api host
func WithHost(host string) func(*Client) error {
	return func(c *Client) error {
		c.host = host
		return nil
	}
}

// WithVersion allows the user to use a custom snapchat ads api version
func WithVersion(version string) func(*Client) error {
	return func(c *Client) error {
		c.version = version
		c.customVersion = true
		return nil
	}
}

// WithHTTPHeaders allows the user to specify custom headers to be used with all requests
func WithHTTPHeaders(headers map[string]string) func(*Client) error {
	return func(c *Client) error {
		c.customHTTPHeaders = headers
		return nil
	}
}

// GetCustomHTTPHeaders returns custom http headers stored by the client
func (cli *Client) GetCustomHTTPHeaders() map[string]string {
	headers := make(map[string]string)
	for k, v := range cli.customHTTPHeaders {
		headers[k] = v
	}
	return headers
}

// WithEnvVars allows retrieves host/version from environment variables
func WithEnvVars(c *Client) error {
	if host := os.Getenv("TWITTER_ADS_HOST"); host != "" {
		if err := WithHost(host)(c); err != nil {
			return err
		}
	}
	if version := os.Getenv("TWITTER_ADS_API_VERSION"); version != "" {
		if err := WithVersion(version)(c); err != nil {
			return err
		}
	}
	return nil
}

// defaultHTTPClient returns an http client with default parameters
func defaultHTTPClient() (*http.Client, error) {
	return &http.Client{
		Timeout: time.Minute,
	}, nil
}
