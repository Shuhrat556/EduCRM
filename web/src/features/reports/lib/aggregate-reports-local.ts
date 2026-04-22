import { getAttendanceSessionsSnapshotSync } from "@/features/attendance/lib/mock-attendance-store";
import { paymentsApi } from "@/features/payments/api/payments-api";
import type { BillingStatus } from "@/features/payments/model/types";
import { scheduleApi } from "@/features/schedule/api/schedule-api";
import { studentsApi } from "@/features/students/api/students-api";
import type { StudentStatus } from "@/features/students/model/types";
import { teachersApi } from "@/features/teachers/api/teachers-api";
import type { ReportsFilters, ReportsSummary } from "@/features/reports/model/types";

function isoDateOnly(iso: string): string {
  return iso.slice(0, 10);
}

function monthsBetweenInclusive(fromYmd: string, toYmd: string): string[] {
  const start = new Date(`${fromYmd}T00:00:00`);
  const end = new Date(`${toYmd}T00:00:00`);
  if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime())) return [];
  const out: string[] = [];
  const cur = new Date(start.getFullYear(), start.getMonth(), 1);
  const endMonth = new Date(end.getFullYear(), end.getMonth(), 1);
  while (cur <= endMonth) {
    const y = cur.getFullYear();
    const m = String(cur.getMonth() + 1).padStart(2, "0");
    out.push(`${y}-${m}`);
    cur.setMonth(cur.getMonth() + 1);
  }
  return out;
}

function emptyByStudentStatus(): Record<StudentStatus, number> {
  return {
    active: 0,
    inactive: 0,
    graduated: 0,
    suspended: 0,
  };
}

function emptyByBillingStatus(): Record<BillingStatus, number> {
  return {
    paid: 0,
    partial: 0,
    unpaid: 0,
    overdue: 0,
    waived: 0,
    no_fee: 0,
  };
}

export async function aggregateReportsLocal(
  filters: ReportsFilters,
): Promise<ReportsSummary> {
  const { dateFrom, dateTo, groupId, teacherId, studentStatus } = filters;

  const studentsRes = await studentsApi.list({
    page: 1,
    pageSize: 2_000,
    search: "",
    status: studentStatus,
    groupId: groupId ?? undefined,
  });

  const roster = studentsRes.items;
  const rosterIds = new Set(roster.map((s) => s.id));

  const byStudentStatus = emptyByStudentStatus();
  let withoutGroup = 0;
  let newEnrollmentsInRange = 0;

  for (const s of roster) {
    byStudentStatus[s.status] += 1;
    if (!s.groupId) withoutGroup += 1;
    const created = isoDateOnly(s.createdAt);
    if (created >= dateFrom && created <= dateTo) newEnrollmentsInRange += 1;
  }

  const teachersRes = await teachersApi.list({
    page: 1,
    pageSize: 500,
    search: "",
    status: "all",
  });

  const teacherItems = teacherId
    ? teachersRes.items.filter((t) => t.id === teacherId)
    : teachersRes.items;

  const byTeacherStatus: Record<string, number> = {};
  for (const t of teacherItems) {
    byTeacherStatus[t.status] = (byTeacherStatus[t.status] ?? 0) + 1;
  }

  let lessonsOnSchedule = 0;
  const groupIdsFromLessons = new Set<string>();
  const lessonMeta = new Map<
    string,
    { groupId: string; teacherId: string }
  >();
  try {
    const lessons = await scheduleApi.listWeek();
    for (const l of lessons) {
      lessonMeta.set(l.id, {
        groupId: l.group.id,
        teacherId: l.teacher.id,
      });
      if (groupId && l.group.id !== groupId) continue;
      if (teacherId && l.teacher.id !== teacherId) continue;
      lessonsOnSchedule += 1;
      groupIdsFromLessons.add(l.group.id);
    }
  } catch {
    /* schedule API optional in some setups */
  }

  const sessions = getAttendanceSessionsSnapshotSync();
  let sessionsRecorded = 0;
  let marksRecorded = 0;
  let present = 0;
  let absent = 0;
  let late = 0;

  for (const [sessionKey, entries] of Object.entries(sessions)) {
    const parts = sessionKey.split("::");
    const lessonId = parts[0];
    const sessionDate = parts[1];
    if (!lessonId || !sessionDate) continue;
    if (sessionDate < dateFrom || sessionDate > dateTo) continue;

    const meta = lessonMeta.get(lessonId);
    if (meta) {
      if (groupId && meta.groupId !== groupId) continue;
      if (teacherId && meta.teacherId !== teacherId) continue;
    }

    let hasAny = false;
    for (const [studentId, row] of Object.entries(entries)) {
      if (!rosterIds.has(studentId)) continue;
      hasAny = true;
      marksRecorded += 1;
      if (row.status === "present") present += 1;
      else if (row.status === "absent") absent += 1;
      else if (row.status === "late") late += 1;
    }
    if (hasAny) sessionsRecorded += 1;
  }

  const denom = present + absent + late;
  const attendanceRate = denom === 0 ? null : present / denom;

  const months = monthsBetweenInclusive(dateFrom, dateTo);
  const byBilling = emptyByBillingStatus();
  let totalExpected = 0;
  let totalPaid = 0;
  let totalOutstanding = 0;
  let overdueCount = 0;
  let receiptRows = 0;

  for (const periodMonth of months) {
    try {
      const res = await paymentsApi.listReceivables({
        periodMonth,
        status: "all",
        search: "",
      });
      for (const row of res.rows) {
        if (!rosterIds.has(row.studentId)) continue;
        if (groupId && row.groupId !== groupId) continue;
        receiptRows += 1;
        totalExpected += row.expectedAmount;
        totalPaid += row.paidAmount;
        if (row.balance > 0) totalOutstanding += row.balance;
        if (row.status === "overdue") overdueCount += 1;
        byBilling[row.status] += 1;
      }
    } catch {
      /* payments API may be unavailable */
    }
  }

  const periodLabel =
    months.length === 0
      ? "—"
      : months.length === 1
        ? months[0]!
        : `${months[0]} — ${months[months.length - 1]}`;

  return {
    generatedAt: new Date().toISOString(),
    currencyCode: "USD",
    source: "aggregated",
    students: {
      totalInRoster: roster.length,
      newEnrollmentsInRange,
      withoutGroup,
      byStatus: byStudentStatus,
    },
    attendance: {
      sessionsRecorded,
      marksRecorded,
      present,
      absent,
      late,
      attendanceRate,
    },
    payments: {
      periodLabel,
      monthsCounted: months.length,
      totalExpected,
      totalPaid,
      totalOutstanding,
      overdueCount,
      receiptRows,
      byBillingStatus: byBilling,
    },
    teachers: {
      total: teacherItems.length,
      byStatus: byTeacherStatus,
      lessonsOnSchedule,
      groupsRepresented: groupIdsFromLessons.size,
    },
  };
}
