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
		case "all":
			workflows := imagr.ParseWorkflows(repoPath)
			jsn, err := json.MarshalIndent(workflows, "", "\t")
			if err != nil {
				log.Println(err)
				w.Write([]byte("There was an error. Check the logs"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(jsn)
		default:
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
			case "DELETE":
				log.Println("Deleting workflow")
				err := os.Remove(wfName)
				if err != nil {
					w.Write([]byte("Could not delete workflow."))
					log.Println(err)
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
