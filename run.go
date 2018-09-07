package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"gopkg.in/ldap.v2"

	"github.com/duosecurity/duo_api_golang"
)

func run(conf DuoLDAPSyncConfig) error {
	l, err := connect(conf.LDAPServers)
	if err != nil {
		return fmt.Errorf("connection to LDAP server(s) failed: %v", err)
	}
	defer l.Close()

	duoAPI := duoapi.NewDuoApi(conf.DuoAPI.Ikey, conf.DuoAPI.Skey, conf.DuoAPI.APIHost, "Duoldapsync", duoapi.SetTimeout(10*time.Second))
	adminAPI := NewAdminAPI(*duoAPI, dryRun)

	// Loop forever sleeping pollTime seconds between iterations.
	ticker := time.NewTicker(time.Second * time.Duration(pollTime))
	done := make(chan bool)

	go tickerLoop(ticker, conf, l, adminAPI, done)
	<-done

	return nil
}

func tickerLoop(ticker *time.Ticker, conf DuoLDAPSyncConfig, ldapConn *ldap.Conn, adminAPI *AdminAPI, done chan bool) {
	for range ticker.C {
		sr, err := enumUsers(ldapConn, conf.LDAPUserSearch)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		if debug {
			log.Printf("LDAP found %d results", len(sr.Entries))
		}

		userSet := UserSet{}
		userSet.addLDAPEntries(sr.Entries, conf.LDAPUserSearch)

		duoUsers, err := adminAPI.Users(url.Values{})
		if err != nil {
			log.Printf("Duo Users Enumeration Fail, %s", err)
			continue
		} else if duoUsers.Stat != "OK" {
			log.Printf("Duo API returned status when attemping user enumeration: %s", duoUsers.Stat)
			continue
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
	}
	done <- true
}
