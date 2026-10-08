# Portable design-system mechanism

Pinned source: https://github.com/byteowlz/design-system at
`0968afc89e0919ce856a4d99351bc475b26a87ee`, implementation shadcn-ts v0.1.0.
`third_party/design-system/src` is byte-identical source, imported by relative
paths. No checkout path, npm publication or local service is needed to build.
The package's public engine is retained intact (including types its entrypoint
imports); only the four relevant spec documents are copied, not playground,
node_modules, compiled bundles or unrelated effect specs.

Pipeline: `normalizeScheme -> mapSchemeToTokens/applyScheme`; closed 16 roles,
Base16 widening, house Base24 schemes, identity and R1 radius come from that
engine. Components use role variables. No consumer color-to-role table exists.

`third_party/hashes.json` locks vendored bytes. `tests/theme.test.ts` verifies
them and compares all reference dark/light tokens, Base16 widening and radius
at 0/8px against a golden generated from the pinned upstream source.
Golden is an observation, not a second implementation.
Change the source pin, preserved metadata, hashes and reference golden together
on an explicit reviewed upgrade; never patch a mapping to fit a screenshot.
Copied source license metadata is preserved in package.upstream.json;
LICENSE documents MIT text/attribution. Source has no separate license file.

## Omarchy (opt-in, data only)

`theme = "omarchy"`, optionally `omarchy_file` or `--omarchy-file`, reads
flat TOML strings, at most 64 KiB. Default is
`~/.local/state/omarchy/current/theme/colors.toml`. A missing/invalid file is
an explicit failure; desktop retains the house default and reports it.

Slot input mapping is the existing templates webapp canonical/legacy mapping,
extracted into data-only `third_party/omarchy/mapping.json`, not newly invented.
Reference source is pinned to templates commit
`6c6214ec1acd3fb0a8bdf8ca805196d564c57593`.
The original script is archived solely for provenance/drift tests; never run it.
Canonical input must contain all 24 keys; legacy ANSI must fill all 16 keys.
Legacy's unspecified slots remain the selected house scheme: this is an
explicit **partial ANSI overlay**, not a claimed full Omarchy/Base24 conversion.
Mode defaults dark unless data declares light; accent metadata is not a new role.
Malformed colors cannot inject CSS. No shell scripts, theme installation,
generated personal CSS, watcher, source or executable content is used.

The Go appearance operation returns parsed raw data; TypeScript validates
palette completeness/colors before atomically applying through the engine.
Appearance data inspection is not a resolved-token parity claim.
