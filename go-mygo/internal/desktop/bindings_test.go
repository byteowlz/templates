package desktop

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/byteowlz/{{project_name}}/internal/config"
	"github.com/byteowlz/{{project_name}}/internal/service"
)

func TestFixtureParityAndSafety(t *testing.T) {
	svc := service.New(config.Default())
	desktop := Desktop{Service: svc}
	ctx := context.Background()
	snapshot := svc.Execute(ctx, service.Request{Operation: "snapshot"})
	data, err := json.Marshal(snapshot.Result.Snapshot)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		request service.Request
		code    string
	}{
		{service.Request{Operation: "snapshot"}, ""},
		{service.Request{Operation: "validate", Input: string(data)}, ""},
		{service.Request{Operation: "validate", Input: strings.ReplaceAll(string(data), "fixture-v1", "stale")}, "STALE_REVISION"},
		{service.Request{Operation: "validate", Input: `{"unknown":true}`}, "INVALID_INPUT"},
		{service.Request{Operation: "validate", Input: string(data) + " {}"}, "INVALID_INPUT"},
		{service.Request{Operation: "validate", Input: strings.Repeat("x", 65537)}, "INVALID_INPUT"},
		{service.Request{Operation: "apply", Input: string(data)}, "AUTHORITY_DENIED"},
		{service.Request{Operation: "unknown"}, "INVALID_OPERATION"},
	}
	for _, tc := range cases {
		gui, headless := desktop.Execute(ctx, tc.request), svc.Execute(ctx, tc.request)
		g, _ := json.Marshal(gui)
		h, _ := json.Marshal(headless)
		if string(g) != string(h) {
			t.Fatalf("parity: %s != %s", g, h)
		}
		if tc.code == "" && !gui.OK {
			t.Fatalf("expected success: %+v", gui)
		}
		if tc.code != "" && (gui.OK || gui.Error.Code != tc.code) {
			t.Fatalf("wanted %s: %+v", tc.code, gui)
		}
	}
	after, _ := json.Marshal(svc.Execute(ctx, service.Request{Operation: "snapshot"}))
	before, _ := json.Marshal(snapshot)
	if string(before) != string(after) {
		t.Fatal("denied mutation changed snapshot")
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if got := desktop.Execute(cancelled, service.Request{Operation: "snapshot"}); got.Error.Code != "CANCELLED" {
		t.Fatal(got)
	}
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			got := desktop.Execute(ctx, service.Request{Operation: "snapshot"})
			if !got.OK {
				t.Error(got)
			}
		})
	}
	wg.Wait()
}
