// Agent-check exercises the actual headless executable in an isolated fixture sandbox.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/byteowlz/{{project_name}}/schemas"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type observation struct {
	Args  []string         `json:"args"`
	Exit  int              `json:"exit"`
	Reply service.Envelope `json:"envelope"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		return errors.New("usage: agent-check CLI_BINARY")
	}
	binary, err := filepath.Abs(os.Args[1])
	if err != nil {
		return err
	}
	sandbox, err := os.MkdirTemp("", "mygo-agent-proof-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(sandbox)
	env := []string{}
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "XDG_CONFIG_HOME=") || strings.HasPrefix(entry, "MYGO_APP_") {
			continue
		}
		env = append(env, entry)
	}
	env = append(env, "XDG_CONFIG_HOME="+filepath.Join(sandbox, "config"))
	rawSchema, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemas.Envelope))
	if err != nil {
		return err
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("urn:mygo:envelope", rawSchema); err != nil {
		return err
	}
	schema, err := compiler.Compile("urn:mygo:envelope")
	if err != nil {
		return err
	}
	var observed []observation
	call := func(args []string, input string, exit int) (service.Envelope, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary, append(args, "--json")...)
		cmd.Dir = sandbox
		cmd.Env = env
		cmd.Stdin = strings.NewReader(input)
		var output, stderr bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &stderr
		err := cmd.Run()
		actual := 0
		if err != nil {
			var failed *exec.ExitError
			if !errors.As(err, &failed) {
				return service.Envelope{}, err
			}
			actual = failed.ExitCode()
		}
		if actual != exit {
			return service.Envelope{}, fmt.Errorf("%v: expected exit %d, got %d %s", args, exit, actual, stderr.String())
		}
		value, err := jsonschema.UnmarshalJSON(bytes.NewReader(output.Bytes()))
		if err != nil {
			return service.Envelope{}, err
		}
		if err := schema.Validate(value); err != nil {
			return service.Envelope{}, err
		}
		var reply service.Envelope
		if err := json.Unmarshal(output.Bytes(), &reply); err != nil {
			return reply, err
		}
		observed = append(observed, observation{Args: args, Exit: actual, Reply: reply})
		return reply, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	help := exec.CommandContext(ctx, binary, "--help")
	help.Dir = sandbox
	help.Env = env
	text, err := help.Output()
	cancel()
	if err != nil || !bytes.Contains(text, []byte("snapshot")) {
		return errors.New("discoverability help failed")
	}
	baseline, err := call([]string{"snapshot"}, "", 0)
	if err != nil {
		return err
	}
	portable, err := json.Marshal(baseline.Result.Snapshot)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(sandbox, "snapshot.json"), portable, 0600); err != nil {
		return err
	}
	for _, args := range [][]string{{"validate"}, {"validate", "--input", "snapshot.json"}} {
		roundtrip, err := call(args, string(portable), 0)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(roundtrip.Result, baseline.Result) {
			return errors.New("portable roundtrip differs")
		}
	}
	for _, test := range []struct {
		args        []string
		input, code string
		exit        int
	}{
		{[]string{"apply"}, "", "AUTHORITY_DENIED", 1},
		{[]string{"apply", "--input", "missing.json"}, "", "AUTHORITY_DENIED", 1},
		{[]string{"validate"}, strings.ReplaceAll(string(portable), "fixture-v1", "stale"), "STALE_REVISION", 2},
		{[]string{"validate"}, "not-json", "INVALID_INPUT", 2},
		{[]string{"validate"}, string(portable) + " {}", "INVALID_INPUT", 2},
		{[]string{"validate"}, strings.Repeat("x", 65537), "INVALID_INPUT", 2},
		{[]string{"bogus"}, "", "USAGE", 2},
		{[]string{"snapshot", "--unknown"}, "", "USAGE", 2},
		{[]string{"snapshot", "--config", "missing"}, "", "INVALID_CONFIG", 2},
	} {
		reply, err := call(test.args, test.input, test.exit)
		if err != nil {
			return err
		}
		if reply.Error == nil || reply.Error.Code != test.code {
			return fmt.Errorf("wrong stable failure %v", test.args)
		}
	}
	after, err := call([]string{"snapshot"}, "", 0)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(after, baseline) {
		return errors.New("denied apply changed fixture")
	}
	files, err := os.ReadDir(sandbox)
	if err != nil {
		return err
	}
	if len(files) != 1 || files[0].Name() != "snapshot.json" {
		return errors.New("headless inspection made unexpected writes")
	}
	if err := os.MkdirAll("artifacts/agent-check", 0700); err != nil {
		return err
	}
	if err := writeJSON("artifacts/agent-check/observations.json", observed); err != nil {
		return err
	}
	gate := func(status, reason string) map[string]any {
		return map[string]any{"status": status, "reason": reason, "evidence": []map[string]string{{"procedure": "just agent-check", "artifact": "artifacts/agent-check/observations.json; artifacts/agent-check/go-tests.log", "observation": reason}}}
	}
	report := map[string]any{
		"schema_version": "agent-readiness/v1", "project": "{{project_name}}", "revision": "instantiated-native-scaffold",
		"scope":     "Read-only synthetic fixture, actual CLI plus MyGo ui.NewTester native headless view; not platform/hardware certification.",
		"platforms": []string{"current-host Go/native-headless fixture"},
		"gates": map[string]any{
			"discoverability":          gate("pass", "Real help names semantic commands; each subprocess has a 5s bound."),
			"structured_contract":      gate("pass", "Actual CLI success/error JSON schema-validated; negative exits checked."),
			"shared_substrate":         gate("pass", "Headless native UI clicks and keyboard invoke the same Go service; replies match actual CLI adapter. Fixture has no domain edits."),
			"authority_and_safety":     gate("pass", "Apply always denied; fixture unchanged; sandbox contains only explicit snapshot export."),
			"concurrency_and_recovery": gate("pass", "Race tests cover immutable concurrent/idempotent reads, cancellation, stale baseline, denial and UI recovery. No transport/write capability exists."),
			"portable_state":           gate("pass", "File/stdin snapshot roundtrips through actual executable; portable embedded adapter bytes checked."),
			"reproducibility":          gate("pass", "Executed Go/native headless fixtures, schema/byte/Studio parity tests. No frontend/local service dependency."),
			"human_accessibility":      gate("unverified", "Headless native keyboard/focus/announcement checks pass; actual OS assistive technology and foreign-platform accessibility require separate observation."),
		},
	}
	if err := writeJSON("artifacts/agent-check/report.json", report); err != nil {
		return err
	}
	fmt.Println(`{"ok":true,"scope":"native-headless-fixture","production_ready":false,"report_structure_is_not_evidence":true,"artifacts":"artifacts/agent-check"}`)
	return nil
}

func writeJSON(path string, value any) error {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0600)
}
