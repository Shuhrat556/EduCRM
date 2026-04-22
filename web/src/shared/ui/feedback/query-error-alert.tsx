import { AlertCircle } from "lucide-react";

import { ApiError } from "@/shared/api/types";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/shared/ui/components/alert";
import { Button } from "@/shared/ui/components/button";
import { cn } from "@/shared/lib/cn";

function messageFromUnknown(error: unknown): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "Something went wrong. Try again.";
}

export type QueryErrorAlertProps = {
  error: unknown;
  title?: string;
  onRetry?: () => void;
  retryLabel?: string;
  className?: string;
};

export function QueryErrorAlert({
  error,
  title = "Could not load data",
  onRetry,
  retryLabel = "Retry",
  className,
}: QueryErrorAlertProps) {
  return (
    <Alert variant="destructive" className={cn("border-destructive/40", className)}>
      <AlertCircle className="h-4 w-4" aria-hidden />
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <span className="text-sm">{messageFromUnknown(error)}</span>
        {onRetry ? (
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="shrink-0 border-destructive/40 bg-background text-destructive hover:bg-destructive/10 hover:text-destructive"
            onClick={onRetry}
          >
            {retryLabel}
          </Button>
        ) : null}
      </AlertDescription>
    </Alert>
  );
}
