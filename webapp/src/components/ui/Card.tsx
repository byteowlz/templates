import { cn } from "@/lib/utils";
import type { ComponentProps } from "react";

/** Card — elevated surface role exemplar. */
export function Card({ className, ...props }: ComponentProps<"section">) {
	return (
		<section
			className={cn(
				"rounded-[var(--radius-lg)] border border-border bg-card text-card-foreground",
				className,
			)}
			{...props}
		/>
	);
}

export function CardHeader({ className, ...props }: ComponentProps<"header">) {
	return (
		<header
			className={cn("flex flex-col gap-1 border-b border-border px-5 py-4", className)}
			{...props}
		/>
	);
}

export function CardTitle({ className, ...props }: ComponentProps<"h2">) {
	return <h2 className={cn("text-sm font-semibold tracking-wide", className)} {...props} />;
}

export function CardContent({ className, ...props }: ComponentProps<"div">) {
	return <div className={cn("px-5 py-4", className)} {...props} />;
}
