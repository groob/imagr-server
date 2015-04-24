package imagr

import (
	"crypto/sha512"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"howett.net/plist"
)

var Workflows []Workflow

type WorkflowComponent struct {
	Type string `plist:"type" json:"type"`
	URL  string `plist:"url" json:"url"`
}

type Workflow struct {
	Name        string              `plist:"name" json:"name"`
	Description string              `plist:"description" json:"description"`
	Components  []WorkflowComponent `plist:"components" json:"components"`
}

type ImagrConfig struct {
	Password  string     `plist:"password", json:"password"`
	Workflows []Workflow `plist:"workflows" json:"workflows"`
}

// Decodes plist into struct
func (p *Workflow) DecodePlist(f *os.File) error {
	decoder := plist.NewDecoder(f)
	err := decoder.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

// Decodes plist into struct
func (p *ImagrConfig) DecodePlist(f *os.File) error {
	decoder := plist.NewDecoder(f)
	err := decoder.Decode(p)
	if err != nil {
		return err
	}
	return nil
}

// Encode a plist and write to file
func (p *ImagrConfig) EncodePlist(f *os.File) error {
	encoder := plist.NewEncoder(f)
	encoder.Indent("  ")
	err := encoder.Encode(p)
	if err != nil {
		return err
	}
	return nil
}

func EncodePassword(p string) string {
	pass512 := sha512.New()
	pass512.Write([]byte("password"))
	password := fmt.Sprintf("%x", pass512.Sum(nil))
	return password
}

func isDirectory(path string) (bool, error) {
	// return true if path is a directory
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func walkpath(path string, f os.FileInfo, err error) error {
	if fileInfo, _ := isDirectory(path); fileInfo == false {
		log.Printf("Parsing workflow: %s\n", path)
		workflow, err := parseWorkflow(path)
		if err != nil {
			return err
		}
		Workflows = append(Workflows, workflow)
	}
	return nil
}

func parseWorkflow(path string) (Workflow, error) {
	var workflow Workflow
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = workflow.DecodePlist(f)
	if err != nil {
		log.Fatal(err)
	}
	return workflow, nil
}

func ParseWorkflows(repoPath string) (workflows []Workflow) {
	workflowPath := fmt.Sprintf("%v/workflows", repoPath) // repo location
	filepath.Walk(workflowPath, walkpath)
	return Workflows
}
