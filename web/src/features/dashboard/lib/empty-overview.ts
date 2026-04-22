import type { DashboardOverview } from "@/features/dashboard/model/types";

export const emptyDashboardOverview: DashboardOverview = {
  stats: {
    totalStudents: 0,
    totalTeachers: 0,
    activeGroups: 0,
    overduePayments: 0,
  },
  todayLessons: [],
  recentPayments: [],
  attendanceSummary: [],
  weeklyActivity: [],
};
