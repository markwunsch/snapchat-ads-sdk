package snapchat

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// FundingSourceService provides functions for interacting with snapchat funding sources
type FundingSourceService service

// FundingSource represents a funding source in the snapchat ads api
type FundingSource struct {
	// Id is the unique id associated with the funding source
	Id string `json:"id"`
	// Type is the type of funding source
	Type string `json:"type"`
	// Status is the status of the funding source
	Status string `json:"status"`
	// BudgetSpentMicro is the total budget spent (micro-currency)
	BudgetSpentMicro int64 `json:"budget_spent_micro"`
	// Currency is the type of currency associated with the funding source
	Currency string `json:"currency"`
	// TotalBudgetMicro is the total budget (micro-currency)
	TotalBudgetMicro int64 `json:"total_budget_micro"`
	// AvailableCreditMicro is the amount of credit available (micro-currency)
	AvailableCreditMicro int64 `json:"available_credit_micro"`
	// CardType is the type of credit card associated with this funding source
	CardType string `json:"card_type"`
	// Name is the name of the funding source
	Name string `json:"name"`
	// Last4 is the last four digits of the credit card associated with this funding source
	Last4 string `json:"last_4"`
	// ExpirationYear is the expiration year of the credit card associated with this funding source
	ExpirationYear string `json:"expiration_year"`
	// ExpirationMonth is the expiration year of the credit card associated with this funding source
	ExpirationMonth string `json:"expiration_month"`
	// DailySpendLimitMicro is the daily spend limit for the credit card associated with this funding source (micro-currency)
	DailySpendLimitMicro int64 `json:"daily_spend_limit_micro"`
	// DailySpendLimitMicro is the currency of the  daily spend limit for the credit card associated with this funding source
	DailySpendLimitCurrency string `json:"daily_spend_limit_currency"`
	// ValueMicro is the value of the coupon (micro-currency)
	ValueMicro int64 `json:"value_micro"`
	// StartDate is the start date of the coupon
	StartDate time.Time `json:"start_date"`
	// EndDate is the end date of the coupon
	EndDate time.Time `json:"end_date"`
	// Email is the email associated with the funding source
	Email string `json:"email"`
	// CreatedAt is the time when the funding source was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the funding source was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// GetFundingSourcesResponse struct { is the response object returned when getting funding sources
type GetFundingSourcesResponse struct {
	RequestStatus  string                   `json:"request_status"`
	RequestId      string                   `json:"request_id"`
	FundingSources []*FundingSourceResponse `json:"fundingsources"`
}

// FundingSourceResponse is the object for a single funding source response
type FundingSourceResponse struct {
	SubRequestStatus string        `json:"sub_request_status"`
	FundingSource    FundingSource `json:"fundingsource"`
}

// Get retrieves a specific funding source
func (fnd *FundingSourceService) Get(ctx context.Context, fundingSourceId string) (*FundingSource, error) {
	path := fmt.Sprintf(`fundingsources/%s`, fundingSourceId)
	req, err := fnd.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetFundingSourcesResponse)
	err = fnd.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.FundingSources) > 0 {
			return &c.FundingSources[0].FundingSource, nil
		}
		return nil, fmt.Errorf("no funding sources found with funding source id: %s", fundingSourceId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get funding source with id %s): %s`, fundingSourceId, c.RequestStatus)
}

// List retrieves all funding sources associated with the specified organization
func (fnd *FundingSourceService) List(ctx context.Context, organizationId string) ([]*FundingSource, error) {
	path := fmt.Sprintf(`organizations/%s/funding-sources`, organizationId)
	req, err := fnd.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetFundingSourcesResponse)
	err = fnd.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.FundingSources) > 0 {
			return getFundingSourcesFromResponse(c.FundingSources), nil
		}
		return nil, fmt.Errorf("no funding sources found for organization id: %s", organizationId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list funding sources for organization id %s): %s`, organizationId, c.RequestStatus)
}

func getFundingSourcesFromResponse(list []*FundingSourceResponse) []*FundingSource {
	var results []*FundingSource
	for _, val := range list {
		if strings.ToLower(val.SubRequestStatus) == "success" {
			results = append(results, &val.FundingSource)
		}
	}
	return results
}
