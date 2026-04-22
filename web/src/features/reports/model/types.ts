import type { BillingStatus } from "@/features/payments/model/types";
import type { StudentStatus } from "@/features/students/model/types";

export type ReportsFilters = {
  dateFrom: string;
  dateTo: string;
  groupId: string | null;
  teacherId: string | null;
  studentStatus: StudentStatus | "all";
};

export type StudentReportSection = {
  totalInRoster: number;
  newEnrollmentsInRange: number;
  withoutGroup: number;
  byStatus: Record<StudentStatus, number>;
};

export type AttendanceReportSection = {
  sessionsRecorded: number;
  marksRecorded: number;
  present: number;
  absent: number;
  late: number;
  /** Present / (present + absent + late), null when no marks. */
  attendanceRate: number | null;
};

export type PaymentReportSection = {
  periodLabel: string;
  monthsCounted: number;
  totalExpected: number;
  totalPaid: number;
  totalOutstanding: number;
  overdueCount: number;
  receiptRows: number;
  byBillingStatus: Record<BillingStatus, number>;
};

export type TeacherReportSection = {
  total: number;
  byStatus: Record<string, number>;
  lessonsOnSchedule: number;
  groupsRepresented: number;
};

export type ReportsSummary = {
  generatedAt: string;
  currencyCode: string;
  source: "api" | "aggregated";
  students: StudentReportSection;
  attendance: AttendanceReportSection;
  payments: PaymentReportSection;
  teachers: TeacherReportSection;
};
