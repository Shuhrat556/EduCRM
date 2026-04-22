import { ArrowLeft, Loader2, RotateCcw, Save, UserCheck } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";

import { ATTENDANCE_STATUS_OPTIONS } from "@/features/attendance/lib/attendance-status";
import { isoDateToday } from "@/features/attendance/lib/iso-date";
import {
  useAttendanceSessionQuery,
  useSaveAttendanceBatchMutation,
} from "@/features/attendance/hooks/use-attendance";
import type {
  AttendanceBatchRowPayload,
  AttendanceStatus,
  AttendanceStudentRow,
} from "@/features/attendance/model/types";
import { weekdayLabel } from "@/features/schedule/lib/weekdays";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
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
import { Skeleton } from "@/shared/ui/components/skeleton";
import { cn } from "@/shared/lib/cn";

type DraftRow = {
  studentId: string;
  fullName: string;
  status: AttendanceStatus;
  comment: string;
  lessonGradeInput: string;
  weeklyRatingInput: string;
};

function toDraft(r: AttendanceStudentRow): DraftRow {
  return {
    studentId: r.studentId,
    fullName: r.studentName,
    status: r.status,
    comment: r.comment ?? "",
    lessonGradeInput:
      r.lessonGrade === null || r.lessonGrade === undefined
        ? ""
        : String(r.lessonGrade),
    weeklyRatingInput:
      r.weeklyRating === null || r.weeklyRating === undefined
        ? ""
        : String(r.weeklyRating),
  };
}

function draftsKey(rows: DraftRow[]) {
  return JSON.stringify(rows);
}

function draftsToPayload(
  rows: DraftRow[],
):
  | { ok: true; rows: AttendanceBatchRowPayload[] }
  | { ok: false; message: string } {
  const out: AttendanceBatchRowPayload[] = [];
  for (const d of rows) {
    const g = d.lessonGradeInput.trim();
    let lessonGrade: number | null = null;
    if (g !== "") {
      const n = Number(g);
      if (!Number.isFinite(n) || n < 0 || n > 100) {
        return {
          ok: false,
          message: `Grade for ${d.fullName} must be between 0 and 100.`,
        };
      }
      lessonGrade = Math.round(n);
    }
    const w = d.weeklyRatingInput.trim();
    let weeklyRating: number | null = null;
    if (w !== "") {
      const n = Number(w);
      if (!Number.isInteger(n) || n < 1 || n > 5) {
        return {
          ok: false,
          message: `Weekly rating for ${d.fullName} must be 1–5 or empty.`,
        };
      }
      weeklyRating = n;
    }
    out.push({
      studentId: d.studentId,
      status: d.status,
      comment: d.comment.trim() || null,
      lessonGrade,
      weeklyRating,
    });
  }
  return { ok: true, rows: out };
}

