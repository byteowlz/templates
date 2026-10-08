package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/byteowlz/{{project_name}}/internal/config"
)

// Request is shared by typed desktop IPC and semantic CLI commands.
type Request struct {
	Operation string `json:"operation"`
	Input     string `json:"input"`
}
type Item struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}
type Snapshot struct {
	Revision string `json:"revision"`
	Items    []Item `json:"items"`
}
type Appearance struct {
	Theme  string            `json:"theme"`
	Radius int               `json:"radius"`
	Colors map[string]string `json:"colors"`
}
type Result struct {
	Snapshot   *Snapshot      `json:"snapshot,omitempty"`
	Appearance *Appearance    `json:"appearance,omitempty"`
	Config     *config.Config `json:"config,omitempty"`
}
type Failure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Envelope struct {
	Version string   `json:"version"`
	OK      bool     `json:"ok"`
	Result  *Result  `json:"result,omitempty"`
	Error   *Failure `json:"error,omitempty"`
}

// Service owns the example substrate. It has no hardware/network/write adapter.
type Service struct {
	mu  sync.RWMutex
	cfg config.Config
}

func New(cfg config.Config) *Service { return &Service{cfg: cfg} }

// Configure is host startup-only; the bound Desktop facade does not expose it.
func (s *Service) Configure(cfg config.Config) { s.mu.Lock(); s.cfg = cfg; s.mu.Unlock() }

func Fail(code, message string) Envelope {
	return Envelope{Version: "1", OK: false, Error: &Failure{Code: code, Message: message}}
}
func success(result Result) Envelope {
	return Envelope{Version: "1", OK: true, Result: &result}
}

// Execute dispatches the complete read-only example contract.
func (s *Service) Execute(ctx context.Context, request Request) Envelope {
	if ctx.Err() != nil {
		return Fail("CANCELLED", "operation cancelled")
	}
	if len(request.Input) > 65536 {
		return Fail("INVALID_INPUT", "input exceeds 64 KiB")
	}
	switch request.Operation {
	case "snapshot":
		if request.Input != "" {
			return Fail("INVALID_INPUT", "snapshot accepts no input")
		}
		return success(Result{Snapshot: fixture()})
	case "validate":
		return validate(request.Input)
	case "apply":
		return Fail("AUTHORITY_DENIED", "starter is read-only; no write adapter exists")
	case "appearance":
		return s.appearance()
	default:
		return Fail("INVALID_OPERATION", "use snapshot, validate, apply or appearance")
	}
}

func fixture() *Snapshot {
	return &Snapshot{Revision: "fixture-v1", Items: []Item{{ID: "sample-a", Label: "Synthetic sample A"}, {ID: "sample-b", Label: "Synthetic sample B"}}}
}

func validate(input string) Envelope {
	var snapshot Snapshot
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&snapshot); err != nil {
		return Fail("INVALID_INPUT", "expected a snapshot JSON object")
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return Fail("INVALID_INPUT", "expected exactly one JSON object")
	}
	if snapshot.Revision != "fixture-v1" {
		return Fail("STALE_REVISION", "expected baseline fixture-v1")
	}
	expected := fixture()
	if len(snapshot.Items) != len(expected.Items) {
		return Fail("INVALID_INPUT", "fixture items do not match")
	}
	for i, item := range snapshot.Items {
		if item != expected.Items[i] {
			return Fail("INVALID_INPUT", "fixture items do not match")
		}
	}
	return success(Result{Snapshot: &snapshot})
}
