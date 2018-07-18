package main

import (
	"fmt"
	"log"
	"net/url"

	ldap "gopkg.in/ldap.v2"
)

// User represents the attributes of a user found in LDAP and records if the user has been found in Duo.
type User struct {
	FullName  string
	Email     string
	FirstName string
	LastName  string

	LDAP        bool // Found in LDAP
	Duo         bool // Found in Duo
	NeedsUpdate bool // Indicates LDAP attributes are different that what is in Duo, and the Duo user needs to be updated.
}

// UserSet is a map of Users indexed by username
type UserSet map[string]*User

// AddLDAPEntries iterates through the results of an LDAP search, adding found users to the UserSet.
func (u UserSet) AddLDAPEntries(entries []*ldap.Entry, ldapUserSearch *LDAPUserSearch) {
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

		u[user].FullName = fullName
		u[user].Email = email
		u[user].FirstName = firstName
		u[user].LastName = lastName
	}
}

// AddDuoResults iterates over a UsersResult from the Duo Admin API and marks the Duo attribute in a User in the UserSet
// to show that the user already exist in Duo.
func (u UserSet) AddDuoResults(result *UsersResponse) {
	for _, user := range result.Response {
		if _, ok := u[user.Username]; ok {
			u[user.Username].Duo = true
		} else {
			u[user.Username] = &User{Duo: true}
		}
	}
}

func (u UserSet) CreateDuoUsers(api *AdminAPI, dryRun bool) []error {
	var errs []error

	for name, user := range u {
		if user.Duo == false && !dryRun {
			err := CreateDuoUser(api, name, user)
			if err != nil {
				errs = append(errs, err)
			}
		} else if user.Duo == false && dryRun {
			log.Printf("Creating Duo user: %s", name)
		}
	}
	return errs
}

func CreateDuoUser(api *AdminAPI, name string, user *User) error {

	params := user.URLValues(name)
	params.Set("status", "active")
	//params.Set("notes", ...) TODO: Add lastUpdated
	result, err := api.CreateUser(params)

	if err != nil {
		return fmt.Errorf("CreateUser failed: %s when attempting to create user: %s", err, name)
	} else if result.Stat != "OK" {
		return fmt.Errorf("CreateUser Duo API returned non-ok status when attemping to create user: %s with message: %s", name, result.Message)
	}

	return nil
}

// URLValues transforms User's attributes into url.Values
func (u *User) URLValues(username string) url.Values {
	params := url.Values{}
	params.Set("username", username)
	params.Set("realname", u.FullName)
	params.Set("email", u.Email)
	params.Set("firstname", u.FirstName)
	params.Set("lastname", u.LastName)

	return params
}
