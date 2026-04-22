import { cn } from "@/shared/lib/cn";

export type SummaryMetricTileProps = {
  label: string;
  value: string;
  hint?: string;
  className?: string;
};

/** Compact KPI cell for analytics-style layouts (reports, dashboards). */
export function SummaryMetricTile({
  label,
  value,
  hint,
  className,
}: SummaryMetricTileProps) {
  return (
    <div
      className={cn(
        "rounded-lg border border-border/60 bg-muted/20 px-3 py-2.5",
        className,
      )}
    >
      <p className="text-[11px] font-medium uppercase tracking-wide text-muted-foreground">
        {label}
      </p>
      <p className="text-xl font-semibold tabular-nums tracking-tight text-foreground">
        {value}
      </p>
      {hint ? (
        <p className="text-xs text-muted-foreground/90">{hint}</p>
      ) : null}
    </div>
  );
}
