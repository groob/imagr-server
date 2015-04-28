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
	err := config.UpdateConfig(repoPath)
	if err != nil {
		log.Println("Failed to update imagr_config.plist")
		log.Println(err)
	}
}

func main() {
	if *cmd.ServeCmd == true {
		server.Serve(repoPath)
	}
}
