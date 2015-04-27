package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/groob/imagr-server/imagr"
)

func wfHandler(repoPath string) http.HandlerFunc {
	// Handler
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/v1/workflows/"):]
		switch key {
		case "":
			w.Write([]byte("No Workflow specified"))
		case "all":
			var workflows []string
			workflowPath := fmt.Sprintf("%v/workflows/", repoPath)
			files, _ := ioutil.ReadDir(workflowPath)
			for _, f := range files {
				if strings.SplitN(f.Name(), ".", 2)[1] == "plist" {
					workflows = append(workflows, f.Name())
				}
			}
			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			err := encoder.Encode(&workflows)
			if err != nil {
				log.Println(err)
			}
			jsn, err := json.MarshalIndent(workflows, "", "\t")
			if err != nil {
				log.Println(err)
				w.Write([]byte("There was an error. Check the logs"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(jsn)
		default:
			var config imagr.ImagrConfig
			wfName := fmt.Sprintf("%v/workflows/%v.plist", repoPath, key)
			switch r.Method {
			case "GET":
				workflow, err := imagr.ParseWorkflow(wfName)
				if err != nil {
					w.Write([]byte("Workflow does not exist."))
					log.Println(err)
					break
				}
				jsn, err := json.MarshalIndent(workflow, "", "\t")
				if err != nil {
					log.Println(err)
					w.Write([]byte("There was an error. Check the logs"))
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(jsn)
			case "PUT":
				var workflow imagr.Workflow
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&workflow)
				if err != nil {
					log.Println(err)
				}
				err = workflow.Save(wfName)
				if err != nil {
					w.Write([]byte("Workflow could not be saved."))
				}
				err = config.UpdateConfig(repoPath)
				if err != nil {
					log.Println("Failed to update config.")
				}
			case "DELETE":
				log.Println("Deleting workflow")
				err := os.Remove(wfName)
				if err != nil {
					w.Write([]byte("Could not delete workflow."))
					log.Println(err)
				}
				err = config.UpdateConfig(repoPath)
				if err != nil {
					log.Println("Failed to update config.")
				}
			default:
				w.Write([]byte("Method not supported."))
			}
		}
	}
}
func Serve(repoPath string) {
	http.Handle("/", http.FileServer(http.Dir(repoPath)))
	http.HandleFunc("/v1/workflows/", wfHandler(repoPath))
	http.ListenAndServe(":3000", nil)
}
