# ADR 0001: MyGo system-webview shared substrate

Status: chosen for the system-webview template request; publication subject to
mode confirmation. MyGo fully native GPU UI is a different adapter/template,
not a switch that quietly embeds a webview.

Use MyGo v0.3.2, native Go service, generated TypeScript IPC and Bun/Vite.
Vanilla TypeScript suffices for two demo operations; React would add framework
weight without a state/composition need. MyGo's build tool embeds built assets
and packages each host platform. Its Go 1.27.1 requirement is declared in the
module and resolved by Go's automatic toolchain without machine-wide changes.

CLI and desktop are thin adapters. Configuration belongs to the Go host and
must not be exposed as remotely callable lifecycle methods. Generation has no
startup config writes: binding registration precedes App.Run, configuration
initialization lives in WhenReady.

Use pinned source for the real design-system engine, not a copied role table.
Source pin and hash/golden tests make drift visible. third_party avoids Go's
special top-level vendor semantics. Native GPU mode would require a matching
native design-system adapter with its own parity/render evidence.

No HTTP/MCP/model layer, HID writes, Three.js or automatic updater is added.
The fixture demo is explicitly unapproved; production composition belongs in
a separately selected Studio record. Platform signing and native runtime proof
remain downstream gates, not consequences of successful compilation.