export function LessonAttendancePage() {
  const base = usePortalBase();
  const { lessonId: lessonIdParam } = useParams<{ lessonId: string }>();
  const lessonId = lessonIdParam ?? null;
  const [searchParams, setSearchParams] = useSearchParams();
  const sessionDate =
    searchParams.get("date")?.trim() || isoDateToday();

  const { data, isLoading, isError, isFetching } = useAttendanceSessionQuery(
    lessonId,
    sessionDate,
  );
  const saveMut = useSaveAttendanceBatchMutation();

  const [rows, setRows] = useState<DraftRow[]>([]);
  const [baseline, setBaseline] = useState("");
  const [clientError, setClientError] = useState<string | null>(null);

  useEffect(() => {
    if (data?.rows) {
      const next = data.rows.map(toDraft);
      setRows(next);
      setBaseline(draftsKey(next));
      setClientError(null);
    }
  }, [data]);

  const dirty = useMemo(
    () => (baseline ? draftsKey(rows) !== baseline : false),
    [rows, baseline],
  );

  function updateRow(studentId: string, patch: Partial<DraftRow>) {
    setRows((prev) =>
      prev.map((r) => (r.studentId === studentId ? { ...r, ...patch } : r)),
    );
  }

  function markAllPresent() {
    setRows((prev) =>
      prev.map((r) => ({ ...r, status: "present" as const })),
    );
  }

  function resetDrafts() {
    if (!data?.rows) return;
    const next = data.rows.map(toDraft);
    setRows(next);
    setBaseline(draftsKey(next));
    setClientError(null);
  }

  function setDateParam(d: string) {
    setSearchParams({ date: d }, { replace: true });
  }

  async function handleSave() {
    if (!lessonId) return;
    setClientError(null);
    const parsed = draftsToPayload(rows);
    if (!parsed.ok) {
      setClientError(parsed.message);
      return;
    }
    try {
      await saveMut.mutateAsync({
        lessonId,
        sessionDate,
        rows: parsed.rows,
      });
    } catch (e) {
      setClientError(
        e instanceof Error ? e.message : "Could not save attendance.",
      );
    }
  }

  if (!lessonId) {
    return (
      <p className="text-sm text-muted-foreground">Missing lesson in URL.</p>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <Button variant="ghost" size="sm" asChild className="gap-2 px-2">
          <Link to={`${base}/attendance`}>
            <ArrowLeft className="h-4 w-4" />
            Attendance
          </Link>
        </Button>
      </div>

      {isError ? (
        <p className="text-sm text-destructive">
          Could not load this lesson or roster. Check the link or schedule.
        </p>
      ) : null}

      {isLoading || !data ? (
        <div className="space-y-4">
          <Skeleton className="h-10 w-full max-w-md" />
          <Skeleton className="h-64 w-full" />
        </div>
      ) : (
        <>
          <Card className="border-border/80">
            <CardHeader className="pb-3">
              <CardTitle className="text-lg">
                {data.lesson.title?.trim() || "Lesson"}
              </CardTitle>
              <CardDescription className="flex flex-wrap gap-x-4 gap-y-1">
                <span>{weekdayLabel(data.lesson.weekday)}</span>
                <span className="tabular-nums">
                  {data.lesson.startTime}–{data.lesson.endTime}
                </span>
                <span>{data.lesson.group.name}</span>
                <span className="text-muted-foreground">
                  {data.lesson.teacher.fullName}
                </span>
                <span className="text-muted-foreground">
                  {data.lesson.room.name}
                </span>
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-end">
              <div className="space-y-2">
                <Label htmlFor="roster-date">Session date</Label>
                <Input
                  id="roster-date"
                  type="date"
                  value={sessionDate}
                  onChange={(e) => setDateParam(e.target.value)}
                  className="w-auto tabular-nums"
                />
              </div>
              <div className="flex flex-wrap gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="gap-2"
                  onClick={markAllPresent}
                  disabled={rows.length === 0}
                >
                  <UserCheck className="h-4 w-4" />
                  Mark all present
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="gap-2"
                  disabled={!dirty}
                  onClick={resetDrafts}
                >
                  <RotateCcw className="h-4 w-4" />
                  Reset
                </Button>
                <Button
                  type="button"
                  size="sm"
                  className="gap-2"
                  disabled={!dirty || saveMut.isPending || rows.length === 0}
                  onClick={() => void handleSave()}
                >
                  {saveMut.isPending ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin" />
                      Saving…
                    </>
                  ) : (
                    <>
                      <Save className="h-4 w-4" />
                      Save all
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>

          {clientError ? (
            <p className="text-sm font-medium text-destructive">{clientError}</p>
          ) : null}

          <div
            className={cn(
              "rounded-xl border border-border/80 bg-card",
              isFetching && "opacity-90 transition-opacity",
            )}
          >
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="hover:bg-muted/50">
                    <TableHead className="sticky left-0 z-20 min-w-[140px] bg-card pl-4 font-semibold shadow-[1px_0_0_hsl(var(--border))]">
                      Student
                    </TableHead>
                    <TableHead className="min-w-[120px] whitespace-nowrap font-semibold">
                      Attendance
                    </TableHead>
                    <TableHead className="min-w-[200px] font-semibold">
                      Comment
                    </TableHead>
                    <TableHead className="min-w-[88px] font-semibold">
                      Grade
                    </TableHead>
                    <TableHead className="min-w-[130px] pr-4 font-semibold">
                      Weekly rating
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="py-12 text-center text-sm text-muted-foreground"
                      >
                        No students in this group. Add enrollments in Students
                        or Groups first.
                      </TableCell>
                    </TableRow>
                  ) : (
                    rows.map((row) => (
                      <TableRow key={row.studentId} className="group">
                        <TableCell className="sticky left-0 z-10 bg-card align-middle font-medium shadow-[1px_0_0_hsl(var(--border))] group-hover:bg-muted/40">
                          {row.fullName}
                        </TableCell>
                        <TableCell className="align-middle">
                          <Select
                            value={row.status}
                            onValueChange={(v) =>
                              updateRow(row.studentId, {
                                status: v as AttendanceStatus,
                              })
                            }
                          >
                            <SelectTrigger className="h-9 w-[124px]">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              {ATTENDANCE_STATUS_OPTIONS.map((o) => (
                                <SelectItem key={o.value} value={o.value}>
                                  {o.label}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </TableCell>
                        <TableCell className="align-middle">
                          <Input
                            className="h-9 min-w-[180px] text-sm"
                            placeholder="Note"
                            value={row.comment}
                            onChange={(e) =>
                              updateRow(row.studentId, {
                                comment: e.target.value,
                              })
                            }
                          />
                        </TableCell>
                        <TableCell className="align-middle">
                          <Input
                            type="number"
                            min={0}
                            max={100}
                            className="h-9 w-20 tabular-nums"
                            placeholder="—"
                            value={row.lessonGradeInput}
                            onChange={(e) =>
                              updateRow(row.studentId, {
                                lessonGradeInput: e.target.value,
                              })
                            }
                          />
                        </TableCell>
                        <TableCell className="align-middle pr-4">
                          <Select
                            value={
                              row.weeklyRatingInput === ""
                                ? "__none__"
                                : row.weeklyRatingInput
                            }
                            onValueChange={(v) =>
                              updateRow(row.studentId, {
                                weeklyRatingInput:
                                  v === "__none__" ? "" : v,
                              })
                            }
                          >
                            <SelectTrigger className="h-9 w-[120px]">
                              <SelectValue placeholder="Optional" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="__none__">—</SelectItem>
                              {[1, 2, 3, 4, 5].map((n) => (
                                <SelectItem key={n} value={String(n)}>
                                  {n} / 5
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </div>

          <p className="text-xs text-muted-foreground">
            {dirty
              ? "You have unsaved changes — click Save all to store this session."
              : rows.length > 0
                ? "All changes saved for this session date."
                : null}
          </p>
        </>
      )}
    </div>
  );
}
