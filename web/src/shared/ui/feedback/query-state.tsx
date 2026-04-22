import type { QueryObserverResult } from "@tanstack/react-query";
import type { ReactNode } from "react";

import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { AppError } from "@/shared/ui/feedback/app-error";
import { AppLoading } from "@/shared/ui/feedback/app-loading";
import { ApiError } from "@/shared/api/types";

type QueryStateProps<T> = {
  query: Pick<
    QueryObserverResult<T, unknown>,
    "isPending" | "isError" | "error" | "isFetching" | "refetch" | "data"
  >;
  loadingLabel?: string;
  empty?: boolean;
  emptyTitle?: string;
  emptyDescription?: string;
  children: (data: NonNullable<T>) => ReactNode;
};

function errorMessage(error: unknown): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "An unexpected error occurred.";
}

export function QueryState<T>({
  query,
  loadingLabel,
  empty,
  emptyTitle = "No data",
  emptyDescription,
  children,
}: QueryStateProps<T>) {
  if (query.isPending || (query.isFetching && query.data === undefined)) {
    return <AppLoading label={loadingLabel} fullPage />;
  }

  if (query.isError) {
    return (
      <AppError
        message={errorMessage(query.error)}
        onRetry={() => void query.refetch()}
      />
    );
  }

  const data = query.data;
  if (data === undefined || data === null || empty) {
    return (
      <AppEmpty title={emptyTitle} description={emptyDescription} />
    );
  }

  return <>{children(data)}</>;
}
