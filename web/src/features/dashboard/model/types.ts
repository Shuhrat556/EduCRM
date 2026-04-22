export type PaymentStatus = "paid" | "pending" | "failed" | "overdue";

export type DashboardStats = {
  totalStudents: number;
  totalTeachers: number;
  activeGroups: number;
  overduePayments: number;
};

export type TodayLesson = {
  id: string;
  startsAt: string;
  endsAt: string;
  subject: string;
  room: string;
  groupName: string;
  teacherName: string;
};

export type RecentPayment = {
  id: string;
  paidAt: string;
  payerName: string;
  amountCents: number;
  currency: string;
  status: PaymentStatus;
  reference?: string;
};

export type AttendanceSummaryRow = {
  id: string;
  date: string;
  presentCount: number;
  absentCount: number;
  lateCount: number;
  attendanceRate: number;
};

/** Simple series for CSS bar / future chart. */
export type ChartPoint = {
  label: string;
  value: number;
};

export type DashboardOverview = {
  stats: DashboardStats;
  todayLessons: TodayLesson[];
  recentPayments: RecentPayment[];
  attendanceSummary: AttendanceSummaryRow[];
  /** Optional trend for chart placeholder — enrollment or revenue. */
  weeklyActivity: ChartPoint[];
};
