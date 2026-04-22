import { AlertTriangle } from "lucide-react";

import { Button } from "@/shared/ui/components/button";
import { cn } from "@/shared/lib/cn";

type AppErrorProps = {
  title?: string;
  message: string;
  onRetry?: () => void;
  className?: string;
};

export function AppError({
  title = "Something went wrong",
  message,
  onRetry,
  className,
}: AppErrorProps) {
  return (
    <div
      role="alert"
      className={cn(
        "flex flex-col items-center justify-center gap-4 rounded-xl border border-destructive/30 bg-destructive/5 px-6 py-10 text-center",
        className,
      )}
    >
      <span className="flex h-12 w-12 items-center justify-center rounded-full bg-destructive/15">
        <AlertTriangle
          className="h-6 w-6 text-destructive"
          aria-hidden
        />
      </span>
      <div className="space-y-1">
        <h3 className="text-base font-semibold text-foreground">{title}</h3>
        <p className="max-w-md text-sm text-muted-foreground">{message}</p>
      </div>
      {onRetry ? (
        <Button type="button" variant="outline" onClick={onRetry}>
          Try again
        </Button>
      ) : null}
    </div>
  );
}
