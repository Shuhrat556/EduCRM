import {
  BarChart3,
  CalendarRange,
  ClipboardCheck,
  FileSpreadsheet,
  FileText,
  GraduationCap,
  Users,
  Wallet,
} from "lucide-react";
import { useMemo, useState } from "react";

import { useGroupsListQuery } from "@/features/groups/hooks/use-groups";
import { studentStatusLabel } from "@/features/students/lib/student-status";
import type { StudentStatus } from "@/features/students/model/types";
import { useReportsSummary } from "@/features/reports/hooks/use-reports-summary";
import type { ReportsFilters } from "@/features/reports/model/types";
import { useTeachersListQuery } from "@/features/teachers/hooks/use-teachers";
import { teacherStatusLabel } from "@/features/teachers/lib/teacher-status";
import type { TeacherStatus } from "@/features/teachers/model/types";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/components/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/components/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/shared/ui/components/tooltip";
import { formatMoneyUnits } from "@/shared/lib/format/money-units";
import { formatPercentRate } from "@/shared/lib/format/percent";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";
import { SummaryMetricTile } from "@/shared/ui/data/summary-metric-tile";
import { PageHeader } from "@/shared/ui/layout/page-header";

const STUDENT_STATUS_OPTIONS: { value: StudentStatus | "all"; label: string }[] =
  [
    { value: "all", label: "All statuses" },
    { value: "active", label: studentStatusLabel("active") },
    { value: "inactive", label: studentStatusLabel("inactive") },
    { value: "graduated", label: studentStatusLabel("graduated") },
    { value: "suspended", label: studentStatusLabel("suspended") },
  ];

function defaultDateRange() {
  const to = new Date();
  const from = new Date(to);
  from.setDate(from.getDate() - 30);
  const ymd = (d: Date) => d.toISOString().slice(0, 10);
  return { from: ymd(from), to: ymd(to) };
}

