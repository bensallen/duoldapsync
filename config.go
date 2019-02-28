package main

import (
	"os"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

// LDAPServer is a LDAP server
type LDAPServer struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	StartTLS     bool   `json:"start_tls"`
	BindDN       string `json:"bind_dn"`
	BindPassword string `json:"bind_password"`
}

// LDAPUserSearch is the config attributes to search for users in the LDAP tree
type LDAPUserSearch struct {
	BaseDN              string `json:"base_dn"`
	UserFilter          string `json:"user_filter"`
	UserAttr            string `json:"user_attr"` // Unique attribute to match an individual user
	GroupMembershipAttr string `json:"group_membership_attr"`
	EmailAttr           string `json:"email_attr"`
	FullNameAttr        string `json:"full_name_attr"`
	FirstNameAttr       string `json:"first_name_attr"`
	LastNameAttr        string `json:"last_name_attr"`
}

// LDAPGroupSearch is the config attributes to search for groups in the LDAP tree
type LDAPGroupSearch struct {
	BaseDN      string `json:"base_dn"`
	GroupFilter string `json:"group_filter"`
	GroupAttr   string `json:"group_attr"` // Attribute to match an individual group
}

// DuoAPI is the config attributes to access and control behavior with the Duo Admin API
type DuoAPI struct {
	DeleteUsers        bool   `json:"delete_users"`
	SendEnrollEmail    bool   `json:"send_enroll_email"`
	MaxDeleteUsers     int    `json:"max_delete_users"`
	EnrollValidSeconds int    `json:"enroll_valid_seconds"`
	Ikey               string `json:"ikey"`
	Skey               string `json:"skey"`
	APIHost            string `json:"api_host"`
	HTTPProxy          string `json:"http_proxy"`
}

// DuoLDAPSyncConfig is overall configuration struct for duoldapsync
type DuoLDAPSyncConfig struct {
	LDAPServers     []*LDAPServer
	LDAPUserSearch  *LDAPUserSearch
	LDAPGroupSearch *LDAPGroupSearch
	DuoAPI          *DuoAPI
}

func loadConfig(path string) (DuoLDAPSyncConfig, error) {
	// Create new config
	conf := config.NewConfig()

	c := DuoLDAPSyncConfig{}

	// Load file source
	f := file.WithPath(path)
	s := file.NewSource(f)

	if err := conf.Load(s); err != nil {
		return c, err
	}
	defer conf.Close()

	if err := conf.Get("servers").Scan(&c.LDAPServers); err != nil {
		return c, err
	}

	if err := conf.Get("user_search").Scan(&c.LDAPUserSearch); err != nil {
		return c, err
	}

	if err := conf.Get("group_search").Scan(&c.LDAPGroupSearch); err != nil {
		return c, err
	}

	if err := conf.Get("duo_api").Scan(&c.DuoAPI); err != nil {
		return c, err
	}

	if c.DuoAPI.HTTPProxy != "" {
		os.Setenv("HTTPS_PROXY", c.DuoAPI.HTTPProxy)
		os.Setenv("HTTP_PROXY", c.DuoAPI.HTTPProxy)
	}

	return c, nil
}
