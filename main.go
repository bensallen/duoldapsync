package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/duosecurity/duo_api_golang"

	"github.com/spf13/pflag"
)

var configPath string
var debug bool
var dryRun bool

func init() {
	pflag.StringVarP(&configPath, "config", "f", "config.json", "Path to configuration file")
	pflag.BoolVarP(&debug, "debug", "d", false, "Enable debug output")
	pflag.BoolVarP(&dryRun, "dryrun", "n", false, "Dry-run mode, don't actually create or update users in Duo")
}

func main() {

	pflag.Parse()

	c, err := loadConfig(configPath)
	if err != nil {
		log.Printf("loadConfig error: %v\n", err)
		os.Exit(1)
	}

	l, err := connect(c.LDAPServers)
	if err != nil {
		log.Printf("Connection to LDAP server(s) failed: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	sr, err := enumUsers(l, c.LDAPUserSearch)
	if err != nil {
		log.Printf("Error running search on LDAP: %v\n", err)
		os.Exit(1)
	}

	if len(sr.Entries) == 0 {
		log.Printf("LDAP search returned no results")
		os.Exit(1)
	} else {
		if debug {
			log.Printf("LDAP found %d results", len(sr.Entries))
		}
	}

	us := UserSet{}
	us.AddLDAPEntries(sr.Entries, c.LDAPUserSearch)

	api := duoapi.NewDuoApi(c.DuoAPI.Ikey, c.DuoAPI.Skey, c.DuoAPI.APIHost, "Duoldapsync", duoapi.SetTimeout(10*time.Second))

	a := NewAdminAPI(*api)
	dUsers, err := a.Users(url.Values{})
	if err != nil {
		log.Printf("Duo Users Enumeration Fail, %s", err)
		os.Exit(1)
	} else if dUsers.Stat != "OK" {
		log.Printf("Duo API returned status when attemping user enumeration: %s", dUsers.Stat)
		os.Exit(1)
	}

	us.AddDuoResults(dUsers)

	errs := us.CreateDuoUsers(a, dryRun)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("Error while creating Duo User, %s", err)
		}
	}
	//for user, userSet := range us {
	//	fmt.Printf("User %s: %#v\n", user, *userSet)
	//}
}