export function ReportsPage() {
  const { from: defaultFrom, to: defaultTo } = useMemo(
    () => defaultDateRange(),
    [],
  );
  const [dateFrom, setDateFrom] = useState(defaultFrom);
  const [dateTo, setDateTo] = useState(defaultTo);
  const [groupId, setGroupId] = useState<string>("all");
  const [teacherId, setTeacherId] = useState<string>("all");
  const [studentStatus, setStudentStatus] = useState<StudentStatus | "all">(
    "all",
  );

  const filters: ReportsFilters = useMemo(
    () => ({
      dateFrom,
      dateTo,
      groupId: groupId === "all" ? null : groupId,
      teacherId: teacherId === "all" ? null : teacherId,
      studentStatus,
    }),
    [dateFrom, dateTo, groupId, teacherId, studentStatus],
  );

  const { data: groupsData } = useGroupsListQuery({
    page: 1,
    pageSize: 200,
    search: "",
    status: "all",
  });
  const { data: teachersData } = useTeachersListQuery({
    page: 1,
    pageSize: 200,
    search: "",
    status: "all",
  });

  const { data, isPending, isError, error, refetch, isFetching } =
    useReportsSummary(filters);

  const billingRows = Object.entries(data?.payments.byBillingStatus ?? {}).filter(
    ([, n]) => n > 0,
  );

  const billingLabels: Record<string, string> = {
    paid: "Paid",
    partial: "Partial",
    unpaid: "Unpaid",
    overdue: "Overdue",
    waived: "Waived",
    no_fee: "No fee",
  };

  return (
    <TooltipProvider delayDuration={200}>
      <div className="space-y-8">
        <PageHeader
          title="Reports"
          description="Operational summaries for students, attendance, tuition, and staff. Adjust filters to narrow the cohort and reporting period."
          meta={
            <>
              {data?.source === "aggregated" ? (
                <Badge variant="secondary" className="font-normal">
                  Aggregated
                </Badge>
              ) : data?.source === "api" ? (
                <Badge variant="outline" className="font-normal">
                  API
                </Badge>
              ) : null}
              {isFetching && !isPending ? (
                <Badge variant="outline" className="font-normal">
                  Refreshing…
                </Badge>
              ) : null}
            </>
          }
          actions={
            <div className="flex w-full flex-wrap gap-2 lg:justify-end">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="outline" size="sm" className="gap-1.5" disabled>
                    <FileText className="h-4 w-4 opacity-80" aria-hidden />
                    Export PDF
                  </Button>
                </TooltipTrigger>
                <TooltipContent>PDF export — coming soon</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="outline" size="sm" className="gap-1.5" disabled>
                    <FileSpreadsheet className="h-4 w-4 opacity-80" aria-hidden />
                    Export Excel
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Excel export — coming soon</TooltipContent>
              </Tooltip>
            </div>
          }
        />

        <Card className="border-border/80">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <CalendarRange className="h-4 w-4 text-primary" aria-hidden />
              Filters
            </CardTitle>
            <CardDescription>
              Date range drives attendance sessions and billing periods; roster
              uses group, teacher context, and student status.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-6">
              <div className="space-y-2">
                <Label htmlFor="reports-from">From</Label>
                <Input
                  id="reports-from"
                  type="date"
                  value={dateFrom}
                  onChange={(e) => setDateFrom(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="reports-to">To</Label>
                <Input
                  id="reports-to"
                  type="date"
                  value={dateTo}
                  onChange={(e) => setDateTo(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Group</Label>
                <Select value={groupId} onValueChange={setGroupId}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="All groups" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All groups</SelectItem>
                    {(groupsData?.items ?? []).map((g) => (
                      <SelectItem key={g.id} value={g.id}>
                        {g.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Teacher</Label>
                <Select value={teacherId} onValueChange={setTeacherId}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="All teachers" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All teachers</SelectItem>
                    {(teachersData?.items ?? []).map((t) => (
                      <SelectItem key={t.id} value={t.id}>
                        {t.fullName}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2 md:col-span-2 xl:col-span-2">
                <Label>Status</Label>
                <Select
                  value={studentStatus}
                  onValueChange={(v) =>
                    setStudentStatus(v as StudentStatus | "all")
                  }
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Student status" />
                  </SelectTrigger>
                  <SelectContent>
                    {STUDENT_STATUS_OPTIONS.map((o) => (
                      <SelectItem key={o.value} value={o.value}>
                        {o.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <p className="text-xs text-muted-foreground">
                  Enrollment status for student and tuition cohorts.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {isError ? (
          <QueryErrorAlert
            error={error}
            title="Could not load reports"
            onRetry={() => void refetch()}
          />
        ) : null}

        {!isError ? (
        <div className="grid gap-6 lg:grid-cols-2">
          {/* Students */}
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <GraduationCap className="h-4 w-4 text-primary" aria-hidden />
                Student summary
              </CardTitle>
              <CardDescription>
                Roster matching filters; new enrollments use created date in
                range.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-3 sm:grid-cols-3">
                <SummaryMetricTile
                  label="In roster"
                  value={
                    isPending
                      ? "…"
                      : (data?.students.totalInRoster ?? 0).toLocaleString()
                  }
                />
                <SummaryMetricTile
                  label="New in range"
                  value={
                    isPending
                      ? "…"
                      : (
                          data?.students.newEnrollmentsInRange ?? 0
                        ).toLocaleString()
                  }
                  hint="By enrollment date"
                />
                <SummaryMetricTile
                  label="No group"
                  value={
                    isPending
                      ? "…"
                      : (data?.students.withoutGroup ?? 0).toLocaleString()
                  }
                />
              </div>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Students</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(
                    [
                      "active",
                      "inactive",
                      "graduated",
                      "suspended",
                    ] as StudentStatus[]
                  ).map((s) => (
                    <TableRow key={s}>
                      <TableCell>{studentStatusLabel(s)}</TableCell>
                      <TableCell className="text-right tabular-nums">
                        {isPending
                          ? "…"
                          : (
                              data?.students.byStatus[s] ?? 0
                            ).toLocaleString()}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {/* Attendance */}
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <ClipboardCheck className="h-4 w-4 text-primary" aria-hidden />
                Attendance summary
              </CardTitle>
              <CardDescription>
                Lesson sessions with marks in the date range; scoped to roster
                and schedule filters.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-3 sm:grid-cols-3">
                <SummaryMetricTile
                  label="Sessions"
                  value={
                    isPending
                      ? "…"
                      : (
                          data?.attendance.sessionsRecorded ?? 0
                        ).toLocaleString()
                  }
                />
                <SummaryMetricTile
                  label="Marks"
                  value={
                    isPending
                      ? "…"
                      : (data?.attendance.marksRecorded ?? 0).toLocaleString()
                  }
                />
                <SummaryMetricTile
                  label="Present rate"
                  value={
                    isPending ? "…" : formatPercentRate(data?.attendance.attendanceRate ?? null)
                  }
                  hint="Present ÷ all outcomes"
                />
              </div>
              <div className="grid grid-cols-3 gap-3">
                <SummaryMetricTile
                  label="Present"
                  value={
                    isPending
                      ? "…"
                      : (data?.attendance.present ?? 0).toLocaleString()
                  }
                  className="border-emerald-500/15 bg-emerald-500/5"
                />
                <SummaryMetricTile
                  label="Absent"
                  value={
                    isPending
                      ? "…"
                      : (data?.attendance.absent ?? 0).toLocaleString()
                  }
                  className="border-rose-500/15 bg-rose-500/5"
                />
                <SummaryMetricTile
                  label="Late"
                  value={
                    isPending ? "…" : (data?.attendance.late ?? 0).toLocaleString()
                  }
                  className="border-amber-500/15 bg-amber-500/5"
                />
              </div>
            </CardContent>
          </Card>

          {/* Payments */}
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <Wallet className="h-4 w-4 text-primary" aria-hidden />
                Payment summary
              </CardTitle>
              <CardDescription>
                Receivable rows per billing month overlapping the range (
                {isPending ? "…" : data?.payments.periodLabel ?? "—"}).
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                <SummaryMetricTile
                  label="Expected"
                  value={
                    isPending
                      ? "…"
                      : formatMoneyUnits(
                          data?.payments.totalExpected ?? 0,
                          data?.currencyCode ?? "USD",
                        )
                  }
                />
                <SummaryMetricTile
                  label="Collected"
                  value={
                    isPending
                      ? "…"
                      : formatMoneyUnits(
                          data?.payments.totalPaid ?? 0,
                          data?.currencyCode ?? "USD",
                        )
                  }
                  className="border-emerald-500/15 bg-emerald-500/5"
                />
                <SummaryMetricTile
                  label="Outstanding"
                  value={
                    isPending
                      ? "…"
                      : formatMoneyUnits(
                          data?.payments.totalOutstanding ?? 0,
                          data?.currencyCode ?? "USD",
                        )
                  }
                />
                <SummaryMetricTile
                  label="Overdue rows"
                  value={
                    isPending
                      ? "…"
                      : (data?.payments.overdueCount ?? 0).toLocaleString()
                  }
                  hint={`${isPending ? "…" : data?.payments.monthsCounted ?? 0} month(s)`}
                />
              </div>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Billing status</TableHead>
                    <TableHead className="text-right">Rows</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {billingRows.map(([k, n]) => (
                    <TableRow key={k}>
                      <TableCell>{billingLabels[k] ?? k}</TableCell>
                      <TableCell className="text-right tabular-nums">
                        {isPending ? "…" : n.toLocaleString()}
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isPending && billingRows.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={2}
                        className="text-center text-muted-foreground"
                      >
                        No receivable rows in this scope / period
                      </TableCell>
                    </TableRow>
                  ) : null}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {/* Teachers */}
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base font-semibold">
                <Users className="h-4 w-4 text-primary" aria-hidden />
                Teacher summary
              </CardTitle>
              <CardDescription>
                Directory slice and schedule footprint (weekly template).
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-3 sm:grid-cols-3">
                <SummaryMetricTile
                  label="Teachers"
                  value={
                    isPending
                      ? "…"
                      : (data?.teachers.total ?? 0).toLocaleString()
                  }
                />
                <SummaryMetricTile
                  label="Lessons"
                  value={
                    isPending
                      ? "…"
                      : (
                          data?.teachers.lessonsOnSchedule ?? 0
                        ).toLocaleString()
                  }
                  hint="After filters"
                />
                <SummaryMetricTile
                  label="Groups on schedule"
                  value={
                    isPending
                      ? "…"
                      : (
                          data?.teachers.groupsRepresented ?? 0
                        ).toLocaleString()
                  }
                />
              </div>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Teachers</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {Object.entries(data?.teachers.byStatus ?? {}).map(
                    ([statusKey, n]) => (
                      <TableRow key={statusKey}>
                        <TableCell>
                          {teacherStatusLabel(statusKey as TeacherStatus)}
                        </TableCell>
                        <TableCell className="text-right tabular-nums">
                          {isPending ? "…" : n.toLocaleString()}
                        </TableCell>
                      </TableRow>
                    ),
                  )}
                  {!isPending &&
                  Object.keys(data?.teachers.byStatus ?? {}).length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={2}
                        className="text-center text-muted-foreground"
                      >
                        No teachers in this slice
                      </TableCell>
                    </TableRow>
                  ) : null}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
        ) : null}

        <Card className="border-dashed border-border/80 bg-muted/10">
          <CardContent className="flex flex-col gap-2 py-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <BarChart3 className="h-4 w-4 shrink-0 opacity-80" aria-hidden />
              <span>
                Need a saved layout or scheduled exports? PDF and Excel actions
                above are placeholders until the export service is connected.
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    </TooltipProvider>
  );
}
