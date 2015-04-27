package main

import (
	"fmt"
	"log"
	"os"

	"github.com/groob/imagr-server/cmd"
	"github.com/groob/imagr-server/imagr"
	"github.com/groob/imagr-server/server"
)

var (
	repoPath string
	password string
	Imagr    imagr.ImagrConfig
)

func init() {
	password := os.Getenv("IMAGR_PASSWORD")
	if password == "" {
		log.Fatal("IMAGR_PASSWORD not set")
	}
	repoPath = *cmd.RepoPathCmd
	Imagr.Password = imagr.EncodePassword(password)
	Imagr.Workflows = imagr.ParseWorkflows(repoPath)
	genConfig()
}

func genConfig() {
	configFile := fmt.Sprintf("%v/imagr_config.plist", repoPath)
	f, err := os.Create(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = Imagr.EncodePlist(f)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if *cmd.ServeCmd == true {
		server.Serve(repoPath)
	}
}
