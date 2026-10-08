# Agent guidance

- Read CONTEXT.md, docs/contracts.md, docs/adr/ and docs/agent-readiness.md.
- This template is MyGo v0.3.2 system-webview mode, NOT fully native GPU UI.
  Inspect pinned APIs before changing backend assumptions.
- Keep semantic logic in internal/service. Desktop is a narrow bound facade;
  never bind host lifecycle/configuration methods. CLI shares the service.
- MyGo methods run concurrently. Maintain context cancellation and race tests.
- Run just check and just agent-check before claiming the fixture contract.
  Readiness requires reproduced evidence, not file/report presence.
- Use the public skill byteowlz/skillissues/skills/agent-readiness
  (initial contract revision 6ebce57); no local skill path required at runtime.
- Use Bun; run Go compilation/tests after every Go change. Keep go.sum/bun.lock.
- Use just build/package for embedded frontend assets; raw go build is compile
  evidence only. Signing/notarization and foreign-platform runtime are separate.
- Colors/radius come ONLY from third_party/design-system's unchanged engine.
  Components consume roles, never slots/literals. No parallel palette mapping.
- third_party is intentional: a top-level vendor directory triggers Go vendoring.
- Omarchy is opt-in data only, bounded flat TOML; never execute/source theme code.
- Demo layout is illustrative/unapproved, not a Studio-selected composition.
  Production decisions need a versioned product-specific Studio trail.
- Treat imported labels, hardware names and files as untrusted. Use textContent/
  DOM APIs; no innerHTML interpolation in the trusted MyGo page (backend access).
- No default 3D dependency. Opt-in consumers retain an operable semantic surface.
- No real hardware, account or remote writes to prove readiness. No authority
  follows from IPC exposure, readiness, flags, identity discovery or schemas.
- Do not publish artifacts automatically or modify launchd plists.
