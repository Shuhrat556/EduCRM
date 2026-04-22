import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { type ReactNode, useState } from "react";

import { isQueryDevtoolsEnabled } from "@/shared/config/env";

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        retry: (failureCount, error) => {
          if (
            typeof error === "object" &&
            error !== null &&
            "status" in error &&
            (error as { status: number }).status === 401
          ) {
            return false;
          }
          return failureCount < 2;
        },
      },
    },
  });
}

export function QueryProvider({ children }: { children: ReactNode }) {
  const [client] = useState(makeQueryClient);

  return (
    <QueryClientProvider client={client}>
      {children}
      {isQueryDevtoolsEnabled() ? (
        <ReactQueryDevtools buttonPosition="bottom-left" />
      ) : null}
    </QueryClientProvider>
  );
}
