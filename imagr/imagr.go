package imagr

import (
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

var (
	password string
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
	BlessTarget   bool                `plist:"bless_target,omitempty" json:"bless_target,omitempty"`
}

type ImagrConfig struct {
	Password  string     `plist:"password", json:"password"`
	Workflows []Workflow `plist:"workflows" json:"workflows"`
}

// Decodes plist into struct
func (p *Workflow) DecodePlist(r io.ReadSeeker) error {
	return plist.NewDecoder(r).Decode(p)
}

func (p *Workflow) EncodePlist(w io.Writer) error {
	encoder := plist.NewEncoder(w)
	encoder.Indent("  ")
	return encoder.Encode(p)
}

// Decodes plist into struct
func (p *ImagrConfig) DecodePlist(r io.ReadSeeker) error {
	return plist.NewDecoder(r).Decode(p)
}

// Encode a plist and write to file
func (p *ImagrConfig) EncodePlist(w io.Writer) error {
	encoder := plist.NewEncoder(w)
	encoder.Indent("  ")
	return encoder.Encode(p)
}

func EncodePassword(p string) string {
	pass512 := sha512.New()
	pass512.Write([]byte(p))
	password := fmt.Sprintf("%x", pass512.Sum(nil))
	return password
}

func ParseWorkflow(path string) (Workflow, error) {
	var workflow Workflow
	f, err := os.Open(path)
	if err != nil {
		return workflow, err
	}
	defer f.Close()
	basename, _ := f.Stat()
	workflow.ID = strings.TrimSuffix(basename.Name(),
		filepath.Ext(basename.Name())) // Get ID from FileName
	return workflow, workflow.DecodePlist(f)
}

func ParseWorkflows(repoPath string) (workflows []Workflow) {
	workflows = []Workflow{}                              // reset slice
	workflowPath := fmt.Sprintf("%v/workflows", repoPath) // repo location
	files, err := ioutil.ReadDir(workflowPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		//only attempt to read plist files in the workflows directory.
		if !f.IsDir() && filepath.Ext(f.Name()) == ".plist" {
			wfPath := fmt.Sprintf("%v/%v", workflowPath, f.Name())
			log.Println("Reading workflow " + f.Name())
			workflow, err := ParseWorkflow(wfPath)
			if err != nil {
				log.Println("Could not parse workflow " + f.Name())
			}
			workflows = append(workflows, workflow)
		}
	}
	return workflows
}

func (w *Workflow) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return w.EncodePlist(f)
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
	return i.EncodePlist(f)
}
