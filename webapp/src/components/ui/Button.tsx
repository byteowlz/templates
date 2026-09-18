import { cn } from "@/lib/utils";
import type { ComponentProps } from "react";

/**
 * Button — the template's component exemplar. Tokens only, 44px touch target,
 * visible focus, disabled state. Add more primitives as needed; keep them
 * this lean and token-pure.
 */
export function Button({ className, variant = "default", ...props }: ComponentProps<"button"> & {
	variant?: "default" | "primary" | "danger" | "ghost";
}) {
	return (
		<button
			className={cn(
				"inline-flex min-h-11 items-center justify-center gap-2 rounded-[var(--radius-md)] border px-4",
				"transition-colors duration-150 select-none",
				"disabled:pointer-events-none disabled:opacity-45",
				{
					default:
						"border-border bg-surface text-foreground hover:bg-accent hover:border-ring",
					primary:
						"border-primary bg-primary text-primary-foreground font-semibold hover:brightness-110",
					danger: "border-transparent bg-transparent text-danger hover:bg-danger/10",
					ghost: "border-transparent bg-transparent text-muted-foreground hover:bg-accent hover:text-foreground",
				}[variant],
				className,
			)}
			{...props}
		/>
	);
}
