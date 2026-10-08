# Native MyGo agent guidance

- Read CONTEXT.md, docs/contracts.md, docs/adr/ and docs/agent-readiness.md.
- This is MyGo v0.3.2 native GPU ui.View. No HTML/webview, Bun frontend or IPC.
  Inspect pinned docs/ui and source before changing API assumptions.
- Native view and semantic CLI share internal/service, not parallel logic.
  Do not preserve illustrative fixture behavior in a production consumer.
- ui.Element belongs to one build pass. Keep model data/handles, not Elements,
  in app state; background work updates native model through Window.Update.
- Run just check and just agent-check before claiming fixture proof. The public
  byteowlz/skillissues/skills/agent-readiness skill is a rubric, not dependency.
  File presence/report validation/build success are not evidence verification.
- Keep Go module/checksum pins and portable relative adapter replacement.
  third_party/mygo-go is unchanged canonical source; never create an alternate
  color mapping in the app. Upgrade source pin/hash/tests together.
- A top-level vendor directory changes Go's module behavior; use third_party.
- Components read closed roles and R1 tiers. Native vocabulary exceptions,
  sRGB/wide-gamut/high-contrast limits are explicit in adapter docs.
- Omarchy is opt-in bounded data only. Never source/execute/install a theme.
- The starter is unapproved composition. Consumers need a versioned
  product-specific Studio record and actual native evidence.
- No default 3D/HID/model/network dependency. Agent access is not authority.
  Never touch real hardware/accounts/shared state for fixture proof.
- Labels/imports/hardware names are untrusted. Native text stays data.
- No release upload or global toolchain changes. Never touch launchd plists.
