import { Button } from "@/components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";

/**
 * Seed surface: an Operate-mode shell (see DESIGN.md). Replace the body with
 * the tool's actual work surface; keep the token discipline and the shell
 * anatomy (header / main / status) — it is what impeccable audits expect.
 */
export function App() {
	return (
		<div className="mx-auto flex min-h-dvh max-w-5xl flex-col px-4 py-6 sm:px-6">
			<header className="flex flex-wrap items-baseline gap-x-3 gap-y-1 border-b border-border pb-4">
				<h1 className="text-base font-bold tracking-wide lowercase">{{project_name}}</h1>
				<p className="text-xs text-muted-foreground">what this workbench does — one line</p>
				<span
					aria-hidden="true"
					className="ml-auto inline-block size-2 rounded-full bg-success"
					title="service: up"
				/>
			</header>

			<main className="flex flex-1 flex-col gap-4 py-6">
				<Card>
					<CardHeader>
						<CardTitle>Get started</CardTitle>
					</CardHeader>
					<CardContent className="flex flex-col items-start gap-3">
						<p className="max-w-prose text-sm text-muted-foreground">
							Empty state. State clearly what to do first, in one sentence, with
							one primary action. Delete this card once the real surface exists.
						</p>
						<div className="flex flex-wrap gap-2">
							<Button variant="primary">Primary action</Button>
							<Button>Secondary action</Button>
						</div>
					</CardContent>
				</Card>
			</main>
		</div>
	);
}
