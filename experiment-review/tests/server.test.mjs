import test from "node:test";
import assert from "node:assert/strict";
import { once } from "node:events";
import { createReviewServer, media } from "../server.mjs";

test("curated media is available; private paths are denied", async (t) => {
  const server = createReviewServer().listen(0, "127.0.0.1");
  await once(server, "listening");
  t.after(() => server.close());
  const address = server.address();
  assert.ok(address && typeof address !== "string");
  const origin = `http://127.0.0.1:${address.port}`;
  const get = (path, options) => fetch(origin + path, options);
  const page = await get("/");
  assert.equal(page.status, 200);
  assert.match(await page.text(), /synthetic/i);
  assert.equal((await get("/health")).status, 200);
  for (const path of Object.keys(media))
    assert.equal((await get(path, { method: "HEAD" })).status, 200, path);
  const gif = await get("/media/method-a.gif");
  assert.equal(gif.headers.get("content-type"), "image/gif");
  assert.equal(
    Buffer.from(await gif.arrayBuffer())
      .subarray(0, 6)
      .toString(),
    "GIF89a",
  );
  for (const path of [
    "/server.mjs",
    "/.pi/data/private",
    "/data/corpus/asset.png",
    "/demo/method-a.gif",
    "/media/%2e%2e%2fserver.mjs",
    "/media/method-a.gif/../../server.mjs",
    "/media/method-a.gif%00",
    "/media/unknown.gif",
  ])
    assert.equal((await get(path)).status, 404, path);
  assert.equal(
    (await get("/media/method-a.gif", { method: "POST" })).status,
    405,
  );
});
