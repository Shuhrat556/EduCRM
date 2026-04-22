export const queryKeys = {
  root: ["educrm"] as const,
  auth: {
    me: () => [...queryKeys.root, "auth", "me"] as const,
  },
  dashboard: {
    overview: () => [...queryKeys.root, "dashboard", "overview"] as const,
  },
  students: {
    all: () => [...queryKeys.root, "students"] as const,
    lists: () => [...queryKeys.root, "students", "list"] as const,
    list: (params: {
      page: number;
      pageSize: number;
      search: string;
      status: string;
      groupId?: string;
    }) => [...queryKeys.students.lists(), params] as const,
    details: () => [...queryKeys.root, "students", "detail"] as const,
    detail: (id: string) =>
      [...queryKeys.students.details(), id] as const,
  },
  groups: {
    all: () => [...queryKeys.root, "groups"] as const,
    lists: () => [...queryKeys.groups.all(), "list"] as const,
    list: (params: {
      page: number;
      pageSize: number;
      search: string;
      status: string;
    }) => [...queryKeys.groups.lists(), params] as const,
    details: () => [...queryKeys.groups.all(), "detail"] as const,
    detail: (id: string) => [...queryKeys.groups.details(), id] as const,
    options: () => [...queryKeys.groups.all(), "options"] as const,
  },
  teachers: {
    all: () => [...queryKeys.root, "teachers"] as const,
    lists: () => [...queryKeys.root, "teachers", "list"] as const,
    list: (params: {
      page: number;
      pageSize: number;
      search: string;
      status: string;
    }) => [...queryKeys.teachers.lists(), params] as const,
    details: () => [...queryKeys.root, "teachers", "detail"] as const,
    detail: (id: string) => [...queryKeys.teachers.details(), id] as const,
  },
  subjects: {
    all: () => [...queryKeys.root, "subjects"] as const,
    lists: () => [...queryKeys.subjects.all(), "list"] as const,
    list: (params: { page: number; pageSize: number; search: string }) =>
      [...queryKeys.subjects.lists(), params] as const,
    details: () => [...queryKeys.subjects.all(), "detail"] as const,
    detail: (id: string) => [...queryKeys.subjects.details(), id] as const,
    options: () => [...queryKeys.subjects.all(), "options"] as const,
  },
  rooms: {
    all: () => [...queryKeys.root, "rooms"] as const,
    lists: () => [...queryKeys.rooms.all(), "list"] as const,
    list: (params: { page: number; pageSize: number; search: string }) =>
      [...queryKeys.rooms.lists(), params] as const,
    details: () => [...queryKeys.rooms.all(), "detail"] as const,
    detail: (id: string) => [...queryKeys.rooms.details(), id] as const,
  },
  schedule: {
    all: () => [...queryKeys.root, "schedule"] as const,
    lessons: () => [...queryKeys.schedule.all(), "lessons"] as const,
    lesson: (id: string) =>
      [...queryKeys.schedule.all(), "lesson", id] as const,
  },
  attendance: {
    all: () => [...queryKeys.root, "attendance"] as const,
    session: (lessonId: string, sessionDate: string) =>
      [...queryKeys.attendance.all(), "session", lessonId, sessionDate] as const,
  },
  reports: {
    all: () => [...queryKeys.root, "reports"] as const,
    summary: (filters: {
      dateFrom: string;
      dateTo: string;
      groupId: string | null;
      teacherId: string | null;
      studentStatus: string;
    }) => [...queryKeys.reports.all(), "summary", filters] as const,
  },
  aiAnalytics: {
    all: () => [...queryKeys.root, "ai-analytics"] as const,
    snapshot: () => [...queryKeys.aiAnalytics.all(), "snapshot"] as const,
  },
} as const;
