package snapchat

import (
	"context"
	"time"
)

type UserService service

type User struct {
	// Id is the unique id associated with the user
	Id string `json:"id"`
	// Email is the email associated with the user
	Email string `json:"email"`
	// OrganizationId is the id of the organization associated with the user
	OrganizationId string `json:"organization_id"`
	// DisplayName is the display name of the user
	DisplayName string `json:"display_name"`
	// CreatedAt is the time when the user was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time when the user was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

type getAuthenticatedUserResponse struct {
	RequestId string `json:"request_id"`
	Me        *User  `json:"me"`
}

// GetAuthenticatedUser returns the user that is currently authenticated
func (usr *UserService) GetAuthenticatedUser(ctx context.Context) (*User, error) {
	path := `me`
	req, err := usr.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := new(getAuthenticatedUserResponse)
	err = usr.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}
	return a.Me, nil
}
