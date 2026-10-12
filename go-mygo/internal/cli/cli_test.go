package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/byteowlz/{{project_name}}/internal/service"
)

func TestCLIEnvelopeAndNoWrites(t *testing.T) {
	t.Chdir(t.TempDir())
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", home)
	for _, key := range []string{"MYGO_APP_THEME", "MYGO_APP_RADIUS", "MYGO_APP_OMARCHY_FILE"} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
	for _, tc := range []struct {
		args  []string
		input string
		exit  int
		code  string
	}{
		{[]string{"snapshot", "--json"}, "", 0, ""},
		{[]string{"validate", "--json"}, "invalid", 2, "INVALID_INPUT"},
		{[]string{"apply", "--json"}, "", 1, "AUTHORITY_DENIED"},
		{[]string{"bogus", "--json"}, "", 2, "USAGE"},
		{[]string{"snapshot", "--nonsense", "--json"}, "", 2, "USAGE"},
		{[]string{"snapshot", "--config", "missing", "--json"}, "", 2, "INVALID_CONFIG"},
	} {
		var out, errOut bytes.Buffer
		exit := Run(tc.args, bytes.NewBufferString(tc.input), &out, &errOut)
		var reply service.Envelope
		if err := json.Unmarshal(out.Bytes(), &reply); err != nil {
			t.Fatalf("%s: %v", out.String(), err)
		}
		if exit != tc.exit || (tc.code != "" && (reply.Error == nil || reply.Error.Code != tc.code)) {
			t.Fatalf("%v: exit %d %+v", tc.args, exit, reply)
		}
	}
	entries, _ := os.ReadDir(home)
	if len(entries) != 0 {
		t.Fatal("headless inspection created config/state")
	}
}

func TestDesktopParserMatchesFlags(t *testing.T) {
	options, err := DesktopOptions([]string{"--theme", "light", "--radius", "8"})
	if err != nil || *options.Theme != "light" || *options.Radius != 8 {
		t.Fatalf("%+v %v", options, err)
	}
}
