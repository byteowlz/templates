import { cpSync, mkdtempSync, readFileSync, readdirSync, statSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";

const source = resolve(import.meta.dir,"..");
const target = mkdtempSync(join(tmpdir(),"mygo-scaffold-"));
const name = "mygo-proof";
cpSync(source,target,{recursive:true,filter:(path) => !/(^|\/)(node_modules|build|dist|artifacts|\.mygo|\.git)(\/|$)/.test(path)});
for(const relative of readdirSync(target,{recursive:true})) {
 const file = join(target,String(relative));
 if(!statSync(file).isFile()) continue;
 const text = readFileSync(file,"utf8");
 if(text.includes("{{project_name}}")) writeFileSync(file,text.replaceAll("{{project_name}}",name));
}
const env: Record<string, string | undefined> = {...process.env, XDG_CONFIG_HOME:join(target,"isolated-config")};
for(const key of ["MYGO_APP_THEME","MYGO_APP_RADIUS","MYGO_APP_OMARCHY_FILE"]) delete env[key];
const commands = [
 ["bun","install","--frozen-lockfile"],["go","mod","download"],
 ["just","check"],["just","build"],["just","agent-check"],
];
for(const command of commands) {
 const run = Bun.spawnSync(command,{cwd:target,env,stdout:"inherit",stderr:"inherit",timeout:120000});
 if(run.exitCode!==0) throw new Error(command.join(" ")+" failed; scaffold retained at "+target);
}
console.log(JSON.stringify({ok:true,scaffold:target,name,commands,limits:["native render not exercised","full signing/notarization not verified"]}));
