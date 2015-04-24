package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/groob/imagr-server/imagr"
)

var (
	repoPath string
	password string
	Imagr    imagr.ImagrConfig
)

func init() {
	flag.Parse()
	password := os.Getenv("IMAGR_PASSWORD")
	if password == "" {
		log.Fatal("IMAGR_PASSWORD not set")
	}
	repoPath = fmt.Sprintf("%v", flag.Arg(0)) // file location
	if repoPath == "" {
		log.Fatal("Please specify a path to Imagr repo")
	}
	Imagr.Password = imagr.EncodePassword(password)
	Imagr.Workflows = imagr.ParseWorkflows(repoPath)
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
	genConfig()
	log.Println("Starting server on port 3000")
	http.Handle("/", http.FileServer(http.Dir(repoPath)))
	http.ListenAndServe(":3000", nil)
}
