import { Fragment } from "react";
import { Link, useLocation } from "react-router-dom";

import {
  labelForPortalNavPath,
  type PortalNavItem,
} from "@/modules/shared/layout/nav-types";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/shared/ui/components/breadcrumb";

function titleCaseSegment(segment: string): string {
  if (!segment) return segment;
  return segment
    .split("-")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

function portalBaseFromPath(pathname: string): string {
  const parts = pathname.split("/").filter(Boolean);
  const seg = parts[0] ?? "student";
  return `/${seg}`;
}

export function usePortalBreadcrumbs(
  navItems: PortalNavItem[],
  dashboardPath: string,
): { label: string; to?: string }[] {
  const { pathname } = useLocation();
  const normalized = pathname.replace(/\/$/, "") || dashboardPath;
  const portalBase = portalBaseFromPath(normalized);

  if (normalized === dashboardPath) {
    return [{ label: "Dashboard", to: dashboardPath }];
  }

  const teacherDetail = new RegExp(
    `^${portalBase.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}/teachers/[^/]+$`,
  );
  if (teacherDetail.test(normalized)) {
    return [
      { label: "Dashboard", to: dashboardPath },
      { label: "Teachers", to: `${portalBase}/teachers` },
      { label: "Teacher profile" },
    ];
  }

  const groupDetail = new RegExp(
    `^${portalBase.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}/groups/[^/]+$`,
  );
  if (groupDetail.test(normalized)) {
    return [
      { label: "Dashboard", to: dashboardPath },
      { label: "Groups", to: `${portalBase}/groups` },
      { label: "Group details" },
    ];
  }

  const lessonAttendance = new RegExp(
    `^${portalBase.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}/attendance/lesson/[^/]+$`,
  );
  if (lessonAttendance.test(normalized)) {
    return [
      { label: "Dashboard", to: dashboardPath },
      { label: "Attendance", to: `${portalBase}/attendance` },
      { label: "Lesson roster" },
    ];
  }

  const registered = labelForPortalNavPath(normalized, navItems);
  if (registered) {
    return [{ label: "Dashboard", to: dashboardPath }, { label: registered }];
  }

  const strip = normalized.startsWith(portalBase)
    ? normalized.slice(portalBase.length)
    : normalized;
  const pathParts = strip.split("/").filter(Boolean);
  const crumbs: { label: string; to?: string }[] = [
    { label: "Dashboard", to: dashboardPath },
  ];
  let acc = portalBase;
  for (let i = 0; i < pathParts.length; i++) {
    acc += `/${pathParts[i]}`;
    const isLast = i === pathParts.length - 1;
    const fromNav = navItems.find((n) => n.to === acc);
    const label = fromNav?.label ?? titleCaseSegment(pathParts[i]!);
    crumbs.push(isLast ? { label } : { label, to: acc });
  }
  return crumbs;
}

export function PortalBreadcrumbs({
  navItems,
  dashboardPath,
  className,
}: {
  navItems: PortalNavItem[];
  dashboardPath: string;
  className?: string;
}) {
  const crumbs = usePortalBreadcrumbs(navItems, dashboardPath);

  return (
    <Breadcrumb className={className}>
      <BreadcrumbList className="flex-nowrap">
        {crumbs.map((c, i) => (
          <Fragment key={`${c.label}-${i}`}>
            {i > 0 ? (
              <BreadcrumbSeparator className="shrink-0 [&>svg]:text-muted-foreground/70" />
            ) : null}
            <BreadcrumbItem className="min-w-0 shrink">
              {c.to ? (
                <BreadcrumbLink asChild>
                  <Link to={c.to} className="truncate">
                    {c.label}
                  </Link>
                </BreadcrumbLink>
              ) : (
                <BreadcrumbPage className="truncate">{c.label}</BreadcrumbPage>
              )}
            </BreadcrumbItem>
          </Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
