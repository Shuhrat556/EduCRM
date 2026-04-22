import { Bell, Plug, UserRound, Wrench } from "lucide-react";
import { useEffect, useState } from "react";

import { formatRoleLabel, useAuth } from "@/modules/auth";
import { Badge } from "@/shared/ui/components/badge";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { Input } from "@/shared/ui/components/input";
import { Label } from "@/shared/ui/components/label";
import { Separator } from "@/shared/ui/components/separator";
import { cn } from "@/shared/lib/cn";
import { PageHeader } from "@/shared/ui/layout/page-header";

export function SettingsPage() {
  const { user, hasRole } = useAuth();
  const isSuperAdmin = hasRole("super_admin");

  const [displayName, setDisplayName] = useState(user?.displayName ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [profileNote, setProfileNote] = useState("");

  useEffect(() => {
    setDisplayName(user?.displayName ?? "");
    setPhone(user?.phone ?? "");
  }, [user?.displayName, user?.phone]);

  return (
    <div className="space-y-8">
      <PageHeader
        title="Settings"
        description="Profile details for your account. System configuration is available to Super Admins as a roadmap-only preview."
      />

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-6 lg:col-span-2">
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <UserRound className="h-4 w-4 text-primary" aria-hidden />
                Profile
              </CardTitle>
              <CardDescription>
                How you appear in EduCRM. Email is managed by your administrator
                or identity provider.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2 sm:col-span-2">
                  <Label htmlFor="settings-display-name">Display name</Label>
                  <Input
                    id="settings-display-name"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                    autoComplete="name"
                  />
                </div>
                <div className="space-y-2 sm:col-span-2">
                  <Label htmlFor="settings-email">Email</Label>
                  <Input
                    id="settings-email"
                    value={user?.email ?? ""}
                    disabled
                    readOnly
                    className="bg-muted/40"
                  />
                </div>
                <div className="space-y-2 sm:col-span-2">
                  <Label htmlFor="settings-phone">Phone</Label>
                  <Input
                    id="settings-phone"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    autoComplete="tel"
                    placeholder="Optional"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="settings-bio">Internal note</Label>
                <textarea
                  id="settings-bio"
                  value={profileNote}
                  onChange={(e) => setProfileNote(e.target.value)}
                  placeholder="Visible only on profile once directory sync is enabled."
                  rows={3}
                  className={cn(
                    "flex min-h-[80px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
                    "resize-none",
                  )}
                />
              </div>
              <Separator />
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-xs text-muted-foreground">
                  Saving will call your API when{" "}
                  <code className="rounded bg-muted px-1 py-0.5 text-[11px]">
                    PATCH /auth/me
                  </code>{" "}
                  (or equivalent) is implemented.
                </p>
                <Button type="button" variant="secondary" disabled>
                  Save profile
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <Bell className="h-4 w-4 text-primary" aria-hidden />
                Notifications
              </CardTitle>
              <CardDescription>
                Placeholder for email and in-app notification preferences.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="rounded-lg border border-dashed border-border/80 bg-muted/15 px-4 py-8 text-center text-sm text-muted-foreground">
                Notification settings — coming soon
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="space-y-6">
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Signed-in as
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <p className="font-medium text-foreground">
                {user?.displayName ?? "—"}
              </p>
              <p className="text-xs text-muted-foreground">{user?.email}</p>
              <div className="flex flex-wrap gap-1 pt-1">
                {user?.roles?.map((r) => (
                  <Badge key={r} variant="secondary" className="text-[10px]">
                    {formatRoleLabel(r)}
                  </Badge>
                ))}
              </div>
            </CardContent>
          </Card>

          {isSuperAdmin ? (
            <Card className="border-border/80">
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center gap-2 text-base font-semibold">
                  <Wrench className="h-4 w-4 text-primary" aria-hidden />
                  System
                  <Badge variant="outline" className="ml-1 font-normal">
                    Super Admin
                  </Badge>
                </CardTitle>
                <CardDescription>
                  Institution-wide controls — wire to your admin API when ready.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="rounded-lg border border-border/60 bg-muted/20 p-3">
                  <div className="flex items-center gap-2 text-sm font-medium">
                    <Plug className="h-4 w-4 text-muted-foreground" aria-hidden />
                    Integrations
                  </div>
                  <p className="mt-1 text-xs text-muted-foreground">
                    SIS, accounting, SMS — placeholders only.
                  </p>
                </div>
                <div className="rounded-lg border border-border/60 bg-muted/20 p-3">
                  <div className="text-sm font-medium">Branding & domain</div>
                  <p className="mt-1 text-xs text-muted-foreground">
                    Logo, support URL, custom domain — placeholders only.
                  </p>
                </div>
                <div className="rounded-lg border border-border/60 bg-muted/20 p-3">
                  <div className="text-sm font-medium">Data retention</div>
                  <p className="mt-1 text-xs text-muted-foreground">
                    Export schedules and deletion policies — placeholders only.
                  </p>
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card className="border-dashed border-border/80 bg-muted/10">
              <CardContent className="py-4 text-sm text-muted-foreground">
                System settings are restricted to{" "}
                <span className="font-medium text-foreground">Super Admin</span>{" "}
                accounts.
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
