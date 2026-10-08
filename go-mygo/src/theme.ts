import { applyScheme, oqtoDark, oqtoLight, type Scheme, type Slots } from "../third_party/design-system/src/index";
import mapping from "../third_party/omarchy/mapping.json";
import type { Appearance } from "./mygo";

// Data adapter only: all role/radius mapping stays in the shared engine.
export function omarchyScheme(colors: Record<string, string>): Scheme {
 const mode = colors.mode ?? colors.theme_type ?? "dark";
 if (mode !== "dark" && mode !== "light") throw new Error("Omarchy mode must be dark or light");
 const base = mode === "light" ? oqtoLight : oqtoDark;
 const canonical = "lighter_background" in colors || "darker_background" in colors;
 const selected = canonical ? mapping.canonical : mapping.legacy;
 const slots: Slots = { ...base.slots };
 let count = 0;
 for (const [slot, key] of Object.entries(selected)) {
  const value = colors[key];
  if (value === undefined) continue;
  if (!/^#[0-9a-fA-F]{6}$/.test(value)) throw new Error("Omarchy palette requires six-digit hex colors");
  slots[slot as keyof Slots] = value;
  count++;
 }
 if (count === 0) throw new Error("Omarchy data has no supported palette keys");
 if (canonical && count !== 24) throw new Error("Canonical Omarchy palette must fill all 24 slots");
 if (!canonical && count !== Object.keys(mapping.legacy).length) throw new Error("ANSI Omarchy palette must fill all 16 reference keys");
 return { ...base, id: "omarchy-data", name: "Omarchy data", slots };
}

export function setAppearance(appearance: Appearance): void {
 const scheme = appearance.theme === "omarchy" ? omarchyScheme(appearance.colors) :
  appearance.theme === "light" ? oqtoLight : oqtoDark;
 applyScheme(scheme, { radius: appearance.radius + "px" });
 document.documentElement.style.colorScheme = scheme.mode;
}
