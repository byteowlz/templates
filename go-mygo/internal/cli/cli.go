// Package cli is a thin semantic adapter; business rules stay in service.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/service"
	"github.com/byteowlz/{{project_name}}/schemas"
	"github.com/spf13/cobra"
)

type flags struct {
	path, theme, omarchy string
	radius               int
	json                 bool
}

func (f *flags) add(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&f.path, "config", "", "explicit TOML (overrides environment)")
	cmd.PersistentFlags().StringVar(&f.theme, "theme", "", "dark, light or omarchy")
	cmd.PersistentFlags().IntVar(&f.radius, "radius", 0, "radius dial in pixels, 0..32")
	cmd.PersistentFlags().StringVar(&f.omarchy, "omarchy-file", "", "read-only colors.toml data")
	cmd.PersistentFlags().BoolVar(&f.json, "json", false, "versioned JSON envelope, including failures")
}

func (f *flags) options(cmd *cobra.Command) config.Options {
	o := config.Options{Path: f.path}
	if cmd.Flags().Changed("theme") {
		o.Theme = &f.theme
	}
	if cmd.Flags().Changed("radius") {
		o.Radius = &f.radius
	}
	if cmd.Flags().Changed("omarchy-file") {
		o.OmarchyFile = &f.omarchy
	}
	return o
}

// DesktopOptions uses the same parser and precedence as the headless adapter.
func DesktopOptions(args []string) (config.Options, error) {
	var f flags
	cmd := &cobra.Command{Use: "{{project_name}}", SilenceErrors: true, SilenceUsage: true}
	f.add(cmd)
	if err := cmd.ParseFlags(args); err != nil {
		return config.Options{}, err
	}
	if len(cmd.Flags().Args()) != 0 {
		return config.Options{}, fmt.Errorf("desktop accepts flags only")
	}
	return f.options(cmd), nil
}

// Run returns 0 success, 2 usage/config/input, 1 denied/runtime/cancellation.
// In JSON mode all failures go to stdout as exactly one envelope.
func Run(args []string, in io.Reader, out, errOut io.Writer) int {
	var f flags
	var reply *service.Envelope
	root := &cobra.Command{
		Use: "{{project_name}}ctl", Version: "0.1.0",
		Short:        "Read-only synthetic desktop service; no hardware or network access",
		Example:      "{{project_name}}ctl snapshot --json\n{{project_name}}ctl validate --input snapshot.json --json\n{{project_name}}ctl apply --json # always denied",
		SilenceUsage: true, SilenceErrors: true,
	}
	root.SetArgs(args)
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	f.add(root)
	for _, operation := range []string{"snapshot", "validate", "apply", "appearance"} {
		cmd := &cobra.Command{Use: operation, Args: cobra.NoArgs, Short: "Shared service: " + operation,
			Example: "{{project_name}}ctl " + operation + " --json"}
		var inputPath string
		if operation == "validate" || operation == "apply" {
			cmd.Flags().StringVar(&inputPath, "input", "-", "snapshot file or - for stdin (64 KiB limit)")
		}
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(f.options(cmd), false)
			if err != nil {
				e := service.Fail("INVALID_CONFIG", err.Error())
				reply = &e
				return nil
			}
			input := ""
			if operation == "validate" {
				input, err = readInput(inputPath, in)
				if err != nil {
					e := service.Fail("INVALID_INPUT", err.Error())
					reply = &e
					return nil
				}
			}
			// apply never reads stdin or files: denied before any potential I/O.
			e := service.New(cfg).Inspect(context.Background(), service.Request{Operation: operation, Input: input})
			reply = &e
			return nil
		}
		root.AddCommand(cmd)
	}
	schema := &cobra.Command{Use: "schema", Args: cobra.NoArgs, Short: "Print JSON Schema for versioned envelopes",
		RunE: func(_ *cobra.Command, _ []string) error { _, err := out.Write(schemas.Envelope); return err }}
	root.AddCommand(schema)
	configCmd := &cobra.Command{Use: "config", Short: "Inspect effective configuration or explicitly initialize defaults"}
	configCmd.AddCommand(&cobra.Command{Use: "schema", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error { _, err := out.Write(schemas.Config); return err }})
	for _, action := range []string{"show", "init"} {
		configCmd.AddCommand(&cobra.Command{Use: action, Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			if action == "init" {
				path := f.path
				if path == "" {
					var err error
					path, err = config.GlobalPath()
					if err != nil {
						return err
					}
				}
				if err := config.Init(path); err != nil {
					return err
				}
			}
			cfg, err := config.Load(f.options(cmd), false)
			if err != nil {
				e := service.Fail("INVALID_CONFIG", err.Error())
				reply = &e
				return nil
			}
			e := service.Envelope{Version: "1", OK: true, Result: &service.Result{Config: &cfg}}
			reply = &e
			return nil
		}})
	}
	root.AddCommand(configCmd)
	if err := root.Execute(); err != nil {
		e := service.Fail("USAGE", err.Error())
		reply = &e
	}
	if reply == nil {
		return 0
	} // help/version/schema intentionally have their own format.
	// Detect JSON even when Cobra stopped at an invalid argument before parsing it.
	for _, arg := range args {
		if arg == "--json" || arg == "--json=true" {
			f.json = true
		}
	}
	if f.json {
		if err := json.NewEncoder(out).Encode(reply); err != nil {
			return 1
		}
	} else if reply.OK {
		if err := json.NewEncoder(out).Encode(reply.Result); err != nil {
			return 1
		}
	} else {
		fmt.Fprintf(errOut, "%s: %s\n", reply.Error.Code, reply.Error.Message)
	}
	if reply.OK {
		return 0
	}
	switch reply.Error.Code {
	case "USAGE", "INVALID_INPUT", "INVALID_CONFIG", "STALE_REVISION", "INVALID_OPERATION":
		return 2
	default:
		return 1
	}
}

func readInput(path string, in io.Reader) (string, error) {
	reader := in
	if path != "-" {
		f, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("cannot open snapshot file")
		}
		defer f.Close()
		reader = f
	} else if f, ok := in.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil || stat.Mode()&os.ModeCharDevice != 0 {
			return "", fmt.Errorf("provide --input FILE or pipe snapshot JSON")
		}
	}
	data, err := io.ReadAll(io.LimitReader(reader, 65537))
	if err != nil {
		return "", fmt.Errorf("cannot read input")
	}
	if len(data) > 65536 {
		return "", fmt.Errorf("input exceeds 64 KiB")
	}
	return strings.TrimSpace(string(data)), nil
}
