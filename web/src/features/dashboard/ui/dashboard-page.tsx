import {
  ArrowRight,
  Calendar,
  ClipboardList,
  CreditCard,
  GraduationCap,
  LayoutGrid,
  TrendingUp,
  UserRound,
  Users,
} from "lucide-react";
import { Link } from "react-router-dom";

import { useAuth } from "@/modules/auth";
import { useDashboardOverview } from "@/features/dashboard/hooks/use-dashboard-overview";
import {
  formatCurrency,
  formatPercent,
  formatRelativeTime,
  formatShortDate,
  formatTime,
} from "@/features/dashboard/lib/format";
import type { PaymentStatus } from "@/features/dashboard/model/types";
import { StatCard } from "@/features/dashboard/ui/stat-card";
import { SimpleBarChart } from "@/features/dashboard/ui/simple-bar-chart";
import { adminPortalNav } from "@/modules/admin/admin-nav";
import { superAdminNav } from "@/modules/super-admin/super-admin-nav";
import type { PortalNavItem } from "@/modules/shared/layout/nav-types";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { Badge } from "@/shared/ui/components/badge";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { ScrollArea } from "@/shared/ui/components/scroll-area";
import { Separator } from "@/shared/ui/components/separator";
import { Skeleton } from "@/shared/ui/components/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/components/table";
import { cn } from "@/shared/lib/cn";

function paymentBadgeVariant(
  status: PaymentStatus,
): "default" | "secondary" | "destructive" | "outline" {
  switch (status) {
    case "paid":
      return "secondary";
    case "pending":
      return "outline";
    case "failed":
      return "destructive";
    case "overdue":
      return "destructive";
    default:
      return "outline";
  }
}

function paymentStatusLabel(status: PaymentStatus) {
  return status.charAt(0).toUpperCase() + status.slice(1);
}

const QUICK_ACTION_BLURBS: Record<string, string> = {
  Students: "Manage enrollment and student records.",
  Teachers: "Staff profiles and assignments.",
  Admins: "Institution administrator accounts.",
  Groups: "Cohorts, classes, and rosters.",
  Subjects: "Curriculum and course catalog.",
  Rooms: "Facilities and capacity.",
  Schedule: "Timetable and room bookings.",
  Attendance: "Roll call and daily status.",
  Payments: "Tuition and invoices.",
  Reports: "Exports and operational insights.",
  Settings: "Workspace preferences and access.",
};

