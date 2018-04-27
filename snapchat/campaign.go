package snapchat

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// CampaignService has methods for interacting with campaigns in snapchat
type CampaignService service

// Campaign has a business objective and organizes Ad Squads
type Campaign struct {
	client *Client
	// Id is the id representing a single campaign
	Id string `json:"id"`
	// AdAccountId is the Ad Account ID that this campaign is under
	AdAccountId string `json:"ad_account_id"`
	// Name is the name of the campaign
	Name string `json:"name"`
	// Status is the status of the campaign (ACTIVE, PAUSED)
	Status string `json:"status"`
	// MeasurementSpec contains the apps to be tracked for this campaign
	MeasurementSpec CampaignMeasurementSpec `json:"measurement_spec"`
	// EndTime is the time when the campaign will end
	EndTime time.Time `json:"end_time"`
	// StartTime is the time when the campaign will start
	StartTime time.Time `json:"start_time"`
	// CreatedAt is the time when the campaign was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the campaign was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// DailyBudgetMicro is the daily spend cap for the campaign(micro-currency)
	DailyBudgetMicro int64 `json:"daily_budget_micro"`
	// LifetimeSpendCapMicro is the lifetime spend cap for the campaign (microcurrency)
	LifetimeSpendCapMicro int64 `json:"lifetime_spend_cap_micro"`
}

// GetCampaignsResponse is the response object returned when getting campaigns
type GetCampaignsResponse struct {
	RequestStatus string              `json:"request_status"`
	RequestId     string              `json:"request_id"`
	Campaigns     []*CampaignResponse `json:"campaigns"`
}

// CampaignResponse is the object for a single campaign response
type CampaignResponse struct {
	SubRequestStatus string   `json:"sub_request_status"`
	Campaign         Campaign `json:"campaign"`
}

// CampaignMeasurementSpec contains the apps to be tracked for this campaign
type CampaignMeasurementSpec struct {
	// IosAppId is the ios app id of the app to use for measurement tracking
	IosAppId string `json:"ios_app_id"`
	// AndroidAppUrl is the android app url of the app to use for measurement tracking
	AndroidAppUrl string `json:"android_app_url"`
}

// Get retrieves a specific campaign
func (cmp *CampaignService) Get(ctx context.Context, campaignId string) (*Campaign, error) {
	path := fmt.Sprintf(`campaigns/%s`, campaignId)
	req, err := cmp.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetCampaignsResponse)
	err = cmp.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.Campaigns) > 0 {
			return &c.Campaigns[0].Campaign, nil
		}
		return nil, fmt.Errorf("no campaigns found with campaign id: %s", campaignId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get campaign): %s`, c.RequestStatus)
}

// List retrieves all campaigns within a specified ad account
func (cmp *CampaignService) List(ctx context.Context, adAccountId string) ([]*Campaign, error) {
	path := fmt.Sprintf(`adaccounts/%s/campaigns`, adAccountId)
	req, err := cmp.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetCampaignsResponse)
	err = cmp.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.Campaigns) > 0 {
			return getCampaignsFromResponse(c.Campaigns), nil
		}
		return nil, fmt.Errorf("no campaigns found for ad account id: %s", adAccountId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list campaigns): %s`, c.RequestStatus)
}

// Delete deletes a specific campaign
func (cmp *CampaignService) Delete(ctx context.Context, campaignId string) error {
	path := fmt.Sprintf(`campaigns/%s`, campaignId)
	req, err := cmp.client.createRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	c := new(GetCampaignsResponse)
	err = cmp.client.do(ctx, req, c)
	if err != nil {
		return err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		return nil
	}
	return fmt.Errorf(`non-success status returned from snapchat api (delete campaign): %s`, c.RequestStatus)
}

func getCampaignsFromResponse(list []*CampaignResponse) []*Campaign {
	var results []*Campaign
	for _, val := range list {
		if strings.ToLower(val.SubRequestStatus) == "success" {
			results = append(results, &val.Campaign)
		}
	}
	return results
}
