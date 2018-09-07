package main

import (
	"log"
	"os"

	"github.com/pkg/profile"
	"github.com/spf13/pflag"
)

var configPath string
var debug bool
var dryRun bool
var pollTime int
var profileOut string

func init() {
	pflag.StringVarP(&configPath, "config", "f", "config.json", "Path to configuration file")
	pflag.BoolVarP(&debug, "debug", "d", false, "Enable debug output")
	pflag.BoolVarP(&dryRun, "dryrun", "n", false, "Dry-run mode, don't actually create or delete users in Duo")
	pflag.IntVarP(&pollTime, "poll", "p", 600, "Number of seconds to wait between polling LDAP and Duo for changes")
	pflag.StringVarP(&profileOut, "profile", "P", "", "Enable cpu, mem, or block profiling")
}

func main() {
	pflag.Parse()

	if profileOut != "" {
		switch profileOut {
		case "cpu":
			defer profile.Start(profile.CPUProfile).Stop()
		case "mem":
			defer profile.Start(profile.MemProfile).Stop()
		case "block":
			defer profile.Start(profile.BlockProfile).Stop()
		default:
			// do nothing
		}
	}

	conf, err := loadConfig(configPath)
	if err != nil {
		log.Printf("loadConfig error: %v\n", err)
		os.Exit(1)
	}

	if err := run(conf); err != nil {
		log.Printf("Run error: %v\n", err)
		os.Exit(1)
	}
}
