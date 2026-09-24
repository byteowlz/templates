# Experiment review — visual contract

Origin: the pxltr experiment review UI. This template copies the **layout and
visual grammar**, not its experiment data, licensed media, or conclusions.

- **Mode:** Experience meets Read. Motion leads; evidence and caveats follow.
- **Palette:** dark background, elevated comparison panels, sunken viewing
  stages. Byteowlz semantic roles (`background`, `surface`, `surface-sunken`,
  `foreground`, `muted-foreground`, `primary`, `secondary`, `accent`, `warning`,
  `border`, `ring`) in `index.html`; no component-specific color tokens.
- **Shape:** one radius dial at 0. Square surfaces and controls.
- **Hierarchy:** large direct question; an explicit rights notice; matched GIF
  pair; source/native comparison; compact method/metrics/rights/artifacts grid;
  one plain-language conclusion; quieter evidence below.
- **Media:** pixel outputs at native 1× or integer nearest-neighbor ×4. Smooth
  source stills are not pixelated. Media paths are curated in `server.mjs`.
- **Responsive:** comparisons remain side-by-side on a phone; descriptive
  metadata stacks. Minimum 44px touch toggle with visible keyboard focus.
- **Boundary:** this is not a corpus explorer or archive. No client or licensed
  material is included. Template examples are synthetic and named as such.
