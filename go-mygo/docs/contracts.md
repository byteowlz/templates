# Shared example contract

`internal/service.Service.Execute(context, Request)` is authoritative.
`internal/desktop.Desktop.Execute` delegates unchanged; only this facade is
bound with MyGo. `internal/cli` parses flags/files and invokes the same service.

| Operation | Input | Result / side effects |
|---|---|---|
| snapshot | none | immutable synthetic revision + items; no domain writes |
| validate | snapshot JSON | exact roundtrip validation; no writes |
| apply | ignored | AUTHORITY_DENIED before reading an input file/stdin; no adapter |
| appearance | none | config and optional read-only TOML colors (raw data, not resolved tokens) |

Snapshot: `{revision:"fixture-v1",items:[{id,label},...]}`, exactly two ordered
synthetic items. Validation denies unknown fields, trailing JSON, mismatched
items, stale revisions and input >64 KiB. All service state is immutable except
host-only startup configuration. Reads are idempotent and thread-safe; no
conflict resolution, retries, writes or partial-success semantics are claimed.
Cancellation is checked at operation entry. File reads are local and bounded
in size; a CLI stdin producer must terminate its input.

## Replies and errors

`schemas/envelope.schema.json` is JSON Schema draft-07.
Success: `{version:"1",ok:true,result:{snapshot|appearance|config:...}}`.
Failure: `{version:"1",ok:false,error:{code,message}}`.
Exactly one result/error. Desktop resolves service failures as these same
envelopes instead of reducing structured failures to MyGo error strings.
Transport/runtime errors can still reject; callers must distinguish those.

Stable codes: USAGE, INVALID_CONFIG, INVALID_INPUT, INVALID_OPERATION,
STALE_REVISION, AUTHORITY_DENIED, CANCELLED, INVALID_THEME, THEME_UNAVAILABLE.
CLI `--json`: one envelope on stdout even on errors, no banners/logs.
Help/version are human text; schema commands are raw schema JSON.
Without `--json`, results are JSON data and errors are brief stderr text.

Exit codes: 0 success, 2 usage/config/input/stale validation, 1 denied/runtime.
No network, hidden model calls, hardware enumeration or automatic mutation.
`config init` is the sole explicit CLI file-creation action; it creates a
0600 default without overwrite. Desktop first-run default creation and MyGo's
own webview profile are startup side effects, not headless audit behavior.

## Consumer replacement

Replace the fixture with the consumer's authoritative service, not another
agent-only implementation. Mutation needs explicit intent, authority,
target/baseline validation, preview, backup/recovery, serialized writes and
truthful partial failures/readback. Readiness is never a grant.
Return partial results with structured errors when needed; do not discard
progress when an underlying apply returns both a result and error.

Do not pass imported names/labels/JSON through innerHTML: trusted-page XSS can
invoke the backend. Use textContent/DOM APIs and bounded validation.
