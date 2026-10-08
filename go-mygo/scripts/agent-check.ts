import { mkdtempSync, mkdirSync, readdirSync, writeFileSync, rmSync, readFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { resolve, join } from "node:path";
import assert from "node:assert/strict";
import Ajv from "ajv";
import { Desktop, type Request, type Envelope } from "../src/mygo";
import schema from "../schemas/envelope.schema.json";

const binary = resolve(process.argv[2] ?? "build/{{project_name}}ctl");
const sandbox = mkdtempSync(join(tmpdir(), "mygo-agent-check-"));
const validate = new Ajv({strict:false}).compile(schema);
const observations: unknown[] = [];
const env: Record<string, string | undefined> = { ...process.env, XDG_CONFIG_HOME: join(sandbox,"config") };
delete env.MYGO_APP_THEME;
delete env.MYGO_APP_RADIUS;
delete env.MYGO_APP_OMARCHY_FILE;
function cli(args: string[], input = "", expected = 0): Envelope {
 const run = Bun.spawnSync([binary,...args,"--json"], {
  cwd:sandbox, env, stdin:Buffer.from(input), stdout:"pipe", stderr:"pipe", timeout:5000,
 });
 assert.equal(run.exitCode,expected,new TextDecoder().decode(run.stderr));
 const reply: Envelope = JSON.parse(new TextDecoder().decode(run.stdout));
 assert(validate(reply),JSON.stringify(validate.errors));
 observations.push({args,exit:run.exitCode,envelope:reply});
 return reply;
}
try {
 const help = Bun.spawnSync([binary,"--help"],{cwd:sandbox,env,timeout:5000});
 assert.equal(help.exitCode,0);
 assert(new TextDecoder().decode(help.stdout).includes("snapshot"));
 const first = cli(["snapshot"]);
 const portable = JSON.stringify(first.result!.snapshot);
 writeFileSync(join(sandbox,"snapshot.json"),portable);
 assert.deepEqual(cli(["validate","--input","snapshot.json"]).result,first.result);
 assert.deepEqual(cli(["validate"],portable).result,first.result);
 assert.equal(cli(["apply"],"",1).error!.code,"AUTHORITY_DENIED");
 assert.equal(cli(["apply","--input","missing.json"],"",1).error!.code,"AUTHORITY_DENIED");
 assert.equal(cli(["validate"],portable.replace("fixture-v1","stale"),2).error!.code,"STALE_REVISION");
 assert.equal(cli(["validate"],"not-json",2).error!.code,"INVALID_INPUT");
 assert.equal(cli(["validate"],portable + " {}",2).error!.code,"INVALID_INPUT");
 assert.equal(cli(["validate"],"x".repeat(65537),2).error!.code,"INVALID_INPUT");
 assert.equal(cli(["bogus"],"",2).error!.code,"USAGE");
 assert.equal(cli(["snapshot","--config","missing"],"",2).error!.code,"INVALID_CONFIG");
 // Adapter harness, NOT a real webview: exercise generated method/argument dispatch.
 const globals = globalThis as typeof globalThis & {mygo?: unknown};
 globals.mygo = { call: async (method:string,request:Request) => {
  assert.equal(method,"Desktop.Execute");
  const exit = request.operation === "apply" ? 1 : 0;
  return cli([request.operation],request.input,exit);
 }};
 try {
  assert.deepEqual(await Desktop.execute({operation:"snapshot",input:""}),first);
  assert.equal((await Desktop.execute({operation:"apply",input:portable})).error!.code,"AUTHORITY_DENIED");
 } finally { delete globals.mygo; }
 assert.deepEqual(cli(["snapshot"]),first); // denied apply is idempotent/no state change.
 assert.deepEqual(readdirSync(sandbox),["snapshot.json"]); // no default config, profile or writes.
 const drift = Bun.spawnSync(["bun","test","tests"],{cwd:process.cwd(),timeout:30000});
 assert.equal(drift.exitCode,0,new TextDecoder().decode(drift.stderr));
 const artifact = "artifacts/agent-check/observations.json";
 mkdirSync("artifacts/agent-check",{recursive:true});
 writeFileSync(artifact,JSON.stringify(observations,null,2)+"\n");
 const evidence = (observation:string) => [{procedure:"just agent-check",artifact,observation}];
 const gate = (status:string,reason:string,observation=reason) => ({status,reason,evidence:evidence(observation)});
 const report = {
  schema_version:"agent-readiness/v1",project:"{{project_name}}",revision:process.env.TEMPLATE_REVISION ?? "instantiated-scaffold",
  scope:"Read-only synthetic fixture; Go CLI, Go Desktop facade, generated TypeScript dispatch harness. Not a native IPC or production/hardware audit.",
  platforms:[process.platform + "/" + process.arch + " fixture"],
  gates:{
   discoverability:gate("pass","Help describes read-only commands and JSON examples; exercised with a 5s process bound."),
   structured_contract:gate("pass","Actual executable JSON success/usage/input/stale/denied envelopes validated against schema."),
   shared_substrate:gate("unverified","Go facade parity tests pass; generated dispatch harness passes. An open native UI/IPC transport has not been exercised by this gate."),
   authority_and_safety:gate("pass","Fixture scope only: apply always denied, no adapter/authority; snapshots unchanged and sandbox filesystem unchanged except explicit export."),
   concurrency_and_recovery:gate("pass","Immutable reads idempotent; race-enabled Go tests cover concurrency/cancellation. No write, transport or partial-failure capability exists in scope."),
   portable_state:gate("pass","Snapshot file and stdin roundtrips use actual CLI; no developer checkout/services required."),
   reproducibility:gate("pass","Go tests, real executable, generated dispatch and theme/provenance tests executed. Platform/hardware claims excluded."),
   human_accessibility:gate("unverified","Native keyboard, focus, names/status and actual rendering require a manual or native accessibility test; source controls alone are not proof."),
  },
  };
 writeFileSync("artifacts/agent-check/report.json",JSON.stringify(report,null,2)+"\n");
 console.log(JSON.stringify({ok:true,scope:"fixture-contract",report:"artifacts/agent-check/report.json",observations:artifact,
  production_ready:false,native_ipc_verified:false,report_structure_is_not_evidence:true}));
} finally { rmSync(sandbox,{recursive:true,force:true}); }
