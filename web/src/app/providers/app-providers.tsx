import { type ReactNode } from "react";

import { ThemeProvider } from "@/app/providers/theme-context";
import { AuthProvider } from "@/modules/auth";
import { QueryProvider } from "@/app/providers/query-provider";

export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <QueryProvider>
      <ThemeProvider>
        <AuthProvider>{children}</AuthProvider>
      </ThemeProvider>
    </QueryProvider>
  );
}
