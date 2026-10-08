# Required agent-readiness gate

Contract: [byteowlz/skillissues/skills/agent-readiness](https://github.com/byteowlz/skillissues/tree/main/skills/agent-readiness),
initial revision `6ebce57`. The skill is a rubric, not a runtime dependency.

`just agent-check` executes race-enabled Go CLI/desktop facade parity and
cancellation/concurrency tests, builds the real CLI, then runs file/stdin
roundtrips, schema validation, usage/stale/invalid/denied paths and no-write
sandbox checks. It also executes the real generated TypeScript call dispatch
against a **test transport harness** and verifies engine/mapping provenance.
This harness is not a native IPC/webview test.

Artifacts: `artifacts/agent-check/observations.json` and `report.json`.
A zero process exit means fixture assertions passed; it does **not** mean all
readiness gates, native UI or a consumer passed. The report retains
`shared_substrate` and `human_accessibility` as unverified for actual desktop
runtime. Report-schema validation checks structure only and never verifies
the observations.

| Gate | Fixture evidence / limit |
|---|---|
| discoverability | real help and documented commands; 5s process bounds |
| structured_contract | actual executable outputs schema-validated, including failures |
| shared_substrate | Go facade parity + generated dispatch; open native UI still unverified |
| authority_and_safety | no mutation adapter; denied apply leaves state/files unchanged |
| concurrency_and_recovery | race/concurrent/cancelled reads, idempotence; no transport or partial writes in scope |
| portable_state | file/stdin snapshot roundtrip in isolated cwd/XDG tree |
| reproducibility | executed fixtures and hashes; no machine-local build dependencies |
| human_accessibility | semantic source controls present; native keyboard/AX/focus verification still required |

A consumer cannot inherit an agent-ready label. Replace the fixture with the
domain's state/actions, reproduce required gates and record pass/fail/unverified/
not_applicable per scope, exact revision/platform, procedure and observed artifact.
Read-only audits must not mutate real hardware/accounts/shared state. Native
render, IPC, platform differences and hardware have separate evidence.
Agent access is not authority; mutation approval cannot be derived from readiness.
