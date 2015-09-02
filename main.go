package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/groob/imagr-server/imagr"
	"github.com/groob/imagr-server/server"
)

var (
	password string
	repoPath string
	serveCmd bool
	config   imagr.ImagrConfig
)

func init() {
	flag.StringVar(&repoPath, "repo", "/imagr_repo", "path to imagr repo")
	flag.BoolVar(&serveCmd, "serve", false, "serve the repo over http")
	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
	password := os.Getenv("IMAGR_PASSWORD")
	if password == "" {
		log.Fatal("IMAGR_PASSWORD not set")
	}
	err := config.UpdateConfig(repoPath)
	if err != nil {
		log.Println("Failed to update imagr_config.plist")
		log.Println(err)
	}
}

func serveWeb() {
	log.Fatal(http.ListenAndServe(":3001",
		http.FileServer(http.Dir("web"))))
}
func main() {
	//	go serveWeb()
	if serveCmd {
		server.Serve(repoPath)
	}
}
