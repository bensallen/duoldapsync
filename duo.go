package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
)

// PostUsersResult represents the response from the POST /admin/v1/users endpoint
type PostUsersResult struct {
	duoapi.StatResult
	Response admin.User
}

// CreateUser creates a new Duo user via the Duo Admin Client
// See https://duo.com/docs/adminapi#create-user
func CreateUser(client *admin.Client, params url.Values, dryRun bool) (*PostUsersResult, error) {
	if !dryRun {
		_, body, err := client.SignedCall("POST", "/admin/v1/users", params, duoapi.UseTimeout)
		if err != nil {
			return nil, err
		}

		ret := &PostUsersResult{}
		if err = json.Unmarshal(body, ret); err != nil {
			return nil, err
		}
		return ret, nil
	}

	return &PostUsersResult{duoapi.StatResult{Stat: "OK"}, admin.User{}}, nil
}

// DeleteUser deletes a Duo user via the Duo Admin Client
// See https://duo.com/docs/adminapi#delete-user
func DeleteUser(client *admin.Client, userID string, dryRun bool) (*duoapi.StatResult, error) {
	if !dryRun {
		path := fmt.Sprintf("/admin/v1/users/%s", userID)
		_, body, err := client.SignedCall("DELETE", path, nil, duoapi.UseTimeout)
		if err != nil {
			return nil, err
		}

		ret := &duoapi.StatResult{}
		if err = json.Unmarshal(body, ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
	return &duoapi.StatResult{Stat: "OK"}, nil
}

// EnrollUser enrolls a user via the Duo Admin Client with user name username and email
// address email and send them an enrollment email that expires after valid_secs seconds.
// See https://duo.com/docs/adminapi#enroll-user
func EnrollUser(client *admin.Client, params url.Values, dryRun bool) (*duoapi.StatResult, error) {
	if !dryRun {
		_, body, err := client.SignedCall("POST", "/admin/v1/users/enroll", params, duoapi.UseTimeout)
		if err != nil {
			return nil, err
		}

		ret := &duoapi.StatResult{}
		if err = json.Unmarshal(body, ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
	return &duoapi.StatResult{Stat: "OK"}, nil
}
