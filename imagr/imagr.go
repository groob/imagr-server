package imagr

import (
	"crypto/sha512"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

var (
	Workflows []Workflow
	password  string
)

func init() {
	password = os.Getenv("IMAGR_PASSWORD")
}

type WorkflowComponent struct {
	Type string `plist:"type" json:"type"`
	URL  string `plist:"url" json:"url"`
}

type Workflow struct {
	ID            string              `plist:"-" json:"id"`
	Name          string              `plist:"name" json:"name"`
	Description   string              `plist:"description" json:"description"`
	Components    []WorkflowComponent `plist:"components" json:"components"`
	RestartAction string              `plist:"restart_action,omitempty" json:"restart_action,omitempty"`
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

func (p *Workflow) EncodePlist(f *os.File) error {
	encoder := plist.NewEncoder(f)
	encoder.Indent("  ")
	err := encoder.Encode(p)
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

func wfWalkFn(path string, f os.FileInfo, err error) error {
	if fileInfo, _ := isDirectory(path); fileInfo == false {
		if filepath.Ext(path) == ".plist" {
			log.Printf("Parsing workflow: %s\n", path)
			workflow, err := ParseWorkflow(path)
			if err != nil {
				return err
			}
			Workflows = append(Workflows, workflow)
		}
	}
	return nil
}

func ParseWorkflow(path string) (Workflow, error) {
	var workflow Workflow
	f, err := os.Open(path)
	if err != nil {
		return Workflow{}, err
	}
	defer f.Close()
	basename, _ := f.Stat()
	workflow.ID = strings.TrimSuffix(basename.Name(),
		filepath.Ext(basename.Name())) // Get ID from FileName
	err = workflow.DecodePlist(f)
	if err != nil {
		return Workflow{}, err
	}
	return workflow, nil
}

func ParseWorkflows(repoPath string) (workflows []Workflow) {
	Workflows = []Workflow{}                              // reset slice
	workflowPath := fmt.Sprintf("%v/workflows", repoPath) // repo location
	filepath.Walk(workflowPath, wfWalkFn)
	return Workflows
}

func (w *Workflow) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = w.EncodePlist(f)
	if err != nil {
		return err
	}
	return nil
}

func (i *ImagrConfig) UpdateConfig(repoPath string) error {
	i.Password = EncodePassword(password)
	i.Workflows = ParseWorkflows(repoPath)
	configFile := fmt.Sprintf("%v/imagr_config.plist", repoPath)
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()
	err = i.EncodePlist(f)
	if err != nil {
		return err
	}
	return nil
}
