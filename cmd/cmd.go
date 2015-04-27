package cmd

import (
	"flag"
	"os"
)

var RepoPathCmd = flag.String("repo", "/imagr_repo", "path to imagr repo")
var ServeCmd = flag.Bool("serve", false, "serve the repo over http")

func init() {
	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
}
