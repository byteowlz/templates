import { test, expect } from "bun:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { createHash } from "node:crypto";
import hashes from "../third_party/hashes.json";
import mapping from "../third_party/omarchy/mapping.json";
import golden from "./fixtures/reference-tokens.json";
import { mapSchemeToTokens, normalizeScheme, oqtoDark, oqtoLight, nordBase16, radiusVars } from "../third_party/design-system/src/index";
import { omarchyScheme } from "../src/theme";

test("vendored snapshot and data mapping have no drift", () => {
 for(const [path,hash] of Object.entries(hashes)) {
  expect(createHash("sha256").update(readFileSync(path)).digest("hex")).toBe(hash);
 }
 const source = readFileSync("third_party/omarchy/reference.mjs","utf8");
 for(const [name,table] of [["CANONICAL",mapping.canonical],["LEGACY_ANSI",mapping.legacy]] as const) {
  const block = source.match(new RegExp("const "+name+" = \\{([\\s\\S]*?)\\};"))![1]!;
  const expected=Object.fromEntries([...block.matchAll(/(base[0-9A-F]{2}): "([^"]+)"/g)].map(x=>[x[1],x[2]]));
  assert.deepEqual(table,expected);
 }
});

test("reference Base24 roles, Base16 widening, R1 dial parity", () => {
 const actual = {dark:mapSchemeToTokens(normalizeScheme(oqtoDark)),light:mapSchemeToTokens(normalizeScheme(oqtoLight)),
  base16:normalizeScheme(nordBase16),sharp:radiusVars("0px"),round:radiusVars("8px")};
 assert.deepEqual(actual,golden);
});

test("Omarchy data overlays only reference slots; malformed input fails before apply", () => {
 const canonical=Object.fromEntries(Object.entries(mapping.canonical).map(([slot,key])=>[key,oqtoDark.slots[slot as keyof typeof oqtoDark.slots]!]));
 expect(omarchyScheme(canonical).slots).toEqual(oqtoDark.slots);
 const ansi=Object.fromEntries(Object.entries(mapping.legacy).map(([slot,key])=>[key,oqtoDark.slots[slot as keyof typeof oqtoDark.slots]!]));
 expect(omarchyScheme(ansi).slots).toEqual(oqtoDark.slots);
 const malformed: Record<string,string>[] = [{}, {color0:"#ffffff"}, {...canonical,red:"url(https://invalid.test)"}, {...canonical,mode:"invalid"}];
 for(const bad of malformed) {
  expect(()=>omarchyScheme(bad)).toThrow();
 }
 const partial={...canonical}; delete partial.red;
 expect(()=>omarchyScheme(partial)).toThrow();
});
