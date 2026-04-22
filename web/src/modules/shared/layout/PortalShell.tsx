import { useEffect, useState, type ReactNode } from "react";
import { Outlet } from "react-router-dom";

import type { PortalNavItem } from "@/modules/shared/layout/nav-types";
import { PortalHeader } from "@/modules/shared/layout/PortalHeader";
import { PortalSidebar } from "@/modules/shared/layout/PortalSidebar";
import { useMediaQuery } from "@/shared/lib/hooks/useMediaQuery";
import { AdminShellProvider, useAdminShell } from "@/modules/shared/layout/admin-shell-context";
import { cn } from "@/shared/lib/cn";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/shared/ui/components/sheet";

export type PortalShellProps = {
  navItems: PortalNavItem[];
  dashboardPath: string;
  settingsPath: string;
  changePasswordPath: string;
  logoutRedirect: string;
  children?: ReactNode;
};

function PortalFrame({
  navItems,
  dashboardPath,
  settingsPath,
  changePasswordPath,
  logoutRedirect,
  children,
}: PortalShellProps) {
  const isDesktop = useMediaQuery("(min-width: 768px)");
  const [mobileOpen, setMobileOpen] = useState(false);
  const { sidebarCollapsed } = useAdminShell();

  useEffect(() => {
    if (isDesktop) setMobileOpen(false);
  }, [isDesktop]);

  return (
    <div className="flex min-h-dvh w-full bg-muted/20">
      {isDesktop ? (
        <aside
          className={cn(
            "fixed inset-y-0 left-0 z-30 hidden flex-col border-r border-sidebar-border bg-sidebar shadow-sm transition-[width] duration-200 ease-in-out md:flex",
            sidebarCollapsed ? "w-[4.5rem]" : "w-64",
          )}
        >
          <PortalSidebar
            variant="desktop"
            items={navItems}
            dashboardPath={dashboardPath}
          />
        </aside>
      ) : (
        <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
          <SheetContent
            side="left"
            className="w-[min(100vw-2rem,20rem)] border-sidebar-border bg-sidebar p-0"
          >
            <SheetHeader className="sr-only">
              <SheetTitle>Navigation</SheetTitle>
            </SheetHeader>
            <PortalSidebar
              variant="mobile"
              items={navItems}
              dashboardPath={dashboardPath}
              onNavigate={() => setMobileOpen(false)}
            />
          </SheetContent>
        </Sheet>
      )}

      <div
        className={cn(
          "flex min-w-0 flex-1 flex-col transition-[padding] duration-200 ease-in-out",
          isDesktop && (sidebarCollapsed ? "md:pl-[4.5rem]" : "md:pl-64"),
        )}
      >
        <PortalHeader
          onMenuClick={() => setMobileOpen(true)}
          navItems={navItems}
          dashboardPath={dashboardPath}
          settingsPath={settingsPath}
          changePasswordPath={changePasswordPath}
          logoutRedirect={logoutRedirect}
        />
        <main className="flex-1 overflow-x-hidden">
          <div className="mx-auto h-full w-full max-w-[1600px] px-4 py-5 sm:px-5 sm:py-6 md:px-6 md:py-8">
            {children ?? <Outlet />}
          </div>
        </main>
      </div>
    </div>
  );
}

export function PortalShell(props: PortalShellProps) {
  return (
    <AdminShellProvider>
      <PortalFrame {...props} />
    </AdminShellProvider>
  );
}
