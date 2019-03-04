package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/ldap.v2"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/admin"
)

func run(conf DuoLDAPSyncConfig, dryRun bool) error {
	l, err := connect(conf.LDAPServers)
	if err != nil {
		return fmt.Errorf("connection to LDAP server(s) failed: %v", err)
	}
	defer l.Close()

	duoAPI := duoapi.NewDuoApi(conf.DuoAPI.Ikey, conf.DuoAPI.Skey, conf.DuoAPI.APIHost, "Duoldapsync", duoapi.SetTimeout(10*time.Second))
	client := admin.New(*duoAPI)

	// Loop forever sleeping pollTime seconds between iterations.
	ticker := time.NewTicker(time.Second * time.Duration(pollTime))
	done := make(chan bool)

	go tickerLoop(ticker, conf, l, client, dryRun, done)

	// Wait for tickerLoop to exit
	<-done
	close(done)

	return nil
}

func tickerLoop(ticker *time.Ticker, conf DuoLDAPSyncConfig, ldapConn *ldap.Conn, client *admin.Client, dryRun bool, done chan bool) {

	// MaxDeleteUsers needs to be 1 or greater to make sense. Disable DeleteUsers
	// to disable user deletion instead of trying to set MaxDeleteUsers to 0.
	maxDeleteUsers := 1
	if conf.DuoAPI.MaxDeleteUsers > 0 {
		maxDeleteUsers = conf.DuoAPI.MaxDeleteUsers
	}

	for range ticker.C {
		sr, err := enumUsers(ldapConn, conf.LDAPUserSearch)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		if debug {
			log.Printf("LDAP found %d results", len(sr.Entries))
		}

		// Skip the rest of the loop so we avoid deleting all Duo users accidently
		if len(sr.Entries) == 0 {
			log.Print("WARNING no LDAP results found, skipping")
			continue
		}

		userSet := UserSet{}
		userSet.addLDAPEntries(sr.Entries, conf.LDAPUserSearch)

		duoUsers, err := client.GetUsers()
		if err != nil {
			log.Printf("Duo Users Enumeration Fail, %s", err)
			continue
		} else if duoUsers.Stat != "OK" {
			log.Printf("Duo API returned status when attemping user enumeration: %s", duoUsers.Stat)
			continue
		}

		userSet.addDuoResults(duoUsers)

		usersDeleteOverflow := false
		usersDelete := make([]*User, 0, maxDeleteUsers)

		for _, user := range userSet {
			if !user.Duo {
				if debug {
					log.Printf("Creating Duo user: %s", user.Username)
				}
				err := user.duoCreate(client, dryRun)
				if err != nil {
					log.Printf("Duo User Creation Failed, %s", err)
					break
				}
				if conf.DuoAPI.SendEnrollEmail {
					if debug {
						log.Printf("Enrolling Duo user: %s", user.Username)
					}
					err := user.duoEnroll(client, conf.DuoAPI.EnrollValidSeconds, dryRun)
					if err != nil {
						log.Printf("Duo User Enrollment Failed, %s", err)
					}
				}
			} else if user.Duo && !user.LDAP && conf.DuoAPI.DeleteUsers {
				if len(usersDelete) < cap(usersDelete) {
					usersDelete = append(usersDelete, user)
				} else {
					usersDeleteOverflow = true
				}
			}
		}

		// Clear usersDelete if overflown (eg. gone over MaxDeleteUsers). Otherwise delete users.
		if usersDeleteOverflow {
			log.Printf("WARNING more users to delete than the configured DuoAPI.MaxDeleteUsers setting of %d allows, no users will be deleted", maxDeleteUsers)
			usersDelete = make([]*User, maxDeleteUsers)
		} else if len(usersDelete) > 0 {
			deleteUsers(client, usersDelete, dryRun)
		}

		// Allow user deletes to retry next tick.
		usersDeleteOverflow = false
	}

	// Tell run() tickerLoop is done
	done <- true
}

func deleteUsers(client *admin.Client, users []*User, dryRun bool) {
	for _, user := range users {
		if debug {
			log.Printf("Deleting Duo user: %s", user.Username)
		}
		resp, err := DeleteUser(client, user.DuoUserID, dryRun)
		if err != nil {
			log.Printf("Duo User Delete Fail, %s", err)
		} else if resp.Stat != "OK" {
			log.Printf("Duo API returned status %d when attemping to delete user %s", resp.Code, user.Username)
		}
	}
}
