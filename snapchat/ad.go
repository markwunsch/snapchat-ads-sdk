package snapchat

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AdService provides functions for interacting with snapchat ads
type AdService service

// Ad represents an ad in the snapchat ads api
type Ad struct {
	// Id is the id representing a single ad
	Id string `json:"id"`
	// AdSquadId is the id of the ad squad that the ad is under
	AdSquadId string `json:"ad_squad_id"`
	// CreativeId is the id of the creative associated with the ad
	CreativeId string `json:"creative_id"`
	// Name is the name of the ad
	Name string `json:"name"`
	// Status is the status of the ad
	Status string `json:"status"`
	// ReviewStatus is the status of the ad's review process
	ReviewStatus string `json:"review_status"`
	// ReviewStatusReason will contain a reason for rejection if an ad was rejected
	ReviewStatusReason string `json:"review_status_reason"`
	// Type is the type of the ad
	Type string `json:"type"`
	// CreatedAt is the time when the campaign was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the campaign was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAdsResponse is the response object returned when getting ads
type GetAdsResponse struct {
	RequestStatus string        `json:"request_status"`
	RequestId     string        `json:"request_id"`
	Ads           []*AdResponse `json:"ads"`
}

// AdResponse is the object for a single ad response
type AdResponse struct {
	SubRequestStatus string `json:"sub_request_status"`
	Ad               Ad     `json:"ad"`
}

// Get is used to get the specific ad associated with the provided ad id
func (ad *AdService) Get(ctx context.Context, adId string) (*Ad, error) {
	path := fmt.Sprintf(`ads/%s`, adId)
	req, err := ad.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := new(GetAdsResponse)
	err = ad.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(a.RequestStatus) == "success" {
		if len(a.Ads) > 1 {
			return &a.Ads[0].Ad, nil
		}
		return nil, fmt.Errorf("no ads found with ad id: %s", adId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get ad with id %s): %s`, adId, a.RequestStatus)
}

// ListByAdSquad will return all of the ads associated with the given ad squad id
func (ad *AdService) ListByAdSquad(ctx context.Context, adSquadId string) ([]*Ad, error) {
	path := fmt.Sprintf(`adsquads/%s/ads`, adSquadId)
	req, err := ad.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetAdsResponse)
	err = ad.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.Ads) > 0 {
			return getAdsFromResponse(c.Ads), nil
		}
		return nil, fmt.Errorf("no ads found for ad squad id: %s", adSquadId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list ads for ad squad with id %s): %s`, adSquadId, c.RequestStatus)
}

// ListByAdAccount will return all of the ads associated with the given ad account id
func (ad *AdService) ListByAdAccount(ctx context.Context, adAccountId string) ([]*Ad, error) {
	path := fmt.Sprintf(`adaccounts/%s/ads`, adAccountId)
	req, err := ad.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetAdsResponse)
	err = ad.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.Ads) > 0 {
			return getAdsFromResponse(c.Ads), nil
		}
		return nil, fmt.Errorf("no ads found for ad account id: %s", adAccountId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list ads for ad account with id %s): %s`, adAccountId, c.RequestStatus)
}

// Delete deletes a specific ad
func (ad *AdService) Delete(ctx context.Context, adId string) error {
	path := fmt.Sprintf(`ads/%s`, adId)
	req, err := ad.client.createRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	c := new(GetAdsResponse)
	err = ad.client.do(ctx, req, c)
	if err != nil {
		return err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		return nil
	}
	return fmt.Errorf(`non-success status returned from snapchat api (delete ad with id %s): %s`, adId, c.RequestStatus)
}

func getAdsFromResponse(list []*AdResponse) []*Ad {
	var results []*Ad
	for _, val := range list {
		if strings.ToLower(val.SubRequestStatus) == "success" {
			results = append(results, &val.Ad)
		}
	}
	return results
}
