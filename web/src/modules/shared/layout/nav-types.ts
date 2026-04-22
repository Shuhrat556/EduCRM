import type { LucideIcon } from "lucide-react";

import type { Role } from "@/modules/auth/model/types";

export type PortalNavItem = {
  label: string;
  to: string;
  icon: LucideIcon;
  roles?: Role[];
};

export function labelForPortalNavPath(
  pathname: string,
  items: PortalNavItem[],
): string | undefined {
  const normalized = pathname.replace(/\/$/, "") || "";
  const match = items.find((item) => item.to === normalized);
  return match?.label;
}
