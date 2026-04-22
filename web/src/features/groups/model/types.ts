export type GroupStatus = "draft" | "active" | "paused" | "completed" | "archived";

export type ScheduleSlot = {
  weekday: string;
  startTime: string;
  endTime: string;
  subjectName: string;
  roomName: string;
};

export interface Group {
  id: string;
  name: string;
  teacherId: string | null;
  subjectId: string | null;
  roomId: string | null;
  /** Whole currency units (e.g. monthly tuition). */
  monthlyFee: number;
  /** ISO date `YYYY-MM-DD`. */
  startDate: string;
  endDate: string | null;
  status: GroupStatus;
  schedulePreview: ScheduleSlot[];
  createdAt: string;
  updatedAt: string;
  teacher: { id: string; fullName: string } | null;
  subject: { id: string; name: string } | null;
  room: { id: string; name: string; capacity: number } | null;
}

export type GroupOption = { id: string; name: string };

export type GroupsListParams = {
  page: number;
  pageSize: number;
  search: string;
  status: GroupStatus | "all";
};

export type GroupsListResponse = {
  items: Group[];
  total: number;
  page: number;
  pageSize: number;
};

export type GroupCreatePayload = {
  name: string;
  teacherId?: string | null;
  subjectId?: string | null;
  roomId?: string | null;
  monthlyFee: number;
  startDate: string;
  endDate?: string | null;
  status: GroupStatus;
  schedulePreview?: ScheduleSlot[];
};

export type GroupUpdatePayload = Partial<GroupCreatePayload>;
