# Native scaffold evidence

Mode selected by user: fully native MyGo v0.3.2 ui.View. No HTML/webview,
typed IPC/Bun frontend or generic 3D dependency. Unselected webview checkpoint:
c5f71ab on feat/mygo-webview-prototype. Native feature starts at templates
origin/main 570caaa, separate from that prototype.

Canonical adapter: design-system feat/mygo-native-adapter, pin
c31fdc8ea0f685db15dcd5a9103f1d981db6c3de. Its just check passes real compile,
race tests and exact Studio sRGB/Base24/closed-role/R1/Omarchy provenance parity.
Finite radius-tier overflow is regression-tested without an arbitrary dial cap.

## Reproduction

Run scripts/scaffold-check.sh from the raw template directory. It creates a
real project-name-replaced temporary scaffold, preserves executable modes,
uses isolated XDG config and runs just install/check/build/agent-check.
No global machine settings are modified. Go 1.27.1 resolves automatically.
Go-only schema/fixture/headless native tests; no frontend dependency is needed.

Observed on macOS arm64:

- Go compilation and go vet pass, including native host and semantic CLI.
- just check and just agent-check pass, including actual CLI subprocess schema
  success/error validation, no-write sandbox and snapshot stdin/file roundtrip.
- Native ui.NewTester clicks/read/validate/denied/recovery, keyboard focus/Enter
  and announcements match CLI replies; dark/light/narrow renders exist.
- just build creates the real native macOS app bundle (about 10.6 MB in the
  initial measured build); no frontend build or hidden webview.
- just package on the final scaffold produces the local macOS app and a
  3.2 MB DMG. This is packaging proof, not notarized distribution approval.
- just native-proof opens native ui.View, asserts Page()==nil and captures
  dark/light content; frame stats independently report drawn on the GPU.
- Packaged native app was separately launched with isolated XDG config.
  Actual macOS AX exposes named buttons, heading and fixture/status text.
  Denied apply and recovery were exercised through native AX actions.
  Window-only OS captures contain only the synthetic fixture, not user data.
- macOS tiling resized the native window to 1009×1131 points; capture metadata
  must retain actual size, not pretend it stayed at requested 900×640.

CapturePage for native content is a CPU re-render: its PNGs are **not**
presented as GPU screenshots. Actual OS window captures + frame-stats logs
are separate evidence. Headless PNGs likewise are software-test output.
A source/compiled fixture is not product composition or hardware proof.

One independently instantiated gate run encountered its 5s help-process
deadline under concurrent load. Direct help subsequently completed in 0.316s
and the same executable gate passed unchanged. The bound was not relaxed.
Reproduction commands fail closed on timeouts; file/report presence is not pass.

## Change review

Root ripwire quality-delta against 570caaa: gating=0. New-code complexity/
verbosity flags remain in the CLI command dispatcher, config layering, fixture
proof driver and scaffolder; these are sequential dispatch/assertion tables.
Test entrypoints are reported dead by the name-based lens, despite explicit
Go test execution. No existing template source changed. The example method is
named Inspect to describe its read-only scope and avoid an unrelated template's
Execute symbol confusing graph resolution. Contract check reports a new symbol;
its native/CLI callers have executable fixture parity coverage. Full owned diff
and canonical snapshot digests were inspected; whitespace checks pass.

## Limits

Linux/Windows runtime, native assistive-technology behavior beyond inspected
macOS AX, high-contrast/wide-P3 behavior, signing/notarization/distribution and
actual HID are unverified. No live device was enumerated or written. No claims
of mutation, transport recovery or partial writes are inherited from this
immutable fixture. Headless audit does not create config/profile state;
desktop first-run defaults are an explicitly documented startup side effect.

The generic demo remains illustrative/unapproved. Vial's Studio B composition
selection does not turn it into a production keyboard UI. Consumers must use
their own core service/config/state, authority and domain readiness proof.
