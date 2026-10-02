# Scheme provenance

The two JSON files here are **example look data** for the mechanism demo, not a
house look. They are a starting palette — roughly the green-tinted oqto look —
chosen so the `{{project_name}}-design` crate has something to resolve and the
sample app renders a coherent dark/light pair out of the box.

This template's intended shape:

- The **mechanism** (`{{project_name}}-design`) is the reusable, target-agnostic
  piece: parse a scheme, resolve closed roles, emit a gpui-component theme set.
- Each **tool wears its own colors and identity**. When you build a real app,
  replace `example-*.json` with your scheme files (loaded from
  `design-system/spec/schema.json` or hand-authored) and set the `Identity` that
  expresses your tool (fonts, radius dial, shadows).

The concrete palette values here were adapted from the open byteowlz design
work for the web reference; they carry no third-party identifiers. Treat them as
disposable example values.