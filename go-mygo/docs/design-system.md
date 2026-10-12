# Canonical native mechanism

Canonical home: byteowlz/design-system impls/mygo-go.
Pin: `c31fdc8ea0f685db15dcd5a9103f1d981db6c3de`, branch feat/mygo-native-adapter.
third_party/mygo-go is a byte-identical portable source snapshot with its MIT
license, reference pins, fixtures and tests. The relative Go replacement is
intentional; there is no absolute local-path dependency or hidden web engine.

third_party/mygo-go.source.json records every adapter byte and source revision.
TestAdapterSnapshotDrift verifies the pin/hashes; the adapter's tests assert
all copied Studio sRGB slots/16 closed roles, native vocabulary fields, Base16
widening/alias and R1 at 0/4/8/32. See its README/PROVENANCE/VALIDATION for exact
inputs and limitations. No app-owned palette table or native default-color
patch exists. Scheme/role data are embedded by Go.

Native field vocabulary is nonportable but role-derived: shared widget surface
and action/status colors, selection/focus/inverse/scrollbar. Controls use SM;
container uses LG. Exact authored colors are resolved using Studio's f32/clipped
sRGB policy; wide-P3 remapping is deliberately not claimed.

Omarchy: theme=omarchy optionally selects omarchy_file/--omarchy-file.
Empty path selects ~/.local/state/omarchy/current/theme/colors.toml.
Only bounded regular flat string TOML is read. Adapter uses the existing
canonical/legacy mapping with provenance, never theme scripts or install hooks.
Legacy ANSI is a documented partial house overlay, not invented full Base24
colors. Missing/malformed selected data fails startup; no silent fallback.

Upgrade the canonical adapter commit, complete snapshot manifest and golden/
drift tests together; never patch only template styling to match a screenshot.
