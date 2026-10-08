# Required scope-aware agent gate

Public rubric: [byteowlz/skillissues/skills/agent-readiness](https://github.com/byteowlz/skillissues/tree/main/skills/agent-readiness)
(initial contract 6ebce57). Not a runtime dependency or guessed score.

just agent-check executes native headless ui.NewTester input/reply/keyboard/
announcement parity against the shared CLI adapter, race/concurrency/cancel/
stale/denied tests, canonical adapter parity/digests, and an actual compiled CLI
in a private temporary cwd/XDG sandbox. It schema-validates real successful
and failed outputs, tests file/stdin snapshot roundtrip and verifies no hidden
writes. Read/validate/apply recovery is tested through native controls.

Artifacts: artifacts/agent-check/{go-tests.log,observations.json,report.json}
and artifacts/headless/*.png. Zero exit means fixture assertions passed.
Human accessibility stays unverified for actual OS assistive-technology and
foreign-platform coverage. Immutable fixture has no domain edits, transport,
partial writes or authority grants; these do not become inherited features.

The report is a conservative scoped claim, not certification. Structural
report validation does not run procedures or verify artifacts. Reproduce them.
Native OS screenshots/AX and GPU logs are separate evidence from headless
software renders and CapturePage CPU content re-renders.

Consumers replace fixture proof with domain-specific gates (discoverability,
structured/shared state, safety, concurrency/recovery, portable state,
reproducibility and human accessibility). Record pass/fail/unverified/
not_applicable per requirement with exact revision/platform/procedure/artifact.
Never touch a real device/account/shared state solely to prove readiness.
