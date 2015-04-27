package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
			jsn, err := json.MarshalIndent(imagr.ParseWorkflows(repoPath), "", "\t")
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
				w.Write([]byte(key + "\n")) //return the UUID created.
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
				w.Write([]byte(key + "\n")) //return the UUID deleted.
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
