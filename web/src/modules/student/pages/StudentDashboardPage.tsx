import { BookOpen, CalendarDays, GraduationCap } from "lucide-react";
import { Link } from "react-router-dom";

import { useAuth } from "@/modules/auth";
import { usePermissions } from "@/modules/auth/hooks/usePermissions";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";

export function StudentDashboardPage() {
  const { user } = useAuth();
  const { can } = usePermissions();
  const base = usePortalBase();

  return (
    <div className="space-y-8">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight md:text-3xl">
          Student dashboard
        </h1>
        <p className="max-w-2xl text-sm text-muted-foreground md:text-base">
          Welcome back{user?.displayName ? `, ${user.displayName}` : ""}. View
          your academic progress and schedule.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <Card className="border-border/80">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base">
              <BookOpen className="h-4 w-4 text-primary" aria-hidden />
              Grades
            </CardTitle>
            <CardDescription>
              {can("student:view_own_grades")
                ? "Course results and instructor feedback."
                : "Grades access is unavailable for this account."}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="secondary" asChild>
              <Link to={`${base}/grades`}>View my grades</Link>
            </Button>
          </CardContent>
        </Card>

        <Card className="border-border/80">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base">
              <CalendarDays className="h-4 w-4 text-primary" aria-hidden />
              Schedule
            </CardTitle>
            <CardDescription>
              Your lessons and room assignments for this term.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="secondary" asChild>
              <Link to={`${base}/schedule`}>View schedule</Link>
            </Button>
          </CardContent>
        </Card>

        <Card className="border-border/80 md:col-span-2 xl:col-span-1">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base">
              <GraduationCap className="h-4 w-4 text-primary" aria-hidden />
              Tips
            </CardTitle>
            <CardDescription>
              Need help? Contact your advisor or administrator through official
              channels.
            </CardDescription>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            Course materials and announcements will appear here when your
            school enables integrations.
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
