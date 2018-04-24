package snapchat

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// OrganizationService provides functions for interacting with snapchat organizations
type OrganizationService service

// Organization represents an organization in the snapchat ads api
type Organization struct {
	// Id is the unique id associated with the organization
	Id string `json:"id"`
	// Name is the name of the organization
	Name string `json:"name"`
	// AddressLine1 is the street address associated with the organization
	AddressLine1 string `json:"address_line_1"`
	// Locality is the region/city associated with the account
	Locality string `json:"locality"`
	// AdministrativeDistrictLevel1 is the region/state associated with the organization
	AdministrativeDistrictLevel1 string `json:"administrative_district_level_1"`
	// County is the country associated with the organization
	Country string `json:"country"`
	// PostalCode is the postal code associated with the organization
	PostalCode string `json:"postal_code"`
	// Type is the type of organization
	Type string `json:"type"`
	// CreatedAt is the time when the organization was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the organization was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// GetOrganizationsResponse is the response object for calls to get organizations
type GetOrganizationsResponse struct {
	RequestStatus string                  `json:"request_status"`
	RequestId     string                  `json:"request_id"`
	Organizations []*OrganizationResponse `json:"organizations"`
}

// OrganizationResponse is the individual organization object in the response for calls to get organizations
type OrganizationResponse struct {
	SubRequestStatus string       `json:"sub_request_status"`
	Organization     Organization `json:"organization"`
}

// Get returns a single organization associated with the provided organization id
func (org *OrganizationService) Get(ctx context.Context, organizationId string) (*Organization, error) {
	path := fmt.Sprintf(`organizations/%s`, organizationId)
	req, err := org.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := new(GetOrganizationsResponse)
	err = org.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(a.RequestStatus) == "success" {
		if len(a.Organizations) >= 1 {
			if strings.ToLower(a.Organizations[0].SubRequestStatus) == "success" {
				return &a.Organizations[0].Organization, nil
			}
			return nil, fmt.Errorf(`non-success status returned from snapchat api (get organization): %s`, a.RequestStatus)
		}
		return nil, fmt.Errorf("no organizations found with organization id: %s", organizationId)
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get organization): %s`, a.RequestStatus)
}

// List returns all organizations associated with the authenticated user
func (org *OrganizationService) List(ctx context.Context) ([]*Organization, error) {
	path := "me/organizations"
	req, err := org.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	c := new(GetOrganizationsResponse)
	err = org.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(c.RequestStatus) == "success" {
		if len(c.Organizations) > 0 {
			return getOrganizationsFromResponse(c.Organizations), nil
		}
		return nil, fmt.Errorf("no organizations found")
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (list organizations): %s`, c.RequestStatus)
}

// getOrganizationsFromResponse returns the organization objects in an OrganizationResponse object
func getOrganizationsFromResponse(list []*OrganizationResponse) []*Organization {
	var results []*Organization
	for _, val := range list {
		if strings.ToLower(val.SubRequestStatus) == "success" {
			results = append(results, &val.Organization)
		}
	}
	return results
}
