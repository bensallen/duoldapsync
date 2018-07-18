package main

import (
	"encoding/json"
	"net/url"

	"github.com/duosecurity/duo_api_golang"
)

// AdminAPI is a Duo Admin API object
type AdminAPI struct {
	duoapi.DuoApi
}

// NewAdminAPI builds a new Duo Admin API object.
// Argument api is a duoapi.DuoApi object used to make the Duo Rest API calls.
// Example: authapi.NewAdminAPI(*duoapi.NewDuoApi(ikey,skey,host,userAgent,duoapi.SetTimeout(10*time.Second)))
func NewAdminAPI(api duoapi.DuoApi) *AdminAPI {
	return &AdminAPI{api}
}

// StatResponse is the standard status response from all endpoints. On success, Stat is 'OK'.
// On error, Stat is 'FAIL', and Code, Message, and Message_Detail contain error information.
type StatResponse struct {
	Stat          string
	Code          int32
	Message       string
	MessageDetail string
}

// GroupResponse represents a group that the user is a member of in the response as part of the /admin/v1/users endpoint
type GroupResponse struct {
	Desc string
	Name string
}

// PhoneResponse represents a phone device response as part of the /admin/v1/users endpoint
type PhoneResponse struct {
	Activated        bool
	Capabilities     []string
	Extension        string
	Fingerprint      string
	LastSeen         string
	Name             string
	Number           string
	PhoneID          string
	Platform         string
	Postdelay        string
	Predelay         string
	SmsPasscodesSent bool
	Tampered         string
	Type             string
}

// TokenResponse represents a security token response as part of the /admin/v1/users endpoint
type TokenResponse struct {
	Serial  string
	TokenID string
	Type    string
}

// UserResponse represents one response from the /admin/v1/users endpoint
type UserResponse struct {
	Alias1            string
	Alias2            string
	Alias3            string
	Alias4            string
	DesktopTokens     []TokenResponse
	Created           int
	Email             string
	Firstname         string
	Groups            []GroupResponse
	LastDirectorySync int
	LastLogin         int
	LastName          string
	Notes             string
	Phones            []PhoneResponse
	Realname          string
	Status            string
	Tokens            []TokenResponse
	U2ftokens         []TokenResponse
	UserID            string
	Username          string
}

// UsersResponse represents the response from the GET /admin/v1/users endpoint
type UsersResponse struct {
	StatResponse
	Response []UserResponse
}

// Users enumerates all existing users via the Duo Admin API
func (api *AdminAPI) Users(params url.Values) (*UsersResponse, error) {
	_, body, err := api.SignedCall("GET", "/admin/v1/users", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	ret := &UsersResponse{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// UsersResponse represents the response from the POST /admin/v1/users endpoint
type CreateUserResponse struct {
	StatResponse
	Response UserResponse
}

// CreateUser creates a new Duo user via the Duo Admin API
func (api *AdminAPI) CreateUser(params url.Values) (*CreateUserResponse, error) {

	_, body, err := api.SignedCall("POST", "/admin/v1/users", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	ret := &CreateUserResponse{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
