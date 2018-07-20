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

	conf, err := loadConfig(configPath)
	if err != nil {
		log.Printf("loadConfig error: %v\n", err)
		os.Exit(1)
	}

	l, err := connect(conf.LDAPServers)
	if err != nil {
		log.Printf("Connection to LDAP server(s) failed: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	sr, err := enumUsers(l, conf.LDAPUserSearch)
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

	userSet := UserSet{}
	userSet.addLDAPEntries(sr.Entries, conf.LDAPUserSearch)

	duoAPI := duoapi.NewDuoApi(conf.DuoAPI.Ikey, conf.DuoAPI.Skey, conf.DuoAPI.APIHost, "Duoldapsync", duoapi.SetTimeout(10*time.Second))

	adminAPI := NewAdminAPI(*duoAPI, dryRun)
	duoUsers, err := adminAPI.Users(url.Values{})
	if err != nil {
		log.Printf("Duo Users Enumeration Fail, %s", err)
		os.Exit(1)
	} else if duoUsers.Stat != "OK" {
		log.Printf("Duo API returned status when attemping user enumeration: %s", duoUsers.Stat)
		os.Exit(1)
	}

	userSet.addDuoResults(duoUsers)

	for _, user := range userSet {
		if user.Duo == false {
			if debug {
				log.Printf("Creating Duo user: %s", user.Username)
			}
			err := user.duoCreate(adminAPI)
			if err != nil {
				log.Printf("Duo User Creation Failed, %s", err)
				break
			}
			if conf.DuoAPI.SendEnrollEmail {
				if debug {
					log.Printf("Enrolling Duo user: %s", user.Username)
				}
				err := user.duoEnroll(adminAPI, conf.DuoAPI.EnrollValidSeconds)
				if err != nil {
					log.Printf("Duo User Enrollment Failed, %s", err)
				}
			}
		} else if user.Duo == true && user.LDAP == false && conf.DuoAPI.DeleteUsers == true {
			// Cleanup Duo Accounts
			if debug {
				log.Printf("Deleting Duo user: %s", user.Username)
			}
			resp, err := adminAPI.DeleteUser(user.DuoUserID)
			if err != nil {
				log.Printf("Duo User Delete Fail, %s", err)
			} else if resp.Stat != "OK" {
				log.Printf("Duo API returned status %d when attemping to delete user %s", resp.Code, user.Username)
			}
		}
	}
	//for user, userSet := range us {
	//	fmt.Printf("User %s: %#v\n", user, *userSet)
	//}

}
