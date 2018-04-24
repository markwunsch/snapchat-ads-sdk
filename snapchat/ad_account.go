package snapchat

import (
	"context"
	"fmt"
	"strings"
)

// AdAccountService provides functions for interacting with snapchat ad accounts
type AdAccountService service

// AdAccount represents an ad account in the snapchat ads api
type AdAccount struct {
	// Id is the unique id associated with the ad account
	Id string `json:"id"`
	// Name is the name associated with the ad account
	Name string `json:"name"`
	// OrganizationId is the id of the organization the user is associated with
	OrganizationId string `json:"name"`
	// Timezone is the timezone that is associated with the account
	Timezone string `json:"name"`
	// Type is the type of ad account
	Type string `json:"type"`
	// LifetimeSpendCapMicro is the lifetime spend limit for the account
	LifetimeSpendCapMicro int64 `json:"lifetime_spend_cap_micro"`
	// AdvertiserOrganizationId is the organization id of the advertiser selected
	AdvertiserOrganizationId string `json:"advertiser_organization_id"`
	// Advertiser is the name of the advertiser associated with the account
	Advertiser string `json:"advertiser"`
	// Currency is the type of currency associated with the ad account
	Currency string `json:"currency"`
	// FundingSourceIds is a list of funding source ids associated with the ad account
	FundingSourceIds []string `json:"funding_source_ids"`
}

// GetAdAccountsResponse is the response object for calls to get ad accounts
type GetAdAccountsResponse struct {
	// RequestStatus is the status of the get ad account request
	RequestStatus string `json:"request_status"`
	// RequestId is the id associated with the request
	RequestId string `json:"request_id"`
	// AdAccounts is a list of individual ad account responses
	AdAccounts []*AdAccountResponse `json:"adaccounts"`
}

// AdAccountResponse is the individual organization object in the response for calls to get ad accounts
type AdAccountResponse struct {
	// SubRequestStatus is the status of this specific ad account request
	SubRequestStatus string `json:"sub_request_status"`
	// AdAccount is the object representing an ad account
	AdAccount AdAccount `json:"adaccount"`
}

// Get returns a single ad account associated with the provided ad account id
func (ad *AdAccountService) Get(ctx context.Context, adAccountId string) (*AdAccount, error) {
	path := fmt.Sprintf(`adaccounts/%s`, adAccountId)
	req, err := ad.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := new(GetAdAccountsResponse)
	err = ad.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}
	if a.RequestStatus == "success" {
		if len(a.AdAccounts) > 1 {
			return &a.AdAccounts[0].AdAccount, nil
		}
		return nil, fmt.Errorf("No ad accounts found with ad account id: %s", adAccountId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get ad account): %s`, a.RequestStatus)
}

// List returns all ad accounts associated with the provided organization id
func (ad *AdAccountService) List(ctx context.Context, organizationId string) ([]*AdAccount, error) {
	path := fmt.Sprintf(`organizations/%s/adaccounts`, organizationId)
	req, err := ad.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetAdAccountsResponse)
	err = ad.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.AdAccounts) > 0 {
			return getAdAccountsFromResponse(c.AdAccounts), nil
		}
		return nil, fmt.Errorf("no ad accounts found for organization id: %s", organizationId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list ad accounts): %s`, c.RequestStatus)
}

// getAdAccountsFromResponse returns the organization objects in an GetAdAccountsResponse object
func getAdAccountsFromResponse(list []*AdAccountResponse) []*AdAccount {
	var results []*AdAccount
	for _, val := range list {
		if strings.ToLower(val.SubRequestStatus) == "success" {
			results = append(results, &val.AdAccount)
		}
	}
	return results
}
