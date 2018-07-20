package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	ldap "gopkg.in/ldap.v2"
)

// User represents the attributes of a user found in LDAP and records if the user has been found in Duo.
type User struct {
	Username  string
	DuoUserID string
	FullName  string
	Email     string
	FirstName string
	LastName  string

	LDAP        bool // User found in LDAP
	Duo         bool // User found in Duo
	NeedsUpdate bool // Indicates LDAP attributes are different that what is in Duo, and the Duo user needs to be updated.
}

// DuoCreate creates a user via the Duo Admin API
func (u *User) duoCreate(api *AdminAPI) error {
	params, err := u.urlValues()
	if err != nil {
		return fmt.Errorf("URLValues failed: %s when attempting to create user: %s", err, u.Username)
	}
	params.Set("status", "active")
	//params.Set("notes", ...) TODO: Add lastUpdated
	result, err := api.CreateUser(params)

	if err != nil {
		return fmt.Errorf("CreateUser failed: %s when attempting to create user: %s", err, u.Username)
	} else if result.Stat != "OK" {
		return fmt.Errorf("CreateUser Duo API returned non-ok status when attemping to create user: %s with message: %s", u.Username, result.Message)
	}
	return nil
}

// DuoEnroll sends an enrollment email via the Duo Admin API
func (u *User) duoEnroll(api *AdminAPI, enrollValidSecs int) error {
	enrollParams := url.Values{}
	enrollParams.Set("username", u.Username)
	enrollParams.Set("email", u.Email)
	enrollParams.Set("valid_secs", strconv.Itoa(enrollValidSecs))

	result, err := api.EnrollUser(enrollParams)
	if err != nil {
		return fmt.Errorf("CreateUser failed: %s when attempting to create user: %s", err, u.Username)
	} else if result.Stat != "OK" {
		return fmt.Errorf("CreateUser Duo API returned non-ok status when attemping to create user: %s with message: %s", u.Username, result.Message)
	}
	return nil
}

// URLValues transforms User's attributes into url.Values
func (u *User) urlValues() (url.Values, error) {
	params := url.Values{}
	if u.Username != "" {
		params.Set("username", u.Username)
	} else {
		return nil, fmt.Errorf("URLValues requires username argument to be a non-empty string")
	}
	if u.FullName != "" {
		params.Set("realname", u.FullName)
	}
	if u.Email != "" {
		params.Set("email", u.Email)
	}
	if u.FirstName != "" {
		params.Set("firstname", u.FirstName)
	}
	if u.LastName != "" {
		params.Set("lastname", u.LastName)
	}
	return params, nil
}

// UserSet is a map of Users indexed by username
type UserSet map[string]*User

// AddLDAPEntries iterates through the results of an LDAP search, adding found users to the UserSet.
func (u UserSet) addLDAPEntries(entries []*ldap.Entry, ldapUserSearch *LDAPUserSearch) {
	for _, entry := range entries {

		var user string
		var fullName string
		var email string
		var firstName string
		var lastName string

		for _, attr := range entry.Attributes {
			if attr.Name == ldapUserSearch.UserAttr {
				if len(attr.Values) != 0 {
					user = attr.Values[0]
				}
			} else if attr.Name == ldapUserSearch.FullNameAttr {
				if len(attr.Values) != 0 {
					fullName = attr.Values[0]
				}
			} else if attr.Name == ldapUserSearch.EmailAttr {
				if len(attr.Values) != 0 {
					email = attr.Values[0]
				}
			} else if attr.Name == ldapUserSearch.FirstNameAttr {
				if len(attr.Values) != 0 {
					firstName = attr.Values[0]
				}
			} else if attr.Name == ldapUserSearch.LastNameAttr {
				if len(attr.Values) != 0 {
					lastName = attr.Values[0]
				}
			}
		}

		if user == "" {
			log.Printf("Warning: Found DN %s but user attribute %s is an empty string", entry.DN, ldapUserSearch.UserAttr)
			continue
		}

		if _, ok := u[user]; ok {
			u[user].LDAP = true
		} else {
			u[user] = &User{LDAP: true}
		}

		u[user].Username = user
		u[user].FullName = fullName
		u[user].Email = email
		u[user].FirstName = firstName
		u[user].LastName = lastName
	}
}

// AddDuoResults iterates over a UsersResult from the Duo Admin API and marks the Duo attribute in a User in the UserSet
// to show that the user already exist in Duo.
func (u UserSet) addDuoResults(result *UsersResponse) {
	for _, dUser := range result.Response {
		if _, ok := u[dUser.Username]; ok {
			u[dUser.Username].Duo = true
			u[dUser.Username].DuoUserID = dUser.UserID
		} else {
			u[dUser.Username] = &User{Duo: true, Username: dUser.Username, DuoUserID: dUser.UserID}
		}
	}
}
