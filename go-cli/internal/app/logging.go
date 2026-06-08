package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Level represents logging severity in ascending order.
type Level int

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

// LogFormat selects the on-the-wire encoding of log records.
type LogFormat int

const (
	// FormatText renders human-readable, optionally colorized lines.
	FormatText LogFormat = iota
	// FormatJSON renders one JSON object per line (JSON Lines).
	FormatJSON
)

// LogSettings control how the logger behaves.
type LogSettings struct {
	Level       Level
	Format      LogFormat
	Diagnostics bool
	Colorize    bool
	Writers     []io.Writer
	FileHandle  *os.File
}

// Logger is a lightweight structured logger tailored for the template.
type Logger struct {
	settings LogSettings
	mu       sync.Mutex
}

// ConfigureLogger returns a logger configured with the supplied settings.
func ConfigureLogger(settings LogSettings) Logger {
	if len(settings.Writers) == 0 {
		settings.Writers = []io.Writer{io.Discard}
	}
	return Logger{
		settings: settings,
	}
}

// Close releases associated resources (currently only the optional file handle).
func (l Logger) Close() error {
	if l.settings.FileHandle != nil {
		return l.settings.FileHandle.Close()
	}
	return nil
}

func (l Logger) Trace(msg string, args ...any) {
	l.log(LevelTrace, msg, args...)
}

func (l Logger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args...)
}

func (l Logger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	l.log(LevelWarn, msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args...)
}

func (l Logger) log(level Level, msg string, args ...any) {
	if level > l.settings.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	formatted := formatMessage(level, l.settings, msg, args...)
	for _, w := range l.settings.Writers {
		io.WriteString(w, formatted)
		io.WriteString(w, "\n")
	}
}

func formatMessage(level Level, settings LogSettings, msg string, args ...any) string {
	body := fmt.Sprintf(msg, args...)

	if settings.Format == FormatJSON {
		return formatJSON(level, body)
	}

	levelLabel, color := levelAttributes(level)

	if settings.Colorize && color != "" {
		levelLabel = color + levelLabel + resetColor()
	}

	if settings.Diagnostics {
		timestamp := time.Now().Format(time.RFC3339)
		return fmt.Sprintf("[%s] %-5s %s", timestamp, levelLabel, body)
	}

	return fmt.Sprintf("%-5s %s", levelLabel, body)
}

// jsonRecord is the unified cross-language log schema. Field order is fixed so
// output is byte-stable: time (RFC3339), level (lowercase), msg.
type jsonRecord struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func formatJSON(level Level, body string) string {
	encoded, err := json.Marshal(jsonRecord{
		Time:  time.Now().Format(time.RFC3339),
		Level: levelName(level),
		Msg:   body,
	})
	if err != nil {
		// json.Marshal of plain strings cannot fail in practice; fall back to a
		// minimal hand-built record so a log line is never silently dropped.
		return fmt.Sprintf(`{"time":%q,"level":%q,"msg":"<unencodable log message>"}`,
			time.Now().Format(time.RFC3339), levelName(level))
	}
	return string(encoded)
}

func levelName(level Level) string {
	switch level {
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	case LevelTrace:
		return "trace"
	default:
		return "info"
	}
}

func levelAttributes(level Level) (string, string) {
	switch level {
	case LevelError:
		return "ERROR", "\033[31m"
	case LevelWarn:
		return "WARN", "\033[33m"
	case LevelInfo:
		return "INFO", "\033[36m"
	case LevelDebug:
		return "DEBUG", "\033[35m"
	case LevelTrace:
		return "TRACE", "\033[34m"
	default:
		return "INFO", ""
	}
}

func resetColor() string {
	return "\033[0m"
}

func ResolveLogSettings(flags CommonFlags, cfg AppConfig) (LogSettings, error) {
	level := parseLevel(cfg.Logging.Level)

	switch {
	case flags.Trace:
		level = LevelTrace
	case flags.Debug:
		level = LevelDebug
	case flags.Verbose >= 2:
		level = LevelTrace
	case flags.Verbose == 1:
		level = LevelDebug
	case flags.Quiet && level > LevelError:
		level = LevelError
	}

	colorPolicy := flags.Color
	if flags.NoColor {
		colorPolicy = "never"
	}

	colorize := shouldColorize(colorPolicy)

	format := resolveLogFormat(flags.LogFormat, cfg.Logging.Format)
	if format == FormatJSON {
		colorize = false
	}

	writers := []io.Writer{os.Stderr}

	if flags.Quiet {
		writers = []io.Writer{io.Discard}
	}

	var fileHandle *os.File
	if cfg.Logging.File != "" && !flags.DryRun {
		handle, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return LogSettings{}, fmt.Errorf("open log file %s: %w", cfg.Logging.File, err)
		}
		fileHandle = handle
		if !flags.Quiet {
			writers = append(writers, handle)
		} else {
			writers = []io.Writer{handle}
		}
	}

	return LogSettings{
		Level:       level,
		Format:      format,
		Diagnostics: flags.Diagnostics,
		Colorize:    colorize,
		Writers:     writers,
		FileHandle:  fileHandle,
	}, nil
}

// resolveLogFormat decides the encoding from the --log-format flag and the
// configured default. Precedence: an explicit text/json flag wins, then the
// config value, then auto-detection (json when stderr is not a terminal).
func resolveLogFormat(flagValue, configValue string) LogFormat {
	choice := configValue
	if choice == "" {
		choice = "auto"
	}
	if flagValue == "text" || flagValue == "json" {
		choice = flagValue
	}

	switch choice {
	case "text":
		return FormatText
	case "json":
		return FormatJSON
	default: // "auto"
		if isTerminal(os.Stderr) {
			return FormatText
		}
		return FormatJSON
	}
}

func parseLevel(value string) Level {
	switch strings.ToLower(value) {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

func shouldColorize(policy string) bool {
	switch policy {
	case "always":
		return true
	case "never":
		return false
	default:
		return isTerminal(os.Stderr)
	}
}
