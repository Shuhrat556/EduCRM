export type TeacherStatus = "active" | "inactive" | "on_leave";

export type NamedRef = { id: string; name: string };

export interface Teacher {
  id: string;
  fullName: string;
  phone: string;
  email: string;
  status: TeacherStatus;
  groupIds: string[];
  groups: NamedRef[];
  subjectIds: string[];
  subjects: NamedRef[];
  photoUrl: string | null;
  createdAt: string;
  updatedAt: string;
}

export type TeachersListParams = {
  page: number;
  pageSize: number;
  search: string;
  status: TeacherStatus | "all";
};

export type TeachersListResponse = {
  items: Teacher[];
  total: number;
  page: number;
  pageSize: number;
};

export type TeacherCreatePayload = {
  fullName: string;
  phone: string;
  email: string;
  status: TeacherStatus;
  groupIds: string[];
  subjectIds: string[];
  photoUrl?: string | null;
};

export type TeacherUpdatePayload = Partial<TeacherCreatePayload>;
