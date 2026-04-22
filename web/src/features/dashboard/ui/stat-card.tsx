import type { LucideIcon } from "lucide-react";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { Skeleton } from "@/shared/ui/components/skeleton";
import { cn } from "@/shared/lib/cn";

type StatCardProps = {
  title: string;
  value: string | number;
  subtitle?: string;
  icon: LucideIcon;
  loading?: boolean;
  variant?: "default" | "warning";
};

export function StatCard({
  title,
  value,
  subtitle,
  icon: Icon,
  loading,
  variant = "default",
}: StatCardProps) {
  return (
    <Card
      className={cn(
        "overflow-hidden border-border/80 shadow-sm",
        variant === "warning" &&
          "border-amber-500/25 bg-amber-500/[0.03] dark:bg-amber-500/[0.06]",
      )}
    >
      <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
        <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
          {title}
        </CardTitle>
        <span
          className={cn(
            "flex h-9 w-9 items-center justify-center rounded-lg",
            variant === "warning"
              ? "bg-amber-500/15 text-amber-700 dark:text-amber-400"
              : "bg-primary/10 text-primary",
          )}
        >
          <Icon className="h-4 w-4" aria-hidden />
        </span>
      </CardHeader>
      <CardContent className="pt-0">
        {loading ? (
          <div className="space-y-2">
            <Skeleton className="h-8 w-24" />
            <Skeleton className="h-3 w-32" />
          </div>
        ) : (
          <>
            <p className="text-2xl font-semibold tracking-tight tabular-nums text-foreground">
              {value}
            </p>
            {subtitle ? (
              <p className="mt-1 text-xs text-muted-foreground">{subtitle}</p>
            ) : null}
          </>
        )}
      </CardContent>
    </Card>
  );
}
