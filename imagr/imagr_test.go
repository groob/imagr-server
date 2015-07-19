package imagr

import (
	"bytes"
	"testing"
)

var w Workflow
var sampleWorkflow = []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>components</key>
    <array>
      <dict>
        <key>type</key>
        <string>image</string>
        <key>url</key>
        <string>http://imagr/images/BaseImage-10.10.3-14D131.hfs.dmg</string>
      </dict>
    </array>
    <key>description</key>
    <string>Deploys the latest 10.10.x image. Munki tools and local admin account included.</string>
    <key>name</key>
    <string>Yosemite - MunkiTools</string>
    <key>restart_action</key>
    <string>restart</string>
    <key>bless_target</key>
    <false/>
  </dict>
</plist>
`)

func TestEncodePassword(t *testing.T) {
	passDict := map[string]string{
		"password": "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86",
		"foo":      "f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7",
		"bar!":     "56c79f1c6e391260bce4418f48fa72b15d2402f78dcfeab5ad5a0fa9e7826d042f534baa2f61557163dbf2b3a40d4f66936cb84e3fd7304e69fbc8759d60b9f9",
	}
	for key, value := range passDict {
		encoded := EncodePassword(key)
		if encoded != value {
			t.Error("Password not hashed correctly")
		}
	}
}

func TestDecodePlist(t *testing.T) {
	workflow := bytes.NewReader(sampleWorkflow)
	err := w.DecodePlist(workflow)
	if err != nil {
		t.Error(err)
	}
	// Test that names match.
	if w.Name != "Yosemite - MunkiTools" {
		t.Fatal(err)
	}
}

func TestEncodePlist(t *testing.T) {
	out := []byte(``)
	b := bytes.NewBuffer(out)
	err := w.EncodePlist(b)
	if err != nil {
		t.Fatal(err)
	}
}
