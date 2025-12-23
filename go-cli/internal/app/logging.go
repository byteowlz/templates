package app

import (
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

// LogSettings control how the logger behaves.
type LogSettings struct {
	Level       Level
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
		Diagnostics: flags.Diagnostics,
		Colorize:    colorize,
		Writers:     writers,
		FileHandle:  fileHandle,
	}, nil
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
