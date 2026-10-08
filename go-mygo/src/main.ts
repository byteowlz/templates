import "./style.css";
import { isMyGo } from "mygo-runtime";
import { Desktop, type Snapshot } from "./mygo";
import { setAppearance } from "./theme";
import { oqtoDark, applyScheme } from "../third_party/design-system/src/index";

applyScheme(oqtoDark, { radius: "0px" });
const status = document.querySelector<HTMLParagraphElement>("#status")!;
const output = document.querySelector<HTMLPreElement>("#snapshot")!;
const refresh = document.querySelector<HTMLButtonElement>("#refresh")!;
const validate = document.querySelector<HTMLButtonElement>("#validate")!;
let snapshot: Snapshot | null = null;

async function execute(operation: string, input = ""): Promise<void> {
 refresh.disabled = true;
 validate.disabled = true;
 status.textContent = "Reading shared service…";
 try {
  const reply = await Desktop.execute({ operation, input });
  if (!reply.ok) {
   status.textContent = reply.error?.code + ": " + reply.error?.message;
   return;
  }
  snapshot = reply.result?.snapshot ?? null;
  output.textContent = JSON.stringify(snapshot, null, 2);
  status.textContent = operation === "validate" ? "Snapshot validated. No writes performed." : "Read-only fixture loaded. Export via the CLI or validate here.";
 } catch {
  status.textContent = "IPC unavailable. Run the desktop app with just dev, then retry.";
 } finally {
  refresh.disabled = false;
  validate.disabled = snapshot === null;
 }
}

refresh.addEventListener("click", () => void execute("snapshot"));
validate.addEventListener("click", () => void execute("validate", JSON.stringify(snapshot)));
if (isMyGo()) {
 void (async () => {
  try {
   const reply = await Desktop.execute({ operation: "appearance", input: "" });
   if (!reply.ok || !reply.result?.appearance) throw new Error(reply.error?.message ?? "Appearance unavailable");
   setAppearance(reply.result.appearance);
   await execute("snapshot");
  } catch (error) {
   status.textContent = "Theme unavailable; safe house default retained. " + String(error);
   refresh.disabled = false;
  }
 })();
} else {
 status.textContent = "Browser preview has no Go IPC. Run just dev for the real fixture.";
 refresh.disabled = true;
}
