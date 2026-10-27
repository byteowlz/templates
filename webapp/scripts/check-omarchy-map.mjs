#!/usr/bin/env node
/**
 * Drift check: this webapp's Omarchy mapping vs the design-system spec map.
 *
 * The webapp is a scaffold that ships to new repos, so it cannot import
 * design-system. It therefore verifies instead: the CANONICAL and LEGACY_ANSI
 * tables in scripts/omarchy-theme.mjs must stay consistent with
 * design-system/spec/omarchy-map.json (contract: design-system/spec/omarchy.md).
 *
 * Offline by default (skips when the spec map is not on this machine). Runs in
 * CI/`just check` where the workspace exists. Never writes.
 *
 *   node scripts/check-omarchy-map.mjs
 *   DESIGN_SYSTEM_MAP=/path/to/omarchy-map.json node scripts/check-omarchy-map.mjs
 */

import { readFileSync, existsSync } from "node:fs";
import { homedir } from "node:os";
import path from "node:path";

const MAP_CANDIDATES = [
	process.env.DESIGN_SYSTEM_MAP,
	path.join(homedir(), "byteowlz/design-system/spec/omarchy-map.json"),
];

const found = MAP_CANDIDATES.find((p) => p && existsSync(p));
if (!found) {
	console.log("omarchy-map drift: skipped (design-system/spec/omarchy-map.json not found)");
	process.exit(0);
}

const here = new URL(".", import.meta.url).pathname;
const script = readFileSync(path.join(here, "omarchy-theme.mjs"), "utf8");

function tableFromScript(name) {
	const start = script.split(`const ${name} = {`)[1];
	if (!start) throw new Error(`table ${name} not found in omarchy-theme.mjs`);
	const block = start.split("};")[0];
	const out = {};
	for (const m of block.matchAll(/(base[0-9A-F]{2}):\s*"([^"]+)"/g)) out[m[1]] = m[2];
	return out;
}

const spec = JSON.parse(readFileSync(found, "utf8"));
const failures = [];

function compare(label, actual, expected, { partial }) {
	const expectedKeys = Object.keys(expected);
	if (!partial && Object.keys(actual).length !== expectedKeys.length) {
		failures.push(`${label}: ${Object.keys(actual).length} entries, expected ${expectedKeys.length}`);
	}
	for (const [slot, key] of Object.entries(expected)) {
		if (actual[slot] !== key) {
			failures.push(`${label}[${slot}] = ${actual[slot] ?? "(missing)"}, spec ${key}`);
		}
	}
	for (const slot of Object.keys(actual)) {
		if (!(slot in expected)) failures.push(`${label}[${slot}] is not in the spec map`);
	}
}

// The canonical table is total; the legacy ANSI overlay is intentionally partial.
compare("CANONICAL", tableFromScript("CANONICAL"), spec.canonical, { partial: false });
compare("LEGACY_ANSI", tableFromScript("LEGACY_ANSI"), spec.legacyAnsi, { partial: true });

if (failures.length) {
	console.error(`omarchy-map drift against ${found}:`);
	for (const f of failures) console.error(`  ${f}`);
	process.exit(1);
}
console.log(`omarchy-map drift: ok (canonical ${Object.keys(spec.canonical).length}, legacy ${Object.keys(spec.legacyAnsi).length})`);
