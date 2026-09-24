# {{project_name}} — experiment review

A private, mobile-first review-gallery template adapted from pxltr's experiment
review surface. **All supplied visuals are synthetic, locally authored samples**;
there are no H3 clips, corpora, client assets, or third-party uploads in this
template. It is a review surface, not a file browser or an archive.

## Start

```bash
just test
just serve                         # http://127.0.0.1:8792/
REVIEW_TAILNET_IP=<this-host-ip> just serve  # also bind a tailnet interface
```

Set `PORT` if 8792 is in use. Node 18+ is the only runtime dependency.
ImageMagick is needed **only** for `just demo` to regenerate the authored
sample GIFs and SVGs from `scripts/make-demo.mjs`.

## Adapt

1. Replace the synthetic content in `index.html`: observation first, then
   source, method, rights/training status, metrics with denominators, artifacts,
   conclusion **and caveat**. Keep a side-by-side, equal-timing moving comparison.
2. Put only reviewed, redistributable presentation assets in `demo/`, or
   explicitly change the root in `server.mjs`. Add each exact file path to its
   `media` allowlist; never allow arbitrary requested paths or expose a corpus.
3. Update the path-denial/availability tests. Keep images self-hosted and
   `image-rendering:pixelated` with integer nearest-neighbor scales for pixel art.
   Use `image-rendering:auto` for smooth source stills.
4. Update `DESIGN.md` when intentionally changing this visual identity.

Layout: restrained dark role surfaces, sharp radius 0, compact comparison
cards, pale green verdict accent, ochre rights notice, quiet evidence sheets,
44px touch controls, and an explicit native/nearest-neighbor toggle. There is
no dependency on Immich or an external image service. A private archive can
remain the archive; this UI presents a curated argument instead.

See `../templates.json` in the templates repository for catalogue metadata.
If your `byt` version does not recognize the named template yet, use
`byt new myreview --from-git byteowlz/templates --subdir experiment-review`.
