import { useEffect, useState } from "react";
import { Outlet } from "react-router-dom";

import { useMediaQuery } from "@/shared/lib/hooks/useMediaQuery";
import { AdminHeader } from "@/layouts/admin-layout/AdminHeader";
import { AdminShellProvider, useAdminShell } from "@/modules/shared/layout/admin-shell-context";
import { AdminSidebar } from "@/layouts/admin-layout/AdminSidebar";
import { cn } from "@/shared/lib/cn";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/shared/ui/components/sheet";

function AdminLayoutFrame() {
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
          <AdminSidebar variant="desktop" />
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
            <AdminSidebar
              variant="mobile"
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
        <AdminHeader onMenuClick={() => setMobileOpen(true)} />
        <main className="flex-1 overflow-x-hidden">
          <div className="mx-auto h-full w-full max-w-[1600px] px-4 py-5 sm:px-5 sm:py-6 md:px-6 md:py-8">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}

export function AdminLayout() {
  return (
    <AdminShellProvider>
      <AdminLayoutFrame />
    </AdminShellProvider>
  );
}
