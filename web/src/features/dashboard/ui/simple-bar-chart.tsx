import type { ChartPoint } from "@/features/dashboard/model/types";
import { cn } from "@/shared/lib/cn";

type SimpleBarChartProps = {
  data: ChartPoint[];
  /** Axis label for accessibility */
  caption: string;
  className?: string;
  emptyHint?: string;
};

/**
 * Lightweight horizontal bar row — swap for Recharts / ECharts when API is stable.
 */
export function SimpleBarChart({
  data,
  caption,
  className,
  emptyHint = "Chart data will appear when your reporting API is connected.",
}: SimpleBarChartProps) {
  const max = Math.max(1, ...data.map((d) => d.value));

  if (!data.length) {
    return (
      <div
        className={cn(
          "flex min-h-[160px] items-center justify-center rounded-lg border border-dashed border-border bg-muted/20 px-4 text-center text-sm text-muted-foreground",
          className,
        )}
      >
        {emptyHint}
      </div>
    );
  }

  return (
    <div
      className={cn("space-y-3", className)}
      role="img"
      aria-label={caption}
    >
      <p className="sr-only">{caption}</p>
      {data.map((point) => (
        <div key={point.label} className="space-y-1">
          <div className="flex items-center justify-between text-xs">
            <span className="font-medium text-foreground">{point.label}</span>
            <span className="tabular-nums text-muted-foreground">
              {point.value}
            </span>
          </div>
          <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
            <div
              className="h-full rounded-full bg-primary/80 transition-all duration-500"
              style={{ width: `${(point.value / max) * 100}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}
