# ADR 0001: Fully native MyGo substrate

Status: selected mode by user; starter composition unapproved.
MyGo v0.3.2 ui.View is the native GPU presentation. No HTML/webview, frontend
build or IPC remains. Native Go views and the semantic CLI share one service.

The prior tested system-webview work was checkpointed as unselected prototype
c5f71ab on feat/mygo-webview-prototype, not merged into this feature.
Use the actual v0.3.2 docs/ui/source; bind/generate/frontend APIs are irrelevant.

The canonical adapter lives in design-system impls/mygo-go. A pinned unchanged
source snapshot provides portable builds while the owning branch is unmerged.
Relative replacement plus complete hashes/native parity tests avoids a local
checkout dependency or parallel template styling implementation.
Native sRGB resolution intentionally matches pinned Studio outputs; native
toolkit vocabulary/metrics and wide-gamut limitations are explicit.

The fixture demonstrates direct Go parity, bounded stable envelopes and safe
denial, not real hardware support. 3D belongs to an opt-in consumer seam with
an operable semantic alternative. Composition approval and platform/hardware/
accessibility acceptance remain separate from compilation and token mechanism.
