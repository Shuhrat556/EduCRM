import { Fragment } from "react";
import { Link, useLocation } from "react-router-dom";

import { adminNav, labelForAdminPath } from "@/layouts/admin-layout/nav-config";
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

/**
 * Builds trail: Dashboard → … → current page from `pathname` and registered nav.
 */
export function useAdminBreadcrumbs(): { label: string; to?: string }[] {
  const { pathname } = useLocation();
  const normalized = pathname.replace(/\/$/, "") || "/app";

  if (normalized === "/app") {
    return [{ label: "Dashboard", to: "/app" }];
  }

  const teacherDetail = /^\/app\/teachers\/[^/]+$/;
  if (teacherDetail.test(normalized)) {
    return [
      { label: "Dashboard", to: "/app" },
      { label: "Teachers", to: "/app/teachers" },
      { label: "Teacher profile" },
    ];
  }

  const groupDetail = /^\/app\/groups\/[^/]+$/;
  if (groupDetail.test(normalized)) {
    return [
      { label: "Dashboard", to: "/app" },
      { label: "Groups", to: "/app/groups" },
      { label: "Group details" },
    ];
  }

  const lessonAttendance = /^\/app\/attendance\/lesson\/[^/]+$/;
  if (lessonAttendance.test(normalized)) {
    return [
      { label: "Dashboard", to: "/app" },
      { label: "Attendance", to: "/app/attendance" },
      { label: "Lesson roster" },
    ];
  }

  const registered = labelForAdminPath(normalized);
  if (registered) {
    return [{ label: "Dashboard", to: "/app" }, { label: registered }];
  }

  const parts = normalized.replace(/^\/app\/?/, "").split("/").filter(Boolean);
  const crumbs: { label: string; to?: string }[] = [
    { label: "Dashboard", to: "/app" },
  ];
  let acc = "/app";
  for (let i = 0; i < parts.length; i++) {
    acc += `/${parts[i]}`;
    const isLast = i === parts.length - 1;
    const fromNav = adminNav.find((n) => n.to === acc);
    const label = fromNav?.label ?? titleCaseSegment(parts[i]!);
    crumbs.push(
      isLast ? { label } : { label, to: acc },
    );
  }
  return crumbs;
}

export function AdminBreadcrumbs({
  className,
}: {
  className?: string;
}) {
  const crumbs = useAdminBreadcrumbs();

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
