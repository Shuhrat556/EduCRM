import { useMemo } from "react";

import { useAuth } from "@/modules/auth/hooks/useAuth";
import {
  canPermission,
  type Permission,
} from "@/modules/auth/lib/permissions";

export function usePermissions() {
  const { user } = useAuth();
  const roles = user?.roles;

  return useMemo(
    () => ({
      can: (permission: Permission) => canPermission(roles, permission),
    }),
    [roles],
  );
}
