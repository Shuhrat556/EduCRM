import { Loader2 } from "lucide-react";

import { cn } from "@/shared/lib/cn";

type AppLoadingProps = {
  label?: string;
  className?: string;
  fullPage?: boolean;
};

export function AppLoading({
  label = "Loading…",
  className,
  fullPage = false,
}: AppLoadingProps) {
  return (
    <div
      role="status"
      aria-live="polite"
      className={cn(
        "flex flex-col items-center justify-center gap-3 text-muted-foreground",
        fullPage && "min-h-dvh",
        className,
      )}
    >
      <Loader2 className="h-8 w-8 animate-spin text-primary" aria-hidden />
      <p className="text-sm font-medium text-foreground">{label}</p>
    </div>
  );
}
