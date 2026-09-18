import { App } from "@/App";
import "@/index.css";
// Omarchy adapter output. Committed as a no-op passthrough; `just omarchy`
// regenerates it from the machine's active omarchy theme (docs/omarchy.md).
// Kept after index.css so slot overrides win.
import "@/theme/omarchy.generated.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<App />
	</StrictMode>,
);
