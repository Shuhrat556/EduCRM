import { KeyRound, LogOut, Menu, Moon, Settings, Sun } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";

import { useTheme } from "@/app/providers/theme-context";
import { formatRoleLabel, primaryRole, useAuth } from "@/modules/auth";
import { PortalBreadcrumbs } from "@/modules/shared/layout/portal-breadcrumbs";
import type { PortalNavItem } from "@/modules/shared/layout/nav-types";
import { useAdminShell } from "@/modules/shared/layout/admin-shell-context";
import { Avatar, AvatarFallback } from "@/shared/ui/components/avatar";
import { Badge } from "@/shared/ui/components/badge";
import { Button } from "@/shared/ui/components/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/components/dropdown-menu";
import { Separator } from "@/shared/ui/components/separator";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/shared/ui/components/tooltip";

type PortalHeaderProps = {
  onMenuClick: () => void;
  navItems: PortalNavItem[];
  dashboardPath: string;
  settingsPath: string;
  changePasswordPath: string;
  logoutRedirect: string;
};

export function PortalHeader({
  onMenuClick,
  navItems,
  dashboardPath,
  settingsPath,
  changePasswordPath,
  logoutRedirect,
}: PortalHeaderProps) {
  const { user, logout } = useAuth();
  const { resolved, toggle } = useTheme();
  const { sidebarCollapsed, toggleSidebar } = useAdminShell();
  const navigate = useNavigate();

  const initials =
    user?.displayName
      ?.split(/\s+/)
      .map((p) => p[0])
      .join("")
      .slice(0, 2)
      .toUpperCase() ||
    user?.email?.slice(0, 2).toUpperCase() ||
    "?";

  const mainRole = primaryRole(user?.roles);
  const extraRoles =
    user?.roles?.filter((r) => r !== mainRole).slice(0, 2) ?? [];

  return (
    <header className="sticky top-0 z-40 border-b bg-background/85 backdrop-blur-md supports-[backdrop-filter]:bg-background/70">
      <div className="flex h-14 items-center gap-2 px-3 sm:px-4 md:gap-3 md:px-5">
        <TooltipProvider delayDuration={200}>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="md:hidden"
                onClick={onMenuClick}
                aria-label="Open navigation"
              >
                <Menu className="h-5 w-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Menu</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="hidden md:inline-flex"
                onClick={toggleSidebar}
                aria-label={
                  sidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"
                }
              >
                <Menu className="h-5 w-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              {sidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>

        <Separator
          orientation="vertical"
          className="hidden h-6 md:block"
        />

        <div className="min-w-0 flex-1 overflow-hidden py-1">
          <PortalBreadcrumbs
            navItems={navItems}
            dashboardPath={dashboardPath}
          />
        </div>

        <div className="hidden items-center gap-1.5 sm:flex">
          {mainRole ? (
            <Badge variant="secondary" className="font-medium">
              {formatRoleLabel(mainRole)}
            </Badge>
          ) : null}
          {extraRoles.map((r) => (
            <Badge
              key={r}
              variant="outline"
              className="hidden font-normal text-muted-foreground lg:inline-flex"
            >
              {formatRoleLabel(r)}
            </Badge>
          ))}
        </div>

        <Separator orientation="vertical" className="hidden h-6 sm:block" />

        <div className="flex items-center gap-1">
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={toggle}
            aria-label={resolved === "dark" ? "Light mode" : "Dark mode"}
          >
            {resolved === "dark" ? (
              <Sun className="h-5 w-5" />
            ) : (
              <Moon className="h-5 w-5" />
            )}
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                type="button"
                variant="ghost"
                className="relative h-9 gap-2 rounded-full pl-1.5 pr-2 sm:pl-2 sm:pr-3"
                aria-label="Account menu"
              >
                <Avatar className="h-8 w-8 border border-border/60">
                  <AvatarFallback className="text-xs font-medium">
                    {initials}
                  </AvatarFallback>
                </Avatar>
                <span className="hidden max-w-[120px] truncate text-sm font-medium lg:inline">
                  {user?.displayName || user?.email || "Account"}
                </span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-60" align="end" forceMount>
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-2">
                  <div className="space-y-1">
                    <p className="text-sm font-medium leading-none">
                      {user?.displayName || "User"}
                    </p>
                    <p className="text-xs leading-none text-muted-foreground">
                      {user?.email}
                    </p>
                  </div>
                  {user?.roles?.length ? (
                    <div className="flex flex-wrap gap-1">
                      {user.roles.map((r) => (
                        <Badge key={r} variant="secondary" className="text-[10px]">
                          {formatRoleLabel(r)}
                        </Badge>
                      ))}
                    </div>
                  ) : null}
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem asChild>
                <Link
                  to={settingsPath}
                  className="flex cursor-default items-center gap-2"
                >
                  <Settings className="h-4 w-4" />
                  Settings
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link
                  to={changePasswordPath}
                  className="flex cursor-default items-center gap-2"
                >
                  <KeyRound className="h-4 w-4" />
                  Change password
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onSelect={() => {
                  void logout().then(() =>
                    navigate(logoutRedirect, { replace: true }),
                  );
                }}
              >
                <LogOut className="h-4 w-4" />
                Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
}
