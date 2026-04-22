import type { ReactElement } from "react";

import { RoleRoute } from "@/app/router/guards/RoleRoute";
import type { Role } from "@/modules/auth/model/types";

/** Wraps a route element with `RoleRoute` — keeps `dashboard-routes` declarative. */
export function guardedElement(
  roles: Role | Role[],
  element: ReactElement,
): ReactElement {
  return <RoleRoute roles={roles}>{element}</RoleRoute>;
}
