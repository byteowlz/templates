package app

import "fmt"

// CommonFlags capture global CLI options shared by all commands.
type CommonFlags struct {
	ConfigPath     string
	Quiet          bool
	Verbose        int
	Debug          bool
	Trace          bool
	JSON           bool
	YAML           bool
	NoColor        bool
	Color          string
	DryRun         bool
	AssumeYes      bool
	TimeoutSeconds *int
	Parallelism    *int
	NoProgress     bool
	Diagnostics    bool
}

// ValidateColor ensures the color flag uses a supported value.
func (c *CommonFlags) ValidateColor() error {
	switch c.Color {
	case "", "auto", "always", "never":
		return nil
	default:
		return fmt.Errorf("invalid --color value %q (expected auto, always, or never)", c.Color)
	}
}
