# Native/shared example contract

internal/service.Service.Inspect(context, Request) is authoritative.
The native App.View calls it directly; internal/cli parses flags/files and
calls the same Go service. ui.NewTester clicks/keys compare complete replies
against actual CLI adapter outputs, not a mocked presentation service.

| Operation | Input | Result / domain effects |
|---|---|---|
| snapshot | none | immutable synthetic revision + items; no writes |
| validate | snapshot JSON | exact roundtrip validation; no writes |
| apply | ignored | AUTHORITY_DENIED; no input read/write adapter |
| appearance | none | configuration plus optional parsed raw TOML colors |

Snapshot is exactly two ordered fixture items at revision fixture-v1.
Validation rejects unknown fields, trailing JSON, mismatched items, stale
revisions and >64 KiB input. Immutable service configuration is startup-only.
Reads are concurrent/idempotent; cancellation is checked at operation entry.
No write, retry, transport/disconnect or partial-result capability is claimed.
Local reads are size-bounded; stdin producers must terminate.

## Stable replies

schemas/envelope.schema.json (draft-07):
success `{version:"1",ok:true,result:{snapshot|appearance|config:...}}`;
failure `{version:"1",ok:false,error:{code,message}}`.
Exactly one result or error. CLI --json emits one envelope on stdout even for
invalid flags/input. No banners; help/version are text, schema commands raw JSON.
Without --json, results are JSON data; failures are stderr text.

Codes: USAGE, INVALID_CONFIG, INVALID_INPUT, INVALID_OPERATION, STALE_REVISION,
AUTHORITY_DENIED, CANCELLED, INVALID_THEME, THEME_UNAVAILABLE.
Exit 0 success; 2 usage/config/input/stale; 1 denied/runtime.
config init explicitly creates defaults (0600, never overwrite). Native first
run can create defaults; headless audits never do. Native window initialization
is not called by CLI/tests. Appearance returns raw data, not resolved tokens;
the canonical native adapter validates data before applying a theme.

## Consumer replacement

Use the consumer's authoritative service instead of this fixture. Mutation
requires explicit intent/authority/target/baseline, preview, backup/recovery,
serialized transport and truthful partial results/readback. Preserve results
when apply returns both progress and a structured error; never blindly retry
uncertain writes. Readiness is not authorization.

Native untrusted labels remain plain text. Any future trusted webview would
need DOM/textContent safety, but this template has none.
