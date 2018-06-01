package main

import (
	"log"
	"os"

	"github.com/spf13/pflag"
)

var configPath string
var debug bool

func init() {
	pflag.StringVarP(&configPath, "config", "f", "config.json", "Path to configuration file")
	pflag.BoolVarP(&debug, "debug", "d", false, "Enable debug output")
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
		log.Printf("Connection to server(s) failed: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	sr, err := enumUsers(l, c.LDAPUserSearch)
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	if len(sr.Entries) == 0 {
		log.Fatal("No results")
	} else {
		log.Printf("Found %d results", len(sr.Entries))
	}

	for _, entry := range sr.Entries {
		entry.Print()
	}
}
