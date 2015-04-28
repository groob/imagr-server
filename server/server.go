package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/groob/imagr-server/imagr"
)

func Serve(repoPath string) {
	http.Handle("/", http.FileServer(http.Dir(repoPath)))
	http.Handle("/v1/workflows/", &wfHandler{repoPath: repoPath})
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type wfHandler struct {
	repoPath string
}

func (wf *wfHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/v1/workflows/"):]
	switch key {
	case "":
		w.Write([]byte("No workflow selected.\n"))
	case "all":
		wf.serveAll(w, r)
	default:
		wf.serveWorkflow(w, r)
	}
}

func (wf *wfHandler) serveAll(w http.ResponseWriter, r *http.Request) {
	jsn, err := json.MarshalIndent(imagr.ParseWorkflows(wf.repoPath), "", "\t")
	if err != nil {
		log.Println(err)
		w.Write([]byte("There was an error. Check the logs"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsn)
}

func (wf *wfHandler) serveWorkflow(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		wf.serveWorkflowGET(w, r)
	case "PUT":
		wf.serveWorkflowPUT(w, r)
	case "DELETE":
		wf.serveWorkflowDELETE(w, r)
	default:
		w.Write([]byte("Method not supported.\n"))
	}
}

func (wf *wfHandler) serveWorkflowGET(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/v1/workflows/"):]
	wfName := fmt.Sprintf("%v/workflows/%v.plist", wf.repoPath, key)
	workflow, err := imagr.ParseWorkflow(wfName)
	if err != nil {
		w.Write([]byte("Workflow does not exist.\n"))
		log.Println(err)
	}
	jsn, err := json.MarshalIndent(workflow, "", "\t")
	if err != nil {
		log.Println(err)
		w.Write([]byte("There was an error. Check the logs"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsn)
}

func (wf *wfHandler) serveWorkflowPUT(w http.ResponseWriter, r *http.Request) {
	var config imagr.ImagrConfig
	key := r.URL.Path[len("/v1/workflows/"):]
	wfName := fmt.Sprintf("%v/workflows/%v.plist", wf.repoPath, key)
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
	err = config.UpdateConfig(wf.repoPath)
	if err != nil {
		log.Println("Failed to update config.")
	}
	w.Write([]byte(key + "\n")) //return the UUID created.
}

func (wf *wfHandler) serveWorkflowDELETE(w http.ResponseWriter, r *http.Request) {
	var config imagr.ImagrConfig
	key := r.URL.Path[len("/v1/workflows/"):]
	wfName := fmt.Sprintf("%v/workflows/%v.plist", wf.repoPath, key)
	log.Println("Deleting workflow")
	err := os.Remove(wfName)
	if err != nil {
		w.Write([]byte("Could not delete workflow."))
		log.Println(err)
	}
	err = config.UpdateConfig(wf.repoPath)
	if err != nil {
		log.Println("Failed to update config.")
	}
	w.Write([]byte(key + "\n")) //return the UUID deleted.
}
