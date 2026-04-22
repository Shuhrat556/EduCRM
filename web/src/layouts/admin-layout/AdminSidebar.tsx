import { PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { NavLink } from "react-router-dom";

import { useAuth } from "@/modules/auth";
import { useAdminShell } from "@/modules/shared/layout/admin-shell-context";
import { adminNav } from "@/layouts/admin-layout/nav-config";
import { getEnv } from "@/shared/config/env";
import { cn } from "@/shared/lib/cn";
import { Button } from "@/shared/ui/components/button";
import { ScrollArea } from "@/shared/ui/components/scroll-area";
import { Separator } from "@/shared/ui/components/separator";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/shared/ui/components/tooltip";

type AdminSidebarProps = {
  onNavigate?: () => void;
  /** Mobile drawer: always expanded with full labels; hides collapse control. */
  variant?: "desktop" | "mobile";
};

export function AdminSidebar({
  onNavigate,
  variant = "desktop",
}: AdminSidebarProps) {
  const { user, hasRole } = useAuth();
  const { sidebarCollapsed, toggleSidebar } = useAdminShell();
  const title = getEnv().VITE_APP_NAME;

  const isMobileDrawer = variant === "mobile";
  const collapsed = isMobileDrawer ? false : sidebarCollapsed;

  const items = adminNav.filter(
    (item) => !item.roles || hasRole(item.roles),
  );

  return (
    <div className="flex h-full w-full flex-col bg-sidebar text-sidebar-foreground">
      <div
        className={cn(
          "flex h-14 shrink-0 items-center gap-1 border-b border-sidebar-border px-2",
          collapsed && !isMobileDrawer ? "justify-center px-1" : "justify-between pl-3 pr-1",
          isMobileDrawer && "justify-start px-3",
        )}
      >
        {!collapsed || isMobileDrawer ? (
          <span className="truncate text-sm font-semibold tracking-tight">
            {title}
          </span>
        ) : (
          <span className="sr-only">{title}</span>
        )}
        {!isMobileDrawer ? (
          <TooltipProvider delayDuration={0}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 shrink-0 text-sidebar-foreground/90 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                  onClick={toggleSidebar}
                  aria-label={
                    collapsed ? "Expand sidebar" : "Collapse sidebar"
                  }
                >
                  {collapsed ? (
                    <PanelLeftOpen className="h-4 w-4" />
                  ) : (
                    <PanelLeftClose className="h-4 w-4" />
                  )}
                </Button>
              </TooltipTrigger>
              <TooltipContent side="right" className="font-medium">
                {collapsed ? "Expand" : "Collapse"}
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        ) : null}
      </div>

      <ScrollArea className="min-h-0 flex-1">
        <TooltipProvider delayDuration={collapsed && !isMobileDrawer ? 0 : 300}>
          <nav className="flex flex-col gap-0.5 p-2" aria-label="Main">
            {items.map((item) => {
              const navLink = (
                <NavLink
                  key={item.to}
                  to={item.to}
                  end={item.to === "/app"}
                  title={collapsed && !isMobileDrawer ? item.label : undefined}
                  onClick={onNavigate}
                  className={({ isActive }) =>
                    cn(
                      "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                      collapsed && !isMobileDrawer && "justify-center px-2",
                      isActive
                        ? "bg-sidebar-accent text-sidebar-accent-foreground"
                        : "text-sidebar-foreground/85 hover:bg-sidebar-accent/75 hover:text-sidebar-accent-foreground",
                    )
                  }
                >
                  <item.icon
                    className="h-[18px] w-[18px] shrink-0 opacity-95"
                    aria-hidden
                  />
                  {!collapsed || isMobileDrawer ? (
                    <span className="truncate">{item.label}</span>
                  ) : null}
                </NavLink>
              );

              if (!collapsed || isMobileDrawer) return navLink;

              return (
                <Tooltip key={item.to}>
                  <TooltipTrigger asChild>{navLink}</TooltipTrigger>
                  <TooltipContent side="right" className="font-medium">
                    {item.label}
                  </TooltipContent>
                </Tooltip>
              );
            })}
          </nav>
        </TooltipProvider>
      </ScrollArea>

      <Separator className="bg-sidebar-border" />

      <div
        className={cn(
          "shrink-0 px-2 py-3 text-xs text-sidebar-foreground/65",
          collapsed && !isMobileDrawer && "flex justify-center px-0",
        )}
      >
        {(!collapsed || isMobileDrawer) && user ? (
          <p className="truncate px-1" title={user.email}>
            {user.email}
          </p>
        ) : null}
        {collapsed && !isMobileDrawer ? (
          <span className="sr-only">
            {user?.email ? `Signed in as ${user.email}` : "Signed in"}
          </span>
        ) : null}
      </div>
    </div>
  );
}
