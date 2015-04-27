package main

import (
	"log"
	"os"

	"github.com/groob/imagr-server/cmd"
	"github.com/groob/imagr-server/imagr"
	"github.com/groob/imagr-server/server"
)

var (
	repoPath string
	password string
	config   imagr.ImagrConfig
)

func init() {
	password := os.Getenv("IMAGR_PASSWORD")
	if password == "" {
		log.Fatal("IMAGR_PASSWORD not set")
	}
	repoPath = *cmd.RepoPathCmd
	config.UpdateConfig(repoPath)
}

func main() {
	if *cmd.ServeCmd == true {
		server.Serve(repoPath)
	}
}
