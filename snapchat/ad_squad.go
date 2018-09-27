package snapchat

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AdSquadService provides functions for interacting with snapchat ad squads
type AdSquadService service

// AdSquad represents an ad squad in the snapchat ads api
type AdSquad struct {
	// Id is the id representing a single ad squad
	Id string `json:"id"`
	// CampaignId is the id of the campaign this ad squad is associated with
	CampaignId string `json:"campaign_id"`
	// BidMicro the max bid (micro-currency)
	BidMicro int64 `json:"bid_micro"`
	// BillingEvent is the billing event associated with the ad squad
	BillingEvent string `json:"billing_event"`
	// DailyBudgetMicro is the daily spend budget (micro-currency)
	DailyBudgetMicro int64 `json:"daily_budget_micro"`
	// EndTime is the ending time of the ad squad
	EndTime time.Time `json:"end_time"`
	// StartTime is the starting time of the ad squad
	StartTime time.Time `json:"start_time"`
	// Name is the name of the ad squad
	Name string `json:"name"`
	// OptimizationGoal is the optimization goal of the ad squad
	OptimizationGoal string `json:"optimization_goal"`
	// Placement is the placement for the ad squad
	Placement string `json:"placement"`
	// Status is the status of the ad squad
	Status string `json:"status"`
	// Targeting is the targeting spec of the ad squad
	//Targeting string `json:"targeting"` // might not be a string
	// IncludedContentType is a list of content types that will be included in this ad squad
	IncludedContentType []string `json:"included_content_types"`
	// ExcludedContentType is a list of content types that will be excluded in this ad squad
	ExcludedContentType []string `json:"excluded_content_types"`
	// CapAndExclusionConfig is the frequency cap and exclusion spec
	//CapAndExclusionConfig string `json:"cap_and_exclusion_config"` // might not be a string
	// LifetimeBudgetMicro is the lifetime spend budget of the ad squad (micro-currency)
	LifetimeBudgetMicro int64 `json:"lifetime_budget_micro"`
	// AdSchedulingConfig is the schedule for running ads on this adsquad
	//AdSchedulingConfig string `json:"ad_scheduling_config"` // might not be a string
	// Type is the type of ad squad
	Type string `json:"type"`
}

// GetAdSquadsResponse is the response object returned when getting ad squads
type GetAdSquadsResponse struct {
	RequestStatus string             `json:"request_status"`
	RequestId     string             `json:"request_id"`
	AdSquads      []*AdSquadResponse `json:"adsquads"`
}

// AdSquadResponse is the object for a single ad squad response
type AdSquadResponse struct {
	SubRequestStatus string  `json:"sub_request_status"`
	AdSquad          AdSquad `json:"adsquad"`
}

// Get retrieves a specific ad squad
func (adsqd *AdSquadService) Get(ctx context.Context, adSquadId string) (*AdSquad, error) {
	path := fmt.Sprintf(`adsquads/%s`, adSquadId)
	req, err := adsqd.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := new(GetAdSquadsResponse)
	err = adsqd.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(a.RequestStatus) == "success" {
		if len(a.AdSquads) > 1 {
			return &a.AdSquads[0].AdSquad, nil
		}
		return nil, fmt.Errorf("no ad squads found with ad squad id: %s", adSquadId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get ad squad with id %s): %s`, adSquadId, a.RequestStatus)
}

// ListByCampaign retrieves all ad squads associated a specified campaign id
func (adsqd *AdSquadService) ListByCampaign(ctx context.Context, campaignId string) ([]*AdSquad, error) {
	path := fmt.Sprintf(`campaigns/%s/adsquads`, campaignId)
	req, err := adsqd.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetAdSquadsResponse)
	err = adsqd.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.AdSquads) > 0 {
			return getAdSquadsFromResponse(c.AdSquads), nil
		}
		return nil, fmt.Errorf("no ad squads found for campaign id: %s", campaignId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list ad squads for campaign with id %s): %s`, campaignId, c.RequestStatus)
}

// ListByAdAccount retrieves all ad squads associated a specified ad account id
func (adsqd *AdSquadService) ListByAdAccount(ctx context.Context, adAccountId string) ([]*AdSquad, error) {
	path := fmt.Sprintf(`adaccounts/%s/adsquads`, adAccountId)
	req, err := adsqd.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetAdSquadsResponse)
	err = adsqd.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.AdSquads) > 0 {
			return getAdSquadsFromResponse(c.AdSquads), nil
		}
		return nil, fmt.Errorf("no ad squads found for ad account id: %s", adAccountId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list ad squads for ad account with id %s): %s`, adAccountId, c.RequestStatus)
}

// Delete deletes a specific ad squad
func (adsqd *AdSquadService) Delete(ctx context.Context, adSquadId string) error {
	path := fmt.Sprintf(`adsquads/%s`, adSquadId)
	req, err := adsqd.client.createRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	c := new(GetAdSquadsResponse)
	err = adsqd.client.do(ctx, req, c)
	if err != nil {
		return err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		return nil
	}
	return fmt.Errorf(`non-success status returned from snapchat api (delete ad squad with id %s): %s`, adSquadId, c.RequestStatus)
}

func getAdSquadsFromResponse(list []*AdSquadResponse) []*AdSquad {
	var results []*AdSquad
	for _, val := range list {
		if strings.ToLower(val.SubRequestStatus) == "success" {
			results = append(results, &val.AdSquad)
		}
	}
	return results
}
