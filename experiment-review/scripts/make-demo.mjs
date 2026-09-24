// Regenerate the template's own synthetic demonstration images.
// ImageMagick is needed only to regenerate the two GIFs, not to serve them.
import { writeFile, mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { resolve, dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const demo = resolve(dirname(fileURLToPath(import.meta.url)), "../demo");
const colors = {
  ink: "#17292c",
  hat: "#f4c97e",
  face: "#e9b987",
  coat: "#a9dfb2",
  trouser: "#5f989c",
};
const rect = (x, y, w, h, color) =>
  `<rect x="${x}" y="${y}" width="${w}" height="${h}" fill="${color}"/>`;
const svg = (width, height, contents, bg = "none") =>
  `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}"><rect width="100%" height="100%" fill="${bg}"/>${contents}</svg>`;

function figure(raised, baseline = false) {
  return [
    rect(12, 7, 9, 3, colors.ink),
    rect(11, 9, 11, 3, colors.hat),
    rect(13, 12, 7, 6, colors.face),
    rect(12, 18, 9, 12, colors.coat),
    rect(11, 30, 5, 9, colors.trouser),
    rect(17, 30, 5, 9, colors.trouser),
    rect(10, 39, 6, 2, colors.ink),
    rect(17, 39, 6, 2, colors.ink),
    raised
      ? rect(21, baseline ? 14 : 10, 3, 11, colors.coat)
      : rect(21, 20, 3, 10, colors.coat),
    rect(9, 21, 3, 9, colors.coat),
    ...(baseline && raised
      ? [rect(25, 14, 1, 1, colors.hat), rect(26, 18, 1, 1, colors.hat)]
      : []),
  ].join("");
}

await writeFile(
  join(demo, "source.svg"),
  svg(32, 48, figure(false), "#00afb5"),
);
await writeFile(
  join(demo, "diff.svg"),
  svg(
    96,
    48,
    `<g>${figure(false)}</g><g transform="translate(32)">${figure(true)}</g><g transform="translate(64)">${figure(true, true)}</g>`,
    "#0b1215",
  ),
);
await writeFile(
  join(demo, "repair.svg"),
  svg(
    128,
    64,
    `${rect(0, 0, 64, 64, "#00afb5")}<g transform="translate(16,8)">${figure(false)}</g><g transform="translate(80,8)">${figure(true)}</g>`,
    "#0b1215",
  ),
);

const tmp = await mkdtemp(join(tmpdir(), "review-demo-"));
try {
  for (const [name, baseline] of [
    ["method-a", false],
    ["method-b", true],
  ]) {
    const frames = [];
    for (let i = 0; i < 2; i++) {
      const source = join(tmp, `${name}-${i}.svg`);
      const frame = join(tmp, `${name}-${i}.png`);
      await writeFile(source, svg(32, 48, figure(Boolean(i), baseline)));
      const raster = spawnSync(
        "magick",
        ["-background", "none", source, frame],
        { encoding: "utf8" },
      );
      if (raster.status !== 0)
        throw new Error(
          raster.stderr || "Install ImageMagick to regenerate demo GIFs",
        );
      frames.push(frame);
    }
    const gif = spawnSync(
      "magick",
      [
        "-delay",
        "55",
        frames[0],
        "-delay",
        "55",
        frames[1],
        "-loop",
        "0",
        "-layers",
        "Optimize",
        join(demo, `${name}.gif`),
      ],
      { encoding: "utf8" },
    );
    if (gif.status !== 0) throw new Error(gif.stderr || "GIF creation failed");
  }
} finally {
  await rm(tmp, { recursive: true, force: true });
}
