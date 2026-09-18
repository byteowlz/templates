# docs/omarchy.md — following the machine's omarchy theme (optional)

On [omarchy](https://github.com/basecamp/omarchy) machines, this app can wear
the currently active omarchy theme instead of the committed house scheme. The
mapping is data-only: omarchy's `colors.toml` → base24 slots. No omarchy files
are modified, nothing is executed from the theme, and without omarchy the app
is untouched.

## Use

```bash
just omarchy          # regenerate src/theme/omarchy.generated.css (exit 0 if no omarchy)
just omarchy watch    # keep following theme switches while you work
```

Then reload the dev server (vite hot-reloads the CSS import).

## How it works

- Reads the active theme, first hit wins:
  1. `$OMARCHY_THEME_FILE`
  2. `~/.local/state/omarchy/current/theme/colors.toml`
  3. `~/.config/omarchy/current/theme/colors.toml` (legacy symlink)
- Maps the canonical semantic palette deterministically onto base24
  (contract: byteowlz govnr ADR-0040, "external theme ecosystems are data
  sources; Omarchy is the first compatibility fixture"):

| slots | omarchy keys |
|---|---|
| `base00`–`base07` | `background`, `lighter_background`, `selection`, `muted`, `dark_foreground`, `foreground`, `light_foreground`, `bright_foreground` |
| `base08`–`base0F` | `red`, `orange`, `yellow`, `green`, `cyan`, `blue`, `magenta`, `brown` |
| `base10`–`base17` | `dark_background`, `darker_background`, `bright_red`, `bright_yellow`, `bright_green`, `bright_cyan`, `bright_blue`, `bright_magenta` |

- `mode` and `accent` are preserved as explicit metadata
  (`--omarchy-mode`, `--omarchy-accent`) — never guessed from slots.
- Legacy themes exposing only ANSI `color0..color15` get a best-effort partial
  mapping (unfilled slots keep the committed house values; the file header
  says so).
- Because roles and Tailwind utilities all reference slots, overriding the
  slots recolors the whole app. The `:root` block in the generated file wins
  over `index.css` (later import); the engine's `applyScheme()` (if vendored)
  still wins over it via inline styles.

## Rules

- `src/theme/omarchy.generated.css` is machine-local and generated: never hand
  edit, never commit personal theme values. Committed state = no-op
  passthrough so the import resolves everywhere.
- Omarchy assets (wallpapers, icons, fonts) are NOT imported — palette data
  only, per the ADR-0040 adapter contract.
- The app must still look correct without omarchy (it does: house scheme).
