import { ShieldOff } from "lucide-react";
import { Link, useLocation, useNavigate } from "react-router-dom";

import {
  formatRoleLabel,
  getDefaultDashboardPath,
  primaryRole,
  resolveLoginPathForRoles,
  useAuth,
} from "@/modules/auth";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";

export function UnauthorizedPage() {
  const { logout, user } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const attempted = (location.state as { from?: string } | null)?.from;
  const safeHome = user?.roles?.length
    ? getDefaultDashboardPath(user.roles)
    : "/login";
  const mainRole = primaryRole(user?.roles);

  return (
    <div className="flex min-h-dvh items-center justify-center bg-muted/30 p-4">
      <Card className="w-full max-w-md text-center shadow-md">
        <CardHeader className="items-center space-y-4">
          <span className="flex h-14 w-14 items-center justify-center rounded-full bg-destructive/10">
            <ShieldOff className="h-7 w-7 text-destructive" aria-hidden />
          </span>
          <div>
            <CardTitle className="text-xl">Access denied</CardTitle>
            <CardDescription className="mt-2">
              Your account does not have permission to view this area. Contact an
              administrator if you believe this is a mistake.
            </CardDescription>
            {mainRole ? (
              <p className="mt-3 text-xs text-muted-foreground">
                Signed in as{" "}
                <span className="font-medium text-foreground">
                  {formatRoleLabel(mainRole)}
                </span>
                {attempted ? (
                  <>
                    {" "}
                    · Blocked path:{" "}
                    <code className="rounded bg-muted px-1 py-0.5 text-[11px]">
                      {attempted}
                    </code>
                  </>
                ) : null}
              </p>
            ) : null}
          </div>
        </CardHeader>
        <CardContent className="flex flex-col gap-2 sm:flex-row sm:justify-center">
          <Button asChild variant="default">
            <Link to={safeHome} replace>
              Back to your workspace
            </Link>
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              void logout().then(() =>
                navigate(resolveLoginPathForRoles(user?.roles), {
                  replace: true,
                }),
              );
            }}
          >
            Sign out
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