export function DashboardPage() {
  const { user } = useAuth();
  const base = usePortalBase();
  const quickSource =
    base === "/super-admin" ? superAdminNav : adminPortalNav;
  const quickActionLinks: PortalNavItem[] = quickSource
    .filter((item: PortalNavItem) => !item.to.endsWith("/dashboard"))
    .slice(0, 8);
  const { data, isPending } = useDashboardOverview();

  const stats = data?.stats;
  const loading = isPending;

  const hasLessons = (data?.todayLessons.length ?? 0) > 0;
  const hasPayments = (data?.recentPayments.length ?? 0) > 0;
  const hasAttendance = (data?.attendanceSummary.length ?? 0) > 0;
  const hasChartData = (data?.weeklyActivity.length ?? 0) > 0;

  const isEmptySnapshot =
    !loading &&
    stats &&
    stats.totalStudents === 0 &&
    stats.totalTeachers === 0 &&
    stats.activeGroups === 0 &&
    stats.overduePayments === 0 &&
    !hasLessons &&
    !hasPayments &&
    !hasAttendance &&
    !hasChartData;

  return (
    <div className="space-y-8">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold tracking-tight text-foreground md:text-3xl">
            Dashboard
          </h1>
          <p className="max-w-2xl text-sm text-muted-foreground md:text-base">
            Welcome back
            {user?.displayName ? `, ${user.displayName}` : ""}. Here is a
            snapshot of your institution today.
          </p>
        </div>
        {import.meta.env.DEV && import.meta.env.VITE_DASHBOARD_DEMO === "true" ? (
          <Badge variant="outline" className="w-fit shrink-0">
            Demo data (VITE_DASHBOARD_DEMO)
          </Badge>
        ) : null}
      </div>

      {isEmptySnapshot ? (
        <Card className="border-dashed bg-muted/15">
          <CardHeader>
            <CardTitle className="text-base">Connect your dashboard API</CardTitle>
            <CardDescription>
              No overview data yet. Implement{" "}
              <code className="rounded bg-muted px-1 py-0.5 text-xs">
                GET /dashboard/overview
              </code>{" "}
              or set{" "}
              <code className="rounded bg-muted px-1 py-0.5 text-xs">
                VITE_DASHBOARD_DEMO=true
              </code>{" "}
              for sample data.
            </CardDescription>
          </CardHeader>
        </Card>
      ) : null}

      {/* Stats */}
      <section aria-labelledby="stats-heading">
        <h2 id="stats-heading" className="sr-only">
          Key statistics
        </h2>
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <StatCard
            title="Total students"
            value={(stats?.totalStudents ?? 0).toLocaleString()}
            icon={GraduationCap}
            loading={loading}
            subtitle="Enrolled learners"
          />
          <StatCard
            title="Total teachers"
            value={(stats?.totalTeachers ?? 0).toLocaleString()}
            icon={UserRound}
            loading={loading}
            subtitle="Active teaching staff"
          />
          <StatCard
            title="Active groups"
            value={(stats?.activeGroups ?? 0).toLocaleString()}
            icon={Users}
            loading={loading}
            subtitle="Classes & cohorts"
          />
          <StatCard
            title="Overdue payments"
            value={(stats?.overduePayments ?? 0).toLocaleString()}
            icon={CreditCard}
            loading={loading}
            variant="warning"
            subtitle="Invoices requiring follow-up"
          />
        </div>
      </section>

      {/* Today lessons + weekly chart */}
      <div className="grid gap-6 lg:grid-cols-5">
        <Card className="border-border/80 lg:col-span-3">
          <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-2 space-y-0 pb-4">
            <div>
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <Calendar className="h-4 w-4 text-primary" aria-hidden />
                Today&apos;s lessons
              </CardTitle>
              <CardDescription>
                Scheduled sessions for{" "}
                {new Date().toLocaleDateString(undefined, {
                  weekday: "long",
                  month: "long",
                  day: "numeric",
                })}
              </CardDescription>
            </div>
            {base === "/teacher" ? (
              <Button variant="outline" size="sm" asChild>
                <Link to={`${base}/schedule`}>
                  Full schedule
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            ) : null}
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? (
              <div className="space-y-3">
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            ) : hasLessons ? (
              <ScrollArea className="max-h-[min(340px,55vh)] pr-3">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Time</TableHead>
                      <TableHead>Subject</TableHead>
                      <TableHead className="hidden sm:table-cell">Group</TableHead>
                      <TableHead className="hidden md:table-cell">Room</TableHead>
                      <TableHead className="hidden lg:table-cell">Teacher</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data!.todayLessons.map((row) => (
                      <TableRow key={row.id}>
                        <TableCell className="whitespace-nowrap font-medium tabular-nums">
                          {formatTime(row.startsAt)} – {formatTime(row.endsAt)}
                        </TableCell>
                        <TableCell>{row.subject}</TableCell>
                        <TableCell className="hidden sm:table-cell">
                          {row.groupName}
                        </TableCell>
                        <TableCell className="hidden md:table-cell">
                          {row.room}
                        </TableCell>
                        <TableCell className="hidden max-w-[180px] truncate lg:table-cell">
                          {row.teacherName}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </ScrollArea>
            ) : (
              <AppEmpty
                icon={Calendar}
                title="No lessons today"
                description="When your schedule API is linked, today’s classes will show here."
                className="border-0 bg-transparent py-10"
              />
            )}
          </CardContent>
        </Card>

        <Card className="border-border/80 lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <TrendingUp className="h-4 w-4 text-primary" aria-hidden />
              Weekly activity
            </CardTitle>
            <CardDescription>
              Lesson sessions & check-ins (sample metric — adjust when API is
              ready).
            </CardDescription>
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? (
              <div className="space-y-4 pt-2">
                {Array.from({ length: 7 }).map((_, i) => (
                  <Skeleton key={i} className="h-8 w-full" />
                ))}
              </div>
            ) : (
              <SimpleBarChart
                data={data?.weeklyActivity ?? []}
                caption="Weekly activity levels by day"
              />
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recent payments + attendance */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="border-border/80">
          <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-2 space-y-0 pb-4">
            <div>
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <CreditCard className="h-4 w-4 text-primary" aria-hidden />
                Recent payments
              </CardTitle>
              <CardDescription>Latest transactions across billing.</CardDescription>
            </div>
            {base === "/student" ? (
              <Button variant="outline" size="sm" asChild>
                <Link to={`${base}/grades`}>
                  Grades
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            ) : null}
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? (
              <div className="space-y-3">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            ) : hasPayments ? (
              <ScrollArea className="max-h-[min(320px,50vh)] pr-3">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>When</TableHead>
                      <TableHead>Payer</TableHead>
                      <TableHead>Amount</TableHead>
                      <TableHead>Status</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data!.recentPayments.map((row) => (
                      <TableRow key={row.id}>
                        <TableCell className="whitespace-nowrap text-muted-foreground">
                          {formatRelativeTime(row.paidAt)}
                        </TableCell>
                        <TableCell className="font-medium">
                          {row.payerName}
                        </TableCell>
                        <TableCell className="tabular-nums">
                          {formatCurrency(row.amountCents, row.currency)}
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={paymentBadgeVariant(row.status)}
                            className="font-normal capitalize"
                          >
                            {paymentStatusLabel(row.status)}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </ScrollArea>
            ) : (
              <AppEmpty
                icon={CreditCard}
                title="No recent payments"
                description="Payment activity will display once billing integration is live."
                className="border-0 bg-transparent py-10"
              />
            )}
          </CardContent>
        </Card>

        <Card className="border-border/80">
          <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-2 space-y-0 pb-4">
            <div>
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <ClipboardList className="h-4 w-4 text-primary" aria-hidden />
                Attendance summary
              </CardTitle>
              <CardDescription>
                Recent daily roll-up — chart area reserved for deeper analytics.
              </CardDescription>
            </div>
            {base === "/teacher" ? (
              <Button variant="outline" size="sm" asChild>
                <Link to={`${base}/attendance`}>
                  Attendance
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            ) : null}
          </CardHeader>
          <CardContent className="space-y-6 pt-0">
            {loading ? (
              <div className="space-y-3">
                <Skeleton className="h-24 w-full" />
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            ) : hasAttendance ? (
              <>
                <div className="rounded-lg border border-border/80 bg-muted/20 p-4">
                  <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                    Attendance rate trend
                  </p>
                  <SimpleBarChart
                    className="mt-3"
                    data={data!.attendanceSummary.slice(0, 7).map((row) => ({
                      label: formatShortDate(row.date),
                      value: Math.round(row.attendanceRate * 100),
                    }))}
                    caption="Daily attendance rate percent"
                    emptyHint="No attendance trend"
                  />
                </div>
                <Separator />
                <ScrollArea className="max-h-[min(220px,40vh)] pr-3">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Date</TableHead>
                        <TableHead className="text-right">Present</TableHead>
                        <TableHead className="hidden sm:table-cell text-right">
                          Absent
                        </TableHead>
                        <TableHead className="text-right">Rate</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {data!.attendanceSummary.map((row) => (
                        <TableRow key={row.id}>
                          <TableCell className="font-medium">
                            {formatShortDate(row.date)}
                          </TableCell>
                          <TableCell className="text-right tabular-nums">
                            {row.presentCount}
                          </TableCell>
                          <TableCell className="hidden sm:table-cell text-right tabular-nums text-muted-foreground">
                            {row.absentCount}
                          </TableCell>
                          <TableCell className="text-right">
                            <span
                              className={cn(
                                "font-medium tabular-nums",
                                row.attendanceRate >= 0.9
                                  ? "text-emerald-600 dark:text-emerald-400"
                                  : row.attendanceRate >= 0.8
                                    ? "text-amber-600 dark:text-amber-400"
                                    : "text-destructive",
                              )}
                            >
                              {formatPercent(row.attendanceRate)}
                            </span>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </ScrollArea>
              </>
            ) : (
              <AppEmpty
                icon={ClipboardList}
                title="No attendance data"
                description="Connect attendance reporting to see daily summaries here."
                className="border-0 bg-transparent py-10"
              />
            )}
          </CardContent>
        </Card>
      </div>

      {/* Quick actions */}
      <section aria-labelledby="quick-actions-heading">
        <div className="mb-4 flex items-center gap-2">
          <LayoutGrid className="h-4 w-4 text-muted-foreground" aria-hidden />
          <h2
            id="quick-actions-heading"
            className="text-lg font-semibold tracking-tight"
          >
            Quick actions
          </h2>
        </div>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {quickActionLinks.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className="group block rounded-xl outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <Card className="h-full border-border/80 transition-all duration-200 group-hover:border-primary/25 group-hover:shadow-md">
                <CardHeader className="pb-2">
                  <div className="flex items-start justify-between gap-2">
                    <span className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary transition-colors group-hover:bg-primary/15">
                      <item.icon className="h-[18px] w-[18px]" aria-hidden />
                    </span>
                    <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5 group-hover:text-primary" />
                  </div>
                  <CardTitle className="pt-2 text-base font-semibold leading-snug">
                    {item.label}
                  </CardTitle>
                </CardHeader>
                <CardContent className="pt-0">
                  <p className="text-xs leading-relaxed text-muted-foreground">
                    {QUICK_ACTION_BLURBS[item.label] ??
                      `Open ${item.label.toLowerCase()}.`}
                  </p>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
}
