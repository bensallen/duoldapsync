package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"

	ldap "gopkg.in/ldap.v2"
)

func connect(servers []*LDAPServer) (*ldap.Conn, error) {
	l := &ldap.Conn{}
	var connErrs []string

	for _, server := range servers {
		l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", server.Address, server.Port))
		if err != nil {
			connErrs = append(connErrs, err.Error())
			continue
		}

		// Reconnect with TLS
		// TODO: verify based on CA
		if server.StartTLS {
			err := l.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				connErrs = append(connErrs, err.Error())
				continue
			}
		}

		// First bind with a read only user
		//err = l.Bind(bindusername, bindpassword)
		//if err != nil {
		//	log.Fatal(err)
		//}

		if debug {
			log.Printf("LDAP connection successful: %v\n", *server)
		}

		return l, nil
	}

	return l, errors.New(strings.Join(connErrs, "\n"))
}

// enumUsers enumerates all users from LDAP
func enumUsers(l *ldap.Conn, c *LDAPUserSearch) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(
		c.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s)", c.UserFilter),
		[]string{c.UserAttr, c.GroupMembershipAttr, c.EmailAttr, c.FullNameAttr, c.FirstNameAttr, c.LastNameAttr},
		nil,
	)

	if debug {
		log.Printf("LDAP executing search: %v\n", searchRequest)
	}

	return l.Search(searchRequest)
}

/*
func groupSearch(group string, l *ldap.Conn, c *LDAPGroupSearch) (*ldap.SearchResult, error) {
	// Search for the given group
	searchRequest := ldap.NewSearchRequest(
		c.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(%s)(%s=%s))", c.GroupFilter, c.GroupAttr, group),
		[]string{"dn"},
		nil,
	)

	if debug {
		log.Printf("LDAP executing search: %v\n", searchRequest)
	}

	return l.Search(searchRequest)
}

// enumGroups enumerates all groups from LDAP
func enumGroups(l *ldap.Conn, c *LDAPGroupSearch) (*ldap.SearchResult, error) {
	// Search for the given group
	searchRequest := ldap.NewSearchRequest(
		c.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s)", c.GroupFilter),
		[]string{"dn"},
		nil,
	)

	if debug {
		log.Printf("LDAP executing search: %v\n", searchRequest)
	}

	return l.Search(searchRequest)
}
*/
