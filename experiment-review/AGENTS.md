# Agent guidance

- This is a curated, read-only review exhibit, not a general file server.
- Read `DESIGN.md` before changing layout or tokens. Use byteowlz roles and
  radius 0; keep integer nearest-neighbor rendering for pixel outputs.
- `server.mjs` serves only exact named paths. Update its allowlist and tests
  together; reject traversal, symlinks outside the media root, corpora and
  private files. Bind localhost by default; tailnet access is explicit.
- Replace sample media and copy before reviewing real experiments. State source,
  rights/training status, method, metrics, status, artifacts and caveats.
- Do not copy licensed experiments into a template or upload review assets to
  third parties. Samples in `demo/` are locally authored and synthetic.
- Run `just test` after editing. Verify desktop and phone in bounded passes.
